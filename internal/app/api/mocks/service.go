// Code generated by MockGen. DO NOT EDIT.
// Source: ./service.go
//
// Generated by this command:
//
//	mockgen -source ./service.go -destination=./mocks/service.go -package=mock_service
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"
	models "rest-wallet/internal/app/models"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockWallet is a mock of Wallet interface.
type MockWallet struct {
	ctrl     *gomock.Controller
	recorder *MockWalletMockRecorder
}

// MockWalletMockRecorder is the mock recorder for MockWallet.
type MockWalletMockRecorder struct {
	mock *MockWallet
}

// NewMockWallet creates a new mock instance.
func NewMockWallet(ctrl *gomock.Controller) *MockWallet {
	mock := &MockWallet{ctrl: ctrl}
	mock.recorder = &MockWalletMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWallet) EXPECT() *MockWalletMockRecorder {
	return m.recorder
}

// GetWalletBalance mocks base method.
func (m *MockWallet) GetWalletBalance(walletID uuid.UUID) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWalletBalance", walletID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWalletBalance indicates an expected call of GetWalletBalance.
func (mr *MockWalletMockRecorder) GetWalletBalance(walletID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWalletBalance", reflect.TypeOf((*MockWallet)(nil).GetWalletBalance), walletID)
}

// UpdateWallet mocks base method.
func (m *MockWallet) UpdateWallet(request models.WalletUpdateRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWallet", request)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateWallet indicates an expected call of UpdateWallet.
func (mr *MockWalletMockRecorder) UpdateWallet(request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWallet", reflect.TypeOf((*MockWallet)(nil).UpdateWallet), request)
}
