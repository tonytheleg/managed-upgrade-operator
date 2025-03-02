// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/openshift/managed-upgrade-operator/pkg/configmanager (interfaces: ConfigManager)
//
// Generated by this command:
//
//	mockgen -destination=mocks/configmanager.go -package=mocks github.com/openshift/managed-upgrade-operator/pkg/configmanager ConfigManager
//
// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	configmanager "github.com/openshift/managed-upgrade-operator/pkg/configmanager"
	gomock "go.uber.org/mock/gomock"
)

// MockConfigManager is a mock of ConfigManager interface.
type MockConfigManager struct {
	ctrl     *gomock.Controller
	recorder *MockConfigManagerMockRecorder
}

// MockConfigManagerMockRecorder is the mock recorder for MockConfigManager.
type MockConfigManagerMockRecorder struct {
	mock *MockConfigManager
}

// NewMockConfigManager creates a new mock instance.
func NewMockConfigManager(ctrl *gomock.Controller) *MockConfigManager {
	mock := &MockConfigManager{ctrl: ctrl}
	mock.recorder = &MockConfigManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfigManager) EXPECT() *MockConfigManagerMockRecorder {
	return m.recorder
}

// Into mocks base method.
func (m *MockConfigManager) Into(arg0 configmanager.ConfigValidator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Into", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Into indicates an expected call of Into.
func (mr *MockConfigManagerMockRecorder) Into(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Into", reflect.TypeOf((*MockConfigManager)(nil).Into), arg0)
}
