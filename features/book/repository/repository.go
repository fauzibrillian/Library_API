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

	if input.Tittle != "" {
		proses.Tittle = input.Tittle
	}
	if input.Author != "" {
		proses.Author = input.Author
	}
	if input.Publisher != "" {
		proses.Publisher = input.Publisher
	}

	if input.Picture != "" {
		proses.Picture = input.Picture
	}

	if err := bq.db.Save(&proses).Error; err != nil {
		return book.Book{}, err
	}

	result := book.Book{
		ID:        proses.ID,
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

// SearchBook implements book.Repository.
func (bq *BookQuery) SearchBook(tittle string, page uint, limit uint) ([]book.Book, uint, error) {
	var books []BookModel
	qry := bq.db.Table("book_models")

	if tittle != "" {
		qry = qry.Where("tittle like ?", "%"+tittle+"%")
		qry = qry.Where("deleted_at IS NULL")
	}

	var totalBooks int64
	if err := qry.Count(&totalBooks).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	qry = qry.Offset(int(offset)).Limit(int(limit))

	if err := qry.Find(&books).Error; err != nil {
		return nil, 0, err
	}

	var result []book.Book
	for _, s := range books {
		result = append(result, book.Book{
			ID:        s.ID,
			Tittle:    s.Tittle,
			Publisher: s.Publisher,
			Author:    s.Author,
			Picture:   s.Picture,
		})
	}

	totalPages := int(totalBooks) / int(limit)
	if totalBooks%int64(limit) != 0 {
		totalPages++
	}

	return result, uint(totalPages), nil
}
