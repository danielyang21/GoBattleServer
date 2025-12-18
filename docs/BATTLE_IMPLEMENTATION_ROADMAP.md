# Battle System Implementation Roadmap

## What We Have ✅

1. **Complete Battle Engine Logic**
   - Damage calculation (Pokemon Showdown accurate)
   - Turn resolution with priority/speed
   - Status effects, weather, terrain
   - Stat stages, entry hazards
   - All battle mechanics implemented

2. **Domain Models**
   - Move, Ability, HeldItem, Battle, BattlePokemon
   - Full type effectiveness chart
   - Battle state machine

3. **Battle Service**
   - In-memory battle state management
   - Wager system
   - Winner determination

## What We Need to Get It Working 🚧

### Phase 1: Database Schema (CRITICAL)

#### 1.1 Update `pokemon_species` table
```sql
-- Add type columns
ALTER TABLE pokemon_species
  ADD COLUMN type1 VARCHAR(50) NOT NULL DEFAULT 'normal',
  ADD COLUMN type2 VARCHAR(50);
```

#### 1.2 Create `moves` table
```sql
CREATE TABLE moves (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  type VARCHAR(50) NOT NULL,
  category VARCHAR(50) NOT NULL,  -- physical, special, status
  power INTEGER,
  accuracy INTEGER,
  pp INTEGER NOT NULL,
  priority INTEGER DEFAULT 0,
  description TEXT,
  -- Store complex data as JSONB for flexibility
  flags JSONB,
  secondary_effect JSONB,
  stat_changes JSONB,
  weather_effect JSONB,
  terrain_effect JSONB,
  status_inflict JSONB,
  created_at TIMESTAMP DEFAULT NOW()
);
```

#### 1.3 Create `user_pokemon_moves` junction table
```sql
CREATE TABLE user_pokemon_moves (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_pokemon_id UUID NOT NULL REFERENCES user_pokemon(id) ON DELETE CASCADE,
  move_id INTEGER NOT NULL REFERENCES moves(id),
  move_slot INTEGER NOT NULL CHECK (move_slot >= 0 AND move_slot <= 3),
  current_pp INTEGER NOT NULL,
  UNIQUE(user_pokemon_id, move_slot),
  created_at TIMESTAMP DEFAULT NOW()
);
```

#### 1.4 Create `battles` table
```sql
CREATE TABLE battles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  player1_id UUID NOT NULL REFERENCES users(id),
  player2_id UUID NOT NULL REFERENCES users(id),
  player1_pokemon_id UUID REFERENCES user_pokemon(id),
  player2_pokemon_id UUID REFERENCES user_pokemon(id),
  wager_amount INTEGER NOT NULL DEFAULT 0,
  status VARCHAR(50) NOT NULL,
  winner_id UUID REFERENCES users(id),
  current_turn INTEGER DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW(),
  started_at TIMESTAMP,
  completed_at TIMESTAMP
);

CREATE INDEX idx_battles_player1 ON battles(player1_id);
CREATE INDEX idx_battles_player2 ON battles(player2_id);
CREATE INDEX idx_battles_status ON battles(status);
```

### Phase 2: Seed Essential Data (CRITICAL)

#### 2.1 Seed Types for Existing Pokemon
Update existing pokemon_species with proper types:
```sql
-- Bulbasaur: Grass/Poison
UPDATE pokemon_species SET type1 = 'grass', type2 = 'poison' WHERE id = 1;
-- Charmander: Fire
UPDATE pokemon_species SET type1 = 'fire', type2 = NULL WHERE id = 4;
-- Squirtle: Water
UPDATE pokemon_species SET type1 = 'water', type2 = NULL WHERE id = 7;
-- Pikachu: Electric
UPDATE pokemon_species SET type1 = 'electric', type2 = NULL WHERE id = 25;
-- etc...
```

#### 2.2 Seed ~50 Essential Moves
Priority moves to seed:
- **Basic attacks**: Tackle, Scratch, Pound
- **Type coverage**: Flamethrower, Surf, Thunderbolt, Ice Beam, Earthquake
- **Priority moves**: Quick Attack, Aqua Jet, Mach Punch
- **Status moves**: Thunder Wave, Will-O-Wisp, Toxic
- **Stat boosters**: Swords Dance, Dragon Dance, Nasty Plot
- **Weather**: Sunny Day, Rain Dance, Sandstorm, Hail
- **Entry hazards**: Stealth Rock, Spikes, Toxic Spikes

#### 2.3 Auto-assign Moves to Pokemon
When a Pokemon is caught/created, assign 2-4 random moves based on their type:
```go
// In gacha service, after creating UserPokemon
assignRandomMoves(pokemon)
```

### Phase 3: Repository Implementation

#### 3.1 Create `postgres_battle.go`
Implement BattleRepository interface with PostgreSQL

#### 3.2 Create `postgres_move.go`
```go
type PostgresMoveRepository interface {
  GetByID(ctx, id) -> Move
  GetByIDs(ctx, []int) -> []Move
  GetByType(ctx, type) -> []Move
  List(ctx) -> []Move
}
```

#### 3.3 Update Battle Service
Add context.Context to all repository calls

### Phase 4: API Layer

#### 4.1 Create Battle Handlers
```go
// HTTP Endpoints
POST   /api/battles/create        // Create challenge
POST   /api/battles/:id/select    // Select Pokemon
POST   /api/battles/:id/action    // Submit move
GET    /api/battles/:id/state     // Get state
POST   /api/battles/:id/forfeit   // Forfeit
```

#### 4.2 Create WebSocket Handler
```go
// Real-time battle updates
/ws/battle/:id
// Messages:
// - battle_start
// - turn_resolved
// - battle_end
```

### Phase 5: Discord Bot Integration

#### 5.1 Add Battle Commands
```
/battle @opponent [wager]     - Challenge player
/accept                       - Accept challenge
/select [pokemon_number]      - Select Pokemon from your box
/move [1-4]                   - Select move during battle
/forfeit                      - Give up
/battlestatus                 - Check current battle state
```

## Minimum Viable Battle System

To get battles working **right now**, we need:

### MUST HAVE (Priority 1)
1. ✅ Battle engine logic (DONE)
2. ⬜ Database migrations (battles, moves tables)
3. ⬜ Seed 20-30 basic moves
4. ⬜ Update pokemon_species with types
5. ⬜ Auto-assign 2 moves to all Pokemon (Tackle + type move)
6. ⬜ Battle repository implementation
7. ⬜ Fix context.Context in battle service
8. ⬜ Basic HTTP endpoints (create, action, state)

### NICE TO HAVE (Priority 2)
9. ⬜ WebSocket real-time updates
10. ⬜ Discord bot commands
11. ⬜ Full 100+ move library
12. ⬜ Battle history/replay system

### FUTURE (Priority 3)
13. ⬜ Abilities system (database + loading)
14. ⬜ Held items system
15. ⬜ Team battles (6v6 with switching)
16. ⬜ Ranked matchmaking
17. ⬜ Battle analytics

## Quick Start Implementation Order

1. **Create migration** for battles + moves tables (10 min)
2. **Seed 20 basic moves** (Tackle, Quick Attack, Flamethrower, Surf, Thunderbolt, etc.) (15 min)
3. **Update pokemon_species** with type1/type2 for existing 150 Pokemon (20 min)
4. **Create battle repository** (20 min)
5. **Fix battle service** context issues (10 min)
6. **Create battle handlers** (30 min)
7. **Test via HTTP** - Create battle, submit moves, see results (10 min)

**Total: ~2 hours to working battle system**

## Testing the Battle System

Once implemented, you can test via curl:

```bash
# 1. Create battle
curl -X POST http://localhost:8080/api/battles/create \
  -d '{"player1_id":"...","player2_id":"...","wager":100}'

# 2. Both players select Pokemon
curl -X POST http://localhost:8080/api/battles/{id}/select \
  -d '{"player_id":"...","pokemon_id":"..."}'

# 3. Submit moves each turn
curl -X POST http://localhost:8080/api/battles/{id}/action \
  -d '{"player_id":"...","action_type":"move","move_index":0}'

# 4. Check battle state
curl http://localhost:8080/api/battles/{id}/state
```

## Architecture Flow

```
Discord Bot / HTTP Client
        ↓
Battle API Handler
        ↓
Battle Service (in-memory state + turn resolution)
        ↓
Battle Repository (persistence)
        ↓
PostgreSQL Database
```

Battle state lives in-memory for performance, persisted to DB on completion.
