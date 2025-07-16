#!/bin/bash

# Go NPM SDK Documentation Verification Script
# This script verifies that the documentation is properly set up and can be built

set -e

echo "🔍 Verifying Go NPM SDK Documentation Setup..."

# Check if we're in the right directory
if [ ! -f "docs/package.json" ]; then
    echo "❌ Error: docs/package.json not found. Please run this script from the project root."
    exit 1
fi

echo "✅ Found docs directory"

# Check Node.js version
if ! command -v node &> /dev/null; then
    echo "❌ Error: Node.js is not installed"
    exit 1
fi

NODE_VERSION=$(node --version)
echo "✅ Node.js version: $NODE_VERSION"

# Check npm
if ! command -v npm &> /dev/null; then
    echo "❌ Error: npm is not installed"
    exit 1
fi

NPM_VERSION=$(npm --version)
echo "✅ npm version: $NPM_VERSION"

# Navigate to docs directory
cd docs

echo "📦 Installing dependencies..."
npm install --silent

echo "🔨 Building documentation..."
npm run build

if [ $? -eq 0 ]; then
    echo "✅ Documentation build successful!"
else
    echo "❌ Documentation build failed!"
    exit 1
fi

# Check if build output exists
if [ -d ".vitepress/dist" ]; then
    echo "✅ Build output directory exists"
    
    # Check for key files
    if [ -f ".vitepress/dist/index.html" ]; then
        echo "✅ Main index.html exists"
    else
        echo "❌ Main index.html missing"
        exit 1
    fi
    
    if [ -f ".vitepress/dist/zh/index.html" ]; then
        echo "✅ Chinese index.html exists"
    else
        echo "❌ Chinese index.html missing"
        exit 1
    fi
    
    # Count total HTML files
    HTML_COUNT=$(find .vitepress/dist -name "*.html" | wc -l)
    echo "✅ Generated $HTML_COUNT HTML files"
    
    # Check for assets
    if [ -d ".vitepress/dist/assets" ]; then
        echo "✅ Assets directory exists"
    else
        echo "⚠️  Warning: Assets directory not found"
    fi
    
else
    echo "❌ Build output directory missing"
    exit 1
fi

echo ""
echo "🎉 Documentation verification completed successfully!"
echo ""
echo "📋 Summary:"
echo "   - Node.js: $NODE_VERSION"
echo "   - npm: $NPM_VERSION"
echo "   - Generated files: $HTML_COUNT HTML files"
echo "   - Build output: docs/.vitepress/dist/"
echo ""
echo "🚀 To preview the documentation locally:"
echo "   cd docs && npm run preview"
echo ""
echo "🌐 To deploy to GitHub Pages:"
echo "   git add . && git commit -m 'Add documentation' && git push"
echo ""
echo "📖 Documentation will be available at:"
echo "   https://scagogogo.github.io/go-npm-sdk/"
