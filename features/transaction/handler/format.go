package handler

import "time"

type TransactionRequest struct {
	BookID int `json:"book_id" form:"book_id"`
}

type TransactionResponse struct {
	ID         int       `json:"transaction_id"`
	BookID     int       `json:"book_id"`
	DateBorrow time.Time `json:"date_borrow"`
}

type SearchTransactionResponse struct {
	ID          int       `json:"transaction_id"`
	UserPicture string    `json:"user_picture" form:"user_picture"`
	UserName    string    `json:"user_name" form:"user_name"`
	TittleBook  string    `json:"tittle_books" form:"tittle_books"`
	PictureBook string    `json:"picture_books" form:"picture_books"`
	DateBorrow  time.Time `json:"date_borrow"`
	DateReturn  time.Time `json:"date_return"`
}

type DateReturnRequest struct {
	ID         int       `json:"transaction_id" form:"transaction_id"`
	DateReturn time.Time `json:"date_return" form:"date_return"`
}

type DateReturnResponse struct {
	ID         int       `json:"transaction_id"`
	UserName   string    `json:"user_name" form:"user_name"`
	TittleBook string    `json:"tittle_books" form:"tittle_books"`
	DateBorrow time.Time `json:"date_borrow"`
	DateReturn time.Time `json:"date_return"`
}
