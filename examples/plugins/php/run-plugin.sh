#!/bin/bash
# Run the PHP plugin using custom PHP installation with built-in gRPC

# Detect PHP binary location (differs on Windows)
if [ -f "bin/php/php.exe" ]; then
    # Windows
    PHP_BIN="bin/php/php.exe"
elif [ -f "bin/php7/bin/php" ]; then
    # Linux/macOS
    PHP_BIN="bin/php7/bin/php"
else
    echo "❌ PHP binary not found"
    echo "   Expected: bin/php7/bin/php (Linux/macOS) or bin/php/php.exe (Windows)"
    echo "   Run ./setup.sh first"
    exit 1
fi

echo "🔍 Using PHP: $PHP_BIN"
$PHP_BIN -v
echo ""

# Check for gRPC extension
echo "🔍 Checking for gRPC extension..."
$PHP_BIN -m | grep grpc
if [ $? -eq 0 ]; then
    echo "✅ gRPC extension found!"
else
    echo "❌ gRPC extension not found"
    exit 1
fi

echo ""
echo "🚀 Starting PHP plugin..."
$PHP_BIN src/HelloPlugin.php

