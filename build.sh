#!/bin/bash

# Generate Swagger documentation
echo "Generating Swagger documentation..."
./generate-docs.sh

# Build Docker images in parallel
echo "Building Docker images..."
docker-compose build --parallel &
BUILD_PID=$!

# Wait for both processes to complete
wait $DOCS_PID
wait $BUILD_PID

echo "Build complete! Run 'docker-compose up' to start the services."