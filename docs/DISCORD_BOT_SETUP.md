# Discord Bot Setup Guide

This guide shows you how to set up and run the Pokemon Gacha Discord bot.

## 🎯 What You'll Need

1. A Discord account
2. A Discord server where you have admin permissions (to invite the bot)
3. The API server running (`go run cmd/api/main.go`)
4. A Discord Bot Token (we'll create this)

---

## Step 1: Create a Discord Bot Application

### 1.1 Go to Discord Developer Portal
Visit: https://discord.com/developers/applications

### 1.2 Create New Application
1. Click **"New Application"** button
2. Give it a name (e.g., "Pokemon Gacha Bot")
3. Click **"Create"**

### 1.3 Create Bot User
1. Go to the **"Bot"** tab on the left
2. Click **"Add Bot"**
3. Click **"Yes, do it!"**
4. Under the bot's username, click **"Reset Token"**
5. Click **"Copy"** to copy your bot token
   - ⚠️ **IMPORTANT:** Save this token securely! You'll need it later.
   - ⚠️ Never share your token publicly (don't commit it to GitHub)

### 1.4 Enable Required Intents
Scroll down to **"Privileged Gateway Intents"** and enable:
- ✅ **SERVER MEMBERS INTENT** (if you want member info)
- ✅ **MESSAGE CONTENT INTENT** (optional)

Click **"Save Changes"**

### 1.5 Configure Bot Permissions
Go to the **"OAuth2"** → **"URL Generator"** tab:

**SCOPES:**
- ✅ `bot`
- ✅ `applications.commands`

**BOT PERMISSIONS:**
- ✅ Send Messages
- ✅ Embed Links
- ✅ Use Slash Commands

### 1.6 Generate Invite URL
1. Copy the generated URL at the bottom
2. Open it in a new tab
3. Select your server from the dropdown
4. Click **"Authorize"**
5. Complete the captcha

✅ Your bot should now appear in your server (offline)

---

## Step 2: Configure Environment Variables

Update your `.env` file with the Discord bot token:

```bash
# Add this line to your .env file
DISCORD_BOT_TOKEN=your_bot_token_here

# API URL (optional, defaults to http://localhost:8080)
API_BASE_URL=http://localhost:8080
```

**Example:**
```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=gobattle
DB_PASSWORD=gobattle_dev
DB_NAME=gobattle_db
DB_SSLMODE=disable

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379

# Server Configuration
SERVER_PORT=8080

# Discord Bot Configuration
DISCORD_BOT_TOKEN=
API_BASE_URL=http://localhost:8080
```

---

## Step 3: Run the Bot

### 3.1 Start the Database (if not running)
```bash
docker-compose up -d
```

### 3.2 Start the API Server
In one terminal:
```bash
go run cmd/api/main.go
```

You should see:
```
🚀 API Server starting on http://localhost:8080
📋 Health check: http://localhost:8080/health
📚 API Base URL: http://localhost:8080/api
```

### 3.3 Start the Discord Bot
In another terminal:
```bash
# Load environment variables
source .env

# Or export manually
export DISCORD_BOT_TOKEN=your_token_here

# Run the bot
go run cmd/bot/main.go
```

You should see:
```
🤖 Starting Discord bot...
📡 API URL: http://localhost:8080
✅ Bot is now running. Registering commands...
✅ Commands registered successfully!
📋 Available commands:
   /daily   - Claim free daily roll (5 Pokemon)
   /roll    - Buy premium rolls with coins
   /balance - Check your coin balance
   /box     - View your Pokemon collection

Press CTRL+C to stop the bot
```

✅ Your bot should now be **online** in Discord!

---

## Step 4: Test the Bot

Go to your Discord server and try these commands:

### 1. `/daily` - Free Daily Roll
Claim your free 5 Pokemon roll (once per 24 hours)

**What you'll see:**
```
🎴 Daily Roll Results!
You received 5 Pokemon:

1. ⚪ Rattata
   Nature: Jolly
   IVs: 54.3% perfect
   Value: 12 coins

2. 🟢 Pikachu
   Nature: Adamant
   IVs: 67.2% perfect
   Value: 85 coins

... (3 more)
```

### 2. `/balance` - Check Your Coins
See how many coins you have

**What you'll see:**
```
💰 Your Balance
You have 1000 coins

💵 Premium Roll: 100 coins per Pokemon
🎁 10-Roll Bonus: 1000 coins (guaranteed Epic+)
```

### 3. `/roll count:10` - Premium Roll
Buy 10 rolls with coins (costs 1000 coins)

**What you'll see:**
```
💎 Premium Roll Results!
You spent 1000 coins and received 10 Pokemon:

[List of 10 Pokemon with stats]

🎁 10-roll bonus: Guaranteed Epic or better!
```

### 4. `/box` - View Your Collection
See all Pokemon you've collected

**What you'll see:**
```
📦 Your Pokemon Collection
Total Pokemon: 15
Collection Value: 2,450 coins

📊 Rarity Breakdown
🟣 Epic: 2
🔵 Rare: 3
🟢 Uncommon: 5
⚪ Common: 5

[Shows top 10 Pokemon]
```

### 5. `/box rarity:epic` - Filter by Rarity
View only Epic Pokemon

---

## 🎨 Rarity Color Legend

- 🔴 **Mythic** - 0.5% drop rate (extremely rare!)
- 🟡 **Legendary** - 2.5% drop rate
- 🟣 **Epic** - 7% drop rate
- 🔵 **Rare** - 15% drop rate
- 🟢 **Uncommon** - 25% drop rate
- ⚪ **Common** - 50% drop rate

---

## 🔧 Troubleshooting

### Bot appears offline
- Check that the bot token is correct in `.env`
- Make sure you're running `go run cmd/bot/main.go` with the token loaded

### Commands don't appear
- Wait 1-2 minutes after starting the bot (Discord caches commands)
- Try kicking and re-inviting the bot
- Make sure you enabled `applications.commands` scope when inviting

### "Failed to get user" error
- Make sure the API server is running on `http://localhost:8080`
- Check that the database is running (`docker ps`)

### "Insufficient coins" error
- Use `/balance` to check your coins
- Use `/daily` to get free Pokemon (refreshes every 24 hours)

### Commands return "Cooldown active"
- Daily roll has a 24-hour cooldown
- Wait until tomorrow or test with a different Discord account

---

## 📁 Project Structure

```
cmd/bot/main.go              # Bot entry point
internal/bot/
  ├── api_client.go           # HTTP client to call API
  └── commands.go             # Discord slash commands
```

---

## 🚀 Running in Production

### 1. Deploy API Server
Deploy the API server to a cloud provider (AWS, GCP, etc.)

### 2. Update API_BASE_URL
```bash
export API_BASE_URL=https://your-api-domain.com
```

### 3. Run Bot as a Service
Use a process manager like `systemd` or run in a Docker container:

```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o bot cmd/bot/main.go
CMD ["./bot"]
```

---

## 🎉 Next Steps

Now that your Discord bot is running, you can:

1. **Add More Commands:**
   - `/trade` - Trade Pokemon with other users
   - `/battle` - Challenge others to battles
   - `/leaderboard` - Show top collectors

2. **Add Features:**
   - Pokemon nicknames
   - Favorite Pokemon
   - Achievement system
   - Daily quests

3. **Improve UI:**
   - Add Pokemon sprite images to embeds
   - Create interactive buttons (buy, sell, etc.)
   - Add pagination for large collections

---

## 🎓 What You Learned

- Creating Discord bot applications
- Slash command registration
- Discord embeds and rich messages
- Calling REST APIs from Discord
- Environment variable management
- Concurrent terminal sessions

---

**Your Discord bot is now live! 🎮 Start collecting Pokemon!**
