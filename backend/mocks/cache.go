// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/contract/cache.go
//
// Generated by this command:
//
//	mockgen -package mocks -source=internal/domain/contract/cache.go -destination=mocks/cache.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockCacheManager is a mock of CacheManager interface.
type MockCacheManager struct {
	ctrl     *gomock.Controller
	recorder *MockCacheManagerMockRecorder
	isgomock struct{}
}

// MockCacheManagerMockRecorder is the mock recorder for MockCacheManager.
type MockCacheManagerMockRecorder struct {
	mock *MockCacheManager
}

// NewMockCacheManager creates a new mock instance.
func NewMockCacheManager(ctrl *gomock.Controller) *MockCacheManager {
	mock := &MockCacheManager{ctrl: ctrl}
	mock.recorder = &MockCacheManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCacheManager) EXPECT() *MockCacheManagerMockRecorder {
	return m.recorder
}

// CleanAll mocks base method.
func (m *MockCacheManager) CleanAll(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CleanAll", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// CleanAll indicates an expected call of CleanAll.
func (mr *MockCacheManagerMockRecorder) CleanAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanAll", reflect.TypeOf((*MockCacheManager)(nil).CleanAll), ctx)
}

// Delete mocks base method.
func (m *MockCacheManager) Delete(ctx context.Context, keys ...string) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range keys {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCacheManagerMockRecorder) Delete(ctx any, keys ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, keys...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCacheManager)(nil).Delete), varargs...)
}

// GetAllKeys mocks base method.
func (m *MockCacheManager) GetAllKeys(ctx context.Context, pattern string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllKeys", ctx, pattern)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllKeys indicates an expected call of GetAllKeys.
func (mr *MockCacheManagerMockRecorder) GetAllKeys(ctx, pattern any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllKeys", reflect.TypeOf((*MockCacheManager)(nil).GetAllKeys), ctx, pattern)
}

// GetExpiration mocks base method.
func (m *MockCacheManager) GetExpiration(ctx context.Context, key string) (time.Duration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExpiration", ctx, key)
	ret0, _ := ret[0].(time.Duration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExpiration indicates an expected call of GetExpiration.
func (mr *MockCacheManagerMockRecorder) GetExpiration(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExpiration", reflect.TypeOf((*MockCacheManager)(nil).GetExpiration), ctx, key)
}

// GetInt mocks base method.
func (m *MockCacheManager) GetInt(ctx context.Context, key string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInt", ctx, key)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInt indicates an expected call of GetInt.
func (mr *MockCacheManagerMockRecorder) GetInt(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInt", reflect.TypeOf((*MockCacheManager)(nil).GetInt), ctx, key)
}

// GetItem mocks base method.
func (m *MockCacheManager) GetItem(ctx context.Context, key string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItem", ctx, key)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItem indicates an expected call of GetItem.
func (mr *MockCacheManagerMockRecorder) GetItem(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItem", reflect.TypeOf((*MockCacheManager)(nil).GetItem), ctx, key)
}

// GetString mocks base method.
func (m *MockCacheManager) GetString(ctx context.Context, key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetString", ctx, key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetString indicates an expected call of GetString.
func (mr *MockCacheManagerMockRecorder) GetString(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetString", reflect.TypeOf((*MockCacheManager)(nil).GetString), ctx, key)
}

// GetStruct mocks base method.
func (m *MockCacheManager) GetStruct(ctx context.Context, key string, data any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStruct", ctx, key, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetStruct indicates an expected call of GetStruct.
func (mr *MockCacheManagerMockRecorder) GetStruct(ctx, key, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStruct", reflect.TypeOf((*MockCacheManager)(nil).GetStruct), ctx, key, data)
}

// Increase mocks base method.
func (m *MockCacheManager) Increase(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Increase", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Increase indicates an expected call of Increase.
func (mr *MockCacheManagerMockRecorder) Increase(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Increase", reflect.TypeOf((*MockCacheManager)(nil).Increase), ctx, key)
}

// SetExpiration mocks base method.
func (m *MockCacheManager) SetExpiration(ctx context.Context, key string, expiration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetExpiration", ctx, key, expiration)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetExpiration indicates an expected call of SetExpiration.
func (mr *MockCacheManagerMockRecorder) SetExpiration(ctx, key, expiration any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetExpiration", reflect.TypeOf((*MockCacheManager)(nil).SetExpiration), ctx, key, expiration)
}

// SetItem mocks base method.
func (m *MockCacheManager) SetItem(ctx context.Context, key string, data []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetItem", ctx, key, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetItem indicates an expected call of SetItem.
func (mr *MockCacheManagerMockRecorder) SetItem(ctx, key, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetItem", reflect.TypeOf((*MockCacheManager)(nil).SetItem), ctx, key, data)
}

// SetItemWithExpiration mocks base method.
func (m *MockCacheManager) SetItemWithExpiration(ctx context.Context, key string, data []byte, expiration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetItemWithExpiration", ctx, key, data, expiration)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetItemWithExpiration indicates an expected call of SetItemWithExpiration.
func (mr *MockCacheManagerMockRecorder) SetItemWithExpiration(ctx, key, data, expiration any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetItemWithExpiration", reflect.TypeOf((*MockCacheManager)(nil).SetItemWithExpiration), ctx, key, data, expiration)
}

// SetString mocks base method.
func (m *MockCacheManager) SetString(ctx context.Context, key, data string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetString", ctx, key, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetString indicates an expected call of SetString.
func (mr *MockCacheManagerMockRecorder) SetString(ctx, key, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetString", reflect.TypeOf((*MockCacheManager)(nil).SetString), ctx, key, data)
}

// SetStringWithExpiration mocks base method.
func (m *MockCacheManager) SetStringWithExpiration(ctx context.Context, key, data string, expiration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetStringWithExpiration", ctx, key, data, expiration)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetStringWithExpiration indicates an expected call of SetStringWithExpiration.
func (mr *MockCacheManagerMockRecorder) SetStringWithExpiration(ctx, key, data, expiration any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStringWithExpiration", reflect.TypeOf((*MockCacheManager)(nil).SetStringWithExpiration), ctx, key, data, expiration)
}

// SetStruct mocks base method.
func (m *MockCacheManager) SetStruct(ctx context.Context, key string, data any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetStruct", ctx, key, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetStruct indicates an expected call of SetStruct.
func (mr *MockCacheManagerMockRecorder) SetStruct(ctx, key, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStruct", reflect.TypeOf((*MockCacheManager)(nil).SetStruct), ctx, key, data)
}

// SetStructWithExpiration mocks base method.
func (m *MockCacheManager) SetStructWithExpiration(ctx context.Context, key string, data any, expiration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetStructWithExpiration", ctx, key, data, expiration)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetStructWithExpiration indicates an expected call of SetStructWithExpiration.
func (mr *MockCacheManagerMockRecorder) SetStructWithExpiration(ctx, key, data, expiration any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStructWithExpiration", reflect.TypeOf((*MockCacheManager)(nil).SetStructWithExpiration), ctx, key, data, expiration)
}
