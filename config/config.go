package config

import (
	"encoding/json"
	"fmt"
	"github.com/struki84/clipt/internal/models"
	"github.com/thanhpk/randstr"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

type AppConfig struct {
	DB    *gorm.DB
	Data  ConfigData
	State State
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
		Data:  configData,
		State: LoadState(),
	}
}

func (config *AppConfig) InitDB() {
	dbPath := "./internal/memory/memory.db"

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Printf("Error connecting to DB: %v", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error getting DB: %v", err)
		return
	}

	sqlDB.Exec("PRAGMA foreign_keys = ON;")
	sqlDB.Exec("PRAGMA journal_mode = WAL;")

	err = config.DB.AutoMigrate(&models.ChatHistory{})
	if err != nil {
		log.Printf("Error migrating DB: %v", err)
		return
	}

	config.DB = db
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

func (config *AppConfig) CurrentSession() string {
	log.Println("Current session:", config.State.CurrentSession)
	config.State = LoadState()
	if config.State.CurrentSession == "" {
		var sessions []models.ChatHistory

		err := config.DB.Find(&sessions).Order("created_at DESC").Error
		if err != nil {
			log.Printf("Error getting sessions: %v", err)
		}

		if len(sessions) > 0 {
			config.State.CurrentSession = sessions[0].SessionID
		} else {
			config.State.CurrentSession = randstr.String(16)
		}

		SetState(config.State)
	}
	return config.State.CurrentSession
}
