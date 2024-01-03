package user

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
	Avatar      string `json:"avatar"`
	Role        string `json:"role"`
}

type Handler interface {
	Login() echo.HandlerFunc
	Register() echo.HandlerFunc
	ResetPassword() echo.HandlerFunc
	UpdateUser() echo.HandlerFunc
}

type Service interface {
	Login(phone string, password string) (User, error)
	Register(newUser User) (User, error)
	ResetPassword(token *jwt.Token, input User) (User, error)
	UpdateUser(token *jwt.Token, input User) (User, error)
}

type Repository interface {
	Login(phone string) (User, error)
	Register(newUser User) (User, error)
	ResetPassword(input User) (User, error)
	UpdateUser(input User) (User, error)
	GetUserByID(userID uint) (*User, error)
}
