# GoBattleServer - Pokemon Gacha + Battle System

A Pokemon gacha game with authentic stat mechanics, marketplace, and battle system built with Go and PostgreSQL.

## Features Implemented

### ✅ Domain Models
- **Pokemon** with authentic 6-stat system (HP, Atk, Def, SpAtk, SpDef, Spd)
- **IVs (Individual Values)** - 0-31 for each stat
- **Natures** - All 25 natures with stat modifiers
- **Rarity System** - Common → Mythic (6 tiers)
- **Users** with coins and daily roll cooldowns

### ✅ Database
- PostgreSQL with connection pooling
- Complete schema with migrations
- Redis for caching (docker-compose)
- Indexes for performance

### ✅ Repositories (Data Access)
- User repository
- Pokemon Species repository
- User Pokemon repository
- Transaction-safe operations

### ✅ Services
- **Gacha Service** with:
  - Free daily rolls (5 Pokemon)
  - Pity system (guaranteed Rare+ on 5th roll)
  - Premium rolls with coins
  - Multi-roll bonus (10 rolls = 1 Epic+)

## Quick Start

### 1. Start Database
```bash
# Start PostgreSQL and Redis
docker-compose up -d

# Verify containers are running
docker ps
```

### 2. Set Up Environment
```bash
# Copy environment template
cp .env.example .env

# (Optional) Edit .env if needed
```

### 3. Verify Build
```bash
# Build the project
go build ./...

# Run tests (when added)
go test ./...
```

## Project Structure

```
GoBattleServer/
├── internal/
│   ├── domain/              # Domain models (pure business logic)
│   │   ├── rarity.go       # Rarity tiers and drop rates
│   │   ├── nature.go       # Pokemon natures (25 types)
│   │   ├── pokemon.go      # Pokemon species & instances
│   │   └── user.go         # User accounts
│   │
│   ├── validators/         # Validation logic (separated from domain)
│   │   ├── rarity_validator.go
│   │   ├── nature_validator.go
│   │   ├── pokemon_validator.go
│   │   └── user_validator.go
│   │
│   ├── database/           # Database connection management
│   │   └── connection.go
│   │
│   ├── repository/         # Data access layer
│   │   ├── interfaces.go
│   │   ├── postgres_user.go
│   │   ├── postgres_pokemon_species.go
│   │   └── postgres_user_pokemon.go
│   │
│   └── service/            # Business logic services
│       └── gacha.go        # Gacha rolling system
│
├── migrations/             # Database schema
│   └── 001_init_schema.sql
│
├── docker-compose.yml      # PostgreSQL + Redis
├── .env.example           # Environment variables template
└── go.mod                 # Dependencies
```

## Database Schema

### Users
- ID, DiscordID, Coins, LastDailyRoll, CreatedAt

### Pokemon Species (Static Reference Data)
- National Dex ID, Name, Rarity
- Base Stats (6 stats)
- Sprite URL, Drop Weight

### User Pokemon (Unique Instances)
- ID, UserID, SpeciesID
- IVs (6 values, 0-31 each)
- Nature (affects stat multipliers)
- Level, Acquired At, Favorite, Nickname

### Market Listings (To Be Implemented)
- Seller, Pokemon, Price, Status

### Battles (To Be Implemented)
- Players, Winner, Wager, Teams

## Next Steps

### Immediate (Required for Testing)
1. **Seed Pokemon Data** - Populate species table with Pokemon
   - Option A: Manual seed with ~20 Pokemon
   - Option B: Fetch from PokeAPI (automated)

2. **Create Main Application** - Wire everything together
   - Initialize database connection
   - Create repository instances
   - Initialize services
   - Build CLI or HTTP API to test gacha

### Short Term
3. **Add Tests** - Unit tests for services and repositories
4. **Market System** - Buy/sell Pokemon
5. **Battle System** - Turn-based combat

### Medium Term
6. **Discord Bot** - Slash commands for gacha
7. **Web Frontend** - React UI for Pokemon box
8. **Leaderboards** - Top collectors/battlers

## Gacha Mechanics

### Daily Roll (Free)
- 5 Pokemon per day
- 24-hour cooldown
- **Pity**: 5th Pokemon guaranteed Rare or better

### Premium Roll (100 coins each)
- Pay coins for additional rolls
- **Bonus**: Every 10th roll guaranteed Epic or better

### Rarity Distribution
```
Common:    50% (0.500)
Uncommon:  25% (0.250)
Rare:      15% (0.150)
Epic:       7% (0.070)
Legendary: 2.5% (0.025)
Mythic:    0.5% (0.005)
```

## Pokemon Stats System

### Calculated Stats Formula
```
HP = floor((2 * BaseHP + IV) * Level / 100) + Level + 10
Other Stats = floor(floor((2 * Base + IV) * Level / 100 + 5) * Nature)
```

### IVs (Individual Values)
- 0-31 per stat (randomly generated)
- Max total: 186 (31 × 6)
- Affects final stat calculation

### Natures
- 25 types (5 neutral, 20 with effects)
- +10% to one stat, -10% to another
- Examples:
  - **Adamant**: +Atk -SpAtk
  - **Jolly**: +Spd -SpAtk
  - **Modest**: +SpAtk -Atk

## Technology Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **Driver**: pgx/v5 (PostgreSQL)
- **Dependencies**: google/uuid

## Environment Variables

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=gobattle
DB_PASSWORD=gobattle_dev
DB_NAME=gobattle_db
DB_SSLMODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379

SERVER_PORT=8080
```

## Contributing

This is a portfolio/learning project demonstrating:
- Clean architecture (domain-driven design)
- Repository pattern
- Service layer
- Transaction safety
- PostgreSQL best practices
- Authentic Pokemon mechanics

## License

MIT
