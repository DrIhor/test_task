// Code generated by mockery (devel). DO NOT EDIT.

package mocks

import (
	context "context"
	csv "encoding/csv"

	mock "github.com/stretchr/testify/mock"

	modelsitems "github.com/DrIhor/test_task/internal/models/items"
)

// ItemSrv is an autogenerated mock type for the ItemSrv type
type ItemSrv struct {
	mock.Mock
}

// AddFromCSV provides a mock function with given fields: _a0, _a1
func (_m *ItemSrv) AddFromCSV(_a0 context.Context, _a1 *csv.Reader) ([]byte, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, *csv.Reader) []byte); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *csv.Reader) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddNewItem provides a mock function with given fields: _a0, _a1
func (_m *ItemSrv) AddNewItem(_a0 context.Context, _a1 modelsitems.Item) (string, error) {
	ret := _m.Called(_a0, _a1)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, modelsitems.Item) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, modelsitems.Item) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteItem provides a mock function with given fields: _a0, _a1
func (_m *ItemSrv) DeleteItem(_a0 context.Context, _a1 string) (bool, error) {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllItems provides a mock function with given fields: _a0
func (_m *ItemSrv) GetAllItems(_a0 context.Context) ([]byte, error) {
	ret := _m.Called(_a0)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context) []byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllItemsAsCSV provides a mock function with given fields: _a0
func (_m *ItemSrv) GetAllItemsAsCSV(_a0 context.Context) ([]byte, error) {
	ret := _m.Called(_a0)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context) []byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetItem provides a mock function with given fields: _a0, _a1
func (_m *ItemSrv) GetItem(_a0 context.Context, _a1 string) ([]byte, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, string) []byte); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateItem provides a mock function with given fields: _a0, _a1
func (_m *ItemSrv) UpdateItem(_a0 context.Context, _a1 string) ([]byte, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, string) []byte); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
