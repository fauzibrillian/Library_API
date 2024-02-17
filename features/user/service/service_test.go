package service_test

import (
	"errors"
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
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
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

	t.Run("Failed Case - Login Required", func(t *testing.T) {
		var userID = uint(2)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})
		input := user.User{ID: uint(2), NewPassword: "newpass"}

		users, err := userService.ResetPassword(token, input)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, users)

		repoMock.AssertExpectations(t)
	})

	t.Run("Failed Case - Wrong ID", func(t *testing.T) {
		var userID = uint(3)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		input := user.User{ID: uint(1), NewPassword: "newpass"}

		users, err := userService.ResetPassword(token, input)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, users)

		repoMock.AssertExpectations(t)
	})

	t.Run("Failed Case - User Not Found", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		input := user.User{ID: uint(1), NewPassword: "newpass"}

		repoMock.On("GetUserByID", uint(1)).Return(nil, errors.New("user tidak ditemukan")).Once()

		users, err := userService.ResetPassword(token, input)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, users)

		repoMock.AssertExpectations(t)
	})

	t.Run("Failed Case - Wrong Password, Admin Required", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		input := user.User{ID: uint(1), Password: "wrongpass"}
		baseUser := user.User{ID: uint(1), Password: "hashpass"}

		repoMock.On("GetUserByID", uint(1)).Return(&baseUser, nil).Once()
		enkrip.On("Compare", baseUser.Password, input.Password).Return(errors.New("password salah")).Once()

		users, err := userService.ResetPassword(token, input)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, users)

		repoMock.AssertExpectations(t)
		enkrip.AssertExpectations(t)
	})

	t.Run("Failed Case - Empty Password", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})

		input := user.User{ID: uint(1), Password: "", NewPassword: "wrongpass"}
		baseUser := user.User{ID: uint(1), Password: ""}

		repoMock.On("GetUserByID", uint(1)).Return(&baseUser, nil).Once()

		users, err := userService.ResetPassword(token, input)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, users)

		repoMock.AssertExpectations(t)

	})

	t.Run("Failed Case - Empty New Password", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		input := user.User{ID: uint(1), Password: "hashpass", NewPassword: "testis123"}
		baseUser := user.User{ID: uint(1), Password: "hashpass", NewPassword: ""}

		repoMock.On("GetUserByID", uint(1)).Return(&baseUser, nil).Once()
		enkrip.On("Compare", baseUser.Password, input.Password).Return(nil).Once()
		enkrip.On("HashPassword", input.NewPassword).Return("", errors.New("masukkan password baru dengan benar")).Once()
		users, err := userService.ResetPassword(token, input)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, users)

		repoMock.AssertExpectations(t)
		enkrip.AssertExpectations(t)
	})

	t.Run("Failed Case - Repository Error", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		input := user.User{ID: uint(1), Password: "hashpass", NewPassword: "testis123"}
		baseUser := user.User{ID: uint(1), Password: "hashpass", NewPassword: ""}

		repoMock.On("GetUserByID", uint(1)).Return(&baseUser, nil).Once()
		repoMock.On("ResetPassword", input).Return(input, errors.New("kesalahan pada database")).Once()

		enkrip.On("HashPassword", input.NewPassword).Return(input.NewPassword, nil).Once()

		users, err := userService.ResetPassword(token, input)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, users)

		repoMock.AssertExpectations(t)
		enkrip.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	repoMock := mocks.NewRepository(t)

	userService := service.New(repoMock, nil)
	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})

		input := user.User{
			ID:     uint(1),
			Name:   "UpdatedUser",
			Email:  "updatedemail@gmail.com",
			Avatar: "https://cloudinary.com",
			Role:   "admin",
		}
		baseUser := user.User{ID: uint(1), Name: "Budi", Email: "budi@gmail.com", Avatar: "https://cloudinary.com", Role: "admin"}
		repoMock.On("GetUserByID", userID).Return(&baseUser, nil).Once()
		repoMock.On("UpdateUser", input).Return(input, nil).Once()

		updatedUser, err := userService.UpdateUser(token, input)
		assert.NoError(t, err)
		assert.Equal(t, user.User{ID: uint(1), Name: "UpdatedUser", Email: "updatedemail@gmail.com", Avatar: "https://cloudinary.com", Role: "admin"}, updatedUser)

		repoMock.AssertExpectations(t)

	})
	t.Run("Failed Case - Login Required", func(t *testing.T) {
		var userID = uint(2)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})
		input := user.User{
			ID:     uint(1),
			Name:   "UpdatedUser",
			Email:  "updatedemail@gmail.com",
			Avatar: "https://cloudinary.com",
			Role:   "admin",
		}
		users, err := userService.UpdateUser(token, input)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, users)

		repoMock.AssertExpectations(t)
	})

	t.Run("Failed Case - Wrong ID", func(t *testing.T) {
		var userID = uint(3)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		input := user.User{
			ID:     uint(1),
			Name:   "UpdatedUser",
			Email:  "updatedemail@gmail.com",
			Avatar: "https://cloudinary.com",
			Role:   "user",
		}

		users, err := userService.UpdateUser(token, input)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, users)

		repoMock.AssertExpectations(t)
	})

	t.Run("Failed Case - User not Found", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		input := user.User{
			ID:     uint(1),
			Name:   "UpdatedUser",
			Email:  "updatedemail@gmail.com",
			Avatar: "https://cloudinary.com",
			Role:   "user",
		}

		repoMock.On("GetUserByID", uint(1)).Return(nil, errors.New("user tidak ditemukan")).Once()

		users, err := userService.UpdateUser(token, input)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, users)

		repoMock.AssertExpectations(t)
	})

	t.Run("Failed Case - Repository Error", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		input := user.User{
			ID:     uint(1),
			Name:   "UpdatedUser",
			Email:  "updatedemail@gmail.com",
			Avatar: "https://cloudinary.com",
			Role:   "admin",
		}
		baseUser := user.User{ID: uint(1), Name: "", Email: "", Avatar: "", Role: ""}

		repoMock.On("GetUserByID", uint(1)).Return(&baseUser, nil).Once()
		repoMock.On("UpdateUser", input).Return(input, errors.New("kesalahan pada database")).Once()

		users, err := userService.UpdateUser(token, input)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, users)

		repoMock.AssertExpectations(t)
	})
}

func TestDeleteUser(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	userService := service.New(repoMock, nil)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})

		repoMock.On("GetUserByID", userID).Return(&user.User{ID: uint(1)}, nil).Once()
		repoMock.On("DeleteUser", userID).Return(nil).Once()
		err := userService.DeleteUser(token, userID)

		assert.Nil(t, err)
		assert.NoError(t, err)
		repoMock.AssertExpectations(t)
	})

	t.Run("Failed Case - Login Required", func(t *testing.T) {
		var userID = uint(2)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})
		err := userService.DeleteUser(token, userID)

		assert.Error(t, err)

	})

	t.Run("Failed Case - Empty Role", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = ""
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		err := userService.DeleteUser(token, userID)

		assert.Error(t, err)

	})

	t.Run("Failed Case - No Permission", func(t *testing.T) {
		var userID = uint(3)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		err := userService.DeleteUser(token, uint(1))

		assert.Error(t, err)

	})

	t.Run("Failed Case - Empty User ID", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		repoMock.On("GetUserByID", userID).Return(nil, errors.New("failed to retrieve the user for deletion")).Once()

		err := userService.DeleteUser(token, userID)

		assert.Error(t, err)

	})

	t.Run("Failed Case - UserID not equal ID User", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		repoMock.On("GetUserByID", userID).Return(&user.User{ID: uint(2), Role: "user"}, nil).Once()

		err := userService.DeleteUser(token, userID)

		assert.Error(t, err)
	})

	t.Run("Failed Case - Repository Error", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})

		repoMock.On("GetUserByID", userID).Return(&user.User{ID: uint(2), Role: "user"}, nil).Once()
		repoMock.On("DeleteUser", userID).Return(errors.New("repository error")).Once()

		err := userService.DeleteUser(token, userID)

		assert.Error(t, err)
	})
}

func TestSearchUser(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	userService := service.New(repoMock, nil)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		expectedUsers := []user.User{{ID: 1, Name: "User1"}, {ID: 2, Name: "User2"}}
		expectedTotalPage := uint(3)

		repoMock.On("SearchUser", userID, "User1", uint(1), uint(10)).Return(expectedUsers, expectedTotalPage, nil).Once()
		users, totalPage, err := userService.SearchUser(token, "User1", uint(1), uint(10))
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Equal(t, expectedTotalPage, totalPage)
	})

	t.Run("Failed Case - Login Required", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})

		users, totalPage, err := userService.SearchUser(token, "User1", uint(1), uint(10))
		assert.Error(t, err)
		assert.Equal(t, []user.User{}, users)
		assert.Equal(t, uint(0), totalPage)
	})

	t.Run("Failed Case - Admin Required", func(t *testing.T) {
		var userID = uint(3)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})

		users, totalPage, err := userService.SearchUser(token, "User1", uint(1), uint(10))
		assert.Error(t, err)
		assert.Equal(t, []user.User{}, users)
		assert.Equal(t, uint(0), totalPage)
	})

	t.Run("Failed Case - Repository Error", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser, "rahasiabanget")
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasiabanget"), nil
		})
		repoMock.On("SearchUser", userID, "User1", uint(1), uint(10)).Return(nil, uint(0), errors.New("repository error")).Once()

		users, totalPage, err := userService.SearchUser(token, "User1", uint(1), uint(10))
		assert.Error(t, err)
		assert.Equal(t, []user.User([]user.User(nil)), users)
		assert.Equal(t, uint(0), totalPage)
		assert.Equal(t, "repository error", err.Error())

		repoMock.AssertExpectations(t)

	})
}
