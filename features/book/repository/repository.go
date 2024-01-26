package repository

import (
	"errors"
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
}

type BookDetail struct {
	gorm.Model
	BookID uint `json:"id_book"`
	RackID uint `json:"id_rack"`
}

type BookQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) book.Repository {
	return &BookQuery{
		db: db,
	}
}

// InsertDetail implements book.Repository.
func (bq *BookQuery) InsertDetail(userID uint, newDetail book.Book, newRack book.Rack) (book.BookDetail, error) {
	var detail book.BookDetail
	detail.BookID = newDetail.ID
	detail.RackID = newRack.ID

	qry := bq.db.Create(&detail)
	if err := qry.Error; err != nil {
		return book.BookDetail{}, err
	}

	return detail, nil
}

// InsertBook implements book.Repository.
func (bq *BookQuery) InsertBook(userID uint, newBook book.Book) (book.Book, error) {
	var inputDB = new(BookModel)
	inputDB.Tittle = newBook.Tittle
	inputDB.Publisher = newBook.Publisher
	inputDB.Author = newBook.Author
	inputDB.Picture = newBook.Picture
	inputDB.Category = newBook.Category

	if err := bq.db.Create(&inputDB).Error; err != nil {
		return book.Book{}, err
	}
	newBook.ID = inputDB.ID

	return newBook, nil
}

// UpdateBook implements book.Repository.
func (bq *BookQuery) UpdateBook(userID uint, bookID uint, input book.Book) (book.Book, error) {
	var proses BookModel
	if err := bq.db.First(&proses, bookID).Error; err != nil {
		return book.Book{}, err
	}

	// Jika tidak ada buku ditemukan
	if proses.ID == 0 {
		err := errors.New("user tidak ditemukan")
		return book.Book{}, err
	}
	if input.Category != "" {
		proses.Category = input.Category
	}
	if input.Tittle != "" {
		proses.Tittle = input.Tittle
	}
	if input.Author != "" {
		proses.Author = input.Author
	}
	if input.Publisher != "" {
		proses.Publisher = input.Publisher
	}
	if input.Category != "" {
		proses.Category = input.Category
	}

	if input.Picture != "" {
		proses.Picture = input.Picture
	}

	if err := bq.db.Save(&proses).Error; err != nil {
		return book.Book{}, err
	}

	result := book.Book{
		ID:        proses.ID,
		Category:  proses.Category,
		Tittle:    proses.Tittle,
		Publisher: proses.Publisher,
		Author:    proses.Author,
		Picture:   proses.Picture,
	}
	return result, nil
}

// DelBook implements book.Repository.
func (bq *BookQuery) DelBook(userID uint, bookID uint) error {
	var prod = new(BookModel)
	if err := bq.db.Where("id", bookID).Find(&prod).Error; err != nil {
		return err
	}

	bq.db.Where("id", bookID).Delete(&prod)
	return nil
}
