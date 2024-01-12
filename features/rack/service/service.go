package service

import (
	"errors"
	"library_api/features/rack"
	"library_api/helper/jwt"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

type RackServices struct {
	repo rack.Repository
}

func New(r rack.Repository) rack.Service {
	return &RackServices{
		repo: r,
	}
}

// AddRack implements rack.Service.
func (rs *RackServices) AddRack(token *golangjwt.Token, newRack rack.Rack) (rack.Rack, error) {
	userId, rolesUser, err := jwt.ExtractToken(token)
	if err != nil {
		return rack.Rack{}, errors.New("Token Error")
	}
	if rolesUser != "admin" {
		return rack.Rack{}, errors.New("unauthorized access: admin role required")
	}

	result, err := rs.repo.AddRack(userId, newRack)
	if err != nil {
		return rack.Rack{}, errors.New("Inputan tidak boleh kosong")
	}

	return result, err
}
