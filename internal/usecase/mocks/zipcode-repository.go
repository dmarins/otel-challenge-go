// Code generated by MockGen. DO NOT EDIT.
// Source: internal/infrastructure/repositories/zipcode-repository.go
//
// Generated by this command:
//
//	mockgen -source=internal/infrastructure/repositories/zipcode-repository.go -destination=internal/usecase/mocks/zipcode-repository.go -typed=true -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	models "github.com/dmarins/otel-challenge-go/internal/domain/models"
	gomock "go.uber.org/mock/gomock"
)

// MockZipcodeRepositoryInterface is a mock of ZipcodeRepositoryInterface interface.
type MockZipcodeRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockZipcodeRepositoryInterfaceMockRecorder
}

// MockZipcodeRepositoryInterfaceMockRecorder is the mock recorder for MockZipcodeRepositoryInterface.
type MockZipcodeRepositoryInterfaceMockRecorder struct {
	mock *MockZipcodeRepositoryInterface
}

// NewMockZipcodeRepositoryInterface creates a new mock instance.
func NewMockZipcodeRepositoryInterface(ctrl *gomock.Controller) *MockZipcodeRepositoryInterface {
	mock := &MockZipcodeRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockZipcodeRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockZipcodeRepositoryInterface) EXPECT() *MockZipcodeRepositoryInterfaceMockRecorder {
	return m.recorder
}

// GetZipcodeInfo mocks base method.
func (m *MockZipcodeRepositoryInterface) GetZipcodeInfo(zipcode string) (*models.Zipcode, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetZipcodeInfo", zipcode)
	ret0, _ := ret[0].(*models.Zipcode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetZipcodeInfo indicates an expected call of GetZipcodeInfo.
func (mr *MockZipcodeRepositoryInterfaceMockRecorder) GetZipcodeInfo(zipcode any) *MockZipcodeRepositoryInterfaceGetZipcodeInfoCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetZipcodeInfo", reflect.TypeOf((*MockZipcodeRepositoryInterface)(nil).GetZipcodeInfo), zipcode)
	return &MockZipcodeRepositoryInterfaceGetZipcodeInfoCall{Call: call}
}

// MockZipcodeRepositoryInterfaceGetZipcodeInfoCall wrap *gomock.Call
type MockZipcodeRepositoryInterfaceGetZipcodeInfoCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockZipcodeRepositoryInterfaceGetZipcodeInfoCall) Return(arg0 *models.Zipcode, arg1 error) *MockZipcodeRepositoryInterfaceGetZipcodeInfoCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockZipcodeRepositoryInterfaceGetZipcodeInfoCall) Do(f func(string) (*models.Zipcode, error)) *MockZipcodeRepositoryInterfaceGetZipcodeInfoCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockZipcodeRepositoryInterfaceGetZipcodeInfoCall) DoAndReturn(f func(string) (*models.Zipcode, error)) *MockZipcodeRepositoryInterfaceGetZipcodeInfoCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}