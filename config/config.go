package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
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

	configFile, err := os.Open("./my_config.json")
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

func (config *AppConfig) AgentLLM() llms.Model {
	llm, err := openai.New(
		openai.WithModel(config.Data.ActiveLLM),
		openai.WithToken(config.Data.OpenAIAPIToken),
	)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return llm
}

func (config *AppConfig) GetTools() []tools.Tool {

	return []tools.Tool{}
}
