package service_test

import (
	"errors"
	"library_api/features/transaction"
	"library_api/features/transaction/mocks"
	"library_api/features/transaction/service"
	golangjwt "library_api/helper/jwt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestBorrow(t *testing.T) {
	repo := mocks.NewRepository(t)
	m := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		input := transaction.Transaction{
			ID:         uint(1),
			BookID:     uint(2),
			DateBorrow: time.Date(2024, time.January, 14, 0, 0, 0, 0, time.UTC),
		}
		repo.On("Borrow", uint(1), uint(2)).Return(input, nil).Once()
		transactions, err := m.Borrow(token, uint(2))
		assert.NoError(t, err)
		assert.Equal(t, input, transactions)

		repo.AssertExpectations(t)

	})

	t.Run("Failed Case - Login Required", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})

		transactions, err := m.Borrow(token, uint(2))
		assert.Error(t, err)
		assert.Equal(t, transaction.Transaction{}, transactions)
	})

	t.Run("Failed Case - Roles Empty", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = ""
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})

		transactions, err := m.Borrow(token, uint(2))
		assert.Error(t, err)
		assert.Equal(t, transaction.Transaction{}, transactions)
	})

	t.Run("Failed Case - Empty Input", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		repo.On("Borrow", uint(1), uint(2)).Return(transaction.Transaction{}, errors.New("repository error")).Once()

		transactions, err := m.Borrow(token, uint(2))
		assert.Error(t, err)
		assert.Equal(t, transaction.Transaction{}, transactions)
	})
}

func TestUpdateReturn(t *testing.T) {
	repo := mocks.NewRepository(t)
	m := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		input := transaction.Transaction{
			ID:         uint(1),
			BookID:     uint(2),
			DateBorrow: time.Date(2024, time.January, 14, 0, 0, 0, 0, time.UTC),
			DateReturn: time.Date(2024, time.January, 16, 0, 0, 0, 0, time.UTC),
		}

		repo.On("UpdateReturn", userID, uint(1), input).Return([]transaction.Transaction{input}, nil).Once()
		transactions, err := m.UpdateReturn(token, uint(1), input)
		assert.NoError(t, err)
		assert.Equal(t, []transaction.Transaction{input}, transactions)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Login Required", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})
		input := transaction.Transaction{
			ID:         uint(1),
			BookID:     uint(2),
			DateBorrow: time.Date(2024, time.January, 14, 0, 0, 0, 0, time.UTC),
			DateReturn: time.Date(2024, time.January, 16, 0, 0, 0, 0, time.UTC),
		}
		transactions, err := m.UpdateReturn(token, uint(1), input)
		assert.Error(t, err)
		assert.Equal(t, []transaction.Transaction{}, transactions)
	})

	t.Run("Failed Case - Roles Empty", func(t *testing.T) {
		t.Run("Failed Case - Login Required", func(t *testing.T) {
			var userID = uint(1)
			var rolesUser = ""
			var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
			var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
				return []byte("rahasiabanget"), nil
			})
			input := transaction.Transaction{
				ID:         uint(1),
				BookID:     uint(2),
				DateBorrow: time.Date(2024, time.January, 14, 0, 0, 0, 0, time.UTC),
				DateReturn: time.Date(2024, time.January, 16, 0, 0, 0, 0, time.UTC),
			}
			transactions, err := m.UpdateReturn(token, uint(1), input)
			assert.Error(t, err)
			assert.Equal(t, []transaction.Transaction{}, transactions)
		})
	})

	t.Run("Failed Case - Admin Required", func(t *testing.T) {
		t.Run("Failed Case - Login Required", func(t *testing.T) {
			var userID = uint(1)
			var rolesUser = "user"
			var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
			var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
				return []byte("rahasiabanget"), nil
			})
			input := transaction.Transaction{
				ID:         uint(1),
				BookID:     uint(2),
				DateBorrow: time.Date(2024, time.January, 14, 0, 0, 0, 0, time.UTC),
				DateReturn: time.Date(2024, time.January, 16, 0, 0, 0, 0, time.UTC),
			}
			transactions, err := m.UpdateReturn(token, uint(1), input)
			assert.Error(t, err)
			assert.Equal(t, []transaction.Transaction{}, transactions)
		})
	})
}

func TestAllTransaction(t *testing.T) {
	repo := mocks.NewRepository(t)
	m := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})

		expectedBooks := []transaction.Transaction{{ID: uint(1), BookID: uint(2)}, {ID: uint(2), BookID: uint(3)}}
		expectedTotalPage := 3

		repo.On("AllTransaction", uint(1), "Harry Potter", uint(1), uint(10)).Return(expectedBooks, expectedTotalPage, nil).Once()
		books, totalPage, err := m.AllTransaction(token, "Harry Potter", uint(1), uint(10))
		assert.NoError(t, err)
		assert.Equal(t, expectedBooks, books)
		assert.Equal(t, expectedTotalPage, totalPage)
	})
	t.Run("Failed Case - Login required", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})
		books, totalPage, err := m.AllTransaction(token, "Harry Potter", uint(1), uint(10))
		assert.Error(t, err)
		assert.Equal(t, []transaction.Transaction{}, books)
		assert.Equal(t, 0, totalPage)
	})

	t.Run("Failed Case - Login required", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		books, totalPage, err := m.AllTransaction(token, "Harry Potter", uint(1), uint(10))
		assert.Error(t, err)
		assert.Equal(t, []transaction.Transaction{}, books)
		assert.Equal(t, 0, totalPage)
	})
}
