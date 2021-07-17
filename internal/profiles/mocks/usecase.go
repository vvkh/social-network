// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	entity "github.com/vvkh/social-network/internal/profiles/entity"
	reflect "reflect"
)

// MockUseCase is a mock of UseCase interface
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// CreateProfile mocks base method
func (m *MockUseCase) CreateProfile(ctx context.Context, firstName, lastName string, age uint8, location, sex, about string) (entity.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProfile", ctx, firstName, lastName, age, location, sex, about)
	ret0, _ := ret[0].(entity.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProfile indicates an expected call of CreateProfile
func (mr *MockUseCaseMockRecorder) CreateProfile(ctx, firstName, lastName, age, location, sex, about interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProfile", reflect.TypeOf((*MockUseCase)(nil).CreateProfile), ctx, firstName, lastName, age, location, sex, about)
}

// GetByID mocks base method
func (m *MockUseCase) GetByID(ctx context.Context, id ...uint64) ([]entity.Profile, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range id {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByID", varargs...)
	ret0, _ := ret[0].([]entity.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockUseCaseMockRecorder) GetByID(ctx interface{}, id ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, id...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockUseCase)(nil).GetByID), varargs...)
}
