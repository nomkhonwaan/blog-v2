// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nomkhonwaan/myblog/pkg/blog (interfaces: Service)

// Package blog is a generated GoMock package.
package blog

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Category mocks base method
func (m *MockService) Category() CategoryRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Category")
	ret0, _ := ret[0].(CategoryRepository)
	return ret0
}

// Category indicates an expected call of Category
func (mr *MockServiceMockRecorder) Category() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Category", reflect.TypeOf((*MockService)(nil).Category))
}

// Post mocks base method
func (m *MockService) Post() PostRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Post")
	ret0, _ := ret[0].(PostRepository)
	return ret0
}

// Post indicates an expected call of Post
func (mr *MockServiceMockRecorder) Post() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Post", reflect.TypeOf((*MockService)(nil).Post))
}
