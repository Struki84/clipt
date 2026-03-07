package main

import (
	"github.com/struki84/clipt"
	"github.com/struki84/clipt/providers"
	"github.com/struki84/clipt/storage"
	"github.com/struki84/clipt/tui/schema"
	"github.com/struki84/clipt/tui/style"
)

func main() {
	dbPath := "./basic.db"

	sqlite := *storage.NewSQLite(dbPath)

	models := []schema.ChatProvider{
		providers.NewOpenAI("gpt-4o", sqlite),
	}

	clipt.Render(
		models,
		clipt.WithStorage(sqlite),
		clipt.WithDebugLog("debug.log"),
		clipt.WithStyle(style.CatppuccinMocha()),
	)
}
