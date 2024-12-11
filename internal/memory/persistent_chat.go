package memory

import (
	"context"
	"errors"
	"log"

	"github.com/struki84/clipt/config"
	"github.com/struki84/clipt/internal/models"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"gorm.io/gorm"
)

var (
	ErrDBConnection     = errors.New("can't connect to database")
	ErrDBMigration      = errors.New("can't migrate database")
	ErrMissingSessionID = errors.New("session id can not be empty")
)

type PersistentChatHistory struct {
	db        *gorm.DB
	records   *models.ChatHistory
	messages  []llms.ChatMessage
	sessionID string
}

var _ schema.ChatMessageHistory = &PersistentChatHistory{}

func NewPersistentChatHistory(config config.AppConfig) *PersistentChatHistory {
	history := &PersistentChatHistory{}
	history.db = config.DB

	return history
}

func (history *PersistentChatHistory) GetSessionID() string {
	return history.sessionID
}

func (history *PersistentChatHistory) SetSession(ID string) {
	log.Println("Setting session ID:", ID)
	history.sessionID = ID
}

func (history *PersistentChatHistory) Messages(context.Context) ([]llms.ChatMessage, error) {
	if history.sessionID == "" {
		return []llms.ChatMessage{}, ErrMissingSessionID
	}

	err := history.db.Where(models.ChatHistory{SessionID: history.sessionID}).Find(&history.records).Error
	if err != nil {
		return nil, err
	}

	history.messages = []llms.ChatMessage{}

	if history.records.ChatHistory != nil {
		for i := range *history.records.ChatHistory {
			msg := (*history.records.ChatHistory)[i]

			if msg.Type == "human" {
				history.messages = append(history.messages, llms.HumanChatMessage{Content: msg.Content})
			}

			if msg.Type == "ai" {
				history.messages = append(history.messages, llms.AIChatMessage{Content: msg.Content})
			}
		}
	}

	return history.messages, nil
}

func (history *PersistentChatHistory) AddMessage(ctx context.Context, message llms.ChatMessage) error {
	if history.sessionID == "" {
		return ErrMissingSessionID
	}

	history.messages = append(history.messages, message)
	bufferString, err := llms.GetBufferString(history.messages, "Human", "AI")
	if err != nil {
		return err
	}

	history.records.SessionID = history.sessionID
	history.records.ChatHistory = history.loadNewMessages()
	history.records.BufferString = bufferString

	err = history.db.Save(&history.records).Error
	if err != nil {
		return err
	}

	return nil
}

func (history *PersistentChatHistory) AddAIMessage(ctx context.Context, message string) error {
	return history.AddMessage(ctx, llms.AIChatMessage{Content: message})
}

func (history *PersistentChatHistory) AddUserMessage(ctx context.Context, message string) error {
	return history.AddMessage(ctx, llms.HumanChatMessage{Content: message})
}

func (history *PersistentChatHistory) SetMessages(ctx context.Context, messages []llms.ChatMessage) error {
	if history.sessionID == "" {
		return ErrMissingSessionID
	}

	history.messages = messages
	bufferString, err := llms.GetBufferString(history.messages, "Human", "AI")
	if err != nil {
		return err
	}

	history.records.SessionID = history.sessionID
	history.records.ChatHistory = history.loadNewMessages()
	history.records.BufferString = bufferString

	err = history.db.Save(&history.records).Error
	if err != nil {
		return err
	}

	return nil
}

func (history *PersistentChatHistory) Clear(context.Context) error {
	history.messages = []llms.ChatMessage{}

	err := history.db.Where(models.ChatHistory{SessionID: history.sessionID}).Delete(&history.records).Error
	if err != nil {
		return err
	}

	return nil
}

func (history *PersistentChatHistory) loadNewMessages() *models.Messages {
	newMsgs := models.Messages{}
	for _, msg := range history.messages {
		newMsgs = append(newMsgs, models.Message{
			Type:    string(msg.GetType()),
			Content: msg.GetContent(),
		})
	}

	return &newMsgs
}
