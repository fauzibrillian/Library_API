package service_test

import (
	"errors"
	"library_api/features/book"
	"library_api/features/book/mocks"
	"library_api/features/book/service"
	golangjwt "library_api/helper/jwt"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAddBook(t *testing.T) {
	repo := mocks.NewRepository(t)
	m := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		input := book.Book{
			ID:        uint(1),
			Tittle:    "Harry Potter",
			Publisher: "Gramedia",
			Author:    "JK. Ronald",
			Picture:   "www.cloudinary.com",
		}
		repo.On("InsertBook", uint(1), input).Return(input, nil).Once()
		products, err := m.AddBook(token, input)

		assert.NoError(t, err, products)
		assert.Equal(t, book.Book{
			ID:        uint(1),
			Tittle:    "Harry Potter",
			Publisher: "Gramedia",
			Author:    "JK. Ronald",
			Picture:   "www.cloudinary.com"}, products)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Login required", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})
		input := book.Book{
			ID:        uint(1),
			Tittle:    "Harry Potter",
			Publisher: "Gramedia",
			Author:    "JK. Ronald",
			Picture:   "www.cloudinary.com",
		}
		books, err := m.AddBook(token, input)
		assert.Error(t, err)
		assert.Equal(t, book.Book{}, books)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Admin Role Required", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		input := book.Book{
			ID:        uint(1),
			Tittle:    "Harry Potter",
			Publisher: "Gramedia",
			Author:    "JK. Ronald",
			Picture:   "www.cloudinary.com",
		}
		products, err := m.AddBook(token, input)
		assert.Error(t, err)
		assert.Equal(t, book.Book{}, products)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Empty Input", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		input := book.Book{
			ID:        uint(1),
			Tittle:    "Harry Potter",
			Publisher: "Gramedia",
			Author:    "JK. Ronald",
			Picture:   "www.cloudinary.com",
		}
		repo.On("InsertBook", uint(1), input).Return(book.Book{}, errors.New("inputan tidak boleh kosong")).Once()

		products, err := m.AddBook(token, input)
		assert.Error(t, err)
		assert.Equal(t, book.Book{}, products)

		repo.AssertExpectations(t)
	})
}

func TestUpdateBook(t *testing.T) {
	repo := mocks.NewRepository(t)
	m := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})

		input := book.Book{
			ID:        uint(1),
			Tittle:    "Harry Potter",
			Publisher: "Gramedia",
			Author:    "JK. Ronald",
			Picture:   "www.cloudinary.com",
		}
		repo.On("UpdateBook", uint(1), uint(1), input).Return(input, nil).Once()
		products, err := m.UpdateBook(token, uint(1), input)

		assert.NoError(t, err, products)
		assert.Equal(t, book.Book{
			ID:        uint(1),
			Tittle:    "Harry Potter",
			Publisher: "Gramedia",
			Author:    "JK. Ronald",
			Picture:   "www.cloudinary.com"}, products)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Login Required", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})

		input := book.Book{
			ID:        uint(1),
			Tittle:    "Harry Potter",
			Publisher: "Gramedia",
			Author:    "JK. Ronald",
			Picture:   "www.cloudinary.com",
		}
		books, err := m.UpdateBook(token, uint(1), input)
		assert.Error(t, err)
		assert.Equal(t, book.Book{}, books)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Admin Role Required", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})

		input := book.Book{
			ID:        uint(1),
			Tittle:    "Harry Potter",
			Publisher: "Gramedia",
			Author:    "JK. Ronald",
			Picture:   "www.cloudinary.com",
		}
		books, err := m.UpdateBook(token, uint(1), input)
		assert.Error(t, err)
		assert.Equal(t, book.Book{}, books)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Empty Input", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})

		input := book.Book{
			ID:        uint(1),
			Tittle:    "Harry Potter",
			Publisher: "Gramedia",
			Author:    "JK. Ronald",
			Picture:   "www.cloudinary.com",
		}
		repo.On("UpdateBook", uint(1), uint(1), input).Return(book.Book{}, errors.New("failed to update the product")).Once()

		books, err := m.UpdateBook(token, uint(1), input)
		assert.Error(t, err)
		assert.Equal(t, book.Book{}, books)

		repo.AssertExpectations(t)
	})
}
func TestDeleteBook(t *testing.T) {
	repo := mocks.NewRepository(t)
	m := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		bookID := uint(1)
		repo.On("DelBook", userID, bookID).Return(nil).Once()
		err := m.DelBook(token, bookID)

		assert.Nil(t, err)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Login Required", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})
		bookID := uint(1)
		err := m.DelBook(token, bookID)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Admin Role Required", func(t *testing.T) {
		var userID = uint(2)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		bookID := uint(1)
		err := m.DelBook(token, bookID)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestSearchBook(t *testing.T) {
	repo := mocks.NewRepository(t)
	m := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		expectedBooks := []book.Book{{ID: 1, Tittle: "Laskar Pelangi"}, {ID: 2, Tittle: "Harry Potter"}}
		expectedTotalPage := uint(3)

		repo.On("SearchBook", "Harry Potter", uint(1), uint(10)).Return(expectedBooks, expectedTotalPage, nil).Once()
		books, totalPage, err := m.SearchBook("Harry Potter", uint(1), uint(10))
		assert.NoError(t, err)
		assert.Equal(t, expectedBooks, books)
		assert.Equal(t, expectedTotalPage, totalPage)
	})

	t.Run("Failed Case - Repository Error", func(t *testing.T) {
		repo.On("SearchBook", "Harry Potter", uint(1), uint(10)).Return(nil, uint(0), errors.New("failed get books data")).Once()
		books, totalPage, err := m.SearchBook("Harry Potter", uint(1), uint(10))
		assert.Error(t, err)
		assert.Equal(t, []book.Book([]book.Book(nil)), books)
		assert.Equal(t, uint(0), totalPage)
		assert.Equal(t, "failed get books data", err.Error())

		repo.AssertExpectations(t)
	})
}
