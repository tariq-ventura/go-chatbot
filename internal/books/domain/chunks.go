package books_domain

import (
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type Chunk struct {
	gorm.Model
	BookID    string `gorm:"index"`
	Page      int
	Content   string          `gorm:"type:text"`
	Embedding pgvector.Vector `gorm:"type:vector(1024)"`
}

type ChunkResult struct {
	Content string
	Page    int
	BookID  string
	Title   string
}
