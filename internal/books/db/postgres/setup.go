package books_postgres

import (
	"context"
	"os"

	"github.com/tariq-ventura/go-chatbot/internal/logs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresClient struct {
	client *gorm.DB
	ctx    context.Context
}

var SetupPostgres = func(ctx context.Context) *PostgresClient {
	dsn := os.Getenv("DB_STRING")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logs.LogError("PostgreSQL connection error", map[string]any{"error": err.Error()})
		panic("app stopped")
	}

	logs.LogInfo("PostgreSQL connected successfully", nil)

	return &PostgresClient{
		client: db,
		ctx:    ctx,
	}
}
