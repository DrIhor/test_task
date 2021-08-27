// Code generated by mockery (devel). DO NOT EDIT.

package mocks

import (
	csv "encoding/csv"

	mock "github.com/stretchr/testify/mock"
)

// ItemCSVServ is an autogenerated mock type for the ItemCSVServ type
type ItemCSVServ struct {
	mock.Mock
}

// AddFromCSV provides a mock function with given fields: rd
func (_m *ItemCSVServ) AddFromCSV(rd *csv.Reader) ([]byte, error) {
	ret := _m.Called(rd)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(*csv.Reader) []byte); ok {
		r0 = rf(rd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*csv.Reader) error); ok {
		r1 = rf(rd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllItemsAsCSV provides a mock function with given fields:
func (_m *ItemCSVServ) GetAllItemsAsCSV() ([]byte, error) {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}