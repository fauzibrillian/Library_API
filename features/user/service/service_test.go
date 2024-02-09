package service_test

import (
	"errors"
	"library_api/config"
	"library_api/features/user"
	"library_api/features/user/mocks"
	"library_api/features/user/service"
	eMock "library_api/helper/enkrip/mocks"
	golangjwt "library_api/helper/jwt"

	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	repo := mocks.NewRepository(t)
	enkrip := eMock.NewHashInterface(t)
	m := service.New(repo, enkrip)

	var newUser = user.User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	t.Run("Success Case", func(t *testing.T) {
		enkrip.On("HashPassword", newUser.Password).Return("password123", nil).Once()

		repo.On("Register", newUser).Return(user.User{ID: 1, Email: newUser.Email}, nil).Once()

		result, err := m.Register(newUser)

		assert.Nil(t, err)
		assert.Equal(t, user.User{ID: 1, Email: newUser.Email}, result)

		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Duplicate Data", func(t *testing.T) {
		enkrip.On("HashPassword", newUser.Password).Return("password123", nil).Once()

		repo.On("Register", newUser).Return(user.User{}, errors.New("duplicate key")).Once()

		_, err := m.Register(newUser)

		assert.Error(t, err)
		assert.Equal(t, "data telah terdaftar pada sistem", err.Error())

		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Error Case - General Error", func(t *testing.T) {
		enkrip.On("HashPassword", newUser.Password).Return("password123", nil).Once()

		repo.On("Register", newUser).Return(user.User{}, errors.New("general error")).Once()

		_, err := m.Register(newUser)

		assert.Error(t, err)
		assert.Equal(t, "terjadi kesalahan pada sistem", err.Error())

		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Empty Name", func(t *testing.T) {
		newUser := user.User{
			Email:    "johndoe@gmail.com",
			Name:     "",
			Password: "johndoe123",
		}
		_, err := m.Register(newUser)

		assert.Error(t, err)
		assert.Equal(t, "name cannot be empty", err.Error())
	})

	t.Run("Error Case - Empty Email", func(t *testing.T) {
		newUser := user.User{
			Email:    "",
			Name:     "John Doe",
			Password: "johndoe123",
		}
		_, err := m.Register(newUser)

		assert.Error(t, err)
		assert.Equal(t, "email cannot be empty", err.Error())
	})

	t.Run("Error Case - Empty Password", func(t *testing.T) {
		newUser := user.User{
			Email:    "johndoe@gmail.com",
			Name:     "John Doe",
			Password: "",
		}
		_, err := m.Register(newUser)

		assert.Error(t, err)
		assert.Equal(t, "password cannot be empty", err.Error())
	})

	t.Run("Error Case - Encryption Failure", func(t *testing.T) {
		enkrip.On("HashPassword", newUser.Password).Return("", errors.New("terdapat masalah saat memproses enkripsi password")).Once()
		_, err := m.Register(newUser)
		assert.Error(t, err)
		assert.Equal(t, "terdapat masalah saat memproses enkripsi password", err.Error())
		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	repo := mocks.NewRepository(t)
	hashMock := eMock.NewHashInterface(t)
	userService := service.New(repo, hashMock)

	email := "test@example.com"
	password := "test_password"
	hashedPassword := "hashed_test_password"

	t.Run("Success Case", func(t *testing.T) {
		repo.On("Login", email).Return(user.User{
			ID:       1,
			Email:    "test@example.com",
			Password: hashedPassword,
		}, nil).Once()

		hashMock.On("Compare", hashedPassword, password).Return(nil).Once()

		result, err := userService.Login(email, password)

		assert.Nil(t, err)
		assert.Equal(t, user.User{
			ID:       1,
			Email:    "test@example.com",
			Password: hashedPassword,
		}, result)

		repo.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})

	t.Run("Error Case - Empty Email and Password", func(t *testing.T) {
		result, err := userService.Login("", "")

		assert.Error(t, err)
		assert.Equal(t, "email and password are required", err.Error())
		assert.Equal(t, user.User{}, result)

		repo.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})

	t.Run("Error Case - User Not Found", func(t *testing.T) {
		repo.On("Login", email).Return(user.User{}, errors.New("not found")).Once()

		result, err := userService.Login(email, password)

		assert.Error(t, err)
		assert.Equal(t, "data tidak ditemukan", err.Error())
		assert.Equal(t, user.User{}, result)

		repo.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})

	t.Run("Error Case - System Error", func(t *testing.T) {
		repo.On("Login", email).Return(user.User{}, errors.New("data tidak ditemukan")).Once()

		result, err := userService.Login(email, password)

		assert.Error(t, err)
		assert.Equal(t, "data tidak ditemukan", err.Error())
		assert.Equal(t, user.User{}, result)

		repo.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})

	t.Run("Error Case - Incorrect Password", func(t *testing.T) {
		repo.On("Login", email).Return(user.User{
			ID:       1,
			Email:    "test@example.com",
			Password: hashedPassword,
		}, nil).Once()

		hashMock.On("Compare", hashedPassword, password).Return(errors.New("password salah")).Once()

		result, err := userService.Login(email, password)

		assert.Error(t, err)
		assert.Equal(t, "password salah", err.Error())
		assert.Equal(t, user.User{}, result)

		repo.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})
}

func TestResetPassword(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	enkrip := eMock.NewHashInterface(t)

	userService := service.New(repoMock, enkrip)
	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.InitConfig().JWT), nil
		})

		input := user.User{ID: uint(1), NewPassword: "newpass"}
		baseUser := user.User{ID: uint(1), Password: "oldpass"}
		repoMock.On("GetUserByID", uint(1)).Return(&baseUser, nil).Once()
		enkrip.On("HashPassword", "newpass").Return("hashednewpass", nil).Once()
		input.NewPassword = "hashednewpass"
		repoMock.On("ResetPassword", input).Return(user.User{ID: uint(1), Password: "hashednewpass"}, nil).Once()

		input.NewPassword = "newpass"
		resetPassword, err := userService.ResetPassword(token, input)

		assert.NoError(t, err)
		assert.Equal(t, user.User{ID: uint(1), Password: "hashednewpass"}, resetPassword)

		repoMock.AssertExpectations(t)
		enkrip.AssertExpectations(t)
	})
}
