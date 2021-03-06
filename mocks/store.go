// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockLockStore is a mock of LockStore interface.
type MockLockStore struct {
	ctrl     *gomock.Controller
	recorder *MockLockStoreMockRecorder
}

// MockLockStoreMockRecorder is the mock recorder for MockLockStore.
type MockLockStoreMockRecorder struct {
	mock *MockLockStore
}

// NewMockLockStore creates a new mock instance.
func NewMockLockStore(ctrl *gomock.Controller) *MockLockStore {
	mock := &MockLockStore{ctrl: ctrl}
	mock.recorder = &MockLockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLockStore) EXPECT() *MockLockStoreMockRecorder {
	return m.recorder
}

// Set mocks base method.
func (m *MockLockStore) Set(ctx context.Context, key, uid string, expiry time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, uid, expiry)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockLockStoreMockRecorder) Set(ctx, key, uid, expiry interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockLockStore)(nil).Set), ctx, key, uid, expiry)
}

// Get mocks base method.
func (m *MockLockStore) Get(ctx context.Context, key string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(string)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockLockStoreMockRecorder) Get(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockLockStore)(nil).Get), ctx, key)
}

// Delete mocks base method.
func (m *MockLockStore) Delete(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockLockStoreMockRecorder) Delete(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockLockStore)(nil).Delete), ctx, key)
}
