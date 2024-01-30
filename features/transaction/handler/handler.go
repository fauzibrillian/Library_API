package handler

import (
	"library_api/features/book"
	"library_api/features/transaction"
	"library_api/features/user/repository"
	"net/http"
	"strconv"
	"strings"

	golangjwt "github.com/golang-jwt/jwt/v5"
	echo "github.com/labstack/echo/v4"
)

type TransactionHandler struct {
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

// AllTransaction implements transaction.Handler.
func (th *TransactionHandler) AllTransaction() echo.HandlerFunc {
	return func(c echo.Context) error {
		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil || page <= 0 {
			page = 1
		}

		limit, err := strconv.Atoi(c.QueryParam("limit"))
		if err != nil || limit <= 0 {
			limit = 10
		}

		name := c.QueryParam("name")
		uintPage := uint(page)
		uintLimit := uint(limit)

		books, totalPage, err := th.s.AllTransaction(c.Get("user").(*golangjwt.Token), name, uintPage, uintLimit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		// slicing data user
		var User []repository.UserModel
		var UserName []string
		var UserPicture []string
		for _, result := range books {
			User = append(User, result.Users...)
		}
		for _, result := range User {
			UserName = append(UserName, result.Name)
			UserPicture = append(UserPicture, result.Avatar)
		}

		// slicing data product
		var Book []book.Book
		var BookTittle []string
		var BookPicture []string

		for _, result := range books {
			Book = append(Book, result.Books...)
		}

		for _, result := range Book {
			BookTittle = append(BookTittle, result.Tittle)
			BookPicture = append(BookPicture, result.Picture)
		}

		// slicing data to response
		var responses []SearchTransactionResponse
		for x, result := range books {
			responses = append(responses, SearchTransactionResponse{
				ID:          int(result.ID),
				UserPicture: UserPicture[x],
				UserName:    UserName[x],
				TittleBook:  BookTittle[x],
				PictureBook: BookPicture[x],
				DateBorrow:  result.DateBorrow,
				DateReturn:  result.DateReturn,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":    "Get Transactions Book Successful",
			"data":       responses,
			"pagination": map[string]interface{}{"page": page, "limit": limit, "total_page": totalPage},
		})
	}
}
