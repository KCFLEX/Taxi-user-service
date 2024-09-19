// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/user.go

// Package mockrepo is a generated GoMock package.
package mockrepo

import (
	context "context"
	reflect "reflect"
	time "time"

	entity "github.com/KCFLEX/Taxi-user-service/internal/repository/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockToken is a mock of Token interface.
type MockToken struct {
	ctrl     *gomock.Controller
	recorder *MockTokenMockRecorder
}

// MockTokenMockRecorder is the mock recorder for MockToken.
type MockTokenMockRecorder struct {
	mock *MockToken
}

// NewMockToken creates a new mock instance.
func NewMockToken(ctrl *gomock.Controller) *MockToken {
	mock := &MockToken{ctrl: ctrl}
	mock.recorder = &MockTokenMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockToken) EXPECT() *MockTokenMockRecorder {
	return m.recorder
}

// GenerateToken mocks base method.
func (m *MockToken) GenerateToken(ctx context.Context, userID string, duration time.Duration) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", ctx, userID, duration)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockTokenMockRecorder) GenerateToken(ctx, userID, duration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockToken)(nil).GenerateToken), ctx, userID, duration)
}

// ParseToken mocks base method.
func (m *MockToken) ParseToken(ctx context.Context, tokenStr string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", ctx, tokenStr)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockTokenMockRecorder) ParseToken(ctx, tokenStr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockToken)(nil).ParseToken), ctx, tokenStr)
}

// ValidateToken mocks base method.
func (m *MockToken) ValidateToken(ctx context.Context, tokenString string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", ctx, tokenString)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateToken indicates an expected call of ValidateToken.
func (mr *MockTokenMockRecorder) ValidateToken(ctx, tokenString interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockToken)(nil).ValidateToken), ctx, tokenString)
}

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockRepository) CreateUser(ctx context.Context, user entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockRepositoryMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockRepository)(nil).CreateUser), ctx, user)
}

// DeleteProfileByID mocks base method.
func (m *MockRepository) DeleteProfileByID(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProfileByID", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProfileByID indicates an expected call of DeleteProfileByID.
func (mr *MockRepositoryMockRecorder) DeleteProfileByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProfileByID", reflect.TypeOf((*MockRepository)(nil).DeleteProfileByID), ctx, id)
}

// GetProfileByID mocks base method.
func (m *MockRepository) GetProfileByID(ctx context.Context, id int) (entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfileByID", ctx, id)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfileByID indicates an expected call of GetProfileByID.
func (mr *MockRepositoryMockRecorder) GetProfileByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfileByID", reflect.TypeOf((*MockRepository)(nil).GetProfileByID), ctx, id)
}

// StoreTokenInRedis mocks base method.
func (m *MockRepository) StoreTokenInRedis(ctx context.Context, userID, token string, expiration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreTokenInRedis", ctx, userID, token, expiration)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreTokenInRedis indicates an expected call of StoreTokenInRedis.
func (mr *MockRepositoryMockRecorder) StoreTokenInRedis(ctx, userID, token, expiration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreTokenInRedis", reflect.TypeOf((*MockRepository)(nil).StoreTokenInRedis), ctx, userID, token, expiration)
}

// UpdateProfileByID mocks base method.
func (m *MockRepository) UpdateProfileByID(ctx context.Context, updateInfo entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfileByID", ctx, updateInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProfileByID indicates an expected call of UpdateProfileByID.
func (mr *MockRepositoryMockRecorder) UpdateProfileByID(ctx, updateInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfileByID", reflect.TypeOf((*MockRepository)(nil).UpdateProfileByID), ctx, updateInfo)
}

// UserExists mocks base method.
func (m *MockRepository) UserExists(ctx context.Context, user entity.User) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserExists", ctx, user)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserExists indicates an expected call of UserExists.
func (mr *MockRepositoryMockRecorder) UserExists(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserExists", reflect.TypeOf((*MockRepository)(nil).UserExists), ctx, user)
}

// UserPhoneExists mocks base method.
func (m *MockRepository) UserPhoneExists(ctx context.Context, user entity.User) (entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserPhoneExists", ctx, user)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserPhoneExists indicates an expected call of UserPhoneExists.
func (mr *MockRepositoryMockRecorder) UserPhoneExists(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserPhoneExists", reflect.TypeOf((*MockRepository)(nil).UserPhoneExists), ctx, user)
}

// ValidateTokenRedis mocks base method.
func (m *MockRepository) ValidateTokenRedis(ctx context.Context, token, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateTokenRedis", ctx, token, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateTokenRedis indicates an expected call of ValidateTokenRedis.
func (mr *MockRepositoryMockRecorder) ValidateTokenRedis(ctx, token, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateTokenRedis", reflect.TypeOf((*MockRepository)(nil).ValidateTokenRedis), ctx, token, userID)
}
