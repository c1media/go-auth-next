#!/bin/bash

# Pre-commit build check script
# Run this before committing to ensure everything builds successfully

set -e

echo "🔍 Running pre-commit checks..."

# Check if we're in the right directory
if [ ! -f "go.mod" ] && [ ! -d "authserver" ]; then
    echo "❌ This script must be run from the project root"
    exit 1
fi

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check required tools
echo "📋 Checking required tools..."
if ! command_exists go; then
    echo "❌ Go is not installed"
    exit 1
fi

if ! command_exists node; then
    echo "❌ Node.js is not installed"
    exit 1
fi

if ! command_exists npm; then
    echo "❌ npm is not installed"
    exit 1
fi

echo "✅ All required tools are available"

# Backend checks
echo ""
echo "🐹 Running backend checks..."
cd authserver

echo "  📥 Downloading Go dependencies..."
go mod download

echo "  🔍 Running go vet..."
go vet ./...

echo "  🧪 Running backend tests..."
go test ./...

echo "  🏗️ Building backend..."
go build -o main cmd/server/main.go
rm -f main

echo "✅ Backend checks passed"

# Frontend checks  
echo ""
echo "⚛️ Running frontend checks..."
cd ../front-end

echo "  📥 Installing dependencies..."
npm ci

echo "  🔍 Running ESLint..."
npm run lint

echo "  🔧 Running type check..."
npm run type-check

echo "  🏗️ Building frontend..."
npm run build

echo "✅ Frontend checks passed"

cd ..

echo ""
echo "🎉 All pre-commit checks passed! You're ready to commit."