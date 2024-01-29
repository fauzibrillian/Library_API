package book

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Book struct {
	ID        uint   `json:"id"`
	Tittle    string `json:"tittle"`
	Publisher string `json:"publisher"`
	Author    string `json:"author"`
	Picture   string `json:"picture"`
}

type Handler interface {
	AddBook() echo.HandlerFunc
	UpdateBook() echo.HandlerFunc
	DeleteBook() echo.HandlerFunc
	SearchBook() echo.HandlerFunc
}

type Service interface {
	AddBook(token *jwt.Token, newBook Book) (Book, error)
	UpdateBook(token *jwt.Token, bookID uint, input Book) (Book, error)
	DelBook(token *jwt.Token, bookID uint) error
	SearchBook(tittle string, page uint, limit uint) ([]Book, uint, error)
}

type Repository interface {
	InsertBook(userID uint, newBook Book) (Book, error)
	UpdateBook(userID uint, bookID uint, input Book) (Book, error)
	DelBook(userID uint, bookID uint) error
	SearchBook(tittle string, page uint, limit uint) ([]Book, uint, error)
}
