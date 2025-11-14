package books_db

import (
	"context"
	"errors"
	"fmt"
	"os"

	books_postgres "github.com/tariq-ventura/go-chatbot/internal/books/db/postgres"
	books_domain "github.com/tariq-ventura/go-chatbot/internal/books/domain"
)

type Database interface {
	InsertBook(result books_domain.Book) error
	InsertChunks(result books_domain.Chunk) error
	BooksMigration(ctx context.Context) error
	ChunksMigration(ctx context.Context) error
	SearchChunks(vec []float32, limit int) ([]string, error)
}

var NewDatabase = func(ctx context.Context) (Database, error) {
	dbType := os.Getenv("DB_CONTEXT")
	switch dbType {
	case "postgres":
		fmt.Print("Using postgresql")
		return books_postgres.SetupPostgres(ctx), nil
	default:
		return nil, errors.New("unsupported database backend")
	}
}
