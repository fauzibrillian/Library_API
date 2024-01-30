package service

import (
	"errors"
	"library_api/features/transaction"
	"library_api/helper/jwt"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

type TransactionServices struct {
	repo transaction.Repository
}

func New(r transaction.Repository) transaction.Service {
	return &TransactionServices{
		repo: r,
	}
}

// Borrow implements transaction.Service.
func (ts *TransactionServices) Borrow(token *golangjwt.Token, BookID uint) (transaction.Transaction, error) {
	userID, rolesUser, err := jwt.ExtractToken(token)
	if err != nil {
		return transaction.Transaction{}, errors.New("user does not exist")
	}
	if rolesUser == "" {
		return transaction.Transaction{}, errors.New("roles user can't empty")
	}
	result, err := ts.repo.Borrow(userID, BookID)
	if err != nil {
		return transaction.Transaction{}, errors.New("Repository Error")
	}
	return result, err
}
