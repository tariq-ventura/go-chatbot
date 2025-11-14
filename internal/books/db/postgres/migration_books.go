package books_postgres

import (
	"context"

	books_domain "github.com/tariq-ventura/go-chatbot/internal/books/domain"
	"github.com/tariq-ventura/go-chatbot/internal/logs"
)

func (db *PostgresClient) BooksMigration(ctx context.Context) error {
	err := db.client.AutoMigrate(books_domain.Book{})

	if err != nil {
		logs.LogError("Error during Books migration in Postgres", map[string]interface{}{"error": err.Error()})
		return err
	}

	logs.LogInfo("Successfully completed Books migration in Postgres", nil)
	return nil
}
