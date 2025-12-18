-- Migration: Create battles and moves tables
-- This migration adds the battle system tables

-- =====================================================
-- 1. Update pokemon_species to include types
-- =====================================================
ALTER TABLE pokemon_species
  ADD COLUMN IF NOT EXISTS type1 VARCHAR(50) NOT NULL DEFAULT 'normal',
  ADD COLUMN IF NOT EXISTS type2 VARCHAR(50);

-- Add index for type filtering
CREATE INDEX IF NOT EXISTS idx_pokemon_species_type1 ON pokemon_species(type1);
CREATE INDEX IF NOT EXISTS idx_pokemon_species_type2 ON pokemon_species(type2);

-- =====================================================
-- 2. Create moves table
-- =====================================================
CREATE TABLE IF NOT EXISTS moves (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  type VARCHAR(50) NOT NULL,
  category VARCHAR(50) NOT NULL,  -- physical, special, status
  power INTEGER,                   -- NULL for status moves
  accuracy INTEGER,                -- 0 means always hits, NULL for status moves
  pp INTEGER NOT NULL,
  priority INTEGER DEFAULT 0,
  crit_ratio INTEGER DEFAULT 0,
  target VARCHAR(50) DEFAULT 'opponent',
  description TEXT,

  -- Store complex effects as JSONB
  flags JSONB DEFAULT '{}',                -- Move flags (contact, sound, etc.)
  secondary_effect JSONB,                  -- Secondary effects
  multi_hit JSONB,                         -- Multi-hit properties
  recoil_percent INTEGER DEFAULT 0,        -- Recoil damage %
  drain_percent INTEGER DEFAULT 0,         -- HP drain %
  heal_percent INTEGER DEFAULT 0,          -- Healing %
  stat_changes JSONB DEFAULT '[]',         -- Stat modifications
  status_inflict JSONB,                    -- Status condition infliction
  weather_effect JSONB,                    -- Weather changes
  terrain_effect JSONB,                    -- Terrain changes
  entry_hazard JSONB,                      -- Entry hazards

  created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for moves
CREATE INDEX IF NOT EXISTS idx_moves_type ON moves(type);
CREATE INDEX IF NOT EXISTS idx_moves_category ON moves(category);
CREATE INDEX IF NOT EXISTS idx_moves_name ON moves(name);

-- =====================================================
-- 3. Create user_pokemon_moves junction table
-- =====================================================
CREATE TABLE IF NOT EXISTS user_pokemon_moves (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_pokemon_id UUID NOT NULL REFERENCES user_pokemon(id) ON DELETE CASCADE,
  move_id INTEGER NOT NULL REFERENCES moves(id),
  move_slot INTEGER NOT NULL CHECK (move_slot >= 0 AND move_slot <= 3),
  current_pp INTEGER NOT NULL,
  max_pp INTEGER NOT NULL,
  UNIQUE(user_pokemon_id, move_slot),
  created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for pokemon moves
CREATE INDEX IF NOT EXISTS idx_user_pokemon_moves_pokemon ON user_pokemon_moves(user_pokemon_id);
CREATE INDEX IF NOT EXISTS idx_user_pokemon_moves_move ON user_pokemon_moves(move_id);

-- =====================================================
-- 4. Create battles table
-- =====================================================
CREATE TABLE IF NOT EXISTS battles (
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
  completed_at TIMESTAMP,

  -- Prevent self-battles
  CHECK (player1_id != player2_id)
);

-- Indexes for battles
CREATE INDEX IF NOT EXISTS idx_battles_player1 ON battles(player1_id);
CREATE INDEX IF NOT EXISTS idx_battles_player2 ON battles(player2_id);
CREATE INDEX IF NOT EXISTS idx_battles_status ON battles(status);
CREATE INDEX IF NOT EXISTS idx_battles_created ON battles(created_at DESC);

-- =====================================================
-- 5. Create abilities table (for future use)
-- =====================================================
CREATE TABLE IF NOT EXISTS abilities (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  description TEXT NOT NULL,
  trigger VARCHAR(50) NOT NULL,  -- on_entry, passive, on_hit, etc.
  effects JSONB DEFAULT '[]',     -- Array of ability effects
  hidden BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_abilities_name ON abilities(name);

-- =====================================================
-- 6. Create held_items table (for future use)
-- =====================================================
CREATE TABLE IF NOT EXISTS held_items (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  description TEXT NOT NULL,
  category VARCHAR(50) NOT NULL,  -- berry, stat_boost, etc.
  trigger VARCHAR(50) NOT NULL,   -- passive, on_hit, etc.
  effects JSONB DEFAULT '[]',     -- Array of item effects
  consumable BOOLEAN DEFAULT FALSE,
  natural BOOLEAN DEFAULT TRUE,   -- Can be obtained naturally
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_held_items_name ON held_items(name);
CREATE INDEX IF NOT EXISTS idx_held_items_category ON held_items(category);

-- =====================================================
-- Comments for documentation
-- =====================================================
COMMENT ON TABLE moves IS 'Stores all Pokemon moves with their properties and effects';
COMMENT ON TABLE user_pokemon_moves IS 'Junction table linking Pokemon to their learned moves';
COMMENT ON TABLE battles IS 'Stores battle records and state';
COMMENT ON TABLE abilities IS 'Stores Pokemon abilities (passive and active effects)';
COMMENT ON TABLE held_items IS 'Stores held items that Pokemon can equip';

COMMENT ON COLUMN moves.flags IS 'JSONB object containing move flags like {contact: true, sound: false}';
COMMENT ON COLUMN moves.secondary_effect IS 'JSONB object for secondary effects with chance and effects';
COMMENT ON COLUMN moves.stat_changes IS 'JSONB array of stat changes [{stat: "attack", stages: 1, target: "self"}]';
