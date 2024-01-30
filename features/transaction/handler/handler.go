package handler

import (
	"library_api/features/book"
	"library_api/features/transaction"
	"net/http"
	"strings"

	golangjwt "github.com/golang-jwt/jwt/v5"
	echo "github.com/labstack/echo/v4"
)

type TransactionHandler struct {
	p book.Handler
	s transaction.Service
}

func New(s transaction.Service) transaction.Handler {
	return &TransactionHandler{
		s: s,
	}
}

// Borrow implements transaction.Handler.
func (th *TransactionHandler) Borrow() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(TransactionRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
			})
		}
		result, err := th.s.Borrow(c.Get("user").(*golangjwt.Token), uint(input.BookID))
		if err != nil {
			c.Logger().Error("terjadi kesalahan", err.Error())
			if strings.Contains(err.Error(), "duplicate") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "dobel input nama",
				})
			}
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "transaction duplicate",
			})
		}

		var response = new(TransactionResponse)
		response.ID = int(result.ID)
		response.BookID = int(result.BookID)
		response.DateBorrow = result.DateBorrow

		return c.JSON(http.StatusCreated, map[string]any{
			"message": "Transaction Borrow created successfully",
			"data":    response,
		})
	}
}
