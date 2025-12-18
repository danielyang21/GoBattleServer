package domain

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/google/uuid"
)

// TurnResolver handles turn resolution logic
type TurnResolver struct {
	damageCalc *DamageCalculator
	rand       *rand.Rand
}

// NewTurnResolver creates a new turn resolver
func NewTurnResolver(source rand.Source) *TurnResolver {
	return &TurnResolver{
		damageCalc: NewDamageCalculator(source),
		rand:       rand.New(source),
	}
}

// ResolveTurn resolves a complete turn of battle
func (tr *TurnResolver) ResolveTurn(state *BattleState) *TurnResolution {
	resolution := &TurnResolution{
		Turn:    state.Turn,
		Actions: []*ResolvedAction{},
	}

	// Check if both players have submitted actions
	if !state.BothPlayersReady() {
		return resolution
	}

	// Determine turn order based on priority and speed
	actions := tr.DetermineTurnOrder(state)

	// Execute each action in order
	for _, action := range actions {
		if action == nil {
			continue
		}

		// Check if Pokemon has fainted
		actor := state.GetPlayer(action.PlayerID)
		if actor == nil || actor.Pokemon.Fainted {
			continue
		}

		// Execute the action
		resolvedAction := tr.ExecuteAction(state, action)
		resolution.Actions = append(resolution.Actions, resolvedAction)

		// Check if battle ended
		if tr.IsBattleOver(state) {
			resolution.BattleEnded = true
			winner := tr.DetermineWinner(state)
			resolution.Winner = &winner
			break
		}
	}

	// End of turn effects (if battle not over)
	if !resolution.BattleEnded {
		tr.ApplyEndOfTurnEffects(state, resolution)

		// Check again if battle ended from end-of-turn effects
		if tr.IsBattleOver(state) {
			resolution.BattleEnded = true
			winner := tr.DetermineWinner(state)
			resolution.Winner = &winner
		}
	}

	// Clear actions and increment turn
	state.ClearActions()
	state.Turn++

	return resolution
}

// DetermineTurnOrder determines the order of actions based on priority and speed
func (tr *TurnResolver) DetermineTurnOrder(state *BattleState) []*BattleAction {
	actions := []*BattleAction{}

	// Add both player actions
	if state.Player1Action != nil {
		actions = append(actions, state.Player1Action)
	}
	if state.Player2Action != nil {
		actions = append(actions, state.Player2Action)
	}

	// Calculate priority and speed for each action
	for _, action := range actions {
		if action.Type == ActionMove && action.Move != nil {
			// Base priority from move
			action.Priority = action.Move.Priority

			// Get player and calculate speed
			player := state.GetPlayer(action.PlayerID)
			if player != nil && player.Pokemon != nil {
				// Calculate speed with stat stages
				speed := player.Pokemon.Stats.Speed
				speedMultiplier := player.Pokemon.StatStages.GetMultiplier(Speed)
				action.Speed = int(float64(speed) * speedMultiplier)

				// Apply paralysis speed reduction (50%)
				if player.Pokemon.Status == StatusParalysis {
					action.Speed = action.Speed / 2
				}

				// Check for priority-boosting abilities (Prankster, Gale Wings, etc.)
				// This would need ability data loaded
				// For now, we'll leave it as is
			}
		} else if action.Type == ActionForfeit {
			// Forfeit has lowest priority
			action.Priority = -8
			action.Speed = 0
		}
	}

	// Sort by priority (descending), then speed (descending), then random for ties
	sort.Slice(actions, func(i, j int) bool {
		// Higher priority goes first
		if actions[i].Priority != actions[j].Priority {
			return actions[i].Priority > actions[j].Priority
		}

		// Higher speed goes first
		if actions[i].Speed != actions[j].Speed {
			return actions[i].Speed > actions[j].Speed
		}

		// Random tiebreaker
		return tr.rand.Intn(2) == 0
	})

	return actions
}

// ExecuteAction executes a single action
func (tr *TurnResolver) ExecuteAction(state *BattleState, action *BattleAction) *ResolvedAction {
	resolved := &ResolvedAction{
		PlayerID:   action.PlayerID,
		ActionType: action.Type,
		Messages:   []string{},
	}

	player := state.GetPlayer(action.PlayerID)
	opponent := state.GetOpponent(action.PlayerID)

	if player == nil || opponent == nil {
		resolved.Failed = true
		resolved.FailReason = "Invalid player"
		return resolved
	}

	switch action.Type {
	case ActionMove:
		return tr.ExecuteMove(state, player, opponent, action)
	case ActionForfeit:
		resolved.Messages = append(resolved.Messages, fmt.Sprintf("%s forfeited the battle!", player.UserID))
		opponent.Pokemon.CurrentHP = 0 // Mark as winner by default logic
		return resolved
	default:
		resolved.Failed = true
		resolved.FailReason = "Unknown action type"
		return resolved
	}
}

// ExecuteMove executes a move action
func (tr *TurnResolver) ExecuteMove(state *BattleState, attacker, defender *BattlePlayer, action *BattleAction) *ResolvedAction {
	resolved := &ResolvedAction{
		PlayerID:   action.PlayerID,
		ActionType: ActionMove,
		Move:       action.Move,
		Target:     defender.UserID,
		Messages:   []string{},
	}

	// Check if Pokemon can move (status conditions)
	if !tr.CanPokemonMove(attacker.Pokemon, resolved) {
		resolved.Failed = true
		return resolved
	}

	// Use the move
	resolved.Messages = append(resolved.Messages,
		fmt.Sprintf("%s used %s!", attacker.Pokemon.Species.Name, action.Move.Name))

	// Decrement PP
	attacker.Pokemon.DecrementPP(action.MoveIndex)

	// Handle status moves separately
	if action.Move.Category == Status {
		return tr.ExecuteStatusMove(state, attacker, defender, action, resolved)
	}

	// Calculate damage
	damageCtx := tr.BuildDamageContext(state, attacker, defender, action.Move)
	damageResult := tr.damageCalc.CalculateDamage(damageCtx)

	// Check if move missed
	if damageResult.Damage == 0 && damageResult.Effectiveness > 0 {
		resolved.Messages = append(resolved.Messages, "But it missed!")
		resolved.Failed = true
		return resolved
	}

	// Apply damage
	defender.Pokemon.TakeDamage(damageResult.Damage)
	resolved.Result = damageResult

	// Add damage message
	damageMsg := fmt.Sprintf("%s took %d damage!", defender.Pokemon.Species.Name, damageResult.Damage)
	if damageResult.IsCritical {
		damageMsg = "A critical hit! " + damageMsg
	}
	resolved.Messages = append(resolved.Messages, damageMsg)

	// Type effectiveness message
	if damageResult.Effectiveness > 1.0 {
		resolved.Messages = append(resolved.Messages, "It's super effective!")
	} else if damageResult.Effectiveness < 1.0 && damageResult.Effectiveness > 0 {
		resolved.Messages = append(resolved.Messages, "It's not very effective...")
	} else if damageResult.Effectiveness == 0 {
		resolved.Messages = append(resolved.Messages, "It doesn't affect the foe!")
		return resolved
	}

	// Check if defender fainted
	if damageResult.Fainted {
		resolved.Messages = append(resolved.Messages, fmt.Sprintf("%s fainted!", defender.Pokemon.Species.Name))
	}

	// Apply recoil damage
	if damageResult.RecoilDamage > 0 {
		attacker.Pokemon.TakeDamage(damageResult.RecoilDamage)
		resolved.Messages = append(resolved.Messages,
			fmt.Sprintf("%s took %d recoil damage!", attacker.Pokemon.Species.Name, damageResult.RecoilDamage))
	}

	// Apply drain/healing
	if damageResult.DrainAmount > 0 {
		healed := attacker.Pokemon.Heal(damageResult.DrainAmount)
		resolved.Messages = append(resolved.Messages,
			fmt.Sprintf("%s restored %d HP!", attacker.Pokemon.Species.Name, healed))
	}

	// Apply secondary effects
	tr.ApplySecondaryEffects(attacker, defender, action.Move, resolved)

	// Log to battle state
	state.AddLogEntry("move", resolved.Messages[0], map[string]interface{}{
		"attacker": attacker.UserID,
		"defender": defender.UserID,
		"move":     action.Move.Name,
		"damage":   damageResult.Damage,
	})

	return resolved
}

// ExecuteStatusMove executes a status move
func (tr *TurnResolver) ExecuteStatusMove(state *BattleState, attacker, defender *BattlePlayer, action *BattleAction, resolved *ResolvedAction) *ResolvedAction {
	move := action.Move

	// Apply stat changes to self
	for _, statChange := range move.StatChanges {
		if statChange.Target == "self" {
			actualChange := attacker.Pokemon.StatStages.ApplyChange(statChange.Stat, statChange.Stages)
			if actualChange != 0 {
				direction := "rose"
				if actualChange < 0 {
					direction = "fell"
				}
				resolved.Messages = append(resolved.Messages,
					fmt.Sprintf("%s's %s %s!", attacker.Pokemon.Species.Name, statChange.Stat, direction))
			}
		} else if statChange.Target == "opponent" {
			actualChange := defender.Pokemon.StatStages.ApplyChange(statChange.Stat, statChange.Stages)
			if actualChange != 0 {
				direction := "rose"
				if actualChange < 0 {
					direction = "fell"
				}
				resolved.Messages = append(resolved.Messages,
					fmt.Sprintf("%s's %s %s!", defender.Pokemon.Species.Name, statChange.Stat, direction))
			}
		}
	}

	// Apply status condition
	if move.StatusInflict != nil {
		if tr.TryInflictStatus(defender.Pokemon, move.StatusInflict.Status, 100, resolved) {
			resolved.Messages = append(resolved.Messages,
				fmt.Sprintf("%s was inflicted with %s!", defender.Pokemon.Species.Name, move.StatusInflict.Status))
		}
	}

	// Apply weather
	if move.WeatherEffect != nil {
		state.Weather = move.WeatherEffect.Weather
		state.WeatherTurns = move.WeatherEffect.Duration
		resolved.Messages = append(resolved.Messages, tr.GetWeatherSetMessage(move.WeatherEffect.Weather))
	}

	// Apply terrain
	if move.TerrainEffect != nil {
		state.Terrain = move.TerrainEffect.Terrain
		state.TerrainTurns = move.TerrainEffect.Duration
		resolved.Messages = append(resolved.Messages, tr.GetTerrainSetMessage(move.TerrainEffect.Terrain))
	}

	// Apply entry hazards
	if move.EntryHazard != nil {
		tr.ApplyEntryHazard(state, attacker.UserID, move.EntryHazard, resolved)
	}

	// Apply healing
	if move.HealPercent > 0 {
		healAmount := (attacker.Pokemon.MaxHP * move.HealPercent) / 100
		healed := attacker.Pokemon.Heal(healAmount)
		resolved.Messages = append(resolved.Messages,
			fmt.Sprintf("%s restored %d HP!", attacker.Pokemon.Species.Name, healed))
	}

	return resolved
}

// CanPokemonMove checks if a Pokemon can execute its move this turn
func (tr *TurnResolver) CanPokemonMove(pokemon *BattlePokemon, resolved *ResolvedAction) bool {
	// Check sleep
	if pokemon.Status == StatusSleep {
		if pokemon.StatusTurns > 0 {
			pokemon.StatusTurns--
			resolved.Messages = append(resolved.Messages, fmt.Sprintf("%s is fast asleep!", pokemon.Species.Name))
			if pokemon.StatusTurns == 0 {
				pokemon.Status = StatusNone
				resolved.Messages = append(resolved.Messages, fmt.Sprintf("%s woke up!", pokemon.Species.Name))
			}
			return false
		}
	}

	// Check freeze (20% chance to thaw)
	if pokemon.Status == StatusFreeze {
		if tr.rand.Intn(100) < 20 {
			pokemon.Status = StatusNone
			resolved.Messages = append(resolved.Messages, fmt.Sprintf("%s thawed out!", pokemon.Species.Name))
			return true
		}
		resolved.Messages = append(resolved.Messages, fmt.Sprintf("%s is frozen solid!", pokemon.Species.Name))
		return false
	}

	// Check paralysis (25% chance of full paralysis)
	if pokemon.Status == StatusParalysis {
		if tr.rand.Intn(100) < 25 {
			resolved.Messages = append(resolved.Messages, fmt.Sprintf("%s is fully paralyzed!", pokemon.Species.Name))
			return false
		}
	}

	return true
}

// ApplySecondaryEffects applies secondary effects of a move
func (tr *TurnResolver) ApplySecondaryEffects(attacker, defender *BattlePlayer, move *Move, resolved *ResolvedAction) {
	if move.SecondaryEffect == nil {
		return
	}

	effect := move.SecondaryEffect

	// Roll for secondary effect
	if effect.Chance > 0 && tr.rand.Intn(100) >= effect.Chance {
		return // Effect doesn't trigger
	}

	// Apply stat changes
	for _, statChange := range effect.StatChanges {
		var target *BattlePokemon
		if statChange.Target == "self" {
			target = attacker.Pokemon
		} else {
			target = defender.Pokemon
		}

		actualChange := target.StatStages.ApplyChange(statChange.Stat, statChange.Stages)
		if actualChange != 0 {
			direction := "rose"
			if actualChange < 0 {
				direction = "fell"
			}
			resolved.Messages = append(resolved.Messages,
				fmt.Sprintf("%s's %s %s!", target.Species.Name, statChange.Stat, direction))
		}
	}

	// Apply status
	if effect.StatusInflict != nil {
		tr.TryInflictStatus(defender.Pokemon, effect.StatusInflict.Status, effect.StatusInflict.Chance, resolved)
	}

	// Apply flinch
	if effect.FlinchChance > 0 && tr.rand.Intn(100) < effect.FlinchChance {
		resolved.Messages = append(resolved.Messages, fmt.Sprintf("%s flinched!", defender.Pokemon.Species.Name))
	}
}

// TryInflictStatus attempts to inflict a status condition
func (tr *TurnResolver) TryInflictStatus(pokemon *BattlePokemon, status StatusCondition, chance int, resolved *ResolvedAction) bool {
	// Check if already has status
	if pokemon.Status != StatusNone {
		return false
	}

	// Roll for chance
	if chance < 100 && tr.rand.Intn(100) >= chance {
		return false
	}

	// Apply status
	pokemon.Status = status

	// Set status turns for sleep (1-3 turns)
	if status == StatusSleep {
		pokemon.StatusTurns = 1 + tr.rand.Intn(3)
	}

	return true
}

// ApplyEndOfTurnEffects applies all end-of-turn effects
func (tr *TurnResolver) ApplyEndOfTurnEffects(state *BattleState, resolution *TurnResolution) {
	// Weather effects
	tr.ApplyWeatherEffects(state, resolution)

	// Status damage
	tr.ApplyStatusDamage(state, resolution)

	// Healing effects (Leftovers, Grassy Terrain, etc.)
	tr.ApplyHealingEffects(state, resolution)

	// Decrement weather/terrain turns
	if state.WeatherTurns > 0 {
		state.WeatherTurns--
		if state.WeatherTurns == 0 {
			state.Weather = WeatherNone
			state.AddLogEntry("weather", "The weather returned to normal", nil)
		}
	}

	if state.TerrainTurns > 0 {
		state.TerrainTurns--
		if state.TerrainTurns == 0 {
			state.Terrain = TerrainNone
			state.AddLogEntry("terrain", "The terrain faded", nil)
		}
	}
}

// ApplyWeatherEffects applies damage from weather
func (tr *TurnResolver) ApplyWeatherEffects(state *BattleState, resolution *TurnResolution) {
	if state.Weather == WeatherNone {
		return
	}

	var damageAmount float64
	switch state.Weather {
	case WeatherSandstorm:
		damageAmount = 1.0 / 16.0 // 1/16 max HP
	case WeatherHail, WeatherSnow:
		damageAmount = 1.0 / 16.0
	default:
		return
	}

	// Apply to both players if not immune
	for _, player := range []*BattlePlayer{state.Player1, state.Player2} {
		if player.Pokemon.Fainted {
			continue
		}

		// Check immunity (Rock/Ground/Steel for Sandstorm, Ice for Hail)
		immune := false
		if state.Weather == WeatherSandstorm {
			if player.Pokemon.Species.Type1 == Rock || player.Pokemon.Species.Type1 == Ground || player.Pokemon.Species.Type1 == Steel {
				immune = true
			}
		} else if state.Weather == WeatherHail || state.Weather == WeatherSnow {
			if player.Pokemon.Species.Type1 == Ice {
				immune = true
			}
		}

		if !immune {
			damage := int(float64(player.Pokemon.MaxHP) * damageAmount)
			if damage < 1 {
				damage = 1
			}
			player.Pokemon.TakeDamage(damage)
			resolution.WeatherDamage = append(resolution.WeatherDamage, WeatherDamage{
				PlayerID: player.UserID,
				Weather:  state.Weather,
				Damage:   damage,
			})
		}
	}
}

// ApplyStatusDamage applies damage from status conditions
func (tr *TurnResolver) ApplyStatusDamage(state *BattleState, resolution *TurnResolution) {
	for _, player := range []*BattlePlayer{state.Player1, state.Player2} {
		if player.Pokemon.Fainted {
			continue
		}

		damage := 0
		switch player.Pokemon.Status {
		case StatusBurn:
			damage = player.Pokemon.MaxHP / 16 // 1/16 max HP
		case StatusPoison:
			damage = player.Pokemon.MaxHP / 8 // 1/8 max HP
		case StatusBadlyPoison:
			// Toxic increases each turn
			player.Pokemon.StatusTurns++
			damage = (player.Pokemon.MaxHP * player.Pokemon.StatusTurns) / 16
		}

		if damage > 0 {
			player.Pokemon.TakeDamage(damage)
			resolution.StatusDamage = append(resolution.StatusDamage, StatusDamage{
				PlayerID: player.UserID,
				Status:   player.Pokemon.Status,
				Damage:   damage,
			})
		}
	}
}

// ApplyHealingEffects applies end-of-turn healing
func (tr *TurnResolver) ApplyHealingEffects(state *BattleState, resolution *TurnResolution) {
	for _, player := range []*BattlePlayer{state.Player1, state.Player2} {
		if player.Pokemon.Fainted {
			continue
		}

		// Grassy Terrain healing
		if state.Terrain == TerrainGrassy {
			healAmount := player.Pokemon.MaxHP / 16
			healed := player.Pokemon.Heal(healAmount)
			if healed > 0 {
				resolution.EndOfTurnHeals = append(resolution.EndOfTurnHeals, EndOfTurnHeal{
					PlayerID: player.UserID,
					Source:   "grassy_terrain",
					Amount:   healed,
				})
			}
		}

		// TODO: Add Leftovers, Black Sludge, etc. when items are implemented
	}
}

// ApplyEntryHazard applies an entry hazard
func (tr *TurnResolver) ApplyEntryHazard(state *BattleState, attackerID uuid.UUID, hazard *EntryHazard, resolved *ResolvedAction) {
	var opponentHazards *EntryHazards
	if state.Player1.UserID == attackerID {
		opponentHazards = state.Player2Hazards
	} else {
		opponentHazards = state.Player1Hazards
	}

	switch hazard.HazardType {
	case HazardStealthRock:
		if !opponentHazards.StealthRock {
			opponentHazards.StealthRock = true
			resolved.Messages = append(resolved.Messages, "Pointed stones float in the air around the foe!")
		}
	case HazardSpikes:
		if opponentHazards.Spikes < 3 {
			opponentHazards.Spikes++
			resolved.Messages = append(resolved.Messages, "Spikes were scattered around the foe!")
		}
	case HazardToxicSpikes:
		if opponentHazards.ToxicSpikes < 2 {
			opponentHazards.ToxicSpikes++
			resolved.Messages = append(resolved.Messages, "Poison spikes were scattered around the foe!")
		}
	case HazardStickyWeb:
		if !opponentHazards.StickyWeb {
			opponentHazards.StickyWeb = true
			resolved.Messages = append(resolved.Messages, "A sticky web has been laid out beneath the foe!")
		}
	}
}

// BuildDamageContext builds the context for damage calculation
func (tr *TurnResolver) BuildDamageContext(state *BattleState, attacker, defender *BattlePlayer, move *Move) *DamageContext {
	return &DamageContext{
		Attacker:        attacker.Pokemon,
		AttackerPlayer:  attacker,
		Defender:        defender.Pokemon,
		DefenderPlayer:  defender,
		Move:            move,
		Weather:         state.Weather,
		Terrain:         state.Terrain,
		Turn:            state.Turn,
		// Abilities and items would be loaded here from database
		AttackerAbility: nil,
		AttackerItem:    nil,
		DefenderAbility: nil,
		DefenderItem:    nil,
	}
}

// IsBattleOver checks if the battle has ended
func (tr *TurnResolver) IsBattleOver(state *BattleState) bool {
	return state.Player1.Pokemon.Fainted || state.Player2.Pokemon.Fainted
}

// DetermineWinner determines the winner of the battle
func (tr *TurnResolver) DetermineWinner(state *BattleState) uuid.UUID {
	if state.Player2.Pokemon.Fainted {
		return state.Player1.UserID
	}
	return state.Player2.UserID
}

// GetWeatherSetMessage returns the message for weather being set
func (tr *TurnResolver) GetWeatherSetMessage(weather Weather) string {
	switch weather {
	case WeatherSun:
		return "The sunlight turned harsh!"
	case WeatherRain:
		return "It started to rain!"
	case WeatherSandstorm:
		return "A sandstorm kicked up!"
	case WeatherHail:
		return "It started to hail!"
	case WeatherSnow:
		return "It started to snow!"
	default:
		return ""
	}
}

// GetTerrainSetMessage returns the message for terrain being set
func (tr *TurnResolver) GetTerrainSetMessage(terrain Terrain) string {
	switch terrain {
	case TerrainElectric:
		return "An electric current ran across the battlefield!"
	case TerrainGrassy:
		return "Grass grew to cover the battlefield!"
	case TerrainMisty:
		return "Mist swirled about the battlefield!"
	case TerrainPsychic:
		return "The battlefield got weird!"
	default:
		return ""
	}
}
