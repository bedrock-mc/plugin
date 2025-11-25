# Dragonfly Rust Plugin API

Welcome to the Rust Plugin API for Dragonfly server software. This library provides the tools to build high-performance, safe, and asynchronous plugins using native Rust.

The API is designed to be simple and explicit. You define your plugin's logic by implementing the `PluginEventHandler` trait and register your event subscriptions using the `#[derive(Handler)]` macro.

## Features

- **Asynchronous by Default:** Built on `tokio` and native Rust `async/await`.
- **Type-Safe:** All events and actions are strongly typed, catching bugs at compile time.
- **Simple Subscriptions:** A clean `#[derive(Handler)]` macro handles all event subscription logic.

---

## Quick Start Guide

Here is the complete process for creating a "Hello, World\!" plugin.

### 1\. Create a New Plugin

First, create a new binary crate for your plugin:

```sh
cargo new my_plugin --bin
```

### 2\. Update `Cargo.toml`

Next, add `dragonfly-plugin` and `tokio` to your dependencies.

```toml
[package]
name = "my_plugin"
version = "0.1.0"
edition = "2021"

[dependencies]
# This is the main API library
dragonfly-plugin = "0.1" # Or use a version number

# Tokio is required for the async runtime
tokio = { version = "1", features = ["full"] }
```

### 3\. Write Your Plugin (`src/main.rs`)

This is the complete code for a simple plugin that greets players on join and adds a prefix to their chat messages.

```rust
// --- Import all the necessary items ---
use dragonfly_plugin::{
    event::{EventContext, PluginEventHandler},
    types, // Contains all event data structs
    Handler, // The derive macro
    Plugin,  // The main plugin runner
    Server,
};

// --- 1. Define Your Plugin Struct ---
//
// `#[derive(Handler)]` is the "trigger" that runs the macro.
// `#[derive(Default)]` is common for simple, stateless plugins.
#[derive(Handler, Default)]
//
// `#[subscriptions(...)]` is the "helper attribute" that lists
// all the events this plugin needs to listen for.
#[subscriptions(PlayerJoin, Chat)]
struct MyPlugin;

// --- 2. Implement the Event Handlers ---
//
// This is where all your plugin's logic lives.
// You only need to implement the `async fn` handlers
// for the events you subscribed to.
impl PluginEventHandler for MyPlugin {
    /// This handler runs when a player joins the server.
    async fn on_player_join(
        &self,
        server: &Server,
        event: &mut EventContext<'_, types::PlayerJoinEvent>,
    ) {
        let player_name = &event.data.name;
        println!("Player '{}' has joined.", player_name);

        let welcome_message = format!(
            "Welcome, {}! This server is running a Rust plugin.",
            player_name
        );

        // Use the `server` handle to send actions
        server.send_chat(event.data.player_uuid.clone(), welcome_message).await.ok();
    }

    /// This handler runs when a player sends a chat message.
    async fn on_chat(
        &self,
        _server: &Server, // We don't need the server for this
        event: &mut EventContext<'_, types::ChatEvent>,
    ) {
        let new_message = format!("[Plugin] {}", event.data.message);

        // Use helper methods on the `event` to mutate it
        event.set_message(new_message);
    }
}

// --- 3. Start the Plugin ---
//
// This is the entry point that connects to the server.
#[tokio::main]
async fn main() {
    println!("Starting my-plugin...");

    // Create an instance of your plugin
    let plugin = MyPlugin::default();

    // Run the plugin. This will connect to the server
    // and block forever, processing events.
    Plugin::run(plugin, "127.0.0.1:50051")
        .await
        .expect("Plugin failed to run");
}
```

---

## Writing Event Handlers

All plugin logic is built by implementing functions from the `PluginEventHandler` trait.

### `async fn` and Lifetimes

Because this API uses native `async fn` in traits, you **must** include the anonymous lifetime (`'_`) annotation in the `EventContext` type:

- **Correct:** `event: &mut EventContext<'_, types::ChatEvent>`
- **Incorrect:** `event: &mut EventContext<types::ChatEvent>`

### Reading Event Data

You can read immutable data directly from the event:

```rust,ignore
let player_name = &event.data.name;
println!("Player name: {}", player_name);
```

### Mutating Events

Some events are mutable. The `EventContext` provides helper methods (like `set_message`) to modify the event before the server processes it:

```rust,ignore
event.set_message(format!("New message: {}", event.data.message));
```

### Cancelling Events

You can also cancel compatible events to stop them from happening:

```rust,ignore
event.cancel();
```

## Available Events

The `#[subscriptions(...)]` macro accepts any variant from the `types::EventType` enum.

You can find a complete list of all available event names (e.g., `PlayerJoin`, `Chat`, `BlockBreak`) and their corresponding data structs (e.g., `types::PlayerJoinEvent`) by looking in the `./src/types.rs` file or by consulting the API documentation.
