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
func (bq *BookQuery) InsertBook(userID uint, newBook book.Book) (book.Book, error) {
	var inputDB = new(BookModel)
	inputDB.Tittle = newBook.Tittle
	inputDB.Publisher = newBook.Publisher
	inputDB.Author = newBook.Author
	inputDB.Picture = newBook.Picture
	inputDB.Category = newBook.Category
	inputDB.Stock = int(newBook.Stock)

	if err := bq.db.Create(&inputDB).Error; err != nil {
		return book.Book{}, err
	}
	newBook.ID = inputDB.ID

	return newBook, nil
}
