// SPDX-License-Identifier: Apache-2.0
// Code generated by MockGen. DO NOT EDIT.
// Source: cmd/npm/os_context.go

// Package test is a generated GoMock package.
package test

import (
	bytes "bytes"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockOSContext is a mock of OSContext interface
type MockOSContext struct {
	ctrl     *gomock.Controller
	recorder *MockOSContextMockRecorder
}

// MockOSContextMockRecorder is the mock recorder for MockOSContext
type MockOSContextMockRecorder struct {
	mock *MockOSContext
}

// NewMockOSContext creates a new mock instance
func NewMockOSContext(ctrl *gomock.Controller) *MockOSContext {
	mock := &MockOSContext{ctrl: ctrl}
	mock.recorder = &MockOSContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOSContext) EXPECT() *MockOSContextMockRecorder {
	return m.recorder
}

// RunCommand mocks base method
func (m *MockOSContext) RunCommand(name string, args ...string) (bytes.Buffer, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{name}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RunCommand", varargs...)
	ret0, _ := ret[0].(bytes.Buffer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunCommand indicates an expected call of RunCommand
func (mr *MockOSContextMockRecorder) RunCommand(name interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{name}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunCommand", reflect.TypeOf((*MockOSContext)(nil).RunCommand), varargs...)
}

// RunCommandBytes mocks base method
func (m *MockOSContext) RunCommandBytes(name string, args ...string) ([]byte, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{name}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RunCommandBytes", varargs...)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunCommandBytes indicates an expected call of RunCommandBytes
func (mr *MockOSContextMockRecorder) RunCommandBytes(name interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{name}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunCommandBytes", reflect.TypeOf((*MockOSContext)(nil).RunCommandBytes), varargs...)
}

// RunCommandString mocks base method
func (m *MockOSContext) RunCommandString(name string, args ...string) (string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{name}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RunCommandString", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunCommandString indicates an expected call of RunCommandString
func (mr *MockOSContextMockRecorder) RunCommandString(name interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{name}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunCommandString", reflect.TypeOf((*MockOSContext)(nil).RunCommandString), varargs...)
}

// GetHomeDir mocks base method
func (m *MockOSContext) GetHomeDir() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHomeDir")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHomeDir indicates an expected call of GetHomeDir
func (mr *MockOSContextMockRecorder) GetHomeDir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHomeDir", reflect.TypeOf((*MockOSContext)(nil).GetHomeDir))
}
