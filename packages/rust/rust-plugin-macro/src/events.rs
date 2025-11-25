use quote::quote;
use syn::{
    parse::{Parse, ParseStream},
    punctuated::Punctuated,
    Attribute, Ident, Token,
};

struct SubscriptionsListParser {
    events: Vec<Ident>,
}

impl Parse for SubscriptionsListParser {
    fn parse(input: ParseStream) -> syn::Result<Self> {
        let punctuated_list: Punctuated<Ident, Token![,]> = Punctuated::parse_terminated(input)?;

        let events = punctuated_list.into_iter().collect();

        Ok(Self { events })
    }
}

pub(crate) fn generate_event_subscriptions_inp(
    attr: &Attribute,
    derive_name: &Ident,
) -> proc_macro2::TokenStream {
    let subscriptions = match attr.parse_args::<SubscriptionsListParser>() {
        Ok(list) => list.events,
        Err(e) => {
            return e.to_compile_error();
        }
    };

    let subscription_variants = subscriptions.iter().map(|ident| {
        quote! { types::EventType::#ident }
    });

    quote! {
        impl dragonfly_plugin::EventSubscriptions for #derive_name {
            fn get_subscriptions(&self) -> Vec<types::EventType> {
                vec![
                    #( #subscription_variants ),*
                ]
            }
        }
    }
}
