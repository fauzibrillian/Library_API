package repository

import (
	"errors"
	"library_api/features/user"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Email    string `json:"email" form:"email" gorm:"unique"`
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
func (uq *UserQuery) Login(email string) (user.User, error) {
	var userData = new(UserModel)

	if err := uq.db.Where("email = ?", email).First(userData).Error; err != nil {
		return user.User{}, err
	}

	var result = new(user.User)
	result.ID = userData.ID
	result.Name = userData.Name
	result.Email = userData.Email
	result.Password = userData.Password
	result.Role = userData.Role

	return *result, nil
}

// Register implements user.Repository.
func (uq *UserQuery) Register(newUser user.User) (user.User, error) {
	var inputDB = new(UserModel)
	inputDB.Name = newUser.Name
	inputDB.Email = newUser.Email
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
		Email:    userModel.Email,
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
		Email:    proses.Email,
		Password: proses.Password,
	}

	return result, nil
}

// UpdateUser implements user.Repository.
func (uq *UserQuery) UpdateUser(input user.User) (user.User, error) {
	var proses UserModel
	if err := uq.db.First(&proses, input.ID).Error; err != nil {
		return user.User{}, err
	}

	if proses.ID == 0 {
		err := errors.New("user tidak ditemukan")
		return user.User{}, err
	}

	if input.Name != "" {
		proses.Name = input.Name
	}
	if input.Email != "" {
		proses.Email = input.Email
	}

	if input.Avatar != "" {
		proses.Avatar = input.Avatar
	}

	if err := uq.db.Save(&proses).Error; err != nil {

		return user.User{}, err
	}
	result := user.User{
		ID:     proses.ID,
		Name:   proses.Name,
		Email:  proses.Email,
		Avatar: proses.Avatar,
	}

	return result, nil
}

// DeleteUser implements user.Repository.
func (uq *UserQuery) DeleteUser(userID uint) error {
	var exitingUser UserModel

	if err := uq.db.First(&exitingUser, userID).Error; err != nil {
		return err
	}

	if err := uq.db.Delete(&exitingUser).Error; err != nil {
		return err
	}

	return nil
}

// SearchUser implements user.Repository.
func (uq *UserQuery) SearchUser(userID uint, name string, page uint, limit uint) ([]user.User, uint, error) {
	var users []UserModel
	qry := uq.db.Table("user_models")

	if name != "" {
		qry = qry.Where("name like ?", "%"+name+"%")
		qry = qry.Where("deleted_at IS NULL")
	}
	var totalUser int64
	if err := qry.Count(&totalUser).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	qry = qry.Offset(int(offset)).Limit(int(limit))

	if err := qry.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	var result []user.User
	for _, s := range users {
		result = append(result, user.User{
			ID:     s.ID,
			Name:   s.Name,
			Email:  s.Email,
			Avatar: s.Avatar,
			Role:   s.Role,
		})
	}

	totalPages := int(totalUser) / int(limit)
	if totalUser%int64(limit) != 0 {
		totalPages++
	}

	return result, uint(totalPages), nil
}
