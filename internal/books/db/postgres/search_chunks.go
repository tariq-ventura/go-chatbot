package books_postgres

import (
	"fmt"
	"strings"

	"github.com/pgvector/pgvector-go"
	"github.com/tariq-ventura/go-chatbot/internal/logs"
)

func (pc *PostgresClient) SearchChunks(vec []float32, limit int) ([]string, error) {
	var results []struct {
		Content string `gorm:"column:content"`
	}
	vector := pgvector.NewVector(vec)
	err := pc.client.Raw(`
        SELECT content
        FROM chunks
        ORDER BY embedding <-> $1
        LIMIT $2
    `, vector, limit).Scan(&results).Error

	if err != nil {
		logs.LogError("Error searching chunks in Postgres", map[string]any{"error": err.Error()})
		return nil, err
	}

	texts := make([]string, len(results))
	for i, result := range results {
		texts[i] = result.Content
	}

	logs.LogInfo("Successfully searched chunks in Postgres", map[string]any{"found": len(texts)})

	return texts, nil
}

func formatVector(vec []float32) string {
	strVec := make([]string, len(vec))
	for i, v := range vec {
		strVec[i] = fmt.Sprintf("%f", v)
	}
	return fmt.Sprintf("{%s}", strings.Join(strVec, ","))
}
