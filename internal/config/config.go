package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/google/uuid"
)

type Config struct {
	DbUrl           string    `json:"db_url"`
	CurrentUserName string    `json:"current_user_name"`
	CurrentUserId   uuid.UUID `json:"current_user_id"`
	filePath        string
}

func Read(path string) (*Config, error) {
	if len(path) == 0 {
		path = "~/.gatorconfig.json"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	c.filePath = path
	return &c, nil
}

func (c *Config) SetUser(user string, id uuid.UUID) error {
	if len(user) == 0 {
		return errors.New("Empty User")
	}
	c.CurrentUserName = user
	c.CurrentUserId = id
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}

func write(c Config) error {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile(c.filePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
