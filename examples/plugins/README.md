# Dragonfly Plugin Examples

This directory contains example plugins demonstrating how to create Dragonfly plugins in different languages.

## Why Dragonfly's Plugin System?

| Benefit | Description | Use Case |
|---------|-------------|----------|
| 🌍 **Any Language** | Write plugins in JavaScript, TypeScript, PHP, Python, Rust, C++, or any language with gRPC support | Use the language your team knows best |
| 💰 **Sell Plugins** | Compile to binary (Rust, Go, C++) and distribute without source code | Create commercial plugins |
| 🔥 **Hot Reload** | Edit JS/TS/PHP plugins and see changes instantly - no server restart needed | Develop and debug plugins in real-time |
| 📱 **Remote Control** | Plugins connect over gRPC - run them anywhere (phone app, cloud service, another server) | Build mobile admin apps |
| 📦 **Use Any Library** | Import npm packages on a Go server, use Python ML libraries, etc. | Leverage entire ecosystems |
| ⚡ **Zero Performance Impact** | Plugins run in separate processes - slow/heavy plugin code doesn't affect server TPS | Run intensive tasks without lag |
| 🚀 **High Performance** | Optimized protobuf protocol with optional batching for low latency | Handle 100+ players with movement events |
| 🔒 **Sandboxing** | Control what plugins can access via gRPC permissions | Host untrusted plugins safely |

### Real-World Examples

```bash
# Hot reload: Edit plugin code while server is running
vim plugins/my-plugin.js   # Make changes
# Changes apply immediately - no restart!

# Remote plugin: Control server from your phone
# Plugin runs on your phone, connects to server over internet
phone-app → [gRPC] → Dragonfly Server

# Binary plugin: Sell without source code
rustc plugin.rs --release   # Compile to binary
# Distribute the binary - customers can't see your code
```

## Available Examples

### 1. Node.js Plugin (`node/`)

Simple JavaScript plugin using `@grpc/grpc-js` and `@grpc/proto-loader`.

```bash
cd node/
npm install
```

---

### 2. TypeScript Plugin (`typescript/`)

Type-safe plugin with generated types

```bash
cd typescript/
npm install
npm run generate
```

---

### 3. PHP Plugin (`php/`)

PHP plugin using gRPC extension.

```bash
cd php/
# Requires: php-grpc extension installed
php HelloPlugin.php
```

**Features:**
- ✅ Use existing PHP libraries
- ⚠️ Requires gRPC extension

---
## Quick Start

1. **Choose a language** based on your needs (TypeScript recommended for production)
2. **Follow the setup** in that example's directory
3. **Enable in config** - Edit `plugins/plugins.yaml`:
   ```yaml
   plugins:
     - id: my-plugin
       name: My Plugin
       command: "node"
       args: ["examples/plugins/typescript/dist/index.js"]
       address: "127.0.0.1:50051"
   ```
4. **Run Dragonfly** - The plugin will connect automatically

## Plugin Configuration

Edit `plugins/plugins.yaml` to enable/configure plugins:

```yaml
plugins:
  # Node.js example
  - id: example-node
    name: Example Node Plugin
    command: "node"
    args: ["examples/plugins/node/hello.js"]
    address: "127.0.0.1:50051"
    env:
      NODE_ENV: development

  # TypeScript example
  - id: example-typescript
    name: Example TypeScript Plugin
    command: "node"
    args: ["examples/plugins/typescript/dist/index.js"]
    address: "127.0.0.1:50052"

  # PHP example
  - id: example-php
    name: Example PHP Plugin
    command: "php"
    args: ["examples/plugins/php/HelloPlugin.php"]
    address: "127.0.0.1:50053"
```

## Protocol Documentation

All plugins communicate using the same protobuf protocol defined in `plugin/proto/plugin.proto`.

**Key concepts:**

### 1. Bidirectional Stream

Plugins use a single bidirectional gRPC stream for all communication:

```
Host ←→ Plugin (EventStream)
```

### 2. Message Types

**Host → Plugin:**
- `HostHello` - Initial handshake
- `EventEnvelope` - Game events (join, quit, chat, commands, etc.)
- `HostShutdown` - Server shutting down

**Plugin → Host:**
- `PluginHello` - Register plugin capabilities
- `EventSubscribe` - Subscribe to specific events
- `ActionBatch` - Execute actions (teleport, chat, kick, etc.)
- `EventResult` - Cancel or mutate events

### 3. Event Flow

```
1. Host sends HostHello
2. Plugin responds with PluginHello
3. Plugin sends EventSubscribe (which events to receive)
4. Host sends events as they occur
5. Plugin can respond with:
   - Actions (do something)
   - EventResult (cancel/modify event)
```

### 4. Example Event Types

- `PLAYER_JOIN` - Player connected
- `PLAYER_QUIT` - Player disconnected
- `CHAT` - Player sent chat message
- `COMMAND` - Player executed command
- `BLOCK_BREAK` - Player broke a block
- `WORLD_CLOSE` - World is closing

### 5. Example Actions

- `SendChatAction` - Send message to player
- `TeleportAction` - Teleport player
- `KickAction` - Kick player

## Creating Your Own Plugin

### Minimal Plugin (Node.js)

```javascript
import grpc from '@grpc/grpc-js';
import protoLoader from '@grpc/proto-loader';

const packageDef = protoLoader.loadSync('plugin/proto/plugin.proto');
const proto = grpc.loadPackageDefinition(packageDef).df.plugin;

const server = new grpc.Server();
server.addService(proto.Plugin.service, {
  EventStream: (call) => {
    call.on('data', (msg) => {
      if (msg.hello) {
        call.write({
          pluginId: 'my-plugin',
          hello: {
            name: 'My Plugin',
            version: '1.0.0',
            apiVersion: msg.hello.apiVersion,
          }
        });
      }
    });
  }
});

server.bindAsync('127.0.0.1:50051', 
  grpc.ServerCredentials.createInsecure(), 
  () => console.log('Plugin ready'));
```
