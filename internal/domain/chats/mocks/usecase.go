// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	entity "github.com/vvkh/social-network/internal/domain/chats/entity"
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

// GetUnreadMessagesCount mocks base method
func (m *MockUseCase) GetUnreadMessagesCount(ctx context.Context, profileID uint64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnreadMessagesCount", ctx, profileID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnreadMessagesCount indicates an expected call of GetUnreadMessagesCount
func (mr *MockUseCaseMockRecorder) GetUnreadMessagesCount(ctx, profileID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnreadMessagesCount", reflect.TypeOf((*MockUseCase)(nil).GetUnreadMessagesCount), ctx, profileID)
}

// GetOrCreateChat mocks base method
func (m *MockUseCase) GetOrCreateChat(ctx context.Context, oneProfileID uint64, otherProfileID int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrCreateChat", ctx, oneProfileID, otherProfileID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrCreateChat indicates an expected call of GetOrCreateChat
func (mr *MockUseCaseMockRecorder) GetOrCreateChat(ctx, oneProfileID, otherProfileID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrCreateChat", reflect.TypeOf((*MockUseCase)(nil).GetOrCreateChat), ctx, oneProfileID, otherProfileID)
}

// ListChatMessages mocks base method
func (m *MockUseCase) ListChatMessages(ctx context.Context, chatID uint64) ([]entity.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListChatMessages", ctx, chatID)
	ret0, _ := ret[0].([]entity.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListChatMessages indicates an expected call of ListChatMessages
func (mr *MockUseCaseMockRecorder) ListChatMessages(ctx, chatID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListChatMessages", reflect.TypeOf((*MockUseCase)(nil).ListChatMessages), ctx, chatID)
}

// SendMessage mocks base method
func (m *MockUseCase) SendMessage(ctx context.Context, chatID uint64, authorProfileID int64, message string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", ctx, chatID, authorProfileID, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage
func (mr *MockUseCaseMockRecorder) SendMessage(ctx, chatID, authorProfileID, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockUseCase)(nil).SendMessage), ctx, chatID, authorProfileID, message)
}
