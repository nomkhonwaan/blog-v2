// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nomkhonwaan/myblog/pkg/blog (interfaces: PostRepository)

// Package mock_blog is a generated GoMock package.
package mock_blog

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	blog "github.com/nomkhonwaan/myblog/pkg/blog"
	reflect "reflect"
)

// MockPostRepository is a mock of PostRepository interface
type MockPostRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPostRepositoryMockRecorder
}

// MockPostRepositoryMockRecorder is the mock recorder for MockPostRepository
type MockPostRepositoryMockRecorder struct {
	mock *MockPostRepository
}

// NewMockPostRepository creates a new mock instance
func NewMockPostRepository(ctrl *gomock.Controller) *MockPostRepository {
	mock := &MockPostRepository{ctrl: ctrl}
	mock.recorder = &MockPostRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPostRepository) EXPECT() *MockPostRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockPostRepository) Create(arg0 context.Context, arg1 string) (blog.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(blog.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockPostRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPostRepository)(nil).Create), arg0, arg1)
}

// FindAll mocks base method
func (m *MockPostRepository) FindAll(arg0 context.Context, arg1 blog.PostQuery) ([]blog.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", arg0, arg1)
	ret0, _ := ret[0].([]blog.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll
func (mr *MockPostRepositoryMockRecorder) FindAll(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockPostRepository)(nil).FindAll), arg0, arg1)
}

// FindByID mocks base method
func (m *MockPostRepository) FindByID(arg0 context.Context, arg1 interface{}) (blog.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", arg0, arg1)
	ret0, _ := ret[0].(blog.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID
func (mr *MockPostRepositoryMockRecorder) FindByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockPostRepository)(nil).FindByID), arg0, arg1)
}

// Save mocks base method
func (m *MockPostRepository) Save(arg0 context.Context, arg1 interface{}, arg2 blog.PostQuery) (blog.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1, arg2)
	ret0, _ := ret[0].(blog.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save
func (mr *MockPostRepositoryMockRecorder) Save(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockPostRepository)(nil).Save), arg0, arg1, arg2)
}
