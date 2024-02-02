package repository

import (
	"errors"
	"fmt"
	p "library_api/features/book"
	"library_api/features/transaction"
	"library_api/features/user/repository"
	ur "library_api/features/user/repository"
	"log"
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

// UpdateReturn implements transaction.Repository.
func (tq *TransactionQuery) UpdateReturn(userID uint, transactionID uint, input transaction.Transaction) ([]transaction.Transaction, error) {
	var update TransactionModel
	var User []ur.UserModel
	var result []transaction.Transaction
	var Book []p.Book

	if err := tq.db.First(&update, transactionID).Error; err != nil {
		return []transaction.Transaction{}, err
	}

	if update.ID == 0 {
		err := errors.New("transaction tidak ditemukan")
		return []transaction.Transaction{}, err
	}
	if input.DateReturn != (time.Time{}) {
		update.DateReturn = input.DateReturn
	}

	if err := tq.db.Save(&update).Error; err != nil {
		return []transaction.Transaction{}, err
	}

	tmpUser := new(repository.UserModel)
	if err := tq.db.Table("user_models").Where("id = ?", update.UserID).Find(&tmpUser).Error; err != nil {
		return []transaction.Transaction{}, err
	}
	User = append(User, *tmpUser)

	tmpBook := new(p.Book)
	if err := tq.db.Table("book_models").Where("id = ?", update.BookID).Find(&tmpBook).Error; err != nil {
		return []transaction.Transaction{}, err
	}
	Book = append(Book, *tmpBook)

	// Create a transaction object with user and book details
	resultTransaction := transaction.Transaction{
		ID:         update.ID,
		BookID:     update.BookID,
		DateBorrow: update.CreatedAt,
		DateReturn: update.DateReturn,
		Users:      User,
		Books:      Book,
	}

	result = append(result, resultTransaction)
	return result, nil
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

// AllTransaction implements transaction.Repository.
func (tq *TransactionQuery) AllTransaction(userID uint, name string, page uint, limit uint) ([]transaction.Transaction, int, error) {
	var User []ur.UserModel
	var tm []TransactionModel
	var result []transaction.Transaction

	if name != "" {
		if err := tq.db.Table("user_models").Where("name LIKE ?", "%"+name+"%").Find(&User).Error; err != nil {
			return nil, 0, err
		}

		if len(User) > 0 {
			userIDs := make([]uint, len(User))
			for i, u := range User {
				userIDs[i] = u.ID
			}

			offset := (page - 1) * limit
			query := tq.db.Offset(int(offset)).Limit(int(limit)).Where("user_id IN (?)", userIDs)
			if err := query.Find(&tm).Error; err != nil {
				return nil, 0, err
			}
		}
	} else {
		offset := (page - 1) * limit
		if err := tq.db.Offset(int(offset)).Limit(int(limit)).Find(&tm).Error; err != nil {
			return nil, 0, err
		}
	}

	var Book []p.Book
	for _, resp := range tm {
		// Clear slices before appending data
		User = nil
		Book = nil

		// Fetch user details for the current transaction
		tmp := new(repository.UserModel)
		if err := tq.db.Table("user_models").Where("id = ?", resp.UserID).Find(&tmp).Error; err != nil {
			return nil, 0, err
		}
		User = append(User, *tmp)

		// Fetch book details for the current transaction
		tmpBook := new(p.Book)
		if err := tq.db.Table("book_models").Where("id = ?", resp.BookID).Find(&tmpBook).Error; err != nil {
			return nil, 0, err
		}
		Book = append(Book, *tmpBook)

		// Create a transaction object with user and book details
		results := transaction.Transaction{
			ID:         resp.ID,
			BookID:     resp.BookID,
			DateBorrow: resp.CreatedAt,
			DateReturn: resp.DateReturn,
			Users:      User,
			Books:      Book,
		}

		// Append the transaction to the result slice
		result = append(result, results)
	}

	var totalPage int
	tableNameUser := "transaction_models"
	columnNameUser := "deleted_at"
	queryuser := fmt.Sprintf("SELECT COUNT(*) AS null_count FROM %s WHERE %s IS NULL", tableNameUser, columnNameUser)
	err := tq.db.Raw(queryuser).Scan(&totalPage).Error
	if err != nil {
		log.Fatal(err)
	}

	if totalPage%int(limit) == 0 {
		totalPage = totalPage / int(limit)
	} else {
		totalPage = totalPage / int(limit)
		totalPage++
	}

	return result, totalPage, err
}
