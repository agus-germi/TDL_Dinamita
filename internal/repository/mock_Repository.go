// Code generated by mockery v2.46.3. DO NOT EDIT.

package repository

import (
	context "context"

	entity "github.com/agus-germi/TDL_Dinamita/internal/entity"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MockRepository is an autogenerated mock type for the Repository type
type MockRepository struct {
	mock.Mock
}

// GetReservation provides a mock function with given fields: ctx, userID, tableNumber, date
func (_m *MockRepository) GetReservation(ctx context.Context, userID int64, tableNumber int64, date time.Time) (*entity.Reservation, error) {
	ret := _m.Called(ctx, userID, tableNumber, date)

	if len(ret) == 0 {
		panic("no return value specified for GetReservation")
	}

	var r0 *entity.Reservation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, time.Time) (*entity.Reservation, error)); ok {
		return rf(ctx, userID, tableNumber, date)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, time.Time) *entity.Reservation); ok {
		r0 = rf(ctx, userID, tableNumber, date)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Reservation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64, time.Time) error); ok {
		r1 = rf(ctx, userID, tableNumber, date)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTableByNumber provides a mock function with given fields: ctx, tableNumber
func (_m *MockRepository) GetTableByNumber(ctx context.Context, tableNumber int64) (*entity.Table, error) {
	ret := _m.Called(ctx, tableNumber)

	if len(ret) == 0 {
		panic("no return value specified for GetTableByNumber")
	}

	var r0 *entity.Table
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*entity.Table, error)); ok {
		return rf(ctx, tableNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *entity.Table); ok {
		r0 = rf(ctx, tableNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Table)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, tableNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByEmail provides a mock function with given fields: ctx, email
func (_m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByEmail")
	}

	var r0 *entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*entity.User, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *entity.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserRole provides a mock function with given fields: ctx, userID
func (_m *MockRepository) GetUserRole(ctx context.Context, userID int64) (*entity.UserRole, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetUserRole")
	}

	var r0 *entity.UserRole
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*entity.UserRole, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *entity.UserRole); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.UserRole)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveReservation provides a mock function with given fields: ctx, userID, tableNumber, date
func (_m *MockRepository) RemoveReservation(ctx context.Context, userID int64, tableNumber int64, date time.Time) error {
	ret := _m.Called(ctx, userID, tableNumber, date)

	if len(ret) == 0 {
		panic("no return value specified for RemoveReservation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, time.Time) error); ok {
		r0 = rf(ctx, userID, tableNumber, date)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveTable provides a mock function with given fields: ctx, tableNumber
func (_m *MockRepository) RemoveTable(ctx context.Context, tableNumber int64) error {
	ret := _m.Called(ctx, tableNumber)

	if len(ret) == 0 {
		panic("no return value specified for RemoveTable")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, tableNumber)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveUser provides a mock function with given fields: ctx, email
func (_m *MockRepository) RemoveUser(ctx context.Context, email string) error {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for RemoveUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveUserRole provides a mock function with given fields: ctx, userID
func (_m *MockRepository) RemoveUserRole(ctx context.Context, userID int64) error {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for RemoveUserRole")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveReservation provides a mock function with given fields: ctx, userID, tableNumber, date
func (_m *MockRepository) SaveReservation(ctx context.Context, userID int64, tableNumber int64, date time.Time) error {
	ret := _m.Called(ctx, userID, tableNumber, date)

	if len(ret) == 0 {
		panic("no return value specified for SaveReservation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, time.Time) error); ok {
		r0 = rf(ctx, userID, tableNumber, date)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveTable provides a mock function with given fields: ctx, tableNumber, seats, location, isAvailable
func (_m *MockRepository) SaveTable(ctx context.Context, tableNumber int64, seats int64, location string, isAvailable bool) error {
	ret := _m.Called(ctx, tableNumber, seats, location, isAvailable)

	if len(ret) == 0 {
		panic("no return value specified for SaveTable")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, string, bool) error); ok {
		r0 = rf(ctx, tableNumber, seats, location, isAvailable)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveUser provides a mock function with given fields: ctx, name, passwd, email
func (_m *MockRepository) SaveUser(ctx context.Context, name string, passwd string, email string) error {
	ret := _m.Called(ctx, name, passwd, email)

	if len(ret) == 0 {
		panic("no return value specified for SaveUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, name, passwd, email)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveUserRole provides a mock function with given fields: ctx, userID, roleID
func (_m *MockRepository) SaveUserRole(ctx context.Context, userID int64, roleID int64) error {
	ret := _m.Called(ctx, userID, roleID)

	if len(ret) == 0 {
		panic("no return value specified for SaveUserRole")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) error); ok {
		r0 = rf(ctx, userID, roleID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockRepository creates a new instance of MockRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRepository {
	mock := &MockRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
