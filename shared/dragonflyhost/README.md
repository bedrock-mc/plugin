# dragonflyhost

Transport-neutral helpers for reading Dragonfly host state.

This package is shared by host runtimes that need to translate Dragonfly players, worlds, items, and damage sources into their own protocol or language bridge. It intentionally does not import protobuf, gRPC, PocketMine, or any plugin runtime package.

Keep this package limited to:

- safe player/world snapshot reads
- neutral item, inventory, position, world, and damage structs
- small Dragonfly enum/string normalization helpers

Do not put runtime behavior, IPC, protobuf messages, PMMP classes, command registration, or action application here. Callers should adapt these neutral snapshots into their own transport payloads.
