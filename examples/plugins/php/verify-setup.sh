#!/bin/bash
# Verify PHP plugin setup with custom PHP installation

echo "========================================="
echo "Dragonfly PHP Plugin Setup Verification"
echo "========================================="
echo ""

# Detect PHP binary location
if [ -f "bin/php/php.exe" ]; then
    PHP_BIN="bin/php/php.exe"
    echo "🖥️  Platform: Windows"
elif [ -f "bin/php7/bin/php" ]; then
    PHP_BIN="bin/php7/bin/php"
    echo "🖥️  Platform: Linux/macOS"
else
    PHP_BIN=""
fi

# Check 1: PHP binary exists
echo ""
echo "✓ Check 1: PHP Binary"
if [ -n "$PHP_BIN" ] && [ -f "$PHP_BIN" ]; then
    echo "  ✅ Found: $PHP_BIN"
    VERSION=$($PHP_BIN -v | head -n 1)
    echo "  📦 $VERSION"
else
    echo "  ❌ Not found"
    echo "  💡 Expected: bin/php7/bin/php (Linux/macOS) or bin/php/php.exe (Windows)"
    echo "  💡 Run: ./setup.sh"
    exit 1
fi
echo ""

# Check 2: gRPC extension
echo "✓ Check 2: gRPC Extension"
$PHP_BIN -m | grep -q grpc
if [ $? -eq 0 ]; then
    echo "  ✅ gRPC extension loaded"
else
    echo "  ❌ gRPC extension not found"
    echo "  💡 Your PHP build may not include gRPC"
    exit 1
fi
echo ""

# Check 3: Protobuf extension (optional)
echo "✓ Check 3: Protobuf Extension (optional)"
$PHP_BIN -m | grep -q protobuf
if [ $? -eq 0 ]; then
    echo "  ✅ Protobuf extension loaded (faster performance)"
else
    echo "  ⚠️  Protobuf extension not found (will use pure PHP, slower)"
    echo "  💡 This is OK - plugin will still work"
fi
echo ""

# Check 4: Composer dependencies
echo "✓ Check 4: Composer Dependencies"
if [ -d "vendor" ]; then
    echo "  ✅ Vendor directory exists"
else
    echo "  ⚠️  Vendor directory not found"
    echo "  💡 Run: $PHP_BIN \$(which composer) install"
fi
echo ""

# Check 5: HelloPlugin.php
echo "✓ Check 5: Plugin File"
if [ -f "src/HelloPlugin.php" ]; then
    echo "  ✅ src/HelloPlugin.php exists"
else
    echo "  ❌ src/HelloPlugin.php not found"
    exit 1
fi
echo ""

# Check 6: Test PHP execution
echo "✓ Check 6: PHP Execution Test"
$PHP_BIN -r "echo '  ✅ PHP can execute code';" 2>/dev/null
if [ $? -eq 0 ]; then
    echo ""
else
    echo "  ❌ PHP execution failed"
    exit 1
fi
echo ""

echo "========================================="
echo "🎉 Setup looks good!"
echo "========================================="
echo ""
echo "Next steps:"
echo "  1. Install dependencies: $PHP_BIN \$(which composer) install"
echo "  2. Run plugin: ./run-plugin.sh"
echo "  3. Or start Dragonfly server (it will auto-start the plugin)"
echo ""
echo "Test commands in Minecraft:"
echo "  /cheers"
echo "  !cheer Hello World"
echo ""

