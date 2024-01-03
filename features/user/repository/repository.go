package repository

import (
	"errors"
	"library_api/features/user"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Phone    string `json:"phone" form:"phone" gorm:"unique"`
	Name     string `json:"name" form:"name"`
	Password string `json:"password" form:"password"`
	Avatar   string `json:"avatar" form:"avatar"`
	Role     string `json:"role" form:"role"`
}

type UserQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) user.Repository {
	return &UserQuery{
		db: db,
	}
}

// Login implements user.Repository.
func (uq *UserQuery) Login(phone string) (user.User, error) {
	var userData = new(UserModel)

	if err := uq.db.Where("phone = ?", phone).First(userData).Error; err != nil {
		return user.User{}, err
	}

	var result = new(user.User)
	result.ID = userData.ID
	result.Name = userData.Name
	result.Phone = userData.Phone
	result.Password = userData.Password
	result.Role = userData.Role

	return *result, nil
}

// Register implements user.Repository.
func (uq *UserQuery) Register(newUser user.User) (user.User, error) {
	var inputDB = new(UserModel)
	inputDB.Name = newUser.Name
	inputDB.Phone = newUser.Phone
	inputDB.Password = newUser.Password
	inputDB.Role = "user"

	if err := uq.db.Create(&inputDB).Error; err != nil {
		return user.User{}, err
	}

	newUser.ID = inputDB.ID

	return newUser, nil
}

// GetUserByID implements user.Repository.
func (uq *UserQuery) GetUserByID(userID uint) (*user.User, error) {
	var userModel UserModel
	if err := uq.db.First(&userModel, userID).Error; err != nil {
		return nil, err
	}

	// Jika tidak ada buku ditemukan
	if userModel.ID == 0 {
		err := errors.New("user tidak ditemukan")
		return nil, err
	}

	result := &user.User{
		ID:       userModel.ID,
		Name:     userModel.Name,
		Phone:    userModel.Phone,
		Password: userModel.Password,
	}

	return result, nil
}

// ResetPassword implements user.Repository.
func (uq *UserQuery) ResetPassword(input user.User) (user.User, error) {
	var proses UserModel
	if err := uq.db.First(&proses, input.ID).Error; err != nil {
		return user.User{}, err
	}

	if proses.ID == 0 {
		err := errors.New("user tidak ditemukan")
		return user.User{}, err
	}

	if input.NewPassword != "" {
		proses.Password = input.NewPassword
	}
	if err := uq.db.Save(&proses).Error; err != nil {

		return user.User{}, err
	}

	result := user.User{
		ID:       proses.ID,
		Name:     proses.Name,
		Phone:    proses.Phone,
		Password: proses.Password,
	}

	return result, nil
}
