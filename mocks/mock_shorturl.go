// Code generated by MockGen. DO NOT EDIT.
// Source: C:\Users\Ivan\Documents\PRACTICUM\go-url-shortener\internal\handlers\shorturl\shorturl.go

// Package mock_shorturl is a generated GoMock package.
package mock_shorturl

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStorageURL is a mock of StorageURL interface.
type MockStorageURL struct {
	ctrl     *gomock.Controller
	recorder *MockStorageURLMockRecorder
}

// MockStorageURLMockRecorder is the mock recorder for MockStorageURL.
type MockStorageURLMockRecorder struct {
	mock *MockStorageURL
}

// NewMockStorageURL creates a new mock instance.
func NewMockStorageURL(ctrl *gomock.Controller) *MockStorageURL {
	mock := &MockStorageURL{ctrl: ctrl}
	mock.recorder = &MockStorageURLMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageURL) EXPECT() *MockStorageURLMockRecorder {
	return m.recorder
}

// AddURL mocks base method.
func (m *MockStorageURL) AddURL(url string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddURL", url)
	ret0, _ := ret[0].(string)
	return ret0
}

// AddURL indicates an expected call of AddURL.
func (mr *MockStorageURLMockRecorder) AddURL(url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddURL", reflect.TypeOf((*MockStorageURL)(nil).AddURL), url)
}

// Close mocks base method.
func (m *MockStorageURL) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStorageURLMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorageURL)(nil).Close))
}

// GetURL mocks base method.
func (m *MockStorageURL) GetURL() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURL")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetURL indicates an expected call of GetURL.
func (mr *MockStorageURLMockRecorder) GetURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURL", reflect.TypeOf((*MockStorageURL)(nil).GetURL))
}

// MockDatabase is a mock of Database interface.
type MockDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseMockRecorder
}

// MockDatabaseMockRecorder is the mock recorder for MockDatabase.
type MockDatabaseMockRecorder struct {
	mock *MockDatabase
}

// NewMockDatabase creates a new mock instance.
func NewMockDatabase(ctrl *gomock.Controller) *MockDatabase {
	mock := &MockDatabase{ctrl: ctrl}
	mock.recorder = &MockDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabase) EXPECT() *MockDatabaseMockRecorder {
	return m.recorder
}

// Connect mocks base method.
func (m *MockDatabase) Connect() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect")
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect.
func (mr *MockDatabaseMockRecorder) Connect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockDatabase)(nil).Connect))
}

// Disconnect mocks base method.
func (m *MockDatabase) Disconnect() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Disconnect")
	ret0, _ := ret[0].(error)
	return ret0
}

// Disconnect indicates an expected call of Disconnect.
func (mr *MockDatabaseMockRecorder) Disconnect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disconnect", reflect.TypeOf((*MockDatabase)(nil).Disconnect))
}

// PingContext mocks base method.
func (m *MockDatabase) PingContext() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PingContext")
	ret0, _ := ret[0].(error)
	return ret0
}

// PingContext indicates an expected call of PingContext.
func (mr *MockDatabaseMockRecorder) PingContext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PingContext", reflect.TypeOf((*MockDatabase)(nil).PingContext))
}
