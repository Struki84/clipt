package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type ChatHistory struct {
	gorm.Model
	SessionID    string    `gorm:"type:varchar(256)" json:"sessionID,omitempty"`
	BufferString string    `gorm:"type:text" json:"bufferString,omitempty"`
	ChatHistory  *Messages `json:"chatHistory" gorm:"type:jsonb;column:chatHistory"`
}

type Messages []Message

type Message struct {
	Type    string `json:"type"`
	Content string `json:"text"`
}

// Value implements the driver.Valuer interface, this method allows us to
// customize how we store the Message type in the database.
func (m Messages) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan implements the sql.Scanner interface, this method allows us to
// define how we convert the Message data from the database into our Message type.
func (m *Messages) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, m)
	}
	return errors.New("could not scan type into Message")
}
