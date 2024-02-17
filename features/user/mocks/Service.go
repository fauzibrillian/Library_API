// Code generated by mockery v2.39.2. DO NOT EDIT.

package mocks

import (
	jwt "github.com/golang-jwt/jwt/v5"
	mock "github.com/stretchr/testify/mock"

	user "library_api/features/user"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// DeleteUser provides a mock function with given fields: token, userID
func (_m *Service) DeleteUser(token *jwt.Token, userID uint) error {
	ret := _m.Called(token, userID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*jwt.Token, uint) error); ok {
		r0 = rf(token, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Login provides a mock function with given fields: email, password
func (_m *Service) Login(email string, password string) (user.User, error) {
	ret := _m.Called(email, password)

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (user.User, error)); ok {
		return rf(email, password)
	}
	if rf, ok := ret.Get(0).(func(string, string) user.User); ok {
		r0 = rf(email, password)
	} else {
		r0 = ret.Get(0).(user.User)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(email, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: newUser
func (_m *Service) Register(newUser user.User) (user.User, error) {
	ret := _m.Called(newUser)

	if len(ret) == 0 {
		panic("no return value specified for Register")
	}

	var r0 user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(user.User) (user.User, error)); ok {
		return rf(newUser)
	}
	if rf, ok := ret.Get(0).(func(user.User) user.User); ok {
		r0 = rf(newUser)
	} else {
		r0 = ret.Get(0).(user.User)
	}

	if rf, ok := ret.Get(1).(func(user.User) error); ok {
		r1 = rf(newUser)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResetPassword provides a mock function with given fields: token, input
func (_m *Service) ResetPassword(token *jwt.Token, input user.User) (user.User, error) {
	ret := _m.Called(token, input)

	if len(ret) == 0 {
		panic("no return value specified for ResetPassword")
	}

	var r0 user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(*jwt.Token, user.User) (user.User, error)); ok {
		return rf(token, input)
	}
	if rf, ok := ret.Get(0).(func(*jwt.Token, user.User) user.User); ok {
		r0 = rf(token, input)
	} else {
		r0 = ret.Get(0).(user.User)
	}

	if rf, ok := ret.Get(1).(func(*jwt.Token, user.User) error); ok {
		r1 = rf(token, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchUser provides a mock function with given fields: token, name, page, limit
func (_m *Service) SearchUser(token *jwt.Token, name string, page uint, limit uint) ([]user.User, uint, error) {
	ret := _m.Called(token, name, page, limit)

	if len(ret) == 0 {
		panic("no return value specified for SearchUser")
	}

	var r0 []user.User
	var r1 uint
	var r2 error
	if rf, ok := ret.Get(0).(func(*jwt.Token, string, uint, uint) ([]user.User, uint, error)); ok {
		return rf(token, name, page, limit)
	}
	if rf, ok := ret.Get(0).(func(*jwt.Token, string, uint, uint) []user.User); ok {
		r0 = rf(token, name, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(*jwt.Token, string, uint, uint) uint); ok {
		r1 = rf(token, name, page, limit)
	} else {
		r1 = ret.Get(1).(uint)
	}

	if rf, ok := ret.Get(2).(func(*jwt.Token, string, uint, uint) error); ok {
		r2 = rf(token, name, page, limit)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// UpdateUser provides a mock function with given fields: token, input
func (_m *Service) UpdateUser(token *jwt.Token, input user.User) (user.User, error) {
	ret := _m.Called(token, input)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUser")
	}

	var r0 user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(*jwt.Token, user.User) (user.User, error)); ok {
		return rf(token, input)
	}
	if rf, ok := ret.Get(0).(func(*jwt.Token, user.User) user.User); ok {
		r0 = rf(token, input)
	} else {
		r0 = ret.Get(0).(user.User)
	}

	if rf, ok := ret.Get(1).(func(*jwt.Token, user.User) error); ok {
		r1 = rf(token, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
