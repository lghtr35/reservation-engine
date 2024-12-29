package models

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog"
)

type Configuration struct {
	DbConnectionString string `json:"dbConnectionString"`
	Secret             string `json:"secret"`
	salt               string
}

func (c *Configuration) ReadAndFillSelf(logger zerolog.Logger) error {
	file, err := os.ReadFile("./config.json")
	if err != nil {
		logger.Error().Err(err).Msg("Error reading local JSON config file")
		return err
	}

	// Unmarshal the JSON data
	var config Configuration
	err = json.Unmarshal(file, &c)
	if err != nil {
		logger.Error().Err(err).Msg("Error unmarshalling config JSON")
		return err
	}

	c.DbConnectionString = config.DbConnectionString
	c.salt = "salty-crackers"
	return nil
}

func (c *Configuration) GetSalt() string {
	return c.salt
}
