package rack

import (
	"library_api/features/book/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Rack struct {
	ID          uint   `json:"rack_id"`
	Name        string `json:"name"`
	RackDetails []repository.BookDetail
}

type Handler interface {
	AddRack() echo.HandlerFunc
}

type Service interface {
	AddRack(token *jwt.Token, newRack Rack) (Rack, error)
}

type Repository interface {
	AddRack(userID uint, newRack Rack) (Rack, error)
}
