#!/bin/bash

# Name of your Go file
GO_FILE="botv4.go"

# Output directory
OUTPUT_DIR="build"

# Create output directory if it doesn't exist
mkdir -p $OUTPUT_DIR

# List of target OS and architectures
targets=(
    "windows:amd64"
    "windows:386"
    "linux:amd64"
    "linux:386"
    "linux:arm"
    "linux:arm64"
    "darwin:amd64"
    "darwin:arm64"
    "freebsd:amd64"
    "openbsd:amd64"
    "netbsd:amd64"
)

# Compile for each target
for target in "${targets[@]}"; do
    OS="${target%%:*}"
    ARCH="${target#*:}"
    output="$OUTPUT_DIR/bot-$OS-$ARCH"
    if [ "$OS" = "windows" ]; then
        output="$output.exe"
    fi
    
    echo "Compiling for $OS/$ARCH..."
    GOOS=$OS GOARCH=$ARCH go build -o "$output" $GO_FILE
    
    if [ $? -eq 0 ]; then
        echo "Successfully compiled for $OS/$ARCH"
    else
        echo "Failed to compile for $OS/$ARCH"
    fi
done

echo "Cross-compilation completed."