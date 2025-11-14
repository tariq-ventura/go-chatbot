package books_domain

type Book struct {
	ID     string `gorm:"primaryKey"`
	Title  string `gorm:"index"`
	Author string
	Chunks []Chunk `gorm:"foreignKey:BookID"`
}
