// Code generated by MockGen. DO NOT EDIT.
// Source: ../usecase.go

// Package mock is a generated GoMock package.
package mock

import (
	models "2019_2_Covenant/internal/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockUsecase is a mock of Usecase interface
type MockUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUsecaseMockRecorder
}

// MockUsecaseMockRecorder is the mock recorder for MockUsecase
type MockUsecaseMockRecorder struct {
	mock *MockUsecase
}

// NewMockUsecase creates a new mock instance
func NewMockUsecase(ctrl *gomock.Controller) *MockUsecase {
	mock := &MockUsecase{ctrl: ctrl}
	mock.recorder = &MockUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUsecase) EXPECT() *MockUsecaseMockRecorder {
	return m.recorder
}

// Fetch mocks base method
func (m *MockUsecase) Fetch(count uint64) ([]*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", count)
	ret0, _ := ret[0].([]*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch
func (mr *MockUsecaseMockRecorder) Fetch(count interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockUsecase)(nil).Fetch), count)
}

// GetByID mocks base method
func (m *MockUsecase) GetByID(id uint64) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", id)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockUsecaseMockRecorder) GetByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockUsecase)(nil).GetByID), id)
}

// GetByEmail mocks base method
func (m *MockUsecase) GetByEmail(email string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", email)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail
func (mr *MockUsecaseMockRecorder) GetByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockUsecase)(nil).GetByEmail), email)
}

// GetByNickname mocks base method
func (m *MockUsecase) GetByNickname(nickname string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNickname", nickname)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByNickname indicates an expected call of GetByNickname
func (mr *MockUsecaseMockRecorder) GetByNickname(nickname interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNickname", reflect.TypeOf((*MockUsecase)(nil).GetByNickname), nickname)
}

// Store mocks base method
func (m *MockUsecase) Store(user *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store
func (mr *MockUsecaseMockRecorder) Store(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockUsecase)(nil).Store), user)
}

// UpdateAvatar mocks base method
func (m *MockUsecase) UpdateAvatar(id uint64, avatarPath string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAvatar", id, avatarPath)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAvatar indicates an expected call of UpdateAvatar
func (mr *MockUsecaseMockRecorder) UpdateAvatar(id, avatarPath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAvatar", reflect.TypeOf((*MockUsecase)(nil).UpdateAvatar), id, avatarPath)
}

// Update mocks base method
func (m *MockUsecase) Update(id uint64, nickname, email string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", id, nickname, email)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockUsecaseMockRecorder) Update(id, nickname, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUsecase)(nil).Update), id, nickname, email)
}

// UpdatePassword mocks base method
func (m *MockUsecase) UpdatePassword(id uint64, plainPassword string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePassword", id, plainPassword)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePassword indicates an expected call of UpdatePassword
func (mr *MockUsecaseMockRecorder) UpdatePassword(id, plainPassword interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePassword", reflect.TypeOf((*MockUsecase)(nil).UpdatePassword), id, plainPassword)
}
