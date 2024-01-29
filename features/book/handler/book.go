package handler

import (
	"context"
	"errors"
	"library_api/features/book"
	"library_api/helper/cld"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	golangjwt "github.com/golang-jwt/jwt/v5"
	echo "github.com/labstack/echo/v4"
)

type BookHandler struct {
	s      book.Service
	cl     *cloudinary.Cloudinary
	ct     context.Context
	folder string
}

func New(s book.Service, cld *cloudinary.Cloudinary, ctx context.Context, uploadparam string) book.Handler {
	return &BookHandler{
		s:      s,
		cl:     cld,
		ct:     ctx,
		folder: uploadparam,
	}
}

// AddBook implements book.Handler.
func (bh *BookHandler) AddBook() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(BookRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
			})
		}
		formHeader, err := c.FormFile("picture")
		if err != nil {
			if errors.Is(err, http.ErrMissingFile) {
				inputProcess := &book.Book{
					Tittle:    input.Tittle,
					Publisher: input.Publisher,
					Author:    input.Author,
					Picture:   "",
				}

				result, err := bh.s.AddBook(c.Get("user").(*golangjwt.Token), *inputProcess)

				if err != nil {
					c.Logger().Error("ERROR Register, explain:", err.Error())
					var statusCode = http.StatusInternalServerError
					var message = "terjadi permasalahan ketika memproses data"

					if strings.Contains(err.Error(), "terdaftar") {
						statusCode = http.StatusBadRequest
						message = "data yang diinputkan sudah terdaftar ada sistem"
					}

					return c.JSON(statusCode, map[string]any{
						"message": message,
					})
				}

				var response = new(BookResponse)
				response.ID = result.ID
				response.Tittle = result.Tittle
				response.Publisher = result.Publisher
				response.Author = result.Author
				response.Picture = result.Picture

				return c.JSON(http.StatusCreated, map[string]any{
					"message": "Success Create Book Data",
					"data":    response,
				})

			}
			return c.JSON(
				http.StatusInternalServerError, map[string]any{
					"message": "formheader error",
				})

		}
		formFile, err := formHeader.Open()
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError, map[string]any{
					"message": "formfile error",
				})
		}

		link, err := cld.UploadImage(bh.cl, bh.ct, formFile, bh.folder)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "harap pilih gambar",
					"data":    nil,
				})
			} else {
				return c.JSON(http.StatusInternalServerError, map[string]any{
					"message": "kesalahan pada server",
					"data":    nil,
				})
			}
		}

		var inputProcess = new(book.Book)
		inputProcess.Tittle = input.Tittle
		inputProcess.Publisher = input.Publisher
		inputProcess.Author = input.Author
		inputProcess.Tittle = input.Tittle
		inputProcess.Picture = link

		result, err := bh.s.AddBook(c.Get("user").(*golangjwt.Token), *inputProcess)
		if err != nil {
			c.Logger().Error("ERROR Register Book, explain:", err.Error())

			if strings.Contains(err.Error(), "terdaftar") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "data duplicate",
				})
			}

			if strings.Contains(err.Error(), "tidak memiliki izin") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "tidak memiliki izin",
				})
			} else if strings.Contains(err.Error(), "tidak memiliki izin") {
				return c.JSON(http.StatusForbidden, map[string]any{
					"message": "tidak memiliki izin",
				})
			}

			return c.JSON(http.StatusForbidden, map[string]any{
				"message": "Tidak memiliki izin",
			})
		}

		var response = new(BookResponse)
		response.ID = result.ID
		response.Tittle = result.Tittle
		response.Publisher = result.Publisher
		response.Author = result.Author
		response.Picture = result.Picture

		return c.JSON(http.StatusOK, map[string]any{
			"message": "Success Created Book Data",
			"data":    response,
		})
	}
}

// UpdateBook implements book.Handler.
func (bh *BookHandler) UpdateBook() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(BookPutRequest)
		bookID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID user tidak valid",
				"data":    nil,
			})
		}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
				"data":    nil,
			})
		}

		formHeader, err := c.FormFile("picture")
		if err != nil {
			if errors.Is(err, http.ErrMissingFile) {
				bookID, err := strconv.ParseUint(c.Param("id"), 10, 64)
				if err != nil {
					return c.JSON(http.StatusBadRequest, map[string]interface{}{
						"message": "ID user tidak valid",
					})
				}
				updatedBook := book.Book{
					ID:        input.ID,
					Tittle:    input.Tittle,
					Author:    input.Author,
					Publisher: input.Publisher,
					Picture:   "",
				}

				result, err := bh.s.UpdateBook(c.Get("user").(*golangjwt.Token), uint(bookID), updatedBook)

				if err != nil {
					c.Logger().Error("ERROR Updating Book, explain:", err.Error())
					var statusCode = http.StatusInternalServerError
					var message = "terjadi permasalahan ketika memproses data"

					if strings.Contains(err.Error(), "admin role required") {
						statusCode = http.StatusUnauthorized
						message = "Anda tidak memiliki izin untuk mengupdate produk"
					} else if strings.Contains(err.Error(), "book tidak ditemukan") {
						statusCode = http.StatusNotFound
						message = "Book tidak ditemukan"
					}

					return c.JSON(statusCode, map[string]any{
						"message": message,
					})
				}

				var response = new(BookPutResponse)
				response.ID = result.ID
				response.Tittle = result.Tittle
				response.Author = result.Author
				response.Picture = result.Picture
				response.Publisher = result.Publisher

				return c.JSON(http.StatusOK, map[string]any{
					"message": "Success Updated data",
					"data":    response,
				})
			}
			return c.JSON(
				http.StatusBadRequest, map[string]any{
					"message": "formheader error",
					"data":    nil,
				})

		}

		formFile, err := formHeader.Open()
		if err != nil {
			return c.JSON(
				http.StatusBadRequest, map[string]any{
					"message": "formfile error",
					"data":    nil,
				})
		}
		defer formFile.Close()

		link, err := cld.UploadImage(bh.cl, bh.ct, formFile, bh.folder)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "harap pilih gambar",
					"data":    nil,
				})
			} else {
				return c.JSON(http.StatusInternalServerError, map[string]any{
					"message": "kesalahan pada server",
					"data":    nil,
				})
			}
		}

		updatedBook := book.Book{
			ID:        input.ID,
			Tittle:    input.Tittle,
			Author:    input.Author,
			Publisher: input.Publisher,
			Picture:   link,
		}

		result, err := bh.s.UpdateBook(c.Get("user").(*golangjwt.Token), uint(bookID), updatedBook)

		if err != nil {
			c.Logger().Error("ERROR Updating Book, explain:", err.Error())
			var statusCode = http.StatusInternalServerError
			var message = "terjadi permasalahan ketika memproses data"

			if strings.Contains(err.Error(), "admin role required") {
				statusCode = http.StatusUnauthorized
				message = "Anda tidak memiliki izin untuk mengupdate produk"
			} else if strings.Contains(err.Error(), "Book tidak ditemukan") {
				statusCode = http.StatusNotFound
				message = "Book tidak ditemukan"
			}

			return c.JSON(statusCode, map[string]any{
				"message": message,
			})
		}

		var response = new(BookPutResponse)
		response.ID = result.ID
		response.Tittle = result.Tittle
		response.Publisher = result.Publisher
		response.Author = result.Author
		response.Picture = result.Picture

		return c.JSON(http.StatusCreated, map[string]any{
			"message": "Success Updated Data",
			"data":    response,
		})
	}
}

// DeleteBook implements book.Handler.
func (bh *BookHandler) DeleteBook() echo.HandlerFunc {
	return func(c echo.Context) error {
		bookID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID book tidak valid",
				"data":    nil,
			})
		}

		errDel := bh.s.DelBook(c.Get("user").(*golangjwt.Token), uint(bookID))

		if errDel != nil {
			c.Logger().Error("ERROR Deleting Book, explain:", errDel.Error())
			var statusCode = http.StatusInternalServerError
			var message = "terjadi permasalahan ketika memproses data"

			if strings.Contains(errDel.Error(), "admin role required") {
				statusCode = http.StatusUnauthorized
				message = "Anda tidak memiliki izin untuk menghapus book"
			} else if strings.Contains(errDel.Error(), "book tidak ditemukan") {
				statusCode = http.StatusNotFound
				message = "Book tidak ditemukan"
			}

			return c.JSON(statusCode, map[string]interface{}{
				"message": message,
			})
		}
		return c.JSON(http.StatusOK, map[string]any{
			"message": "Delete Book Success",
		})
	}
}

// SearchBook implements book.Handler.
func (bh *BookHandler) SearchBook() echo.HandlerFunc {
	return func(c echo.Context) error {
		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil || page <= 0 {
			page = 1
		}

		limit, err := strconv.Atoi(c.QueryParam("limit"))
		if err != nil || limit <= 0 {
			limit = 10
		}

		tittle := c.QueryParam("tittle")
		uintPage := uint(page)
		uintLimit := uint(limit)

		books, totalPage, err := bh.s.SearchBook(tittle, uintPage, uintLimit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		var response []SearchBookResponse
		for _, result := range books {
			response = append(response, SearchBookResponse{
				ID:      result.ID,
				Tittle:  result.Tittle,
				Picture: result.Picture,
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":    "Get Book Successful",
			"data":       response,
			"pagination": map[string]interface{}{"page": page, "limit": limit, "total_page": totalPage},
		})
	}
}

// GetBook implements book.Handler.
func (bh *BookHandler) GetBook() echo.HandlerFunc {
	return func(c echo.Context) error {
		bookID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID book tidak valid",
				"data":    nil,
			})
		}

		result, err := bh.s.GetBook(uint(bookID))
		if err != nil {
			c.Logger().Error("Error fetching product: ", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Failed to retrieve product data",
			})
		}
		var response = new(BookResponse)
		response.ID = result.ID
		response.Tittle = result.Tittle
		response.Publisher = result.Publisher
		response.Author = result.Author
		response.Picture = result.Picture

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Get Details Books Successful",
			"data":    response,
		})
	}
}
