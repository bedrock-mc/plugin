use heck::ToSnakeCase;
use proc_macro2::TokenStream;
use quote::{format_ident, quote};
use syn::{
    parse::{Parse, ParseStream},
    punctuated::Punctuated,
    spanned::Spanned,
    Attribute, Data, DeriveInput, Expr, ExprLit, Fields, GenericArgument, Ident, Lit, LitStr, Meta,
    PathArguments, Token, Type, TypePath,
};

struct CommandInfoParser {
    pub name: LitStr,
    pub description: LitStr,
    pub aliases: Vec<LitStr>,
}

impl Parse for CommandInfoParser {
    fn parse(input: ParseStream) -> syn::Result<Self> {
        let metas = Punctuated::<Meta, Token![,]>::parse_terminated(input)?;

        let mut name = None;
        let mut description = None;
        let mut aliases = Vec::new();

        for meta in metas {
            match meta {
                Meta::NameValue(nv) if nv.path.is_ident("name") => {
                    if let Expr::Lit(ExprLit {
                        lit: Lit::Str(s), ..
                    }) = nv.value
                    {
                        name = Some(s);
                    } else {
                        return Err(syn::Error::new(
                            nv.value.span(),
                            "expected string literal for `name`",
                        ));
                    }
                }
                Meta::NameValue(nv) if nv.path.is_ident("description") => {
                    if let Expr::Lit(ExprLit {
                        lit: Lit::Str(s), ..
                    }) = nv.value
                    {
                        description = Some(s);
                    } else {
                        return Err(syn::Error::new(
                            nv.value.span(),
                            "expected string literal for `description`",
                        ));
                    }
                }
                Meta::List(list) if list.path.is_ident("aliases") => {
                    aliases = list
                        .parse_args_with(Punctuated::<LitStr, Token![,]>::parse_terminated)?
                        .into_iter()
                        .collect();
                }
                _ => {
                    return Err(syn::Error::new(
                        meta.span(),
                        "unrecognized command attribute",
                    ));
                }
            }
        }

        Ok(Self {
            name: name.ok_or_else(|| {
                syn::Error::new(input.span(), "missing required attribute `name`")
            })?,
            description: description.unwrap_or_else(|| LitStr::new("", input.span())),
            aliases,
        })
    }
}

#[derive(Default)]
struct SubcommandAttr {
    pub name: Option<LitStr>,
    pub aliases: Vec<LitStr>,
}

fn parse_subcommand_attr(attrs: &[Attribute]) -> syn::Result<SubcommandAttr> {
    let mut out = SubcommandAttr::default();

    for attr in attrs {
        if !attr.path().is_ident("subcommand") {
            continue;
        }

        let metas = attr.parse_args_with(Punctuated::<Meta, Token![,]>::parse_terminated)?;

        for meta in metas {
            match meta {
                // name = "pay"
                Meta::NameValue(nv) if nv.path.is_ident("name") => {
                    if let Expr::Lit(ExprLit {
                        lit: Lit::Str(s), ..
                    }) = nv.value
                    {
                        out.name = Some(s);
                    } else {
                        return Err(syn::Error::new_spanned(
                            nv.value,
                            "subcommand `name` must be a string literal",
                        ));
                    }
                }

                // aliases("give", "send")
                Meta::List(list) if list.path.is_ident("aliases") => {
                    out.aliases = list
                        .parse_args_with(Punctuated::<LitStr, Token![,]>::parse_terminated)?
                        .into_iter()
                        .collect();
                }

                _ => {
                    return Err(syn::Error::new_spanned(
                        meta,
                        "unknown subcommand attribute; expected `name = \"...\"` or `aliases(..)`",
                    ));
                }
            }
        }
    }

    Ok(out)
}

pub fn generate_command_impls(ast: &DeriveInput, attr: &Attribute) -> TokenStream {
    let command_info = match attr.parse_args::<CommandInfoParser>() {
        Ok(info) => info,
        Err(e) => return e.to_compile_error(),
    };

    let cmd_ident = &ast.ident;
    let cmd_name_lit = &command_info.name;
    let cmd_desc_lit = &command_info.description;
    let aliases_lits = &command_info.aliases;

    let shape = match collect_command_shape(ast) {
        Ok(s) => s,
        Err(e) => return e.to_compile_error(),
    };

    let spec_impl = generate_spec_impl(cmd_ident, cmd_name_lit, cmd_desc_lit, aliases_lits, &shape);
    let try_from_impl = generate_try_from_impl(cmd_ident, cmd_name_lit, aliases_lits, &shape);
    let trait_impl = generate_handler_trait(cmd_ident, &shape);
    let exec_impl = generate_execute_impl(cmd_ident, &shape);

    quote! {
        #spec_impl
        #try_from_impl
        #trait_impl
        #exec_impl
    }
}

fn generate_spec_impl(
    cmd_ident: &Ident,
    cmd_name_lit: &LitStr,
    cmd_desc_lit: &LitStr,
    aliases_lits: &[LitStr],
    shape: &CommandShape,
) -> TokenStream {
    let params_tokens = match shape {
        CommandShape::Struct { params } => {
            let specs = params.iter().map(param_to_spec);
            quote! { vec![ #( #specs ),* ] }
        }
        CommandShape::Enum { variants } => {
            // First param: subcommand enum
            let variant_names_iter = variants.iter().map(|v| &v.canonical);
            let subcommand_names: Vec<&LitStr> = variants
                .iter()
                .flat_map(|v| &v.aliases)
                .chain(variant_names_iter)
                .collect();
            let subcommand_spec = quote! {
                dragonfly_plugin::types::ParamSpec {
                    name: "subcommand".to_string(),
                    r#type: dragonfly_plugin::types::ParamType::ParamEnum as i32,
                    optional: false,
                    suffix: String::new(),
                    enum_values: vec![ #( #subcommand_names.to_string() ),* ],
                }
            };

            // Merge all variant params, marking as optional if not present in all variants
            let merged = merge_variant_params(variants);
            let merged_specs = merged.iter().map(param_to_spec);

            quote! { vec![ #subcommand_spec, #( #merged_specs ),* ] }
        }
    };

    quote! {
        impl #cmd_ident {
            pub fn spec() -> dragonfly_plugin::types::CommandSpec {
                dragonfly_plugin::types::CommandSpec {
                    name: #cmd_name_lit.to_string(),
                    description: #cmd_desc_lit.to_string(),
                    aliases: vec![ #( #aliases_lits.to_string() ),* ],
                    params: #params_tokens,
                }
            }
        }
    }
}

fn param_to_spec(p: &ParamMeta) -> TokenStream {
    let name_str = p.field_ident.to_string();
    let param_type_expr = &p.param_type_expr;
    let is_optional = p.is_optional;

    quote! {
        dragonfly_plugin::types::ParamSpec {
            name: #name_str.to_string(),
            r#type: #param_type_expr as i32,
            optional: #is_optional,
            suffix: String::new(),
            enum_values: Vec::new(),
        }
    }
}

/// Merge params from all variants: a param is optional if not present in every variant.
fn merge_variant_params(variants: &[VariantMeta]) -> Vec<ParamMeta> {
    use std::collections::HashMap;

    // Collect all unique param names with their metadata
    let mut seen: HashMap<String, (ParamMeta, usize)> = HashMap::new(); // name -> (meta, count)

    for variant in variants {
        for param in &variant.params {
            let name = param.field_ident.to_string();
            seen.entry(name)
                .and_modify(|(_, count)| *count += 1)
                .or_insert_with(|| (param.clone(), 1));
        }
    }

    let variant_count = variants.len();

    seen.into_iter()
        .map(|(_, (mut meta, count))| {
            // If not present in all variants, mark as optional
            if count < variant_count {
                meta.is_optional = true;
            }
            meta
        })
        .collect()
}

fn generate_try_from_impl(
    cmd_ident: &Ident,
    cmd_name_lit: &LitStr,
    cmd_aliases: &[LitStr],
    shape: &CommandShape,
) -> TokenStream {
    let body = match shape {
        CommandShape::Struct { params } => {
            let field_inits = params.iter().map(|p| struct_field_init(p, 0));
            quote! {
                Ok(Self {
                    #( #field_inits, )*
                })
            }
        }
        CommandShape::Enum { variants } => {
            let match_arms = variants.iter().map(|v| {
                let variant_ident = &v.ident;
                let params = &v.params;

                // Build patterns: "canonical" | "alias1" | "alias2"
                let mut name_lits = Vec::new();
                name_lits.push(&v.canonical);
                for alias in &v.aliases {
                    name_lits.push(alias);
                }

                let subcommand_patterns = quote! { #( #name_lits )|* };

                if params.is_empty() {
                    quote! {
                        #subcommand_patterns => Ok(Self::#variant_ident),
                    }
                } else {
                    let field_inits = params.iter().map(enum_field_init);
                    quote! {
                        #subcommand_patterns => Ok(Self::#variant_ident {
                            #( #field_inits, )*
                        }),
                    }
                }
            });

            quote! {
                let subcommand = event.args.first()
                    .ok_or(dragonfly_plugin::command::CommandParseError::Missing("subcommand"))?
                    .as_str();

                match subcommand {
                    #( #match_arms )*
                    _ => Err(dragonfly_plugin::command::CommandParseError::UnknownSubcommand),
                }
            }
        }
    };

    let mut conditions = Vec::with_capacity(1 + cmd_aliases.len());
    conditions.push(quote! { event.command != #cmd_name_lit });

    for alias in cmd_aliases {
        conditions.push(quote! { && event.command != #alias });
    }

    quote! {
        impl ::core::convert::TryFrom<&dragonfly_plugin::types::CommandEvent> for #cmd_ident {
            type Error = dragonfly_plugin::command::CommandParseError;

            fn try_from(event: &dragonfly_plugin::types::CommandEvent) -> Result<Self, Self::Error> {
                if #(#conditions)* {
                    return Err(dragonfly_plugin::command::CommandParseError::NoMatch);
                }

                #body
            }
        }
    }
}

/// Generate field init for struct commands (args start at index 0).
fn struct_field_init(p: &ParamMeta, offset: usize) -> TokenStream {
    let ident = &p.field_ident;
    let idx = p.index + offset;
    let name_str = ident.to_string();
    let ty = &p.field_ty;

    if p.is_optional {
        quote! {
            #ident: dragonfly_plugin::command::parse_optional_arg::<#ty>(&event.args, #idx, #name_str)?
        }
    } else {
        quote! {
            #ident: dragonfly_plugin::command::parse_required_arg::<#ty>(&event.args, #idx, #name_str)?
        }
    }
}

/// Generate field init for enum variant (args start at index 1, after subcommand).
fn enum_field_init(p: &ParamMeta) -> TokenStream {
    let ident = &p.field_ident;
    let idx = p.index + 1; // +1 because index 0 is the subcommand
    let name_str = ident.to_string();
    let ty = &p.field_ty;

    if p.is_optional {
        quote! {
            #ident: dragonfly_plugin::command::parse_optional_arg::<#ty>(&event.args, #idx, #name_str)?
        }
    } else {
        quote! {
            #ident: dragonfly_plugin::command::parse_required_arg::<#ty>(&event.args, #idx, #name_str)?
        }
    }
}

/// Metadata for a single command parameter (struct field or enum variant field).
#[derive(Clone)]
struct ParamMeta {
    /// Rust field identifier, e.g. `amount`
    field_ident: Ident,
    /// Full Rust type of the field, e.g. `f64` or `Option<String>`
    field_ty: Type,
    /// Expression for ParamType (e.g. `ParamType::ParamFloat`)
    param_type_expr: TokenStream,
    /// Whether this is optional in the spec
    is_optional: bool,
    /// Position index in args (0-based, relative to this variant/struct)
    index: usize,
}

/// Metadata for a single enum variant (subcommand).
struct VariantMeta {
    /// Variant identifier, e.g. `Pay`
    ident: Ident,
    /// name for matching, e.g. `"pay"`
    canonical: LitStr,
    /// Aliases e.g give, donate.
    aliases: Vec<LitStr>,
    /// Parameters (fields) for this variant
    params: Vec<ParamMeta>,
}

/// The shape of a command: either a struct or an enum with variants.
enum CommandShape {
    Struct { params: Vec<ParamMeta> },
    Enum { variants: Vec<VariantMeta> },
}
fn collect_command_shape(ast: &DeriveInput) -> syn::Result<CommandShape> {
    match &ast.data {
        Data::Struct(data) => {
            let params = collect_params_from_fields(&data.fields)?;
            Ok(CommandShape::Struct { params })
        }
        Data::Enum(data) => {
            let mut variants_meta = Vec::new();

            for variant in &data.variants {
                let ident = variant.ident.clone();
                let default_name =
                    LitStr::new(ident.to_string().to_snake_case().as_str(), ident.span());

                let sub_attr = parse_subcommand_attr(&variant.attrs)?;
                let canonical = sub_attr.name.unwrap_or(default_name);

                let aliases = sub_attr.aliases.into_iter().collect::<Vec<_>>();

                let params = collect_params_from_fields(&variant.fields)?;

                variants_meta.push(VariantMeta {
                    ident,
                    canonical,
                    aliases,
                    params,
                });
            }

            if variants_meta.is_empty() {
                return Err(syn::Error::new_spanned(
                    &ast.ident,
                    "enum commands must have at least one variant",
                ));
            }

            Ok(CommandShape::Enum {
                variants: variants_meta,
            })
        }
        Data::Union(_) => Err(syn::Error::new_spanned(
            ast,
            "unions are not supported for #[derive(Command)]",
        )),
    }
}

fn collect_params_from_fields(fields: &Fields) -> syn::Result<Vec<ParamMeta>> {
    let fields = match fields {
        Fields::Named(named) => &named.named,
        Fields::Unit => {
            // No params
            return Ok(Vec::new());
        }
        Fields::Unnamed(_) => {
            // TODO: maybe do support this but typenames are used as param names?
            return Err(syn::Error::new_spanned(
                fields,
                "tuple structs are not supported for commands; use named fields",
            ));
        }
    };

    let mut out = Vec::new();
    for (index, field) in fields.iter().enumerate() {
        let field_ident = field
            .ident
            .clone()
            .expect("command struct fields must be named");
        let field_ty = field.ty.clone();

        let (param_type_expr, is_optional) = get_param_type(&field_ty);

        out.push(ParamMeta {
            field_ident,
            field_ty,
            param_type_expr,
            is_optional,
            index,
        });
    }

    Ok(out)
}

fn get_param_type(ty: &Type) -> (TokenStream, bool) {
    if let Some(inner) = option_inner(ty) {
        let (inner_param, _inner_opt) = get_param_type(inner);
        return (inner_param, true);
    }

    if let Type::Reference(r) = ty {
        return get_param_type(&r.elem);
    }

    if let Type::Path(TypePath { path, .. }) = ty {
        if let Some(seg) = path.segments.last() {
            let ident = seg.ident.to_string();

            // Floats
            if ident == "f32" || ident == "f64" {
                return (
                    quote! { dragonfly_plugin::types::ParamType::ParamFloat },
                    false,
                );
            }

            if matches!(
                ident.as_str(),
                "i8" | "i16"
                    | "i32"
                    | "i64"
                    | "i128"
                    | "u8"
                    | "u16"
                    | "u32"
                    | "u64"
                    | "u128"
                    | "isize"
                    | "usize"
            ) {
                return (
                    quote! { dragonfly_plugin::types::ParamType::ParamInt },
                    false,
                );
            }

            if ident == "bool" {
                return (
                    quote! { dragonfly_plugin::types::ParamType::ParamBool },
                    false,
                );
            }

            if ident == "String" {
                return (
                    quote! { dragonfly_plugin::types::ParamType::ParamString },
                    false,
                );
            }
        }
    }

    (
        quote! { dragonfly_plugin::types::ParamType::ParamString },
        false,
    )
}

fn option_inner(ty: &Type) -> Option<&Type> {
    if let Type::Path(TypePath { path, .. }) = ty {
        if let Some(seg) = path.segments.last() {
            if seg.ident == "Option" {
                if let PathArguments::AngleBracketed(args) = &seg.arguments {
                    for arg in &args.args {
                        if let GenericArgument::Type(inner) = arg {
                            return Some(inner);
                        }
                    }
                }
            }
        }
    }
    None
}

fn generate_handler_trait(cmd_ident: &Ident, shape: &CommandShape) -> TokenStream {
    let trait_ident = format_ident!("{}Handler", cmd_ident);

    match shape {
        CommandShape::Struct { params } => {
            // method name = struct name in snake_case, e.g. Ping -> ping
            let method_ident = format_ident!("{}", cmd_ident.to_string().to_snake_case());

            // args from struct fields if you want
            let args = params.iter().map(|p| {
                let ident = &p.field_ident;
                let ty = &p.field_ty;
                quote! { #ident: #ty }
            });

            quote! {
                #[allow(async_fn_in_trait)]
                pub trait #trait_ident: Send + Sync {
                    async fn #method_ident(
                        &self,
                        ctx: dragonfly_plugin::command::Ctx<'_>,
                        #( #args ),*
                    );
                }
            }
        }
        CommandShape::Enum { variants } => {
            let methods = variants.iter().map(|v| {
                let method_ident = format_ident!("{}", v.ident.to_string().to_snake_case());
                let args = v.params.iter().map(|p| {
                    let ident = &p.field_ident;
                    let ty = &p.field_ty;
                    quote! { #ident: #ty }
                });
                quote! {
                    async fn #method_ident(
                        &self,
                        ctx: dragonfly_plugin::command::Ctx<'_>,
                        #( #args ),*
                    );
                }
            });

            quote! {
                #[allow(async_fn_in_trait)]
                pub trait #trait_ident: Send + Sync {
                    #( #methods )*
                }
            }
        }
    }
}

fn generate_execute_impl(cmd_ident: &Ident, shape: &CommandShape) -> TokenStream {
    let trait_ident = format_ident!("{}Handler", cmd_ident);

    match shape {
        CommandShape::Struct { params } => {
            let method_ident = format_ident!("{}", cmd_ident.to_string().to_snake_case());
            let field_idents: Vec<_> = params.iter().map(|p| &p.field_ident).collect();

            quote! {
                impl #cmd_ident {
                    pub async fn __execute<H: #trait_ident>(
                        self,
                        handler: &H,
                        ctx: dragonfly_plugin::command::Ctx<'_>,
                    ) {
                        let Self { #( #field_idents ),* } = self;
                        handler.#method_ident(ctx, #( #field_idents ),*).await;
                    }
                }
            }
        }
        CommandShape::Enum { variants } => {
            let match_arms = variants.iter().map(|v| {
                let variant_ident = &v.ident;
                let method_ident = format_ident!("{}", v.ident.to_string().to_snake_case());
                let field_idents: Vec<_> = v.params.iter().map(|p| &p.field_ident).collect();

                if field_idents.is_empty() {
                    quote! {
                        Self::#variant_ident => handler.#method_ident(ctx).await,
                    }
                } else {
                    quote! {
                        Self::#variant_ident { #( #field_idents ),* } => {
                            handler.#method_ident(ctx, #( #field_idents ),*).await
                        }
                    }
                }
            });

            quote! {
                impl #cmd_ident {
                    pub async fn __execute<H: #trait_ident>(
                        self,
                        handler: &H,
                        ctx: dragonfly_plugin::command::Ctx<'_>,
                    ) {
                        match self {
                            #( #match_arms )*
                        }
                    }
                }
            }
        }
    }
}
