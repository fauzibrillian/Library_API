// Code generated by mockery v2.39.2. DO NOT EDIT.

package mocks

import (
	echo "github.com/labstack/echo/v4"
	mock "github.com/stretchr/testify/mock"
)

// Handler is an autogenerated mock type for the Handler type
type Handler struct {
	mock.Mock
}

// Delete provides a mock function with given fields:
func (_m *Handler) Delete() echo.HandlerFunc {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 echo.HandlerFunc
	if rf, ok := ret.Get(0).(func() echo.HandlerFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(echo.HandlerFunc)
		}
	}

	return r0
}

// Login provides a mock function with given fields:
func (_m *Handler) Login() echo.HandlerFunc {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 echo.HandlerFunc
	if rf, ok := ret.Get(0).(func() echo.HandlerFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(echo.HandlerFunc)
		}
	}

	return r0
}

// Register provides a mock function with given fields:
func (_m *Handler) Register() echo.HandlerFunc {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Register")
	}

	var r0 echo.HandlerFunc
	if rf, ok := ret.Get(0).(func() echo.HandlerFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(echo.HandlerFunc)
		}
	}

	return r0
}

// ResetPassword provides a mock function with given fields:
func (_m *Handler) ResetPassword() echo.HandlerFunc {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ResetPassword")
	}

	var r0 echo.HandlerFunc
	if rf, ok := ret.Get(0).(func() echo.HandlerFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(echo.HandlerFunc)
		}
	}

	return r0
}

// SearchUser provides a mock function with given fields:
func (_m *Handler) SearchUser() echo.HandlerFunc {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for SearchUser")
	}

	var r0 echo.HandlerFunc
	if rf, ok := ret.Get(0).(func() echo.HandlerFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(echo.HandlerFunc)
		}
	}

	return r0
}

// UpdateUser provides a mock function with given fields:
func (_m *Handler) UpdateUser() echo.HandlerFunc {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for UpdateUser")
	}

	var r0 echo.HandlerFunc
	if rf, ok := ret.Get(0).(func() echo.HandlerFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(echo.HandlerFunc)
		}
	}

	return r0
}

// NewHandler creates a new instance of Handler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *Handler {
	mock := &Handler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
