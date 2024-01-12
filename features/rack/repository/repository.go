package repository

import (
	"library_api/features/book/repository"
	"library_api/features/rack"

	"gorm.io/gorm"
)

type RackModel struct {
	gorm.Model
	Name        string
	RackDetails []repository.BookDetail `gorm:"foreignKey:RackID;"`
}

type RackQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) rack.Repository {
	return &RackQuery{
		db: db,
	}
}

// AddRack implements rack.Repository.
func (rq *RackQuery) AddRack(userID uint, input rack.Rack) (rack.Rack, error) {
	var inputDB = new(RackModel)
	inputDB.Name = input.Name

	rq.db.Create(&inputDB)
	input.ID = inputDB.ID
	return input, nil
}
