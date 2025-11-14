package main

import (
	"context"

	books_db "github.com/tariq-ventura/go-chatbot/internal/books/db"
	"github.com/tariq-ventura/go-chatbot/internal/logs"
	"github.com/tariq-ventura/go-chatbot/internal/router"
)

func RunApp(ctx context.Context) error {
	database, err := books_db.NewDatabase(ctx)
	if err != nil {
		logs.LogError("Database Books connection error", map[string]any{"error": err.Error()})
		return err
	}

	err = database.BooksMigration(ctx)
	if err != nil {
		logs.LogError("Database Books migration error", map[string]any{"error": err.Error()})
		return err
	}

	logs.LogInfo("Database connected and migrated successfully", nil)

	err = database.ChunksMigration(ctx)
	if err != nil {
		logs.LogError("Database Chunks migration error", map[string]any{"error": err.Error()})
		return err
	}

	logs.LogInfo("Database connected and migrated successfully", nil)

	r := &router.Routes{}
	r.Routes = r.SetupRouter()
	r.Run()

	return nil
}

func main() {
	ctx := context.Background()
	if err := RunApp(ctx); err != nil {
		panic("app stopped")
	}
}
