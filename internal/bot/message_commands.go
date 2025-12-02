package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// MessageCommandPrefix is the prefix for message-based commands
var MessageCommandPrefix = "!"

// SetCommandPrefix sets a custom command prefix
func (b *Bot) SetCommandPrefix(prefix string) {
	MessageCommandPrefix = prefix
}

// handleMessageCommand handles traditional message commands (!daily, !roll, etc.)
func (b *Bot) handleMessageCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages
	if m.Author.Bot {
		return
	}

	// Check if message starts with prefix
	if !strings.HasPrefix(m.Content, MessageCommandPrefix) {
		return
	}

	// Remove prefix and split into command and args
	content := strings.TrimPrefix(m.Content, MessageCommandPrefix)
	parts := strings.Fields(content)
	if len(parts) == 0 {
		return
	}

	command := strings.ToLower(parts[0])
	args := parts[1:]

	// Route to appropriate handler
	switch command {
	case "daily":
		b.handleDailyMessage(s, m)
	case "roll":
		b.handleRollMessage(s, m, args)
	case "balance", "bal", "coins":
		b.handleBalanceMessage(s, m)
	case "box", "collection":
		b.handleBoxMessage(s, m, args)
	case "help":
		b.handleHelpMessage(s, m)
	default:
		// Unknown command - ignore
		return
	}
}

// handleDailyMessage handles !daily command
func (b *Bot) handleDailyMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	discordID := m.Author.ID

	// Get or create user
	user, err := b.apiClient.GetOrCreateUser(discordID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Failed to get user: "+err.Error())
		return
	}

	// Perform daily roll
	pokemons, err := b.apiClient.DailyRoll(user.ID)
	if err != nil {
		if strings.Contains(err.Error(), "cooldown") {
			s.ChannelMessageSend(m.ChannelID, "⏰ You've already claimed your daily roll! Come back tomorrow.")
		} else {
			s.ChannelMessageSend(m.ChannelID, "❌ Failed to roll: "+err.Error())
		}
		return
	}

	// Build embed
	embed := &discordgo.MessageEmbed{
		Title:       "🎴 Daily Roll Results!",
		Description: fmt.Sprintf("You received **%d Pokemon**:", len(pokemons)),
		Color:       0x00ff00,
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}

	for i, p := range pokemons {
		rarityEmoji := getRarityEmoji(p.Species.Rarity)
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: fmt.Sprintf("%d. %s %s", i+1, rarityEmoji, p.Species.Name),
			Value: fmt.Sprintf(
				"**Nature:** %s\n**IVs:** %.1f%% perfect\n**Value:** %d coins",
				p.Nature, p.IVPercentage, p.EstimatedValue,
			),
			Inline: false,
		})
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// handleRollMessage handles !roll <count> command
func (b *Bot) handleRollMessage(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "❌ Usage: `!roll <count>` (e.g., `!roll 10`)")
		return
	}

	count, err := strconv.Atoi(args[0])
	if err != nil || count < 1 || count > 10 {
		s.ChannelMessageSend(m.ChannelID, "❌ Count must be a number between 1 and 10")
		return
	}

	discordID := m.Author.ID

	// Get user
	user, err := b.apiClient.GetOrCreateUser(discordID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Failed to get user: "+err.Error())
		return
	}

	// Check if user has enough coins
	cost := count * 100
	if user.Coins < cost {
		s.ChannelMessageSend(m.ChannelID,
			fmt.Sprintf("❌ Insufficient coins! You need %d coins but only have %d.", cost, user.Coins))
		return
	}

	// Perform premium roll
	pokemons, err := b.apiClient.PremiumRoll(user.ID, count)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Failed to roll: "+err.Error())
		return
	}

	// Build embed
	embed := &discordgo.MessageEmbed{
		Title:       "💎 Premium Roll Results!",
		Description: fmt.Sprintf("You spent **%d coins** and received **%d Pokemon**:", cost, len(pokemons)),
		Color:       0xffd700,
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}

	for i, p := range pokemons {
		rarityEmoji := getRarityEmoji(p.Species.Rarity)
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: fmt.Sprintf("%d. %s %s", i+1, rarityEmoji, p.Species.Name),
			Value: fmt.Sprintf(
				"**Nature:** %s | **IVs:** %.1f%%\n**Value:** %d coins",
				p.Nature, p.IVPercentage, p.EstimatedValue,
			),
			Inline: true,
		})
	}

	if count >= 10 {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: "🎁 10-roll bonus: Guaranteed Epic or better!",
		}
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// handleBalanceMessage handles !balance command
func (b *Bot) handleBalanceMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	discordID := m.Author.ID

	// Get user
	user, err := b.apiClient.GetOrCreateUser(discordID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Failed to get user: "+err.Error())
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "💰 Your Balance",
		Description: fmt.Sprintf("You have **%d coins**", user.Coins),
		Color:       0xffd700,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "💵 Premium Roll",
				Value:  "100 coins per Pokemon",
				Inline: true,
			},
			{
				Name:   "🎁 10-Roll Bonus",
				Value:  "1000 coins (guaranteed Epic+)",
				Inline: true,
			},
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// handleBoxMessage handles !box [rarity] command
func (b *Bot) handleBoxMessage(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	discordID := m.Author.ID

	// Get user
	user, err := b.apiClient.GetOrCreateUser(discordID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Failed to get user: "+err.Error())
		return
	}

	// Get Pokemon collection
	pokemons, err := b.apiClient.GetUserPokemon(user.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Failed to get Pokemon: "+err.Error())
		return
	}

	// Filter by rarity if specified
	var rarityFilter string
	if len(args) > 0 {
		rarityFilter = strings.ToLower(args[0])
	}

	filtered := pokemons
	if rarityFilter != "" {
		filtered = make([]Pokemon, 0)
		for _, p := range pokemons {
			if p.Species.Rarity == rarityFilter {
				filtered = append(filtered, p)
			}
		}
	}

	if len(filtered) == 0 {
		s.ChannelMessageSend(m.ChannelID, "You don't have any Pokemon in your collection yet! Use `!daily` to get started.")
		return
	}

	// Calculate total value
	totalValue := 0
	rarityCounts := make(map[string]int)
	for _, p := range pokemons {
		totalValue += p.EstimatedValue
		rarityCounts[p.Species.Rarity]++
	}

	// Build embed
	title := "📦 Your Pokemon Collection"
	if rarityFilter != "" {
		title = fmt.Sprintf("📦 Your %s Pokemon", strings.Title(rarityFilter))
	}

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: fmt.Sprintf("**Total Pokemon:** %d\n**Collection Value:** %d coins", len(pokemons), totalValue),
		Color:       0x3498db,
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}

	// Show rarity breakdown
	rarityBreakdown := ""
	for _, rarity := range []string{"mythic", "legendary", "epic", "rare", "uncommon", "common"} {
		if count := rarityCounts[rarity]; count > 0 {
			rarityBreakdown += fmt.Sprintf("%s **%s**: %d\n", getRarityEmoji(rarity), strings.Title(rarity), count)
		}
	}
	if rarityBreakdown != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "📊 Rarity Breakdown",
			Value:  rarityBreakdown,
			Inline: false,
		})
	}

	// Show top Pokemon (limit to 10)
	limit := 10
	if len(filtered) < limit {
		limit = len(filtered)
	}

	for i := 0; i < limit; i++ {
		p := filtered[i]
		rarityEmoji := getRarityEmoji(p.Species.Rarity)
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: fmt.Sprintf("%s %s", rarityEmoji, p.Species.Name),
			Value: fmt.Sprintf(
				"**Nature:** %s\n**IVs:** %.1f%%\n**Value:** %d coins",
				p.Nature, p.IVPercentage, p.EstimatedValue,
			),
			Inline: true,
		})
	}

	if len(filtered) > limit {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Showing %d of %d Pokemon", limit, len(filtered)),
		}
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// handleHelpMessage shows available commands
func (b *Bot) handleHelpMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "🎮 Pokemon Gacha Bot Commands",
		Description: fmt.Sprintf("Use `%s` prefix for commands", MessageCommandPrefix),
		Color:       0x3498db,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("%sdaily", MessageCommandPrefix),
				Value:  "Claim your free daily roll (5 Pokemon)",
				Inline: false,
			},
			{
				Name:   fmt.Sprintf("%sroll <count>", MessageCommandPrefix),
				Value:  "Buy premium rolls (1-10) with coins\nExample: `!roll 10`",
				Inline: false,
			},
			{
				Name:   fmt.Sprintf("%sbalance", MessageCommandPrefix),
				Value:  "Check your coin balance\nAliases: `!bal`, `!coins`",
				Inline: false,
			},
			{
				Name:   fmt.Sprintf("%sbox [rarity]", MessageCommandPrefix),
				Value:  "View your Pokemon collection\nExample: `!box epic`\nAlias: `!collection`",
				Inline: false,
			},
			{
				Name:   fmt.Sprintf("%shelp", MessageCommandPrefix),
				Value:  "Show this help message",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "You can also use slash commands: /daily, /roll, /balance, /box",
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
