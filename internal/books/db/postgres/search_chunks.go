package books_postgres

import (
	"fmt"
	"strings"

	"github.com/pgvector/pgvector-go"
	books_domain "github.com/tariq-ventura/go-chatbot/internal/books/domain"
	"github.com/tariq-ventura/go-chatbot/internal/logs"
)

func (pc *PostgresClient) SearchChunks(vec []float32, limit int) ([]books_domain.ChunkResult, error) {
	var results []books_domain.ChunkResult

	vector := pgvector.NewVector(vec)

	err := pc.client.Raw(`
        SELECT c.content, c.page, c.book_id, b.title
        FROM chunks c
        INNER JOIN books b ON c.book_id = b.id
        ORDER BY c.embedding <-> $1
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

	return results, nil
}

func formatVector(vec []float32) string {
	strVec := make([]string, len(vec))
	for i, v := range vec {
		strVec[i] = fmt.Sprintf("%f", v)
	}
	return fmt.Sprintf("{%s}", strings.Join(strVec, ","))
}
