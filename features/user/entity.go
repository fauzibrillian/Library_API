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
	Email       string `json:"email"`
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
	Delete() echo.HandlerFunc
	SearchUser() echo.HandlerFunc
}

type Service interface {
	Login(email string, password string) (User, error)
	Register(newUser User) (User, error)
	ResetPassword(token *jwt.Token, input User) (User, error)
	UpdateUser(token *jwt.Token, input User) (User, error)
	DeleteUser(token *jwt.Token, userID uint) error
	SearchUser(token *jwt.Token, name string, page uint, limit uint) ([]User, uint, error)
}

type Repository interface {
	Login(email string) (User, error)
	Register(newUser User) (User, error)
	ResetPassword(input User) (User, error)
	UpdateUser(input User) (User, error)
	GetUserByID(userID uint) (*User, error)
	DeleteUser(userID uint) error
	SearchUser(userID uint, name string, page uint, limit uint) ([]User, uint, error)
}
