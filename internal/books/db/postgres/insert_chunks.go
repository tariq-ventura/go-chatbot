package books_postgres

import (
	books_domain "github.com/tariq-ventura/go-chatbot/internal/books/domain"
	"github.com/tariq-ventura/go-chatbot/internal/logs"
)

func (db *PostgresClient) InsertChunks(result books_domain.Chunk) error {
	insert := db.client.Create(&result)

	if insert.Error != nil {
		logs.LogError("PostgreSQL insert error", map[string]any{"error": insert.Error.Error()})
		return insert.Error
	}

	logs.LogInfo("PostgreSQL insert success in Chunks", map[string]any{"insertedID": result.ID})
	return nil
}
