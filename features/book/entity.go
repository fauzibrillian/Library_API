package book

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Book struct {
	ID          uint         `json:"id"`
	Tittle      string       `json:"tittle"`
	Publisher   string       `json:"publisher"`
	Author      string       `json:"author"`
	Picture     string       `json:"picture"`
	Category    string       `json:"category"`
	Stock       uint         `json:"stock"`
	BookDetails []BookDetail `json:"book_detail" gorm:"many2many:book_detailbook;"`
}

type BookDetail struct {
	ID     uint `json:"id"`
	BookID uint `json:"id_book"`
	RackID uint `json:"id_rack"`
}

type Rack struct {
	ID   uint   `json:"id_rack"`
	Name string `json:"name"`
}

type Handler interface {
	AddBook() echo.HandlerFunc
	AddDetail() echo.HandlerFunc
	UpdateBook() echo.HandlerFunc
	DeleteBook() echo.HandlerFunc
}

type Service interface {
	AddBook(token *jwt.Token, newBook Book) (Book, error)
	AddDetail(token *jwt.Token, newDetail Book, newRack Rack) (BookDetail, error)
	UpdateBook(token *jwt.Token, bookID uint, input Book) (Book, error)
	DelBook(token *jwt.Token, bookID uint) error
}

type Repository interface {
	InsertBook(userID uint, newBook Book) (Book, error)
	InsertDetail(userID uint, newDetail Book, newRack Rack) (BookDetail, error)
	UpdateBook(userID uint, bookID uint, input Book) (Book, error)
	DelBook(userID uint, bookID uint) error
}
