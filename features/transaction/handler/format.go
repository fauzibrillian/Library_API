package handler

import "time"

type TransactionRequest struct {
	BookID int `json:"book_id" form:"book_id"`
}

type TransactionResponse struct {
	ID         int       `json:"transaction_id"`
	BookID     int       `json:"book_id"`
	DateBorrow time.Time `json:"created_at"`
}
