// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nomkhonwaan/myblog/pkg/blog (interfaces: TagRepository)

// Package blog is a generated GoMock package.
package blog

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	reflect "reflect"
)

// MockTagRepository is a mock of TagRepository interface
type MockTagRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTagRepositoryMockRecorder
}

// MockTagRepositoryMockRecorder is the mock recorder for MockTagRepository
type MockTagRepositoryMockRecorder struct {
	mock *MockTagRepository
}

// NewMockTagRepository creates a new mock instance
func NewMockTagRepository(ctrl *gomock.Controller) *MockTagRepository {
	mock := &MockTagRepository{ctrl: ctrl}
	mock.recorder = &MockTagRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTagRepository) EXPECT() *MockTagRepositoryMockRecorder {
	return m.recorder
}

// FindAllByIDs mocks base method
func (m *MockTagRepository) FindAllByIDs(arg0 context.Context, arg1 []primitive.ObjectID) ([]Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByIDs", arg0, arg1)
	ret0, _ := ret[0].([]Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllByIDs indicates an expected call of FindAllByIDs
func (mr *MockTagRepositoryMockRecorder) FindAllByIDs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByIDs", reflect.TypeOf((*MockTagRepository)(nil).FindAllByIDs), arg0, arg1)
}
