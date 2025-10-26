package main

import (
	"github.com/struki84/clipt"
	"github.com/struki84/clipt/providers"
	"github.com/struki84/clipt/storage"
	"github.com/struki84/clipt/tui"
)

func main() {
	p := []tui.ChatProvider{
		providers.NewOpenAI("gpt-4o"),
	}

	dbPath := "./basic.db"

	clipt.Render(p, storage.NewSQLite(dbPath))
}
