package storage

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/struki84/clipt/tui"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/thanhpk/randstr"
)

type Session struct {
	gorm.Model
	SessionID string
	Title     string
	Msgs      Messages `gorm:"type:jsonb;column:msgs"`
}

type Messages []Message

type Message struct {
	Role    string
	Content string
}

func (m Messages) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *Messages) Scan(src interface{}) error {
	var bytes []byte
	switch v := src.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("could not scan type into Messages")
	}
	return json.Unmarshal(bytes, m)
}

type SQLite struct {
	db     *gorm.DB
	path   string
	record Session
}

func NewSQLite(dbPath string) *SQLite {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Printf("Error connecting to DB: %v", err)
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error getting DB: %v", err)
		return nil
	}

	sqlDB.Exec("PRAGMA foreign_keys = ON;")
	sqlDB.Exec("PRAGMA journal_mode = WAL;")

	err = db.AutoMigrate(Session{})
	if err != nil {
		log.Printf("Error migrating DB: %v", err)
		return nil
	}

	return &SQLite{
		db:   db,
		path: dbPath,
	}
}

func (sql SQLite) NewSession() (tui.ChatSession, error) {
	sessionID := randstr.String(8)

	sql.record = Session{
		SessionID: sessionID,
		Title:     "New Session",
		Msgs:      Messages{},
	}

	err := sql.db.Create(&sql.record).Error

	if err != nil {
		return tui.ChatSession{}, fmt.Errorf("Error creating new session, %v", err)
	}

	return tui.ChatSession{
		ID:        sessionID,
		Title:     "New Session",
		Msgs:      []tui.ChatMsg{},
		CreatedAt: sql.record.CreatedAt.Unix(),
	}, nil
}

func (sql SQLite) ListSessions() []tui.ChatSession {
	sessions := []Session{}
	err := sql.db.Find(&sessions).Error
	if err != nil {
		return []tui.ChatSession{}
	}

	list := []tui.ChatSession{}
	for _, session := range sessions {
		msgs := []tui.ChatMsg{}
		for _, msg := range session.Msgs {
			msgs = append(msgs, tui.ChatMsg{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}

		list = append(list, tui.ChatSession{
			ID:        session.SessionID,
			Title:     session.Title,
			Msgs:      msgs,
			CreatedAt: session.CreatedAt.Unix(),
		})
	}

	return list
}

func (sql SQLite) LoadRecentSession() (tui.ChatSession, error) {
	sessions := []Session{}
	err := sql.db.Find(&sessions).Order("created_at, DESC").Error
	if err != nil {
		return tui.ChatSession{}, fmt.Errorf("Error loading recent sessions, %v", err)
	}

	if len(sessions) > 0 {
		session := sessions[0]

		msgs := []tui.ChatMsg{}
		for _, msg := range session.Msgs {
			msgs = append(msgs, tui.ChatMsg{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}

		sql.record = session

		return tui.ChatSession{
			ID:        sql.record.SessionID,
			Title:     sql.record.Title,
			Msgs:      msgs,
			CreatedAt: sql.record.CreatedAt.Unix(),
		}, nil

	}

	sql.record = Session{
		SessionID: randstr.String(8),
		Title:     "New Session",
		Msgs:      Messages{},
	}

	err = sql.db.Save(&sql.record).Error
	if err != nil {
		return tui.ChatSession{}, fmt.Errorf("Error loading recent sessions, %v", err)
	}

	return tui.ChatSession{
		ID:        sql.record.SessionID,
		Title:     sql.record.Title,
		Msgs:      []tui.ChatMsg{},
		CreatedAt: sql.record.CreatedAt.Unix(),
	}, nil
}

func (sql SQLite) LoadSession(sessionID string) (tui.ChatSession, error) {
	err := sql.db.Where("session_id = ?", sessionID).Find(&sql.record).Error
	if err != nil {
		return tui.ChatSession{}, fmt.Errorf("Error loading session: %v", err)
	}

	msgs := []tui.ChatMsg{}
	for _, msg := range sql.record.Msgs {
		msgs = append(msgs, tui.ChatMsg{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	return tui.ChatSession{
		ID:        sql.record.SessionID,
		Title:     sql.record.Title,
		Msgs:      msgs,
		CreatedAt: sql.record.CreatedAt.Unix(),
	}, nil
}

func (sql SQLite) SaveSession(session tui.ChatSession) (tui.ChatSession, error) {
	msgs := Messages{}
	for i, msg := range session.Msgs {
		msgs[i] = Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	sql.record = Session{
		SessionID: session.ID,
		Title:     session.Title,
		Msgs:      msgs,
	}

	err := sql.db.Save(&sql.record).Error
	if err != nil {
		return tui.ChatSession{}, fmt.Errorf("Can't save session, %v", err)
	}

	return session, nil
}

func (sql SQLite) DeleteSession(sessionID string) error {
	err := sql.db.Where("session_id = ?", sessionID).Delete(&sql.record)
	if err != nil {
		return fmt.Errorf("Error deleting session: %v ", err)
	}
	return nil
}

func (sql SQLite) SaveSessionMsg(sessionID string, humanMsg string, aiMsg string) error {
	err := sql.db.Where("session_id = ?", sessionID).Find(&sql.record).Error
	if err != nil {
		return fmt.Errorf("Error loading session: %v", err)
	}

	human := Message{
		Role:    "User",
		Content: humanMsg,
	}

	ai := Message{
		Role:    "AI",
		Content: aiMsg,
	}

	sql.record.Msgs = append(sql.record.Msgs, human)
	sql.record.Msgs = append(sql.record.Msgs, ai)

	err = sql.db.Save(&sql.record).Error
	if err != nil {
		return fmt.Errorf("Can't save session, %v", err)
	}

	return nil
}

func (sql SQLite) LoadSessionMsgs(sessionID string) (string, error) {
	result := []string{}
	err := sql.db.Where("session_id = ?", sessionID).Find(&sql.record).Error
	if err != nil {
		return "", fmt.Errorf("Error loading session: %v", err)
	}

	if sql.record.Msgs != nil {
		for _, msg := range sql.record.Msgs {
			result = append(result, fmt.Sprintf("%s: %s", msg.Role, msg.Content))
		}
	}

	return strings.Join(result, "\n"), nil
}
