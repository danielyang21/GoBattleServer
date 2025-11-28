# API Testing Guide

This guide shows how to test the HTTP API endpoints using `curl`.

## Prerequisites

1. **Start the database:**
   ```bash
   docker-compose up -d
   ```

2. **Start the API server:**
   ```bash
   go run cmd/api/main.go
   ```

   The server will start on http://localhost:8080

## API Endpoints

### 1. Health Check

Check if the server is running:

```bash
curl http://localhost:8080/health
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "status": "healthy"
  }
}
```

---

### 2. Register a New User

Create a new user account:

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "discord_id": "123456789012345678"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "discord_id": "123456789012345678",
    "coins": 1000,
    "created_at": "2025-11-26T15:00:00Z"
  }
}
```

**Save the `id` field - you'll need it for other requests!**

---

### 3. Get User by ID

Retrieve user information:

```bash
# Replace {USER_ID} with the ID from the register response
curl http://localhost:8080/api/users/{USER_ID}
```

**Example:**
```bash
curl http://localhost:8080/api/users/550e8400-e29b-41d4-a716-446655440000
```

---

### 4. Get User by Discord ID

Retrieve user by Discord ID:

```bash
curl http://localhost:8080/api/users/discord/123456789012345678
```

---

### 5. Daily Roll (Free)

Perform a free daily roll (5 Pokemon):

```bash
# Replace {USER_ID} with your user ID
curl -X POST http://localhost:8080/api/gacha/daily-roll \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "{USER_ID}"
  }'
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/gacha/daily-roll \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "count": 5,
    "pokemons": [
      {
        "id": "...",
        "species": {
          "id": 25,
          "name": "Pikachu",
          "rarity": "uncommon"
        },
        "nature": "Jolly",
        "level": 5,
        "ivs": {
          "hp": 23,
          "attack": 31,
          "defense": 15,
          "sp_attack": 8,
          "sp_defense": 22,
          "speed": 29
        },
        "stats": {
          "hp": 22,
          "attack": 12,
          "defense": 9,
          "sp_attack": 10,
          "sp_defense": 11,
          "speed": 13
        },
        "iv_percentage": 68.3,
        "estimated_value": 250
      }
      // ... 4 more Pokemon
    ]
  }
}
```

**Note:** The 5th Pokemon is guaranteed to be Rare or better (pity system)!

**Cooldown:** You can only do this once per 24 hours. If you try again too soon:
```json
{
  "success": false,
  "error": {
    "code": "cooldown_active",
    "message": "daily roll already claimed today"
  }
}
```

---

### 6. Premium Roll (Costs Coins)

Perform a premium roll (costs 100 coins each):

```bash
# Replace {USER_ID} with your user ID
# This example rolls 10 Pokemon (costs 1000 coins)
curl -X POST http://localhost:8080/api/gacha/premium-roll \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "{USER_ID}",
    "count": 10
  }'
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/gacha/premium-roll \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "count": 10
  }'
```

**Bonus:** Every 10th roll is guaranteed Epic or better!

**If insufficient coins:**
```json
{
  "success": false,
  "error": {
    "code": "insufficient_coins",
    "message": "insufficient coins for premium roll"
  }
}
```

---

### 7. Get User's Pokemon Collection

View all Pokemon owned by a user:

```bash
# Replace {USER_ID} with your user ID
curl http://localhost:8080/api/users/{USER_ID}/pokemon
```

**Example:**
```bash
curl http://localhost:8080/api/users/550e8400-e29b-41d4-a716-446655440000/pokemon
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "count": 15,
    "pokemons": [
      {
        "id": "...",
        "species": {
          "id": 6,
          "name": "Charizard",
          "rarity": "epic"
        },
        "nature": "Adamant",
        "level": 5,
        "ivs": { "hp": 31, "attack": 31, ... },
        "stats": { "hp": 25, "attack": 17, ... },
        "iv_percentage": 95.7,
        "estimated_value": 7500
      }
      // ... all your Pokemon
    ]
  }
}
```

---

### 8. Get Specific Pokemon Details

View details of a specific Pokemon:

```bash
# Replace {POKEMON_ID} with the Pokemon's ID
curl http://localhost:8080/api/pokemon/{POKEMON_ID}
```

**Example:**
```bash
curl http://localhost:8080/api/pokemon/650e8400-e29b-41d4-a716-446655440000
```

---

## Testing Flow

Here's a complete testing flow:

```bash
# 1. Check server health
curl http://localhost:8080/health

# 2. Register a new user
USER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"discord_id": "999888777666555444"}')

echo $USER_RESPONSE

# 3. Extract user ID (manual - copy the ID from the response)
USER_ID="<paste-user-id-here>"

# 4. Perform daily roll
curl -X POST http://localhost:8080/api/gacha/daily-roll \
  -H "Content-Type: application/json" \
  -d "{\"user_id\": \"$USER_ID\"}"

# 5. View your Pokemon collection
curl http://localhost:8080/api/users/$USER_ID/pokemon

# 6. Check your user balance
curl http://localhost:8080/api/users/$USER_ID

# 7. Perform a premium roll (if you have coins)
curl -X POST http://localhost:8080/api/gacha/premium-roll \
  -H "Content-Type: application/json" \
  -d "{\"user_id\": \"$USER_ID\", \"count\": 1}"
```

---

## Error Responses

All errors follow this format:

```json
{
  "success": false,
  "error": {
    "code": "error_code",
    "message": "Human readable error message"
  }
}
```

### Common Error Codes

- `bad_request` - Invalid request data
- `not_found` - Resource not found
- `conflict` - Resource already exists (e.g., duplicate Discord ID)
- `cooldown_active` - Daily roll cooldown still active
- `insufficient_coins` - Not enough coins for action
- `internal_server_error` - Server error

---

## Advanced Testing with jq

If you have `jq` installed, you can pretty-print and parse responses:

```bash
# Pretty print response
curl -s http://localhost:8080/health | jq '.'

# Extract specific fields
curl -s http://localhost:8080/api/users/discord/123456789012345678 | jq '.data.coins'

# Count Pokemon in collection
curl -s http://localhost:8080/api/users/$USER_ID/pokemon | jq '.data.count'

# List all Pokemon names
curl -s http://localhost:8080/api/users/$USER_ID/pokemon | jq '.data.pokemons[].species.name'
```

---

## Next Steps

- Test with Postman or Insomnia for a better UI
- Build a Discord bot that calls these endpoints
- Create a React frontend that uses this API
- Add authentication middleware
- Implement the marketplace endpoints
