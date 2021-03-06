// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	entity "github.com/vvkh/social-network/internal/domain/friendship/entity"
	reflect "reflect"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// ListFriends mocks base method
func (m *MockRepository) ListFriends(ctx context.Context, profileID uint64) ([]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFriends", ctx, profileID)
	ret0, _ := ret[0].([]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFriends indicates an expected call of ListFriends
func (mr *MockRepositoryMockRecorder) ListFriends(ctx, profileID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFriends", reflect.TypeOf((*MockRepository)(nil).ListFriends), ctx, profileID)
}

// ListPendingRequests mocks base method
func (m *MockRepository) ListPendingRequests(ctx context.Context, profileID uint64) ([]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPendingRequests", ctx, profileID)
	ret0, _ := ret[0].([]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPendingRequests indicates an expected call of ListPendingRequests
func (mr *MockRepositoryMockRecorder) ListPendingRequests(ctx, profileID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPendingRequests", reflect.TypeOf((*MockRepository)(nil).ListPendingRequests), ctx, profileID)
}

// GetFriendshipStatus mocks base method
func (m *MockRepository) GetFriendshipStatus(ctx context.Context, one, other uint64) (entity.FriendshipStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFriendshipStatus", ctx, one, other)
	ret0, _ := ret[0].(entity.FriendshipStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFriendshipStatus indicates an expected call of GetFriendshipStatus
func (mr *MockRepositoryMockRecorder) GetFriendshipStatus(ctx, one, other interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFriendshipStatus", reflect.TypeOf((*MockRepository)(nil).GetFriendshipStatus), ctx, one, other)
}

// CreateRequest mocks base method
func (m *MockRepository) CreateRequest(ctx context.Context, profileFromID, profileToID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRequest", ctx, profileFromID, profileToID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateRequest indicates an expected call of CreateRequest
func (mr *MockRepositoryMockRecorder) CreateRequest(ctx, profileFromID, profileToID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRequest", reflect.TypeOf((*MockRepository)(nil).CreateRequest), ctx, profileFromID, profileToID)
}

// AcceptRequest mocks base method
func (m *MockRepository) AcceptRequest(ctx context.Context, profileFromID, profileToID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcceptRequest", ctx, profileFromID, profileToID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AcceptRequest indicates an expected call of AcceptRequest
func (mr *MockRepositoryMockRecorder) AcceptRequest(ctx, profileFromID, profileToID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptRequest", reflect.TypeOf((*MockRepository)(nil).AcceptRequest), ctx, profileFromID, profileToID)
}

// DeclineRequest mocks base method
func (m *MockRepository) DeclineRequest(ctx context.Context, profileFromID, profileToID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeclineRequest", ctx, profileFromID, profileToID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeclineRequest indicates an expected call of DeclineRequest
func (mr *MockRepositoryMockRecorder) DeclineRequest(ctx, profileFromID, profileToID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeclineRequest", reflect.TypeOf((*MockRepository)(nil).DeclineRequest), ctx, profileFromID, profileToID)
}

// StopFriendship mocks base method
func (m *MockRepository) StopFriendship(ctx context.Context, profileID, otherProfileID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopFriendship", ctx, profileID, otherProfileID)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopFriendship indicates an expected call of StopFriendship
func (mr *MockRepositoryMockRecorder) StopFriendship(ctx, profileID, otherProfileID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopFriendship", reflect.TypeOf((*MockRepository)(nil).StopFriendship), ctx, profileID, otherProfileID)
}
