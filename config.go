package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type config interface {
	AgentLLM() llms.Model
}

var _ config = &AppConfig{}

type AppConfig struct {
	Data ConfigData
}

type ConfigData struct {
	ActiveLLM      string `json:"ActiveLLM,omitempty"`
	OpenAIAPIToken string `json:"OpenAIAPIToken,omitempty"`
}

func NewConfig() AppConfig {
	var configData ConfigData

	configFile, err := os.Open("my_config.json")
	if err != nil {
		fmt.Println(err)
	}

	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&configData); err != nil {
		fmt.Println(err)
	}

	return AppConfig{
		Data: configData,
	}
}

func (c *AppConfig) AgentLLM() llms.Model {
	llm, err := openai.New(
		openai.WithModel(c.Data.ActiveLLM),
		openai.WithToken(c.Data.OpenAIAPIToken),
	)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return llm
}
