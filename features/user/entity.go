package user

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Handler interface {
	Login() echo.HandlerFunc
	Register() echo.HandlerFunc
}

type Service interface {
	Login(phone string, password string) (User, error)
	Register(newUser User) (User, error)
}

type Repository interface {
	Login(phone string) (User, error)
	Register(newUser User) (User, error)
}
