// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nomkhonwaan/myblog/pkg/storage (interfaces: FileRepository)

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	storage "github.com/nomkhonwaan/myblog/pkg/storage"
	reflect "reflect"
)

// MockFileRepository is a mock of FileRepository interface
type MockFileRepository struct {
	ctrl     *gomock.Controller
	recorder *MockFileRepositoryMockRecorder
}

// MockFileRepositoryMockRecorder is the mock recorder for MockFileRepository
type MockFileRepositoryMockRecorder struct {
	mock *MockFileRepository
}

// NewMockFileRepository creates a new mock instance
func NewMockFileRepository(ctrl *gomock.Controller) *MockFileRepository {
	mock := &MockFileRepository{ctrl: ctrl}
	mock.recorder = &MockFileRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFileRepository) EXPECT() *MockFileRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockFileRepository) Create(arg0 context.Context, arg1 storage.File) (storage.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(storage.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockFileRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockFileRepository)(nil).Create), arg0, arg1)
}

// FindByID mocks base method
func (m *MockFileRepository) FindByID(arg0 context.Context, arg1 interface{}) (storage.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", arg0, arg1)
	ret0, _ := ret[0].(storage.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID
func (mr *MockFileRepositoryMockRecorder) FindByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockFileRepository)(nil).FindByID), arg0, arg1)
}

// FindByPath mocks base method
func (m *MockFileRepository) FindByPath(arg0 context.Context, arg1 string) (storage.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByPath", arg0, arg1)
	ret0, _ := ret[0].(storage.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByPath indicates an expected call of FindByPath
func (mr *MockFileRepositoryMockRecorder) FindByPath(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByPath", reflect.TypeOf((*MockFileRepository)(nil).FindByPath), arg0, arg1)
}

// IsErrorRecordNotFound mocks base method
func (m *MockFileRepository) IsErrorRecordNotFound(arg0 error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsErrorRecordNotFound", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsErrorRecordNotFound indicates an expected call of IsErrorRecordNotFound
func (mr *MockFileRepositoryMockRecorder) IsErrorRecordNotFound(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsErrorRecordNotFound", reflect.TypeOf((*MockFileRepository)(nil).IsErrorRecordNotFound), arg0)
}
