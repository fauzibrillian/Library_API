package repository

import (
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
