package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const defaultConfigPath = "./config.json"

type ConnectorsConfig []ConnectorConfig

type ConnectorConfig struct {
	ClientId   string `json:"clientId"`
	ApiUrl     string `json:"apiUrl"`
	ApiToken   string `json:"apiToken"`
	BufferSize uint64 `json:"bufferSize"`
}

func Load() (ConnectorsConfig, error) {
	file, err := os.Open(defaultConfigPath)
	if err != nil {
		return []ConnectorConfig{}, fmt.Errorf("error when reading config file: %w", err)
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return []ConnectorConfig{}, fmt.Errorf("error when reading config file: %w", err)
	}
	var config ConnectorsConfig
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return []ConnectorConfig{}, fmt.Errorf("error when parsing config file: %w", err)
	}
	return config, nil
}
