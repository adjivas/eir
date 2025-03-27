package sbi

import (
	reflect "reflect"

	context "github.com/adjivas/eir/internal/context"
	processor "github.com/adjivas/eir/internal/sbi/processor"
	factory "github.com/adjivas/eir/pkg/factory"
	gomock "github.com/golang/mock/gomock"
)

// MockEIR is a mock of EIR interface.
type MockEIR struct {
	ctrl     *gomock.Controller
	recorder *MockEIRMockRecorder
}

// MockEIRMockRecorder is the mock recorder for MockEIR.
type MockEIRMockRecorder struct {
	mock *MockEIR
}

// NewMockEIR creates a new mock instance.
func NewMockEIR(ctrl *gomock.Controller) *MockEIR {
	mock := &MockEIR{ctrl: ctrl}
	mock.recorder = &MockEIRMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEIR) EXPECT() *MockEIRMockRecorder {
	return m.recorder
}

// Config mocks base method.
func (m *MockEIR) Config() *factory.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*factory.Config)
	return ret0
}

// Config indicates an expected call of Config.
func (mr *MockEIRMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockEIR)(nil).Config))
}

// Context mocks base method.
func (m *MockEIR) Context() *context.EIRContext {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(*context.EIRContext)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockEIRMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockEIR)(nil).Context))
}

// Processor mocks base method.
func (m *MockEIR) Processor() *processor.Processor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Processor")
	ret0, _ := ret[0].(*processor.Processor)
	return ret0
}

// Processor indicates an expected call of Processor.
func (mr *MockEIRMockRecorder) Processor() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Processor", reflect.TypeOf((*MockEIR)(nil).Processor))
}

// SetLogEnable mocks base method.
func (m *MockEIR) SetLogEnable(enable bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetLogEnable", enable)
}

// SetLogEnable indicates an expected call of SetLogEnable.
func (mr *MockEIRMockRecorder) SetLogEnable(enable interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()

	reflect := reflect.TypeOf((*MockEIR)(nil).SetLogEnable)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLogEnable", reflect, enable)
}

// SetLogLevel mocks base method.
func (m *MockEIR) SetLogLevel(level string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetLogLevel", level)
}

// SetLogLevel indicates an expected call of SetLogLevel.
func (mr *MockEIRMockRecorder) SetLogLevel(level interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	reflect := reflect.TypeOf((*MockEIR)(nil).SetLogLevel)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLogLevel", reflect, level)
}

// SetReportCaller mocks base method.
func (m *MockEIR) SetReportCaller(reportCaller bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetReportCaller", reportCaller)
}

// SetReportCaller indicates an expected call of SetReportCaller.
func (mr *MockEIRMockRecorder) SetReportCaller(reportCaller interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	reflect := reflect.TypeOf((*MockEIR)(nil).SetReportCaller)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetReportCaller", reflect, reportCaller)
}

// Start mocks base method.
func (m *MockEIR) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start.
func (mr *MockEIRMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockEIR)(nil).Start))
}

// Terminate mocks base method.
func (m *MockEIR) Terminate() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Terminate")
}

// Terminate indicates an expected call of Terminate.
func (mr *MockEIRMockRecorder) Terminate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Terminate", reflect.TypeOf((*MockEIR)(nil).Terminate))
}
