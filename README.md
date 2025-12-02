# GoBattleServer - Pokemon Gacha + Battle System

A Pokemon gacha game with authentic stat mechanics, Discord bot, and REST API built with Go and PostgreSQL.

## 🎮 What's Built

### ✅ Complete REST API
- User registration and management
- Gacha rolling system (daily + premium)
- Pokemon collection management
- Full CRUD operations

### ✅ Discord Bot
- **Slash commands:** `/daily`, `/roll`, `/balance`, `/box`
- **Message commands:** `!daily`, `!roll`, `!balance`, `!box`
- Beautiful embeds with Pokemon stats
- Automatic user registration

### ✅ Pokemon System
- Authentic 6-stat system (HP, Atk, Def, SpAtk, SpDef, Spd)
- IVs (Individual Values) - 0-31 for each stat
- 25 Pokemon natures with stat modifiers
- Rarity tiers: Common → Mythic
- Estimated value calculation

### ✅ Gacha Mechanics
- **Free daily rolls** (5 Pokemon, 24hr cooldown)
- **Pity system** (5th roll guaranteed Rare+)
- **Premium rolls** (100 coins each)
- **Multi-roll bonus** (10 rolls = Epic+ guaranteed)

### ✅ Database
- PostgreSQL with migrations
- 150+ Pokemon seeded
- Redis for caching
- Transaction-safe operations

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Discord account (for bot)

### 1. Clone & Setup
```bash
git clone <your-repo>
cd GoBattleServer

# Copy environment template
cp .env.example .env

# Edit .env and add your Discord bot token
# Get token from: https://discord.com/developers/applications
```

### 2. Run Everything
```bash
# Starts database, API server, and Discord bot
./run-bot.sh
```

That's it! 🎉

### Alternative: Run Components Separately

**Terminal 1 - Database:**
```bash
docker-compose up -d
```

**Terminal 2 - API Server:**
```bash
go run cmd/api/main.go
```

**Terminal 3 - Discord Bot:**
```bash
export DISCORD_BOT_TOKEN=your_token
go run cmd/bot/main.go
```

## 📁 Project Structure

```
GoBattleServer/
├── cmd/
│   ├── api/                # REST API server
│   ├── bot/                # Discord bot
│   └── example/            # Example usage
│
├── internal/
│   ├── domain/             # Domain models
│   ├── validators/         # Validation logic
│   ├── database/           # Database connection
│   ├── repository/         # Data access layer
│   ├── service/            # Business logic
│   ├── handler/            # HTTP handlers
│   └── bot/                # Discord bot logic
│
├── migrations/             # Database migrations
├── docs/                   # Documentation
├── docker-compose.yml      # Database setup
├── .env.example            # Environment template
├── .gitignore              # Git ignore rules
└── run-bot.sh              # Startup script
```

## 📚 Documentation

- **[Architecture Overview](docs/ARCHITECTURE.md)** - System design and roadmap
- **[Discord Bot Setup](docs/DISCORD_BOT_SETUP.md)** - How to create and run the bot
- **[Message Commands Guide](docs/MESSAGE_COMMANDS_SETUP.md)** - Using `!` prefix commands
- **[API Testing Guide](docs/API_TESTING_GUIDE.md)** - Testing REST endpoints with curl
- **[HTTP API Summary](docs/HTTP_API_SUMMARY.md)** - API architecture details

## 🎮 Discord Commands

### Slash Commands
- `/daily` - Free daily roll (5 Pokemon)
- `/roll <count>` - Premium roll (1-10 Pokemon)
- `/balance` - Check coin balance
- `/box [rarity]` - View Pokemon collection

### Message Commands
- `!daily` - Free daily roll
- `!roll 10` - Premium roll (specify count)
- `!balance` - Check coins (`!bal`, `!coins`)
- `!box [rarity]` - View collection (`!collection`)
- `!help` - Show all commands

## 🔮 Next Steps

### Phase 2: Marketplace
- List Pokemon for sale
- Buy/sell with coins
- Price history and trends
- Search and filters

### Phase 3: Battle System
- Turn-based combat
- Type effectiveness
- Move system
- Wager battles

### Phase 4: Advanced Features
- Trading between users
- Leaderboards
- Daily quests
- Web frontend

## 🎲 Gacha Mechanics

### Rarity Drop Rates
```
Common:    50%  ⚪
Uncommon:  25%  🟢
Rare:      15%  🔵
Epic:       7%  🟣
Legendary: 2.5% 🟡
Mythic:    0.5% 🔴
```

### Special Systems
- **Pity System:** 5th daily roll guaranteed Rare+
- **10-Roll Bonus:** Guaranteed Epic+ on 10th premium roll
- **IVs:** Each Pokemon has unique stats (0-31 per stat)
- **Natures:** 25 types that modify stats (+10%/-10%)

## 🛠️ Technology Stack

- **Language:** Go 1.21+
- **Database:** PostgreSQL 16 with pgx/v5
- **Cache:** Redis 7
- **Discord:** discordgo
- **Docker:** PostgreSQL + Redis containers

## 📖 Learning Highlights

This project demonstrates:
- Clean architecture (domain-driven design)
- Repository pattern
- Service layer
- Transaction safety
- PostgreSQL best practices
- Authentic Pokemon mechanics

## License

MIT
