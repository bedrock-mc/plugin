use heck::ToPascalCase;
use proc_macro::TokenStream;
use quote::{format_ident, quote};
use syn::{ImplItem, ItemImpl, Type, parse_macro_input};

#[proc_macro_attribute]
pub fn bedrock_plugin(_attr: TokenStream, item: TokenStream) -> TokenStream {
    let input = parse_macro_input!(item as ItemImpl);
    let original_impl = input.clone();

    let self_ty: &Type = &input.self_ty;

    let trait_path = match &input.trait_ {
        Some((_, path, _)) => path,
        None => {
            let msg = "The #[bedrock_plugin] attribute can only be used on an `impl PluginEventHandler for ...` block.";
            return syn::Error::new_spanned(&input.self_ty, msg)
                .to_compile_error()
                .into();
        }
    };

    if trait_path
        .segments
        .last()
        .is_some_and(|s| s.ident != "PluginEventHandler")
    {
        let msg = "The #[bedrock_plugin] attribute must be on an `impl PluginEventHandler for ...` block.";
        return syn::Error::new_spanned(trait_path, msg)
            .to_compile_error()
            .into();
    }

    let mut subscriptions = Vec::new();
    for item in &input.items {
        if let ImplItem::Fn(method) = item {
            let fn_name_str = method.sig.ident.to_string();

            if let Some(event_name_snake) = fn_name_str.strip_prefix("on_") {
                let event_name_pascal = event_name_snake.to_pascal_case();

                let event_type_ident = format_ident!("{}", event_name_pascal);

                subscriptions.push(quote! { types::EventType::#event_type_ident });
            }
        }
    }

    let subscription_impl = quote! {
        impl PluginSubscriptions for #self_ty {
            fn get_subscriptions(&self) -> Vec<types::EventType> {
                vec![
                    #( #subscriptions ),*
                ]
            }
        }
    };

    let output = quote! {
        #original_impl
        #subscription_impl
    };

    output.into()
}
