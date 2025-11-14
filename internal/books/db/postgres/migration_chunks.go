package books_postgres

import (
	"context"

	books_domain "github.com/tariq-ventura/go-chatbot/internal/books/domain"
	"github.com/tariq-ventura/go-chatbot/internal/logs"
)

func (db *PostgresClient) ChunksMigration(ctx context.Context) error {
	err := db.client.AutoMigrate(books_domain.Chunk{})

	if err != nil {
		logs.LogError("Error during Chunks migration in Postgres", map[string]interface{}{"error": err.Error()})
		return err
	}

	logs.LogInfo("Successfully completed Chunks migration in Postgres", nil)
	return nil
}
