mod events;
mod plugin;

use proc_macro::TokenStream;
use quote::quote;
use syn::{parse_macro_input, Attribute, DeriveInput};

use crate::{events::generate_event_subscriptions_inp, plugin::generate_plugin_impl};

#[proc_macro_derive(Plugin, attributes(plugin, events))]
pub fn handler_derive(input: TokenStream) -> TokenStream {
    let ast = parse_macro_input!(input as DeriveInput);

    let derive_name = &ast.ident;

    let info_attr = find_attribute(
        &ast,
        "plugin",
        "Missing `#[plugin(...)]` attribute with metadata.",
    );

    //  generate the code for impling Plugin.
    let plugin_impl = generate_plugin_impl(info_attr, derive_name);

    let subscription_attr = find_attribute(
        &ast,
        "events",
        "Missing #[events(...)] attribute. Please list the events to \
        subscribe to, e.g., #[events(Chat, PlayerJoin)]",
    );

    let event_subscriptions_impl = generate_event_subscriptions_inp(subscription_attr, derive_name);

    quote! {
        #plugin_impl
        #event_subscriptions_impl
    }
    .into()
}

fn find_attribute<'a>(ast: &'a syn::DeriveInput, name: &str, error: &str) -> &'a Attribute {
    ast.attrs
        .iter()
        .find(|a| a.path().is_ident(name))
        .ok_or_else(|| syn::Error::new(ast.ident.span(), error))
        .unwrap()
}
