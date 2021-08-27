// Code generated by mockery (devel). DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/DrIhor/test_task/internal/service/transport/graphQL/graph/model"
)

// QueryResolver is an autogenerated mock type for the QueryResolver type
type QueryResolver struct {
	mock.Mock
}

// GetItem provides a mock function with given fields: ctx, id
func (_m *QueryResolver) GetItem(ctx context.Context, id string) (*model.Item, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.Item
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Item); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Item)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetItems provides a mock function with given fields: ctx
func (_m *QueryResolver) GetItems(ctx context.Context) ([]*model.Item, error) {
	ret := _m.Called(ctx)

	var r0 []*model.Item
	if rf, ok := ret.Get(0).(func(context.Context) []*model.Item); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Item)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
