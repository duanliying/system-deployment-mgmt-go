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
package group

import (
	"commons/errors"
	"commons/results"
	dbmocks "db/mocks"
	msgmocks "messenger/mocks"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

const (
	appId   = "000000000000000000000000"
	agentId = "000000000000000000000001"
	groupId = "000000000000000000000002"
	host    = "192.168.0.1"
	port    = "8888"
)

var (
	agent = map[string]interface{}{
		"id":   agentId,
		"host": host,
		"port": port,
		"apps": []string{appId},
	}
	members = []map[string]interface{}{agent, agent}
	address = map[string]interface{}{
		"host": host,
		"port": port,
	}
	membersAddress = []map[string]interface{}{address, address}
	group          = map[string]interface{}{
		"id":      groupId,
		"members": []string{},
	}

	body                   = `{"description":"description"}`
	respCode               = []int{results.OK, results.OK}
	partialSuccessRespCode = []int{results.OK, results.ERROR}
	errorRespCode          = []int{results.ERROR, results.ERROR}
	invalidRespStr         = []string{`{"invalidJson"}`}
	notFoundError          = errors.NotFound{}
	connectionError        = errors.DBConnectionError{}
)

var controller GroupInterface

func init() {
	controller = GroupController{}
}

func TestCalledCreateGroup_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().CreateGroup().Return(group, nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, res, err := controller.CreateGroup()

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(group, res) {
		t.Errorf("Expected res: %s, actual res: %s", group, res)
	}
}

func TestCalledCreateGroupWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.CreateGroup()

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", err.Error())
	case errors.DBConnectionError:
	}
}

func TestCalledCreateGroupWhenFailedToInsertGroupToDB_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().CreateGroup().Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.CreateGroup()

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetGroup_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroup(groupId).Return(group, nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, res, err := controller.GetGroup(groupId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(group, res) {
		t.Errorf("Expected res: %s, actual res: %s", group, res)
	}
}

func TestCalledGetGroupWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.GetGroup(groupId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", err.Error())
	case errors.DBConnectionError:
	}
}

func TestCalledGetGroupWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroup(groupId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.GetGroup(groupId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetGroups_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groups := []map[string]interface{}{group}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAllGroups().Return(groups, nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, res, err := controller.GetGroups()

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(groups, res["groups"].([]map[string]interface{})) {
		t.Errorf("Expected res: %s, actual res: %s", groups, res["groups"].([]map[string]interface{}))
	}
}

func TestCalledGetGroupsWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.GetGroups()

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", err.Error())
	case errors.DBConnectionError:
	}
}

func TestCalledGetGroupsWhenFailedToGetGroupsFromDB_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAllGroups().Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.GetGroups()

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledJoinGroup_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().JoinGroup(groupId, agentId).Return(nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	agents := `{"agents":["000000000000000000000001"]}`
	code, _, err := controller.JoinGroup(groupId, agents)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledJoinGroupWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	agents := `{"agents":["000000000000000000000001"]}`
	code, _, err := controller.JoinGroup(groupId, agents)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", err.Error())
	case errors.DBConnectionError:
	}
}

func TestCalledJoinGroupWithInvalidRequestBody_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	invalidJsonStr := `{"invalidJson"}`
	code, _, err := controller.JoinGroup(groupId, invalidJsonStr)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InvalidParamError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InvalidParamError", err.Error())
	case errors.InvalidJSON:
	}
}

func TestCalledJoinGroupWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().JoinGroup(groupId, agentId).Return(notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	agents := `{"agents":["000000000000000000000001"]}`
	code, _, err := controller.JoinGroup(groupId, agents)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledLeaveGroup_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().LeaveGroup(groupId, agentId).Return(nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	agents := `{"agents":["000000000000000000000001"]}`
	code, _, err := controller.LeaveGroup(groupId, agents)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledLeaveGroupWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	agents := `{"agents":["000000000000000000000001"]}`
	code, _, err := controller.LeaveGroup(groupId, agents)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", err.Error())
	case errors.DBConnectionError:
	}
}

func TestCalledLeaveGroupWithInvalidRequestBody_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	invalidJsonStr := `{"invalidJson"}`
	code, _, err := controller.LeaveGroup(groupId, invalidJsonStr)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InvalidParamError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InvalidParamError", err.Error())
	case errors.InvalidJSON:
	}
}

func TestCalledLeaveGroupWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().LeaveGroup(groupId, agentId).Return(notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	agents := `{"agents":["000000000000000000000001"]}`
	code, _, err := controller.LeaveGroup(groupId, agents)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledDeleteGroup_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().DeleteGroup(groupId).Return(nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.DeleteGroup(groupId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledDeleteGroupWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.DeleteGroup(groupId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", err.Error())
	case errors.DBConnectionError:
	}
}

func TestCalledDeleteGroupWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().DeleteGroup(groupId).Return(notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.DeleteGroup(groupId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledDeployApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	respStr := []string{`{"id":"000000000000000000000000"}`, `{"id":"000000000000000000000000"}`}
	expectedRes := map[string]interface{}{
		"id": "000000000000000000000000",
	}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembers(groupId).Return(members, nil),
		msgMockObj.EXPECT().DeployApp(membersAddress, body).Return(respCode, respStr),
		dbManagerMockObj.EXPECT().AddAppToAgent(agentId, appId).Return(nil).AnyTimes(),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, res, err := controller.DeployApp(groupId, body)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledDeployAppWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.DeployApp(groupId, body)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", err.Error())
	case errors.DBConnectionError:
	}
}

func TestCalledDeployAppWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembers(groupId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.DeployApp(groupId, body)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledDeployAppWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembers(groupId).Return(members, nil),
		msgMockObj.EXPECT().DeployApp(membersAddress, body).Return(respCode, invalidRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.DeployApp(groupId, body)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InternalServerError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InternalServerError", err.Error())
	case errors.InternalServerError:
	}
}

func TestCalledDeployAppWhenFailedToAddAppIdToDB_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	respStr := []string{`{"id":"000000000000000000000000"}`}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembers(groupId).Return(members, nil),
		msgMockObj.EXPECT().DeployApp(membersAddress, body).Return(respCode, respStr),
		dbManagerMockObj.EXPECT().AddAppToAgent(agentId, appId).Return(notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.DeployApp(groupId, body)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledDeployAppWhenMessengerReturnsPartialSuccess_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	partialSuccessRespStr := []string{`{"id":"000000000000000000000000"}`, `{"message":"errorMsg"}`}
	expectedRes := map[string]interface{}{
		"id": "000000000000000000000000",
		"responses": []map[string]interface{}{
			map[string]interface{}{
				"id":   agentId,
				"code": results.OK,
			},
			map[string]interface{}{
				"id":      agentId,
				"code":    results.ERROR,
				"message": "errorMsg",
			},
		},
	}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembers(groupId).Return(members, nil),
		msgMockObj.EXPECT().DeployApp(membersAddress, body).Return(partialSuccessRespCode, partialSuccessRespStr),
		dbManagerMockObj.EXPECT().AddAppToAgent(agentId, appId).Return(nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, res, err := controller.DeployApp(groupId, body)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.MULTI_STATUS {
		t.Errorf("Expected code: %d, actual code: %d", results.MULTI_STATUS, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledGetApps_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedRes := map[string]interface{}{
		"apps": []map[string]interface{}{{
			"id":      appId,
			"members": []string{agentId, agentId},
		}},
	}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembers(groupId).Return(members, nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, res, err := controller.GetApps(groupId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledGetAppsWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.GetApps(groupId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", err.Error())
	case errors.DBConnectionError:
	}
}

func TestCalledGetAppsWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembers(groupId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.GetApps(groupId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundError", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	respStr := []string{`{"description":"description"}`, `{"description":"description"}`}
	expectedRes := map[string]interface{}{
		"responses": []map[string]interface{}{{
			"description": "description",
			"id":          members[0]["id"],
		},
			{
				"description": "description",
				"id":          members[0]["id"],
			}},
	}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().InfoApp(membersAddress, appId).Return(respCode, respStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, res, err := controller.GetApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledGetAppWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.GetApp(groupId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", err.Error())
	case errors.DBConnectionError:
	}
}

func TestCalledGetAppWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.GetApp(groupId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundError", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetAppWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidRespStr := []string{`{"invalidJson"}`, `{"invalidJson"}`}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().InfoApp(membersAddress, appId).Return(respCode, invalidRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.GetApp(groupId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InternalServerError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InternalServerError", err.Error())
	case errors.InternalServerError:
	}
}

func TestCalledGetAppWhenMessengerReturnsPartialSuccess_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	partialSuccessRespStr := []string{`{"description": "description"}`, `{"message":"errorMsg"}`}
	expectedRes := map[string]interface{}{
		"responses": []map[string]interface{}{
			map[string]interface{}{
				"id":          agentId,
				"code":        results.OK,
				"description": "description",
			},
			map[string]interface{}{
				"id":      agentId,
				"code":    results.ERROR,
				"message": "errorMsg",
			},
		},
	}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().InfoApp(membersAddress, appId).Return(partialSuccessRespCode, partialSuccessRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, res, err := controller.GetApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.MULTI_STATUS {
		t.Errorf("Expected code: %d, actual code: %d", results.MULTI_STATUS, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledUpdateAppInfo_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().UpdateAppInfo(membersAddress, appId, body).Return(respCode, nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.UpdateAppInfo(groupId, appId, body)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledUpdateAppInfoWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.UpdateAppInfo(groupId, appId, body)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", err.Error())
	case errors.DBConnectionError:
	}
}

func TestCalledUpdateAppInfoWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.UpdateAppInfo(groupId, appId, body)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundError", err.Error())
	case errors.NotFound:
	}
}

func TestCalledUpdateAppInfoWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidRespStr := []string{`{"invalidJson"}`, `{"invalidJson"}`}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().UpdateAppInfo(membersAddress, appId, body).Return(respCode, invalidRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.UpdateAppInfo(groupId, appId, body)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InternalServerError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InternalServerError", err.Error())
	case errors.InternalServerError:
	}
}

func TestCalledUpdateAppInfoWhenMessengerReturnsPartialSuccess_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	partialSuccessRespStr := []string{`{"message": "successMsg"}`, `{"message":"errorMsg"}`}
	expectedRes := map[string]interface{}{
		"responses": []map[string]interface{}{
			map[string]interface{}{
				"id":   agentId,
				"code": results.OK,
			},
			map[string]interface{}{
				"id":      agentId,
				"code":    results.ERROR,
				"message": "errorMsg",
			},
		},
	}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().UpdateAppInfo(membersAddress, appId, body).Return(partialSuccessRespCode, partialSuccessRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, res, err := controller.UpdateAppInfo(groupId, appId, body)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.MULTI_STATUS {
		t.Errorf("Expected code: %d, actual code: %d", results.MULTI_STATUS, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledUpdateApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().UpdateApp(membersAddress, appId).Return(respCode, nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.UpdateApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledUpdateAppWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.UpdateApp(groupId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", err.Error())
	case errors.DBConnectionError:
	}
}

func TestCalledUpdateAppWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.UpdateApp(groupId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFoundError", err.Error())
	case errors.NotFound:
	}
}

func TestCalledUpdateAppWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidRespStr := []string{`{"invalidJson"}`, `{"invalidJson"}`}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().UpdateApp(membersAddress, appId).Return(respCode, invalidRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.UpdateApp(groupId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InternalServerError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InternalServerError", err.Error())
	case errors.InternalServerError:
	}
}

func TestCalledUpdateAppWhenMessengerReturnsPartialSuccess_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	partialSuccessRespStr := []string{`{"message": "successMsg"}`, `{"message":"errorMsg"}`}
	expectedRes := map[string]interface{}{
		"responses": []map[string]interface{}{
			map[string]interface{}{
				"id":   agentId,
				"code": results.OK,
			},
			map[string]interface{}{
				"id":      agentId,
				"code":    results.ERROR,
				"message": "errorMsg",
			},
		},
	}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().UpdateApp(membersAddress, appId).Return(partialSuccessRespCode, partialSuccessRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, res, err := controller.UpdateApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.MULTI_STATUS {
		t.Errorf("Expected code: %d, actual code: %d", results.MULTI_STATUS, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledStartApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().StartApp(membersAddress, appId).Return(respCode, nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.StartApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledStartAppWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.StartApp(groupId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", err.Error())
	case errors.DBConnectionError:
	}
}

func TestCalledStartAppWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.StartApp(groupId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledStartAppWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidRespStr := []string{`{"invalidJson"}`, `{"invalidJson"}`}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().StartApp(membersAddress, appId).Return(respCode, invalidRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.StartApp(groupId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InternalServerError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InternalServerError", err.Error())
	case errors.InternalServerError:
	}
}

func TestCalledStartAppWhenMessengerReturnsPartialSuccess_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	partialSuccessRespStr := []string{`{"message": "successMsg"}`, `{"message":"errorMsg"}`}
	expectedRes := map[string]interface{}{
		"responses": []map[string]interface{}{
			map[string]interface{}{
				"id":   agentId,
				"code": results.OK,
			},
			map[string]interface{}{
				"id":      agentId,
				"code":    results.ERROR,
				"message": "errorMsg",
			},
		},
	}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().StartApp(membersAddress, appId).Return(partialSuccessRespCode, partialSuccessRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, res, err := controller.StartApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.MULTI_STATUS {
		t.Errorf("Expected code: %d, actual code: %d", results.MULTI_STATUS, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledStopApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().StopApp(membersAddress, appId).Return(respCode, nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.StopApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledStopAppWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.StopApp(groupId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", err.Error())
	case errors.DBConnectionError:
	}
}

func TestCalledStopAppWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.StopApp(groupId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledStopAppWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidRespStr := []string{`{"invalidJson"}`, `{"invalidJson"}`}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().StopApp(membersAddress, appId).Return(respCode, invalidRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.StopApp(groupId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InternalServerError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InternalServerError", err.Error())
	case errors.InternalServerError:
	}
}

func TestCalledStopAppWhenMessengerReturnsPartialSuccess_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	partialSuccessRespStr := []string{`{"message": "successMsg"}`, `{"message":"errorMsg"}`}
	expectedRes := map[string]interface{}{
		"responses": []map[string]interface{}{
			map[string]interface{}{
				"id":   agentId,
				"code": results.OK,
			},
			map[string]interface{}{
				"id":      agentId,
				"code":    results.ERROR,
				"message": "errorMsg",
			},
		},
	}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().StopApp(membersAddress, appId).Return(partialSuccessRespCode, partialSuccessRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, res, err := controller.StopApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.MULTI_STATUS {
		t.Errorf("Expected code: %d, actual code: %d", results.MULTI_STATUS, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledDeleteApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().DeleteApp(membersAddress, appId).Return(respCode, nil),
		dbManagerMockObj.EXPECT().DeleteAppFromAgent(agentId, appId).Return(nil).AnyTimes(),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.DeleteApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledDeleteAppWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.DeleteApp(groupId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "DBConnectionError", err.Error())
	case errors.DBConnectionError:
	}
}

func TestCalledDeleteAppWhenDBHasNotMatchedGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.DeleteApp(groupId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledDeleteAppWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidRespStr := []string{`{"invalidJson"}`, `{"invalidJson"}`}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().DeleteApp(membersAddress, appId).Return(respCode, invalidRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.DeleteApp(groupId, appId)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InternalServerError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InternalServerError", err.Error())
	case errors.InternalServerError:
	}
}

func TestCalledDeleteAppWhenMessengerReturnsPartialSuccess_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	partialSuccessRespStr := []string{`{"message": "successMsg"}`, `{"message":"errorMsg"}`}
	expectedRes := map[string]interface{}{
		"responses": []map[string]interface{}{
			map[string]interface{}{
				"id":   agentId,
				"code": results.OK,
			},
			map[string]interface{}{
				"id":      agentId,
				"code":    results.ERROR,
				"message": "errorMsg",
			},
		},
	}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetGroupMembersByAppID(groupId, appId).Return(members, nil),
		msgMockObj.EXPECT().DeleteApp(membersAddress, appId).Return(partialSuccessRespCode, partialSuccessRespStr),
		dbManagerMockObj.EXPECT().DeleteAppFromAgent(agentId, appId).Return(nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, res, err := controller.DeleteApp(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.MULTI_STATUS {
		t.Errorf("Expected code: %d, actual code: %d", results.MULTI_STATUS, code)
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}
