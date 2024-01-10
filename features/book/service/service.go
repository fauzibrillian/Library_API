package service

import (
	"errors"
	"library_api/features/book"
	"library_api/helper/jwt"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

type BookServices struct {
	repo book.Repository
}

func New(r book.Repository) book.Service {
	return &BookServices{
		repo: r,
	}
}

// AddBook implements book.Service.
func (bs *BookServices) AddBook(token *golangjwt.Token, newBook book.Book) (book.Book, error) {
	userId, rolesUser, err := jwt.ExtractToken(token)
	if err != nil {
		return book.Book{}, errors.New("Token Error")
	}
	if rolesUser != "admin" {
		return book.Book{}, errors.New("unauthorized access: admin role required")
	}

	result, err := bs.repo.InsertBook(userId, newBook)
	if err != nil {
		return book.Book{}, errors.New("Inputan tidak boleh kosong")
	}

	return result, err
}
