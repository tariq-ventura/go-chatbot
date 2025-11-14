package books_domain

import "gorm.io/gorm"

type Chunk struct {
	gorm.Model
	BookID  string `gorm:"index"`
	Page    int
	Content string `gorm:"type:text"`
}
