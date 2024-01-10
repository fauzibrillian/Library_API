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
	Category  string `json:"category"`
	Stock     uint   `json:"stock"`
}

type Handler interface {
	AddBook() echo.HandlerFunc
}

type Service interface {
	AddBook(token *jwt.Token, newBook Book) (Book, error)
}

type Repository interface {
	InsertBook(userID uint, newBook Book) (Book, error)
}
