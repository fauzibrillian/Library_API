package transaction

import (
	"library_api/features/book"
	"library_api/features/user/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Transaction struct {
	ID         uint
	BookID     uint
	DateBorrow time.Time
	DateReturn time.Time
	Users      []repository.UserModel
	Books      []book.Book
}

type Handler interface {
	Borrow() echo.HandlerFunc
}

type Repository interface {
	Borrow(userID uint, BookID uint) (Transaction, error)
}

type Service interface {
	Borrow(token *jwt.Token, BookID uint) (Transaction, error)
}
