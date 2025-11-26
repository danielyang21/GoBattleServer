package validators

import (
	"errors"

	"github.com/danielyang21/GoBattleServer/internal/domain"
)

var (
	ErrInvalidIV      = errors.New("IV must be between 0 and 31")
	ErrInvalidLevel   = errors.New("level must be between 1 and 100")
	ErrInvalidSpecies = errors.New("invalid species data")
	ErrInvalidRarity  = errors.New("invalid rarity")
	ErrInvalidNature  = errors.New("invalid nature")
	ErrInvalidStats   = errors.New("invalid base stats")
	ErrInvalidWeight  = errors.New("invalid drop weight")
	ErrEmptyName      = errors.New("name cannot be empty")
)

// ValidateIVs checks if all IVs are in valid range (0-31)
func ValidateIVs(ivs domain.IVs) error {
	if ivs.HP < 0 || ivs.HP > 31 ||
		ivs.Attack < 0 || ivs.Attack > 31 ||
		ivs.Defense < 0 || ivs.Defense > 31 ||
		ivs.SpAttack < 0 || ivs.SpAttack > 31 ||
		ivs.SpDefense < 0 || ivs.SpDefense > 31 ||
		ivs.Speed < 0 || ivs.Speed > 31 {
		return ErrInvalidIV
	}
	return nil
}

// ValidateLevel checks if level is valid (1-100)
func ValidateLevel(level int) error {
	if level < 1 || level > 100 {
		return ErrInvalidLevel
	}
	return nil
}

// ValidateUserPokemon checks if a UserPokemon has valid data
func ValidateUserPokemon(p *domain.UserPokemon) error {
	if err := ValidateIVs(p.IVs); err != nil {
		return err
	}

	if !ValidateNature(p.Nature) {
		return ErrInvalidNature
	}

	if err := ValidateLevel(p.Level); err != nil {
		return err
	}

	return nil
}

// ValidatePokemonSpecies checks if species data is valid
func ValidatePokemonSpecies(s *domain.PokemonSpecies) error {
	if s.ID <= 0 {
		return ErrInvalidSpecies
	}

	if s.Name == "" {
		return ErrEmptyName
	}

	if !ValidateRarity(s.Rarity) {
		return ErrInvalidRarity
	}

	if s.BaseHP <= 0 || s.BaseAttack < 0 || s.BaseDefense < 0 ||
		s.BaseSpAttack < 0 || s.BaseSpDefense < 0 || s.BaseSpeed < 0 {
		return ErrInvalidStats
	}

	if s.DropWeight <= 0 {
		return ErrInvalidWeight
	}

	return nil
}
