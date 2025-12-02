#!/bin/bash

echo "🚀 Starting Pokemon Gacha System..."
echo ""

# Load environment variables from .env file
if [ -f .env ]; then
    echo "📝 Loading environment variables..."
    export $(cat .env | grep -v '^#' | grep -v '^$' | xargs)
else
    echo "❌ Error: .env file not found"
    exit 1
fi

# Check if DISCORD_BOT_TOKEN is set
if [ -z "$DISCORD_BOT_TOKEN" ]; then
    echo "❌ Error: DISCORD_BOT_TOKEN not found in .env file"
    exit 1
fi

# Start database
echo "🗄️  Starting database..."
docker-compose up -d
echo "⏳ Waiting for database to be ready..."
sleep 3
echo ""

# Start API server in background
echo "🌐 Starting API server..."
go run cmd/api/main.go > api.log 2>&1 &
API_PID=$!
echo "   API server started (PID: $API_PID)"
echo "   Logs: tail -f api.log"
echo "⏳ Waiting for API to be ready..."
sleep 3
echo ""

# Start Discord bot
echo "🤖 Starting Discord bot..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
go run cmd/bot/main.go

# Cleanup on exit (Ctrl+C)
trap "echo ''; echo '🛑 Shutting down...'; kill $API_PID 2>/dev/null; exit" INT TERM
