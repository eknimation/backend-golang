#!/bin/bash

# Script to set up the development environment

set -e

echo "🚀 Setting up Backend Go API development environment..."

# Create necessary directories
echo "📁 Creating directories..."
mkdir -p .docker/mongo-data
mkdir -p .docker/logs

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose is not installed. Please install docker-compose and try again."
    exit 1
fi

echo "🔨 Building Docker images..."
docker-compose build

echo "🗄️  Starting MongoDB..."
docker-compose up -d mongodb

echo "⏳ Waiting for MongoDB to be ready..."
sleep 10

# Check if MongoDB is healthy
until docker-compose exec mongodb mongosh --eval "db.adminCommand('ping')" --quiet > /dev/null 2>&1; do
    echo "Waiting for MongoDB..."
    sleep 5
done

echo "✅ MongoDB is ready!"

echo "🚀 Starting API service..."
docker-compose up -d api

echo "⏳ Waiting for API to be ready..."
sleep 10

# Check API health
until curl -f http://localhost:5555/health > /dev/null 2>&1; do
    echo "Waiting for API..."
    sleep 5
done

echo "✅ API is ready!"

echo ""
echo "🎉 Setup complete! Your services are running:"
echo "   📊 API: http://localhost:5555"
echo "   📚 Swagger UI: http://localhost:5555/swagger/index.html"
echo "   🗄️  MongoDB: localhost:27017"
echo ""
echo "📋 Useful commands:"
echo "   make logs      - View all logs"
echo "   make logs-api  - View API logs only"
echo "   make status    - Check service status"
echo "   make down      - Stop all services"
echo "   make help      - See all available commands"
echo ""
