package repository

import (
	"library_api/features/transaction"
	"time"

	"gorm.io/gorm"
)

type TransactionModel struct {
	gorm.Model
	BookID     uint
	UserID     uint
	DateReturn time.Time
}

type TransactionQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) transaction.Repository {
	return &TransactionQuery{
		db: db,
	}
}

// Borrow implements transaction.Repository.
func (tq *TransactionQuery) Borrow(userID uint, BookID uint) (transaction.Transaction, error) {
	var inputDB TransactionModel
	inputDB.BookID = BookID
	inputDB.UserID = userID

	if err := tq.db.Create(&inputDB).Error; err != nil {
		return transaction.Transaction{}, err
	}

	var result transaction.Transaction
	result.ID = inputDB.ID
	result.BookID = inputDB.BookID
	result.DateBorrow = inputDB.CreatedAt

	return result, nil
}
