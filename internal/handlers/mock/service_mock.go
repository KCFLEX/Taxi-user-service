// Code generated by MockGen. DO NOT EDIT.
// Source: internal/handlers/user.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/KCFLEX/Taxi-user-service/internal/handlers/models"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CheckTokenInRedis mocks base method.
func (m *MockService) CheckTokenInRedis(ctx context.Context, token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckTokenInRedis", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckTokenInRedis indicates an expected call of CheckTokenInRedis.
func (mr *MockServiceMockRecorder) CheckTokenInRedis(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckTokenInRedis", reflect.TypeOf((*MockService)(nil).CheckTokenInRedis), ctx, token)
}

// DeleteUserProfile mocks base method.
func (m *MockService) DeleteUserProfile(ctx context.Context, userID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserProfile", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserProfile indicates an expected call of DeleteUserProfile.
func (mr *MockServiceMockRecorder) DeleteUserProfile(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserProfile", reflect.TypeOf((*MockService)(nil).DeleteUserProfile), ctx, userID)
}

// GetUserProfile mocks base method.
func (m *MockService) GetUserProfile(ctx context.Context, userID int) (models.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserProfile", ctx, userID)
	ret0, _ := ret[0].(models.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserProfile indicates an expected call of GetUserProfile.
func (mr *MockServiceMockRecorder) GetUserProfile(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserProfile", reflect.TypeOf((*MockService)(nil).GetUserProfile), ctx, userID)
}

// SignIN mocks base method.
func (m *MockService) SignIN(ctx context.Context, user models.UserInfo) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIN", ctx, user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIN indicates an expected call of SignIN.
func (mr *MockServiceMockRecorder) SignIN(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIN", reflect.TypeOf((*MockService)(nil).SignIN), ctx, user)
}

// SignUP mocks base method.
func (m *MockService) SignUP(ctx context.Context, User models.UserInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUP", ctx, User)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignUP indicates an expected call of SignUP.
func (mr *MockServiceMockRecorder) SignUP(ctx, User interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUP", reflect.TypeOf((*MockService)(nil).SignUP), ctx, User)
}

// UpdateUserProfile mocks base method.
func (m *MockService) UpdateUserProfile(ctx context.Context, userID int, updateInfo models.UserInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserProfile", ctx, userID, updateInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserProfile indicates an expected call of UpdateUserProfile.
func (mr *MockServiceMockRecorder) UpdateUserProfile(ctx, userID, updateInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserProfile", reflect.TypeOf((*MockService)(nil).UpdateUserProfile), ctx, userID, updateInfo)
}

// VerifyToken mocks base method.
func (m *MockService) VerifyToken(ctx context.Context, token string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyToken", ctx, token)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyToken indicates an expected call of VerifyToken.
func (mr *MockServiceMockRecorder) VerifyToken(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyToken", reflect.TypeOf((*MockService)(nil).VerifyToken), ctx, token)
}
