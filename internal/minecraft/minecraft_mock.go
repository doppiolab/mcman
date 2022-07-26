// Code generated by MockGen. DO NOT EDIT.
// Source: minecraft.go

// Package minecraft is a generated GoMock package.
package minecraft

import (
	exec "os/exec"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockServer is a mock of Server interface.
type MockServer struct {
	ctrl     *gomock.Controller
	recorder *MockServerMockRecorder
}

// MockServerMockRecorder is the mock recorder for MockServer.
type MockServerMockRecorder struct {
	mock *MockServer
}

// NewMockServer creates a new mock instance.
func NewMockServer(ctrl *gomock.Controller) *MockServer {
	mock := &MockServer{ctrl: ctrl}
	mock.recorder = &MockServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServer) EXPECT() *MockServerMockRecorder {
	return m.recorder
}

// GetProcess mocks base method.
func (m *MockServer) GetProcess() *exec.Cmd {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProcess")
	ret0, _ := ret[0].(*exec.Cmd)
	return ret0
}

// GetProcess indicates an expected call of GetProcess.
func (mr *MockServerMockRecorder) GetProcess() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProcess", reflect.TypeOf((*MockServer)(nil).GetProcess))
}

// PutCommand mocks base method.
func (m *MockServer) PutCommand(cmd string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutCommand", cmd)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutCommand indicates an expected call of PutCommand.
func (mr *MockServerMockRecorder) PutCommand(cmd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutCommand", reflect.TypeOf((*MockServer)(nil).PutCommand), cmd)
}

// Start mocks base method.
func (m *MockServer) Start() (chan string, chan string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(chan string)
	ret1, _ := ret[1].(chan string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Start indicates an expected call of Start.
func (mr *MockServerMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockServer)(nil).Start))
}

// Stop mocks base method.
func (m *MockServer) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockServerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockServer)(nil).Stop))
}
