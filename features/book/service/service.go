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
		return book.Book{}, errors.New("token Error")
	}
	if rolesUser != "admin" {
		return book.Book{}, errors.New("unauthorized access: admin role required")
	}

	result, err := bs.repo.InsertBook(userId, newBook)
	if err != nil {
		return book.Book{}, errors.New("inputan tidak boleh kosong")
	}

	return result, err
}

// UpdateBook implements book.Service.
func (bs *BookServices) UpdateBook(token *golangjwt.Token, bookID uint, input book.Book) (book.Book, error) {
	userId, rolesUser, err := jwt.ExtractToken(token)
	if err != nil {
		return book.Book{}, errors.New("token error")
	}
	if rolesUser == "" {
		return book.Book{}, errors.New("role cannot empty")
	}
	if rolesUser != "admin" {
		return book.Book{}, errors.New("unauthorized access: admin role required")
	}

	result, err := bs.repo.UpdateBook(userId, bookID, input)
	if err != nil {
		return book.Book{}, errors.New("failed to update the product")
	}

	return result, nil
}

// DelBook implements book.Service.
func (bs *BookServices) DelBook(token *golangjwt.Token, bookID uint) error {
	userId, rolesUser, err := jwt.ExtractToken(token)
	if err != nil {
		return errors.New("token error")
	}
	if rolesUser != "admin" {
		return errors.New("unauthorized access: admin role required")
	}
	err = bs.repo.DelBook(userId, bookID)
	return err
}

// SearchBook implements book.Service.
func (bs *BookServices) SearchBook(tittle string, page uint, limit uint) ([]book.Book, uint, error) {
	books, totalPage, err := bs.repo.SearchBook(tittle, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return books, totalPage, err
}
