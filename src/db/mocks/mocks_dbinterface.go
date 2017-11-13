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
// Source: DBInterface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	"db"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCommand is a mock of Command interface
type MockCommand struct {
	ctrl     *gomock.Controller
	recorder *MockCommandMockRecorder
}

// MockCommandMockRecorder is the mock recorder for MockCommand
type MockCommandMockRecorder struct {
	mock *MockCommand
}

// NewMockCommand creates a new mock instance
func NewMockCommand(ctrl *gomock.Controller) *MockCommand {
	mock := &MockCommand{ctrl: ctrl}
	mock.recorder = &MockCommandMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCommand) EXPECT() *MockCommandMockRecorder {
	return m.recorder
}

// AddAgent mocks base method
func (m *MockCommand) AddAgent(host, port, status string) (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "AddAgent", host, port, status)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddAgent indicates an expected call of AddAgent
func (mr *MockCommandMockRecorder) AddAgent(host, port, status interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAgent", reflect.TypeOf((*MockCommand)(nil).AddAgent), host, port, status)
}

// UpdateAgentAddress mocks base method
func (m *MockCommand) UpdateAgentAddress(agent_id, host, port string) error {
	ret := m.ctrl.Call(m, "UpdateAgentAddress", agent_id, host, port)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAgentAddress indicates an expected call of UpdateAgentAddress
func (mr *MockCommandMockRecorder) UpdateAgentAddress(agent_id, host, port interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAgentAddress", reflect.TypeOf((*MockCommand)(nil).UpdateAgentAddress), agent_id, host, port)
}

// UpdateAgentStatus mocks base method
func (m *MockCommand) UpdateAgentStatus(agent_id, status string) error {
	ret := m.ctrl.Call(m, "UpdateAgentStatus", agent_id, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAgentStatus indicates an expected call of UpdateAgentStatus
func (mr *MockCommandMockRecorder) UpdateAgentStatus(agent_id, status interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAgentStatus", reflect.TypeOf((*MockCommand)(nil).UpdateAgentStatus), agent_id, status)
}

// GetAgent mocks base method
func (m *MockCommand) GetAgent(agent_id string) (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetAgent", agent_id)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAgent indicates an expected call of GetAgent
func (mr *MockCommandMockRecorder) GetAgent(agent_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAgent", reflect.TypeOf((*MockCommand)(nil).GetAgent), agent_id)
}

// GetAllAgents mocks base method
func (m *MockCommand) GetAllAgents() ([]map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetAllAgents")
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllAgents indicates an expected call of GetAllAgents
func (mr *MockCommandMockRecorder) GetAllAgents() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllAgents", reflect.TypeOf((*MockCommand)(nil).GetAllAgents))
}

// GetAgentByAppID mocks base method
func (m *MockCommand) GetAgentByAppID(agent_id, app_id string) (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetAgentByAppID", agent_id, app_id)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAgentByAppID indicates an expected call of GetAgentByAppID
func (mr *MockCommandMockRecorder) GetAgentByAppID(agent_id, app_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAgentByAppID", reflect.TypeOf((*MockCommand)(nil).GetAgentByAppID), agent_id, app_id)
}

// AddAppToAgent mocks base method
func (m *MockCommand) AddAppToAgent(agent_id, app_id string) error {
	ret := m.ctrl.Call(m, "AddAppToAgent", agent_id, app_id)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddAppToAgent indicates an expected call of AddAppToAgent
func (mr *MockCommandMockRecorder) AddAppToAgent(agent_id, app_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAppToAgent", reflect.TypeOf((*MockCommand)(nil).AddAppToAgent), agent_id, app_id)
}

// DeleteAppFromAgent mocks base method
func (m *MockCommand) DeleteAppFromAgent(agent_id, app_id string) error {
	ret := m.ctrl.Call(m, "DeleteAppFromAgent", agent_id, app_id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAppFromAgent indicates an expected call of DeleteAppFromAgent
func (mr *MockCommandMockRecorder) DeleteAppFromAgent(agent_id, app_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAppFromAgent", reflect.TypeOf((*MockCommand)(nil).DeleteAppFromAgent), agent_id, app_id)
}

// DeleteAgent mocks base method
func (m *MockCommand) DeleteAgent(agent_id string) error {
	ret := m.ctrl.Call(m, "DeleteAgent", agent_id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAgent indicates an expected call of DeleteAgent
func (mr *MockCommandMockRecorder) DeleteAgent(agent_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAgent", reflect.TypeOf((*MockCommand)(nil).DeleteAgent), agent_id)
}

// CreateGroup mocks base method
func (m *MockCommand) CreateGroup() (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "CreateGroup")
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGroup indicates an expected call of CreateGroup
func (mr *MockCommandMockRecorder) CreateGroup() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroup", reflect.TypeOf((*MockCommand)(nil).CreateGroup))
}

// GetGroup mocks base method
func (m *MockCommand) GetGroup(group_id string) (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetGroup", group_id)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroup indicates an expected call of GetGroup
func (mr *MockCommandMockRecorder) GetGroup(group_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroup", reflect.TypeOf((*MockCommand)(nil).GetGroup), group_id)
}

// GetAllGroups mocks base method
func (m *MockCommand) GetAllGroups() ([]map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetAllGroups")
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllGroups indicates an expected call of GetAllGroups
func (mr *MockCommandMockRecorder) GetAllGroups() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllGroups", reflect.TypeOf((*MockCommand)(nil).GetAllGroups))
}

// GetGroupMembers mocks base method
func (m *MockCommand) GetGroupMembers(group_id string) ([]map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetGroupMembers", group_id)
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupMembers indicates an expected call of GetGroupMembers
func (mr *MockCommandMockRecorder) GetGroupMembers(group_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupMembers", reflect.TypeOf((*MockCommand)(nil).GetGroupMembers), group_id)
}

// GetGroupMembersByAppID mocks base method
func (m *MockCommand) GetGroupMembersByAppID(group_id, app_id string) ([]map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetGroupMembersByAppID", group_id, app_id)
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupMembersByAppID indicates an expected call of GetGroupMembersByAppID
func (mr *MockCommandMockRecorder) GetGroupMembersByAppID(group_id, app_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupMembersByAppID", reflect.TypeOf((*MockCommand)(nil).GetGroupMembersByAppID), group_id, app_id)
}

// JoinGroup mocks base method
func (m *MockCommand) JoinGroup(group_id, agent_id string) error {
	ret := m.ctrl.Call(m, "JoinGroup", group_id, agent_id)
	ret0, _ := ret[0].(error)
	return ret0
}

// JoinGroup indicates an expected call of JoinGroup
func (mr *MockCommandMockRecorder) JoinGroup(group_id, agent_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JoinGroup", reflect.TypeOf((*MockCommand)(nil).JoinGroup), group_id, agent_id)
}

// LeaveGroup mocks base method
func (m *MockCommand) LeaveGroup(group_id, agent_id string) error {
	ret := m.ctrl.Call(m, "LeaveGroup", group_id, agent_id)
	ret0, _ := ret[0].(error)
	return ret0
}

// LeaveGroup indicates an expected call of LeaveGroup
func (mr *MockCommandMockRecorder) LeaveGroup(group_id, agent_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeaveGroup", reflect.TypeOf((*MockCommand)(nil).LeaveGroup), group_id, agent_id)
}

// DeleteGroup mocks base method
func (m *MockCommand) DeleteGroup(group_id string) error {
	ret := m.ctrl.Call(m, "DeleteGroup", group_id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGroup indicates an expected call of DeleteGroup
func (mr *MockCommandMockRecorder) DeleteGroup(group_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGroup", reflect.TypeOf((*MockCommand)(nil).DeleteGroup), group_id)
}

// MockCloser is a mock of Closer interface
type MockCloser struct {
	ctrl     *gomock.Controller
	recorder *MockCloserMockRecorder
}

// MockCloserMockRecorder is the mock recorder for MockCloser
type MockCloserMockRecorder struct {
	mock *MockCloser
}

// NewMockCloser creates a new mock instance
func NewMockCloser(ctrl *gomock.Controller) *MockCloser {
	mock := &MockCloser{ctrl: ctrl}
	mock.recorder = &MockCloserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCloser) EXPECT() *MockCloserMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockCloser) Close() {
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockCloserMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockCloser)(nil).Close))
}

// MockDBManager is a mock of DBManager interface
type MockDBManager struct {
	ctrl     *gomock.Controller
	recorder *MockDBManagerMockRecorder
}

// MockDBManagerMockRecorder is the mock recorder for MockDBManager
type MockDBManagerMockRecorder struct {
	mock *MockDBManager
}

// NewMockDBManager creates a new mock instance
func NewMockDBManager(ctrl *gomock.Controller) *MockDBManager {
	mock := &MockDBManager{ctrl: ctrl}
	mock.recorder = &MockDBManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDBManager) EXPECT() *MockDBManagerMockRecorder {
	return m.recorder
}

// AddAgent mocks base method
func (m *MockDBManager) AddAgent(host, port, status string) (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "AddAgent", host, port, status)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddAgent indicates an expected call of AddAgent
func (mr *MockDBManagerMockRecorder) AddAgent(host, port, status interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAgent", reflect.TypeOf((*MockDBManager)(nil).AddAgent), host, port, status)
}

// UpdateAgentAddress mocks base method
func (m *MockDBManager) UpdateAgentAddress(agent_id, host, port string) error {
	ret := m.ctrl.Call(m, "UpdateAgentAddress", agent_id, host, port)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAgentAddress indicates an expected call of UpdateAgentAddress
func (mr *MockDBManagerMockRecorder) UpdateAgentAddress(agent_id, host, port interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAgentAddress", reflect.TypeOf((*MockDBManager)(nil).UpdateAgentAddress), agent_id, host, port)
}

// UpdateAgentStatus mocks base method
func (m *MockDBManager) UpdateAgentStatus(agent_id, status string) error {
	ret := m.ctrl.Call(m, "UpdateAgentStatus", agent_id, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAgentStatus indicates an expected call of UpdateAgentStatus
func (mr *MockDBManagerMockRecorder) UpdateAgentStatus(agent_id, status interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAgentStatus", reflect.TypeOf((*MockDBManager)(nil).UpdateAgentStatus), agent_id, status)
}

// GetAgent mocks base method
func (m *MockDBManager) GetAgent(agent_id string) (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetAgent", agent_id)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAgent indicates an expected call of GetAgent
func (mr *MockDBManagerMockRecorder) GetAgent(agent_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAgent", reflect.TypeOf((*MockDBManager)(nil).GetAgent), agent_id)
}

// GetAllAgents mocks base method
func (m *MockDBManager) GetAllAgents() ([]map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetAllAgents")
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllAgents indicates an expected call of GetAllAgents
func (mr *MockDBManagerMockRecorder) GetAllAgents() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllAgents", reflect.TypeOf((*MockDBManager)(nil).GetAllAgents))
}

// GetAgentByAppID mocks base method
func (m *MockDBManager) GetAgentByAppID(agent_id, app_id string) (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetAgentByAppID", agent_id, app_id)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAgentByAppID indicates an expected call of GetAgentByAppID
func (mr *MockDBManagerMockRecorder) GetAgentByAppID(agent_id, app_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAgentByAppID", reflect.TypeOf((*MockDBManager)(nil).GetAgentByAppID), agent_id, app_id)
}

// AddAppToAgent mocks base method
func (m *MockDBManager) AddAppToAgent(agent_id, app_id string) error {
	ret := m.ctrl.Call(m, "AddAppToAgent", agent_id, app_id)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddAppToAgent indicates an expected call of AddAppToAgent
func (mr *MockDBManagerMockRecorder) AddAppToAgent(agent_id, app_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAppToAgent", reflect.TypeOf((*MockDBManager)(nil).AddAppToAgent), agent_id, app_id)
}

// DeleteAppFromAgent mocks base method
func (m *MockDBManager) DeleteAppFromAgent(agent_id, app_id string) error {
	ret := m.ctrl.Call(m, "DeleteAppFromAgent", agent_id, app_id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAppFromAgent indicates an expected call of DeleteAppFromAgent
func (mr *MockDBManagerMockRecorder) DeleteAppFromAgent(agent_id, app_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAppFromAgent", reflect.TypeOf((*MockDBManager)(nil).DeleteAppFromAgent), agent_id, app_id)
}

// DeleteAgent mocks base method
func (m *MockDBManager) DeleteAgent(agent_id string) error {
	ret := m.ctrl.Call(m, "DeleteAgent", agent_id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAgent indicates an expected call of DeleteAgent
func (mr *MockDBManagerMockRecorder) DeleteAgent(agent_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAgent", reflect.TypeOf((*MockDBManager)(nil).DeleteAgent), agent_id)
}

// CreateGroup mocks base method
func (m *MockDBManager) CreateGroup() (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "CreateGroup")
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGroup indicates an expected call of CreateGroup
func (mr *MockDBManagerMockRecorder) CreateGroup() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroup", reflect.TypeOf((*MockDBManager)(nil).CreateGroup))
}

// GetGroup mocks base method
func (m *MockDBManager) GetGroup(group_id string) (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetGroup", group_id)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroup indicates an expected call of GetGroup
func (mr *MockDBManagerMockRecorder) GetGroup(group_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroup", reflect.TypeOf((*MockDBManager)(nil).GetGroup), group_id)
}

// GetAllGroups mocks base method
func (m *MockDBManager) GetAllGroups() ([]map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetAllGroups")
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllGroups indicates an expected call of GetAllGroups
func (mr *MockDBManagerMockRecorder) GetAllGroups() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllGroups", reflect.TypeOf((*MockDBManager)(nil).GetAllGroups))
}

// GetGroupMembers mocks base method
func (m *MockDBManager) GetGroupMembers(group_id string) ([]map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetGroupMembers", group_id)
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupMembers indicates an expected call of GetGroupMembers
func (mr *MockDBManagerMockRecorder) GetGroupMembers(group_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupMembers", reflect.TypeOf((*MockDBManager)(nil).GetGroupMembers), group_id)
}

// GetGroupMembersByAppID mocks base method
func (m *MockDBManager) GetGroupMembersByAppID(group_id, app_id string) ([]map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetGroupMembersByAppID", group_id, app_id)
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupMembersByAppID indicates an expected call of GetGroupMembersByAppID
func (mr *MockDBManagerMockRecorder) GetGroupMembersByAppID(group_id, app_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupMembersByAppID", reflect.TypeOf((*MockDBManager)(nil).GetGroupMembersByAppID), group_id, app_id)
}

// JoinGroup mocks base method
func (m *MockDBManager) JoinGroup(group_id, agent_id string) error {
	ret := m.ctrl.Call(m, "JoinGroup", group_id, agent_id)
	ret0, _ := ret[0].(error)
	return ret0
}

// JoinGroup indicates an expected call of JoinGroup
func (mr *MockDBManagerMockRecorder) JoinGroup(group_id, agent_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JoinGroup", reflect.TypeOf((*MockDBManager)(nil).JoinGroup), group_id, agent_id)
}

// LeaveGroup mocks base method
func (m *MockDBManager) LeaveGroup(group_id, agent_id string) error {
	ret := m.ctrl.Call(m, "LeaveGroup", group_id, agent_id)
	ret0, _ := ret[0].(error)
	return ret0
}

// LeaveGroup indicates an expected call of LeaveGroup
func (mr *MockDBManagerMockRecorder) LeaveGroup(group_id, agent_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeaveGroup", reflect.TypeOf((*MockDBManager)(nil).LeaveGroup), group_id, agent_id)
}

// DeleteGroup mocks base method
func (m *MockDBManager) DeleteGroup(group_id string) error {
	ret := m.ctrl.Call(m, "DeleteGroup", group_id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGroup indicates an expected call of DeleteGroup
func (mr *MockDBManagerMockRecorder) DeleteGroup(group_id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGroup", reflect.TypeOf((*MockDBManager)(nil).DeleteGroup), group_id)
}

// Close mocks base method
func (m *MockDBManager) Close() {
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockDBManagerMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDBManager)(nil).Close))
}

// MockDBConnection is a mock of DBConnection interface
type MockDBConnection struct {
	ctrl     *gomock.Controller
	recorder *MockDBConnectionMockRecorder
}

// MockDBConnectionMockRecorder is the mock recorder for MockDBConnection
type MockDBConnectionMockRecorder struct {
	mock *MockDBConnection
}

// NewMockDBConnection creates a new mock instance
func NewMockDBConnection(ctrl *gomock.Controller) *MockDBConnection {
	mock := &MockDBConnection{ctrl: ctrl}
	mock.recorder = &MockDBConnectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDBConnection) EXPECT() *MockDBConnectionMockRecorder {
	return m.recorder
}

// Connect mocks base method
func (m *MockDBConnection) Connect() (db.DBManager, error) {
	ret := m.ctrl.Call(m, "Connect")
	ret0, _ := ret[0].(db.DBManager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Connect indicates an expected call of Connect
func (mr *MockDBConnectionMockRecorder) Connect() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockDBConnection)(nil).Connect))
}

