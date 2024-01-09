package repository

import (
	"library_api/features/book"

	"gorm.io/gorm"
)

type BookModel struct {
	gorm.Model
	Tittle    string
	Publisher string
	Author    string
	Picture   string
	Category  string
	Stock     int
}

type BookQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) book.Repository {
	return &BookQuery{
		db: db,
	}
}

// InsertBook implements book.Repository.
func (*BookQuery) InsertBook(userID uint, newBook book.Book) (book.Book, error) {
	panic("unimplemented")
}
