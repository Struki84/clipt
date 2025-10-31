package main

import (
	"github.com/struki84/clipt"
	"github.com/struki84/clipt/providers"
	"github.com/struki84/clipt/storage"
	"github.com/struki84/clipt/tui"
)

func main() {
	dbPath := "./basic.db"
	s := storage.NewSQLite(dbPath)
	p := []tui.ChatProvider{
		providers.NewOpenAI("gpt-4o", *s),
	}

	clipt.Render(p, storage.NewSQLite(dbPath))
}
