package books_postgres

import (
	"strings"
	"unicode/utf8"

	books_domain "github.com/tariq-ventura/go-chatbot/internal/books/domain"
	"github.com/tariq-ventura/go-chatbot/internal/logs"
)

func (db *PostgresClient) InsertChunks(result books_domain.Chunk) error {

	result.Content = sanitizeUTF8(result.Content)

	result.Content = strings.Map(func(r rune) rune {
		if r < 32 && r != '\n' && r != '\r' && r != '\t' {
			return -1
		}
		return r
	}, result.Content)

	insert := db.client.Create(&result)

	if insert.Error != nil {
		logs.LogError("PostgreSQL insert error", map[string]any{"error": insert.Error.Error()})
		return insert.Error
	}

	logs.LogInfo("PostgreSQL insert success in Chunks", map[string]any{"insertedID": result.ID})
	return nil
}

func sanitizeUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}

	// Reemplaza caracteres inválidos con espacio
	v := make([]rune, 0, len(s))
	for _, r := range s {
		if r == utf8.RuneError {
			v = append(v, ' ')
		} else {
			v = append(v, r)
		}
	}
	return string(v)
}
