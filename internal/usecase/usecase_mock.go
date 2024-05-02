// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package usecase is a generated GoMock package.
package usecase

import (
	context "context"
	domain "homework/internal/domain"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockSensorRepository is a mock of SensorRepository interface.
type MockSensorRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSensorRepositoryMockRecorder
}

// MockSensorRepositoryMockRecorder is the mock recorder for MockSensorRepository.
type MockSensorRepositoryMockRecorder struct {
	mock *MockSensorRepository
}

// NewMockSensorRepository creates a new mock instance.
func NewMockSensorRepository(ctrl *gomock.Controller) *MockSensorRepository {
	mock := &MockSensorRepository{ctrl: ctrl}
	mock.recorder = &MockSensorRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSensorRepository) EXPECT() *MockSensorRepositoryMockRecorder {
	return m.recorder
}

// GetSensorByID mocks base method.
func (m *MockSensorRepository) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSensorByID", ctx, id)
	ret0, _ := ret[0].(*domain.Sensor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSensorByID indicates an expected call of GetSensorByID.
func (mr *MockSensorRepositoryMockRecorder) GetSensorByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSensorByID", reflect.TypeOf((*MockSensorRepository)(nil).GetSensorByID), ctx, id)
}

// GetSensorBySerialNumber mocks base method.
func (m *MockSensorRepository) GetSensorBySerialNumber(ctx context.Context, sn string) (*domain.Sensor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSensorBySerialNumber", ctx, sn)
	ret0, _ := ret[0].(*domain.Sensor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSensorBySerialNumber indicates an expected call of GetSensorBySerialNumber.
func (mr *MockSensorRepositoryMockRecorder) GetSensorBySerialNumber(ctx, sn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSensorBySerialNumber", reflect.TypeOf((*MockSensorRepository)(nil).GetSensorBySerialNumber), ctx, sn)
}

// GetSensors mocks base method.
func (m *MockSensorRepository) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSensors", ctx)
	ret0, _ := ret[0].([]domain.Sensor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSensors indicates an expected call of GetSensors.
func (mr *MockSensorRepositoryMockRecorder) GetSensors(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSensors", reflect.TypeOf((*MockSensorRepository)(nil).GetSensors), ctx)
}

// SaveSensor mocks base method.
func (m *MockSensorRepository) SaveSensor(ctx context.Context, sensor *domain.Sensor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveSensor", ctx, sensor)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveSensor indicates an expected call of SaveSensor.
func (mr *MockSensorRepositoryMockRecorder) SaveSensor(ctx, sensor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveSensor", reflect.TypeOf((*MockSensorRepository)(nil).SaveSensor), ctx, sensor)
}

// MockEventRepository is a mock of EventRepository interface.
type MockEventRepository struct {
	ctrl     *gomock.Controller
	recorder *MockEventRepositoryMockRecorder
}

// MockEventRepositoryMockRecorder is the mock recorder for MockEventRepository.
type MockEventRepositoryMockRecorder struct {
	mock *MockEventRepository
}

// NewMockEventRepository creates a new mock instance.
func NewMockEventRepository(ctrl *gomock.Controller) *MockEventRepository {
	mock := &MockEventRepository{ctrl: ctrl}
	mock.recorder = &MockEventRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventRepository) EXPECT() *MockEventRepositoryMockRecorder {
	return m.recorder
}

// GetLastEventBySensorID mocks base method.
func (m *MockEventRepository) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastEventBySensorID", ctx, id)
	ret0, _ := ret[0].(*domain.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastEventBySensorID indicates an expected call of GetLastEventBySensorID.
func (mr *MockEventRepositoryMockRecorder) GetLastEventBySensorID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastEventBySensorID", reflect.TypeOf((*MockEventRepository)(nil).GetLastEventBySensorID), ctx, id)
}

// GetSensorHistory mocks base method.
func (m *MockEventRepository) GetSensorHistory(ctx context.Context, id int64, start, end time.Time) ([]domain.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSensorHistory", ctx, id, start, end)
	ret0, _ := ret[0].([]domain.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSensorHistory indicates an expected call of GetSensorHistory.
func (mr *MockEventRepositoryMockRecorder) GetSensorHistory(ctx, id, start, end interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSensorHistory", reflect.TypeOf((*MockEventRepository)(nil).GetSensorHistory), ctx, id, start, end)
}

// SaveEvent mocks base method.
func (m *MockEventRepository) SaveEvent(ctx context.Context, event *domain.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveEvent", ctx, event)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveEvent indicates an expected call of SaveEvent.
func (mr *MockEventRepositoryMockRecorder) SaveEvent(ctx, event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveEvent", reflect.TypeOf((*MockEventRepository)(nil).SaveEvent), ctx, event)
}

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// GetUserByID mocks base method.
func (m *MockUserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, id)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockUserRepositoryMockRecorder) GetUserByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUserRepository)(nil).GetUserByID), ctx, id)
}

// SaveUser mocks base method.
func (m *MockUserRepository) SaveUser(ctx context.Context, user *domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveUser indicates an expected call of SaveUser.
func (mr *MockUserRepositoryMockRecorder) SaveUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUser", reflect.TypeOf((*MockUserRepository)(nil).SaveUser), ctx, user)
}

// MockSensorOwnerRepository is a mock of SensorOwnerRepository interface.
type MockSensorOwnerRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSensorOwnerRepositoryMockRecorder
}

// MockSensorOwnerRepositoryMockRecorder is the mock recorder for MockSensorOwnerRepository.
type MockSensorOwnerRepositoryMockRecorder struct {
	mock *MockSensorOwnerRepository
}

// NewMockSensorOwnerRepository creates a new mock instance.
func NewMockSensorOwnerRepository(ctrl *gomock.Controller) *MockSensorOwnerRepository {
	mock := &MockSensorOwnerRepository{ctrl: ctrl}
	mock.recorder = &MockSensorOwnerRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSensorOwnerRepository) EXPECT() *MockSensorOwnerRepositoryMockRecorder {
	return m.recorder
}

// GetSensorsByUserID mocks base method.
func (m *MockSensorOwnerRepository) GetSensorsByUserID(ctx context.Context, userID int64) ([]domain.SensorOwner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSensorsByUserID", ctx, userID)
	ret0, _ := ret[0].([]domain.SensorOwner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSensorsByUserID indicates an expected call of GetSensorsByUserID.
func (mr *MockSensorOwnerRepositoryMockRecorder) GetSensorsByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSensorsByUserID", reflect.TypeOf((*MockSensorOwnerRepository)(nil).GetSensorsByUserID), ctx, userID)
}

// SaveSensorOwner mocks base method.
func (m *MockSensorOwnerRepository) SaveSensorOwner(ctx context.Context, sensorOwner domain.SensorOwner) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveSensorOwner", ctx, sensorOwner)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveSensorOwner indicates an expected call of SaveSensorOwner.
func (mr *MockSensorOwnerRepositoryMockRecorder) SaveSensorOwner(ctx, sensorOwner interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveSensorOwner", reflect.TypeOf((*MockSensorOwnerRepository)(nil).SaveSensorOwner), ctx, sensorOwner)
}
