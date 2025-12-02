package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type APIResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   *APIError       `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type User struct {
	ID            string `json:"id"`
	DiscordID     string `json:"discord_id"`
	Coins         int    `json:"coins"`
	LastDailyRoll string `json:"last_daily_roll,omitempty"`
	CreatedAt     string `json:"created_at"`
}

type Pokemon struct {
	ID             string  `json:"id"`
	Species        Species `json:"species"`
	Nature         string  `json:"nature"`
	Level          int     `json:"level"`
	IVs            IVs     `json:"ivs"`
	Stats          Stats   `json:"stats"`
	IVPercentage   float64 `json:"iv_percentage"`
	EstimatedValue int     `json:"estimated_value"`
}

type Species struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Rarity string `json:"rarity"`
}

type IVs struct {
	HP        int `json:"hp"`
	Attack    int `json:"attack"`
	Defense   int `json:"defense"`
	SpAttack  int `json:"sp_attack"`
	SpDefense int `json:"sp_defense"`
	Speed     int `json:"speed"`
}

type Stats struct {
	HP        int `json:"hp"`
	Attack    int `json:"attack"`
	Defense   int `json:"defense"`
	SpAttack  int `json:"sp_attack"`
	SpDefense int `json:"sp_defense"`
	Speed     int `json:"speed"`
}

func (c *APIClient) RegisterUser(discordID string) (*User, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"discord_id": discordID,
	})

	resp, err := c.httpClient.Post(
		c.baseURL+"/api/users/register",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("%s: %s", apiResp.Error.Code, apiResp.Error.Message)
	}

	var user User
	if err := json.Unmarshal(apiResp.Data, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *APIClient) GetUserByDiscordID(discordID string) (*User, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/api/users/discord/" + discordID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("%s: %s", apiResp.Error.Code, apiResp.Error.Message)
	}

	var user User
	if err := json.Unmarshal(apiResp.Data, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *APIClient) GetOrCreateUser(discordID string) (*User, error) {
	user, err := c.GetUserByDiscordID(discordID)
	if err == nil {
		return user, nil
	}

	return c.RegisterUser(discordID)
}

func (c *APIClient) DailyRoll(userID string) ([]Pokemon, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"user_id": userID,
	})

	resp, err := c.httpClient.Post(
		c.baseURL+"/api/gacha/daily-roll",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("%s: %s", apiResp.Error.Code, apiResp.Error.Message)
	}

	var result struct {
		Pokemons []Pokemon `json:"pokemons"`
		Count    int       `json:"count"`
	}
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, err
	}

	return result.Pokemons, nil
}

func (c *APIClient) PremiumRoll(userID string, count int) ([]Pokemon, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"user_id": userID,
		"count":   count,
	})

	resp, err := c.httpClient.Post(
		c.baseURL+"/api/gacha/premium-roll",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("%s: %s", apiResp.Error.Code, apiResp.Error.Message)
	}

	var result struct {
		Pokemons []Pokemon `json:"pokemons"`
		Count    int       `json:"count"`
	}
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, err
	}

	return result.Pokemons, nil
}

func (c *APIClient) GetUserPokemon(userID string) ([]Pokemon, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/api/users/" + userID + "/pokemon")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("%s: %s", apiResp.Error.Code, apiResp.Error.Message)
	}

	var result struct {
		Pokemons []Pokemon `json:"pokemons"`
		Count    int       `json:"count"`
	}
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, err
	}

	return result.Pokemons, nil
}