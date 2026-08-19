#!/bin/bash

# Navigate to the directory containing the proto files
cd "$(dirname "$0")/../api/proto"

# Run protoc to generate Go files
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       *.proto

# Optionally, you can print a message indicating that the files have been generated
echo "Protobuf files generated successfully."