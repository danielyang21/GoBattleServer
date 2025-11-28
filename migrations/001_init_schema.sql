-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    discord_id VARCHAR(255) UNIQUE NOT NULL,
    coins INTEGER NOT NULL DEFAULT 1000,
    last_daily_roll TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Pokemon species (static reference data from PokeAPI)
CREATE TABLE pokemon_species (
    id INTEGER PRIMARY KEY,  -- National Dex number (1-1025+)
    name VARCHAR(255) NOT NULL,
    rarity VARCHAR(50) NOT NULL,
    base_hp INTEGER NOT NULL,
    base_attack INTEGER NOT NULL,
    base_defense INTEGER NOT NULL,
    base_sp_attack INTEGER NOT NULL,
    base_sp_defense INTEGER NOT NULL,
    base_speed INTEGER NOT NULL,
    sprite_url TEXT,
    drop_weight FLOAT NOT NULL DEFAULT 1.0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (rarity IN ('common', 'uncommon', 'rare', 'epic', 'legendary', 'mythic'))
);

-- User's Pokemon (each instance is unique with IVs and nature)
CREATE TABLE user_pokemon (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    species_id INTEGER NOT NULL REFERENCES pokemon_species(id),

    -- Individual Values (0-31 for each stat)
    iv_hp INTEGER NOT NULL CHECK (iv_hp BETWEEN 0 AND 31),
    iv_attack INTEGER NOT NULL CHECK (iv_attack BETWEEN 0 AND 31),
    iv_defense INTEGER NOT NULL CHECK (iv_defense BETWEEN 0 AND 31),
    iv_sp_attack INTEGER NOT NULL CHECK (iv_sp_attack BETWEEN 0 AND 31),
    iv_sp_defense INTEGER NOT NULL CHECK (iv_sp_defense BETWEEN 0 AND 31),
    iv_speed INTEGER NOT NULL CHECK (iv_speed BETWEEN 0 AND 31),

    -- Nature (affects stat growth)
    nature VARCHAR(50) NOT NULL,

    -- Level (for future leveling system)
    level INTEGER NOT NULL DEFAULT 50 CHECK (level BETWEEN 1 AND 100),

    -- Metadata
    acquired_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_favorite BOOLEAN NOT NULL DEFAULT false,
    nickname VARCHAR(255),

    CHECK (nature IN (
        'hardy', 'docile', 'serious', 'bashful', 'quirky',
        'lonely', 'brave', 'adamant', 'naughty',
        'bold', 'relaxed', 'impish', 'lax',
        'timid', 'hasty', 'jolly', 'naive',
        'modest', 'mild', 'quiet', 'rash',
        'calm', 'gentle', 'sassy', 'careful'
    ))
);

-- Market listings
CREATE TABLE market_listings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    seller_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_pokemon_id UUID NOT NULL REFERENCES user_pokemon(id) ON DELETE CASCADE,
    price INTEGER NOT NULL CHECK (price > 0),
    listed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL DEFAULT 'active',

    CHECK (status IN ('active', 'sold', 'cancelled'))
);

-- Market transactions (history)
CREATE TABLE market_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    listing_id UUID NOT NULL REFERENCES market_listings(id),
    buyer_id UUID NOT NULL REFERENCES users(id),
    seller_id UUID NOT NULL REFERENCES users(id),
    price INTEGER NOT NULL,
    completed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Battle records
CREATE TABLE battles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    player1_id UUID NOT NULL REFERENCES users(id),
    player2_id UUID NOT NULL REFERENCES users(id),
    winner_id UUID REFERENCES users(id),
    wager INTEGER NOT NULL DEFAULT 0,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP
);

-- Battle team composition
CREATE TABLE battle_teams (
    battle_id UUID NOT NULL REFERENCES battles(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    user_pokemon_id UUID NOT NULL REFERENCES user_pokemon(id),
    position INTEGER NOT NULL CHECK (position BETWEEN 1 AND 6),

    PRIMARY KEY (battle_id, user_id, position)
);

-- Indexes for performance
CREATE INDEX idx_users_discord_id ON users(discord_id);
CREATE INDEX idx_user_pokemon_user_id ON user_pokemon(user_id);
CREATE INDEX idx_user_pokemon_species_id ON user_pokemon(species_id);
CREATE INDEX idx_market_listings_status ON market_listings(status);
CREATE INDEX idx_market_listings_seller ON market_listings(seller_id);
CREATE INDEX idx_market_listings_price ON market_listings(price);
CREATE INDEX idx_market_transactions_buyer ON market_transactions(buyer_id);
CREATE INDEX idx_market_transactions_seller ON market_transactions(seller_id);
CREATE INDEX idx_battles_player1 ON battles(player1_id);
CREATE INDEX idx_battles_player2 ON battles(player2_id);
CREATE INDEX idx_pokemon_species_rarity ON pokemon_species(rarity);
