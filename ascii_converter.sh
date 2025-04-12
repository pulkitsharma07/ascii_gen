#!/bin/bash

# Ensure the script exits on the first error
set -e

# Project setup
echo "[INFO] Initializing project and dependencies..."
if [ ! -f "go.mod" ]; then
    go mod init github.com/osama/ascii_gen
fi

# Tidy dependencies
go mod tidy

# Install required Go package
echo "[INFO] Installing dependencies..."
go get github.com/nfnt/resize

# Compile Go program
echo "[INFO] Building the ASCII converter..."
go build -o ascii_converter main.go

# Create output directory if it doesn't exist
OUTPUT_DIR="ascii_output"
mkdir -p "$OUTPUT_DIR"

# Convert all images in the 'input_images' directory
INPUT_DIR="input_images"
echo "[INFO] Processing images from directory: $INPUT_DIR"

if [ ! -d "$INPUT_DIR" ]; then
    echo "[INFO] Creating 'input_images' directory. Please add your images here."
    mkdir -p "$INPUT_DIR"
    exit 0
fi

if ls "$INPUT_DIR"/* &>/dev/null; then
    for image in "$INPUT_DIR"/*; do
        filename=$(basename -- "$image")
        output_file="$OUTPUT_DIR/${filename%.*}.txt"
        echo "[INFO] Converting $filename to ASCII..."
        ./ascii_converter "$image" > "$output_file"
        echo "[INFO] ASCII art saved to $output_file"
    done
else
    echo "[INFO] No images found in $INPUT_DIR. Add images and rerun the script."
fi

echo "[INFO] Done! ASCII art can be found in the $OUTPUT_DIR directory."
