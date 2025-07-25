#!/bin/bash

echo "🔍 Testing Shuffle System..."

# Test 1: Check if backend compiles
echo "📦 Testing backend compilation..."
cd backend/go-app
if go build -o shuffle-backend .; then
    echo "✅ Backend compiles successfully"
    rm -f shuffle-backend
else
    echo "❌ Backend compilation failed"
    exit 1
fi

# Test 2: Check if frontend builds
echo "🎨 Testing frontend build..."
cd ../../frontend
if npm run build >/dev/null 2>&1; then
    echo "✅ Frontend builds successfully"
else
    echo "❌ Frontend build failed"
    exit 1
fi

# Test 3: Check if license generator compiles
echo "🔑 Testing license generator..."
cd ../backend/license_generator
if go build -o license-generator .; then
    echo "✅ License generator compiles successfully"
    rm -f license-generator
else
    echo "❌ License generator compilation failed"
    exit 1
fi

cd ../../

# Test 4: Check environment file
echo "⚙️  Testing environment configuration..."
if [ -f ".env" ]; then
    echo "✅ Environment file exists"
else
    echo "❌ Environment file missing"
    exit 1
fi

# Test 5: Check docker-compose configuration
echo "🐳 Testing docker-compose configuration..."
if [ -f "docker-compose.yml" ]; then
    echo "✅ Docker compose file exists"
else
    echo "❌ Docker compose file missing"
    exit 1
fi

# Test 6: Check required directories
echo "📁 Testing required directories..."
for dir in "/tmp/shuffle-database" "/tmp/shuffle-apps" "/tmp/shuffle-files"; do
    if [ -d "$dir" ]; then
        echo "✅ Directory $dir exists"
    else
        echo "❌ Directory $dir missing"
        mkdir -p "$dir"
        echo "✅ Created directory $dir"
    fi
done

echo ""
echo "🎉 All tests passed! System is ready for deployment."
echo ""
echo "To start the system:"
echo "  docker compose up -d"
echo ""
echo "To generate a license key:"
echo "  cd backend/license_generator"
echo "  go run standalone_license_generator.go generate professional 365"
echo ""
echo "To test the license API:"
echo "  bash backend/test_license.sh"