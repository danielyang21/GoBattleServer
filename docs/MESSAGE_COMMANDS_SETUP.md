# Message Commands Setup (! prefix)

Your bot now supports **BOTH** slash commands (`/daily`) and message commands (`!daily`)!

## ⚠️ Important: Enable Message Content Intent

For message commands to work, you MUST enable the Message Content Intent in Discord Developer Portal:

### Step 1: Go to Discord Developer Portal
https://discord.com/developers/applications

### Step 2: Select Your Bot Application

### Step 3: Go to "Bot" Tab

### Step 4: Scroll to "Privileged Gateway Intents"
Enable:
- ✅ **MESSAGE CONTENT INTENT** (Required!)
- ✅ **SERVER MEMBERS INTENT** (Optional)

### Step 5: Click "Save Changes"

**⚠️ WITHOUT THIS, MESSAGE COMMANDS WON'T WORK!**

---

## Available Commands

Now you can use EITHER format:

### Slash Commands (Original)
- `/daily` - Free daily roll
- `/roll <count>` - Premium roll
- `/balance` - Check coins
- `/box [rarity]` - View collection

### Message Commands (NEW!)
- `!daily` - Free daily roll
- `!roll 10` - Premium roll (specify count)
- `!balance` - Check coins (also: `!bal`, `!coins`)
- `!box` - View collection
- `!box epic` - Filter by rarity
- `!help` - Show all commands

---

## Usage Examples

```
!daily
→ Claims your free 5 Pokemon

!roll 10
→ Buys 10 premium rolls (costs 1000 coins)

!balance
→ Shows your coin balance

!box epic
→ Shows only your Epic Pokemon

!help
→ Shows command list
```

---

## Changing the Prefix

Want to use a different prefix? (e.g., `?daily`, `.daily`, `$daily`)

Edit `cmd/bot/main.go` and add this after creating the bot:

```go
// Create bot
discordBot, err := bot.NewBot(token, apiURL)
if err != nil {
    log.Fatalf("Failed to create bot: %v", err)
}

// Change prefix (add this line)
discordBot.SetCommandPrefix("?")  // Now use ?daily instead of !daily
```

---

## Testing

1. **Enable Message Content Intent** (see above)
2. **Restart your bot:**
   ```bash
   ./run-bot.sh
   ```
3. **Try message commands in Discord:**
   ```
   !daily
   !balance
   !box
   !help
   ```

---

## How It Works

### Message Flow
```
User types: !daily
    ↓
Discord sends message event
    ↓
Bot checks if message starts with "!"
    ↓
Bot parses command and arguments
    ↓
Bot calls API and returns result
```

### Slash vs Message Commands

| Feature | Slash (`/`) | Message (`!`) |
|---------|------------|---------------|
| Discord UI autocomplete | ✅ Yes | ❌ No |
| Permission control | ✅ Built-in | ❌ Manual |
| Argument validation | ✅ Automatic | ❌ Manual |
| Classic feel | ❌ No | ✅ Yes |
| Works in DMs | ✅ Yes | ✅ Yes |
| Requires intent | ❌ No | ✅ Yes |

---

## Aliases

The bot supports command aliases:

- `!balance` = `!bal` = `!coins`
- `!box` = `!collection`

---

## Command Arguments

### `!roll <count>`
Count is required:
```
!roll 1    ✅ Rolls 1 Pokemon
!roll 10   ✅ Rolls 10 Pokemon
!roll      ❌ Error: count required
!roll 99   ❌ Error: max is 10
```

### `!box [rarity]`
Rarity is optional:
```
!box           ✅ Shows all Pokemon
!box epic      ✅ Shows only Epic
!box legendary ✅ Shows only Legendary
```

---

## Troubleshooting

### Message commands don't work
- ✅ Did you enable **Message Content Intent** in Discord Portal?
- ✅ Did you restart the bot after enabling it?
- ✅ Is the bot online?

### Bot responds to other bots
- The bot automatically ignores messages from other bots

### Commands are case-sensitive
- No! All commands work in any case: `!daily`, `!DAILY`, `!Daily`

---

## What's Next?

Now that you have both command types, you can:
- Let users choose their preferred style
- Keep slash commands for new users (easier to discover)
- Keep message commands for power users (faster to type)

Enjoy! 🎮
