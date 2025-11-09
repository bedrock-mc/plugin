#!/bin/bash
set -e

# Generate PHP protobuf files from proto definitions
# Requires: protoc, protoc-gen-grpc-php plugin

PROTO_DIR="../../../plugin/proto/types"
OUT_DIR="./generated"

echo "[php] Generating protobuf PHP files..."

# Create output directory
mkdir -p "$OUT_DIR"

# Generate PHP protobuf and gRPC files
protoc \
  --proto_path="$PROTO_DIR" \
  --php_out="$OUT_DIR" \
  --grpc_out="$OUT_DIR" \
  --plugin=protoc-gen-grpc=/usr/local/bin/grpc_php_plugin \
  "$PROTO_DIR/plugin.proto"

echo "[php] Protobuf generation complete!"
echo "[php] Generated files are in: $OUT_DIR"

