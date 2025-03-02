// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/openshift/managed-upgrade-operator/pkg/machinery (interfaces: Machinery)
//
// Generated by this command:
//
//	mockgen -destination=mocks/machinery.go -package=mocks github.com/openshift/managed-upgrade-operator/pkg/machinery Machinery
//
// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	machinery "github.com/openshift/managed-upgrade-operator/pkg/machinery"
	gomock "go.uber.org/mock/gomock"
	v1 "k8s.io/api/core/v1"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

// MockMachinery is a mock of Machinery interface.
type MockMachinery struct {
	ctrl     *gomock.Controller
	recorder *MockMachineryMockRecorder
}

// MockMachineryMockRecorder is the mock recorder for MockMachinery.
type MockMachineryMockRecorder struct {
	mock *MockMachinery
}

// NewMockMachinery creates a new mock instance.
func NewMockMachinery(ctrl *gomock.Controller) *MockMachinery {
	mock := &MockMachinery{ctrl: ctrl}
	mock.recorder = &MockMachineryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMachinery) EXPECT() *MockMachineryMockRecorder {
	return m.recorder
}

// IsNodeCordoned mocks base method.
func (m *MockMachinery) IsNodeCordoned(arg0 *v1.Node) *machinery.IsCordonedResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNodeCordoned", arg0)
	ret0, _ := ret[0].(*machinery.IsCordonedResult)
	return ret0
}

// IsNodeCordoned indicates an expected call of IsNodeCordoned.
func (mr *MockMachineryMockRecorder) IsNodeCordoned(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNodeCordoned", reflect.TypeOf((*MockMachinery)(nil).IsNodeCordoned), arg0)
}

// IsUpgrading mocks base method.
func (m *MockMachinery) IsUpgrading(arg0 client.Client, arg1 string) (*machinery.UpgradingResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsUpgrading", arg0, arg1)
	ret0, _ := ret[0].(*machinery.UpgradingResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsUpgrading indicates an expected call of IsUpgrading.
func (mr *MockMachineryMockRecorder) IsUpgrading(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsUpgrading", reflect.TypeOf((*MockMachinery)(nil).IsUpgrading), arg0, arg1)
}
