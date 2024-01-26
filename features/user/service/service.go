package service

import (
	"errors"
	"library_api/features/user"
	"library_api/helper/enkrip"
	"library_api/helper/jwt"
	"strings"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	repo user.Repository
	hash enkrip.HashInterface
}

func New(r user.Repository, h enkrip.HashInterface) user.Service {
	return &UserService{
		repo: r,
		hash: h,
	}
}

// Login implements user.Service.
func (us *UserService) Login(email string, password string) (user.User, error) {
	if email == "" || password == "" {
		return user.User{}, errors.New("email and password are required")
	}
	result, err := us.repo.Login(email)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return user.User{}, errors.New("data tidak ditemukan")
		}
		return user.User{}, errors.New("data tidak ditemukan")
	}

	err = us.hash.Compare(result.Password, password)

	if err != nil {
		return user.User{}, errors.New("password salah")
	}

	return result, nil
}

// Register implements user.Service.
func (us *UserService) Register(newUser user.User) (user.User, error) {
	if newUser.Name == "" {
		return user.User{}, errors.New("name cannot be empty")
	}
	if newUser.Email == "" {
		return user.User{}, errors.New("email cannot be empty")
	}
	if newUser.Password == "" {
		return user.User{}, errors.New("password cannot be empty")
	}
	// enkripsi password
	ePassword, err := us.hash.HashPassword(newUser.Password)

	if err != nil {
		return user.User{}, errors.New("terdapat masalah saat memproses enkripsi password")
	}

	newUser.Password = ePassword
	result, err := us.repo.Register(newUser)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return user.User{}, errors.New("data telah terdaftar pada sistem")
		}
		return user.User{}, errors.New("terjadi kesalahan pada sistem")
	}

	return result, nil
}

// ResetPassword implements user.Service.
func (us *UserService) ResetPassword(token *golangjwt.Token, input user.User) (user.User, error) {
	userID, rolesUser, err := jwt.ExtractToken(token)
	if err != nil {
		return user.User{}, errors.New("harap login")
	}

	if userID != input.ID && rolesUser != "admin" {
		return user.User{}, errors.New("id tidak cocok")
	}

	base, err := us.repo.GetUserByID(userID)
	if err != nil {
		return user.User{}, errors.New("user tidak ditemukan")
	}
	if input.Password != "" && rolesUser != "admin" {
		err = us.hash.Compare(base.Password, input.Password)

		if err != nil {
			return user.User{}, errors.New("password salah")
		}
	}

	if input.NewPassword != "" {
		if input.Password == "" && rolesUser != "admin" {
			return user.User{}, errors.New("masukkan password yang lama")
		}
		newpass, err := us.hash.HashPassword(input.NewPassword)
		if err != nil {
			return user.User{}, errors.New("masukkan password baru dengan benar")
		}
		input.NewPassword = newpass
	}

	respons, err := us.repo.ResetPassword(input)
	if err != nil {

		return user.User{}, errors.New("kesalahan pada database")
	}
	return respons, nil
}

// UpdateUser implements user.Service.
func (us *UserService) UpdateUser(token *golangjwt.Token, input user.User) (user.User, error) {
	userID, rolesUser, err := jwt.ExtractToken(token)
	if err != nil {
		return user.User{}, errors.New("harap login")
	}

	if userID != input.ID && rolesUser != "admin" {
		return user.User{}, errors.New("id tidak cocok")
	}
	exitingUser, err := us.repo.GetUserByID(userID)
	if err != nil {
		return user.User{}, errors.New("failed to retrieve the user for deletion")
	}
	if exitingUser.ID != userID {
		return user.User{}, errors.New("you don't have permission to update this user")
	}
	respons, err := us.repo.UpdateUser(input)
	if err != nil {

		return user.User{}, errors.New("kesalahan pada database")
	}
	return respons, nil
}

// HapusUser implements user.Service.
func (us *UserService) DeleteUser(token *golangjwt.Token, userID uint) error {
	userId, rolesUser, err := jwt.ExtractToken(token)
	if err != nil {
		return errors.New("harap login")
	}
	if rolesUser == "" {
		return errors.New("you don't have permission to delete this user")
	}
	if rolesUser != "admin" && userId != userID {
		return errors.New("you don't have permission to delete this user")
	}
	exitingUser, err := us.repo.GetUserByID(userID)
	if err != nil {
		return errors.New("failed to retrieve the user for deletion")
	}
	if exitingUser.ID != userId && rolesUser != "admin" {
		return errors.New("you don't have permission to delete this user")
	}
	err = us.repo.DeleteUser(userID)
	if err != nil {
		return errors.New("failed to delete the user")
	}

	return nil
}
