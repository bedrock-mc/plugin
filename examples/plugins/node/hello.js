// Example Dragonfly plugin implemented in Node.js.
// Requires: npm install @grpc/grpc-js @grpc/proto-loader

import grpc from '@grpc/grpc-js';
import protoLoader from '@grpc/proto-loader';
import { fileURLToPath } from 'url';
import { dirname, resolve } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const PROTO_PATH = resolve(__dirname, '../../../proto/types/plugin.proto');
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});
const dfplugin = grpc.loadPackageDefinition(packageDefinition).df.plugin;

const pluginId = process.env.DF_PLUGIN_ID || 'node-plugin';
const serverAddress = process.env.DF_PLUGIN_SERVER_ADDRESS || '127.0.0.1:50050';
const API_VERSION = 'v1';

/**
 * IMPORTANT: All events MUST receive an eventResult response to avoid timeout warnings.
 * Even if your plugin doesn't modify or cancel an event, send an acknowledgment with cancel: false.
 */

// Connect to Dragonfly as a gRPC client
const client = new dfplugin.Plugin(
  serverAddress,
  grpc.credentials.createInsecure()
);

const call = client.EventStream();

console.log(`[node] connecting to ${serverAddress}...`);

call.write({
  pluginId: pluginId,
  hello: {
    name: 'example-node',
    version: '0.1.0',
    apiVersion: API_VERSION,
    commands: [
      { name: '/hello', description: 'Send a greeting from the Node plugin' },
    ],
  },
});

call.write({
  pluginId: pluginId,
  subscribe: { events: ['PLAYER_JOIN', 'COMMAND', 'CHAT'], observeEvents: [] },
});

call.on('data', (message) => {
  console.log('[node] received message:', JSON.stringify(message, null, 2));
  if (message.hello) {
    console.log('[node] host hello', message.hello);
    if (message.hello.apiVersion !== API_VERSION) {
      console.warn(`[node] API version mismatch: host=${message.hello.apiVersion}, plugin=${API_VERSION}`);
    }
    return;
  }

  if (message.event) {
    handleEvent(call, message.event);
  }
  if (message.shutdown) {
    console.log('[node] host shutdown:', message.shutdown.reason);
    call.end();
  }
});

call.on('end', () => {
  console.log('[node] stream ended');
  process.exit(0);
});

call.on('error', (err) => {
  console.error('[node] stream error:', err);
  process.exit(1);
});

process.on('SIGINT', () => {
  console.log('[node] shutting down...');
  call.end();
});

function handleEvent(call, event) {
  switch (event.type) {
    case 'PLAYER_JOIN': {
      const player = event.playerJoin;
      console.log(`[node] player joined ${player.name}`);
      call.write({
        pluginId: pluginId,
        actions: {
          actions: [
            {
              sendChat: {
                targetUuid: player.playerUuid,
                message: `Welcome to the server, ${player.name}! (from Node)`,
              },
            },
          ],
        },
      });
      // Acknowledge the event
      call.write({
        pluginId: pluginId,
        eventResult: {
          eventId: event.eventId,
          cancel: false,
        },
      });
      break;
    }
    case 'COMMAND': {
      const command = event.command;
      if (command.raw.startsWith('/hello')) {
        call.write({
          pluginId: pluginId,
          actions: {
            actions: [
              {
                sendChat: {
                  targetUuid: command.playerUuid,
                  message: `Hello ${command.name}!`,
                },
              },
            ],
          },
        });
      }
      // Always acknowledge command events
      call.write({
        pluginId: pluginId,
        eventResult: {
          eventId: event.eventId,
          cancel: false,
        },
      });
      break;
    }
    case 'CHAT': {
      const chat = event.chat;
      if (!chat) {
        break;
      }
      if (chat.message.toLowerCase().includes('badword')) {
        call.write({
          pluginId: pluginId,
          eventResult: {
            eventId: event.eventId,
            cancel: true,
          },
        });
        call.write({
          pluginId: pluginId,
          actions: {
            actions: [
              {
                sendChat: {
                  targetUuid: chat.playerUuid,
                  message: "Please keep the chat friendly!",
                },
              },
            ],
          },
        });
        break;
      }
      if (chat.message.startsWith('!shout ')) {
        const updated = chat.message.substring(7).toUpperCase();
        call.write({
          pluginId: pluginId,
          eventResult: {
            eventId: event.eventId,
            chat: { message: updated },
          },
        });
        break;
      }
      // Acknowledge regular chat messages
      call.write({
        pluginId: pluginId,
        eventResult: {
          eventId: event.eventId,
          cancel: false,
        },
      });
      break;
    }
    default:
      console.log('[node] event', event.type);
  }
}
