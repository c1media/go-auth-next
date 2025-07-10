#!/bin/bash

# Docker build script for Go Auth Template

set -e

echo "🐳 Building Go Auth Template Docker Images..."

# Build options
BUILD_TYPE=${1:-"compose"}  # compose, single, or all
TAG=${2:-"latest"}

case $BUILD_TYPE in
  "single")
    echo "📦 Building single container with both services..."
    docker build -t go-auth-template:$TAG .
    echo "✅ Single container built: go-auth-template:$TAG"
    ;;
    
  "compose")
    echo "📦 Building individual services for docker-compose..."
    docker build -f authserver/Dockerfile -t go-auth-backend:$TAG ./authserver
    docker build -f front-end/Dockerfile -t go-auth-frontend:$TAG ./front-end
    echo "✅ Individual services built:"
    echo "   - go-auth-backend:$TAG"
    echo "   - go-auth-frontend:$TAG"
    ;;
    
  "all")
    echo "📦 Building all variants..."
    # Single container
    docker build -t go-auth-template:$TAG .
    # Individual services
    docker build -f authserver/Dockerfile -t go-auth-backend:$TAG ./authserver
    docker build -f front-end/Dockerfile -t go-auth-frontend:$TAG ./front-end
    echo "✅ All containers built"
    ;;
    
  *)
    echo "❌ Invalid build type. Use: single, compose, or all"
    exit 1
    ;;
esac

echo "🎉 Build complete!"
echo ""
echo "Usage:"
echo "  Single container: docker run -p 3000:3000 -p 8080:8080 go-auth-template:$TAG"
echo "  Docker Compose:   docker-compose up"