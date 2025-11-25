use quote::quote;
use syn::{
    parse::{Parse, ParseStream},
    punctuated::Punctuated,
    Attribute, Expr, Ident, Lit, LitStr, Meta, Token,
};

struct PluginInfoParser {
    pub id: LitStr,
    pub name: LitStr,
    pub version: LitStr,
    pub api: LitStr,
}

impl Parse for PluginInfoParser {
    fn parse(input: ParseStream) -> syn::Result<Self> {
        let metas = Punctuated::<Meta, Token![,]>::parse_terminated(input)?;

        let mut id = None;
        let mut name = None;
        let mut version = None;
        let mut api = None;

        for meta in metas {
            match meta {
                Meta::NameValue(nv) => {
                    let key_ident = nv.path.get_ident().ok_or_else(|| {
                        syn::Error::new_spanned(&nv.path, "Expected an identifier (e.g., 'id')")
                    })?;

                    let value_str = match &nv.value {
                        Expr::Lit(expr_lit) => match &expr_lit.lit {
                            // clone out the reference as we will now
                            // be owning it and putting it into our generated code
                            Lit::Str(lit_str) => lit_str.clone(),
                            _ => {
                                return Err(syn::Error::new_spanned(
                                    &nv.value,
                                    "Expected a string literal",
                                ));
                            }
                        },
                        _ => {
                            return Err(syn::Error::new_spanned(
                                &nv.value,
                                "Expected a string literal",
                            ));
                        }
                    };

                    // Store the value
                    if key_ident == "id" {
                        id = Some(value_str);
                    } else if key_ident == "name" {
                        name = Some(value_str);
                    } else if key_ident == "version" {
                        version = Some(value_str);
                    } else if key_ident == "api" {
                        api = Some(value_str);
                    } else {
                        return Err(syn::Error::new_spanned(
                            key_ident,
                            "Unknown key. Expected 'id', 'name', 'version', or 'api'",
                        ));
                    }
                }
                _ => {
                    return Err(syn::Error::new_spanned(
                        meta,
                        "Expected `key = \"value\"` format",
                    ));
                }
            };
        }

        // Validate that all required fields were found
        // We use `input.span()` to point the error at the whole `#[plugin(...)]`
        // attribute if a field is missing.
        let id = id.ok_or_else(|| syn::Error::new(input.span(), "Missing required field 'id'"))?;
        let name =
            name.ok_or_else(|| syn::Error::new(input.span(), "Missing required field 'name'"))?;
        let version = version
            .ok_or_else(|| syn::Error::new(input.span(), "Missing required field 'version'"))?;
        let api =
            api.ok_or_else(|| syn::Error::new(input.span(), "Missing required field 'api'"))?;

        Ok(Self {
            id,
            name,
            version,
            api,
        })
    }
}

pub(crate) fn generate_plugin_impl(
    attr: &Attribute,
    derive_name: &Ident,
) -> proc_macro2::TokenStream {
    let plugin_info = match attr.parse_args::<PluginInfoParser>() {
        Ok(info) => info,
        Err(e) => return e.to_compile_error(),
    };

    // gotta define these outside because of quote rules with the . access.
    let id_lit = &plugin_info.id;
    let name_lit = &plugin_info.name;
    let version_lit = &plugin_info.version;
    let api_lit = &plugin_info.api;

    quote! {
        impl dragonfly_plugin::Plugin for #derive_name {
            fn get_info(&self) -> dragonfly_plugin::PluginInfo<'static> {
                dragonfly_plugin::PluginInfo::<'static> {
                    id: #id_lit,
                    name: #name_lit,
                    version: #version_lit,
                    api_version: #api_lit
                }
            }
            fn get_id(&self) -> &'static str { #id_lit }
            fn get_name(&self) -> &'static str { #name_lit }
            fn get_version(&self) -> &'static str { #version_lit }
            fn get_api_version(&self) -> &'static str { #api_lit }
        }
    }
}
