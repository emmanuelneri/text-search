// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	users "api/pkg/users"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// Scroll provides a mock function with given fields: id
func (_m *Service) Scroll(id string) (*users.UserPaged, error) {
	ret := _m.Called(id)

	var r0 *users.UserPaged
	if rf, ok := ret.Get(0).(func(string) *users.UserPaged); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*users.UserPaged)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Search provides a mock function with given fields: ctx, query
func (_m *Service) Search(ctx context.Context, query string) (*users.UserPaged, error) {
	ret := _m.Called(ctx, query)

	var r0 *users.UserPaged
	if rf, ok := ret.Get(0).(func(context.Context, string) *users.UserPaged); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*users.UserPaged)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
