package storage

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/struki84/clipt/tui"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/thanhpk/randstr"
)

type Session struct {
	gorm.Model
	ID    string
	Title string
	Msgs  *Messages
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
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, m)
	}
	return errors.New("could not scan type into Message")
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

func (sql *SQLite) NewSession() (tui.ChatSession, error) {
	sessionID := randstr.String(8)

	sql.record = Session{
		ID:    sessionID,
		Title: "New Session",
		Msgs:  &Messages{},
	}

	err := sql.db.Create(&sql.record).Error

	if err != nil {
		return tui.ChatSession{}, errors.New(fmt.Sprintf("Error creating session: %v", err))
	}

	return tui.ChatSession{
		ID:        sessionID,
		Title:     "New Session",
		Msgs:      []tui.ChatMsg{},
		CreatedAt: sql.record.CreatedAt.Unix(),
	}, nil
}

func (sql *SQLite) ListSessions() []tui.ChatSession {

	return []tui.ChatSession{}
}

func (sql *SQLite) LoadSession(sessionID string) (tui.ChatSession, error) {
	err := sql.db.Where(Session{ID: sessionID}).Find(&sql.record).Error
	if err != nil {
		return tui.ChatSession{}, errors.New(fmt.Sprintf("Error loading session: %v", err))
	}

	msgs := []tui.ChatMsg{}
	for _, msg := range *sql.record.Msgs {
		msgs = append(msgs, tui.ChatMsg{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	return tui.ChatSession{
		ID:        sql.record.ID,
		Title:     sql.record.Title,
		Msgs:      msgs,
		CreatedAt: sql.record.CreatedAt.Unix(),
	}, nil
}

func (sql *SQLite) SaveSession(session tui.ChatSession) (tui.ChatSession, error) {
	msgs := Messages{}
	for i, msg := range session.Msgs {
		msgs[i] = Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	sql.record = Session{
		ID:    session.ID,
		Title: session.Title,
		Msgs:  &msgs,
	}

	err := sql.db.Save(&sql.record).Error
	if err != nil {
		return tui.ChatSession{}, errors.New(fmt.Sprintf("Error saving session: %v", err))
	}

	return session, nil
}

func (sql *SQLite) DeleteSession(string) error {
	return nil
}
