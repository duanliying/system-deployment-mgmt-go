/*******************************************************************************
 * Copyright 2017 Samsung Electronics All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *******************************************************************************/
// Code generated by MockGen. DO NOT EDIT.
// Source: MessengerInterface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockMessengerInterface is a mock of MessengerInterface interface
type MockMessengerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockMessengerInterfaceMockRecorder
}

// MockMessengerInterfaceMockRecorder is the mock recorder for MockMessengerInterface
type MockMessengerInterfaceMockRecorder struct {
	mock *MockMessengerInterface
}

// NewMockMessengerInterface creates a new mock instance
func NewMockMessengerInterface(ctrl *gomock.Controller) *MockMessengerInterface {
	mock := &MockMessengerInterface{ctrl: ctrl}
	mock.recorder = &MockMessengerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMessengerInterface) EXPECT() *MockMessengerInterfaceMockRecorder {
	return m.recorder
}

// DeployApp mocks base method
func (m *MockMessengerInterface) DeployApp(members []map[string]interface{}, data string) ([]int, []string) {
	ret := m.ctrl.Call(m, "DeployApp", members, data)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].([]string)
	return ret0, ret1
}

// DeployApp indicates an expected call of DeployApp
func (mr *MockMessengerInterfaceMockRecorder) DeployApp(members, data interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeployApp", reflect.TypeOf((*MockMessengerInterface)(nil).DeployApp), members, data)
}

// InfoApp mocks base method
func (m *MockMessengerInterface) InfoApp(members []map[string]interface{}, appId string) ([]int, []string) {
	ret := m.ctrl.Call(m, "InfoApp", members, appId)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].([]string)
	return ret0, ret1
}

// InfoApp indicates an expected call of InfoApp
func (mr *MockMessengerInterfaceMockRecorder) InfoApp(members, appId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InfoApp", reflect.TypeOf((*MockMessengerInterface)(nil).InfoApp), members, appId)
}

// DeleteApp mocks base method
func (m *MockMessengerInterface) DeleteApp(members []map[string]interface{}, appId string) ([]int, []string) {
	ret := m.ctrl.Call(m, "DeleteApp", members, appId)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].([]string)
	return ret0, ret1
}

// DeleteApp indicates an expected call of DeleteApp
func (mr *MockMessengerInterfaceMockRecorder) DeleteApp(members, appId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteApp", reflect.TypeOf((*MockMessengerInterface)(nil).DeleteApp), members, appId)
}

// StartApp mocks base method
func (m *MockMessengerInterface) StartApp(members []map[string]interface{}, appId string) ([]int, []string) {
	ret := m.ctrl.Call(m, "StartApp", members, appId)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].([]string)
	return ret0, ret1
}

// StartApp indicates an expected call of StartApp
func (mr *MockMessengerInterfaceMockRecorder) StartApp(members, appId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartApp", reflect.TypeOf((*MockMessengerInterface)(nil).StartApp), members, appId)
}

// StopApp mocks base method
func (m *MockMessengerInterface) StopApp(members []map[string]interface{}, appId string) ([]int, []string) {
	ret := m.ctrl.Call(m, "StopApp", members, appId)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].([]string)
	return ret0, ret1
}

// StopApp indicates an expected call of StopApp
func (mr *MockMessengerInterfaceMockRecorder) StopApp(members, appId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopApp", reflect.TypeOf((*MockMessengerInterface)(nil).StopApp), members, appId)
}

// UpdateApp mocks base method
func (m *MockMessengerInterface) UpdateApp(members []map[string]interface{}, appId string) ([]int, []string) {
	ret := m.ctrl.Call(m, "UpdateApp", members, appId)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].([]string)
	return ret0, ret1
}

// UpdateApp indicates an expected call of UpdateApp
func (mr *MockMessengerInterfaceMockRecorder) UpdateApp(members, appId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateApp", reflect.TypeOf((*MockMessengerInterface)(nil).UpdateApp), members, appId)
}

// InfoApps mocks base method
func (m *MockMessengerInterface) InfoApps(member []map[string]interface{}) ([]int, []string) {
	ret := m.ctrl.Call(m, "InfoApps", member)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].([]string)
	return ret0, ret1
}

// InfoApps indicates an expected call of InfoApps
func (mr *MockMessengerInterfaceMockRecorder) InfoApps(member interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InfoApps", reflect.TypeOf((*MockMessengerInterface)(nil).InfoApps), member)
}

// UpdateAppInfo mocks base method
func (m *MockMessengerInterface) UpdateAppInfo(member []map[string]interface{}, appId, data string) ([]int, []string) {
	ret := m.ctrl.Call(m, "UpdateAppInfo", member, appId, data)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].([]string)
	return ret0, ret1
}

// UpdateAppInfo indicates an expected call of UpdateAppInfo
func (mr *MockMessengerInterfaceMockRecorder) UpdateAppInfo(member, appId, data interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAppInfo", reflect.TypeOf((*MockMessengerInterface)(nil).UpdateAppInfo), member, appId, data)
}
