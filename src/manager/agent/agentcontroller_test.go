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
package agent

import (
	"commons/errors"
	"commons/results"
	dbmocks "db/mocks"
	"github.com/golang/mock/gomock"
	msgmocks "messenger/mocks"
	"reflect"
	"testing"
)

const (
	status  = "connected"
	appId   = "000000000000000000000000"
	agentId = "000000000000000000000001"
	host    = "127.0.0.1"
	port    = "48098"
)

var (
	agent = map[string]interface{}{
		"id":   agentId,
		"host": host,
		"port": port,
		"apps": []string{},
	}
	address = []map[string]interface{}{
		map[string]interface{}{
			"host": host,
			"port": port,
		}}
	body             = `{"description":"description"}`
	respCode         = []int{results.OK}
	errorRespCode    = []int{results.ERROR}
	respStr          = []string{`{"response":"response"}`}
	invalidRespStr   = []string{`{"invalidJson"}`}
	notFoundError    = errors.NotFound{}
	connectionError  = errors.DBConnectionError{}
	invalidJsonError = errors.InvalidJSON{}
)

var controller AgentInterface

func init() {
	controller = AgentController{}
}

func TestCalledAddAgentWithValidBody_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	body := `{"ip":"127.0.0.1"}`
	expectedRes := map[string]interface{}{
		"id": "000000000000000000000001",
	}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().AddAgent(host, port, status).Return(agent, nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, res, err := controller.AddAgent(body)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(res, expectedRes) {
		t.Error()
	}
}

func TestCalledAddAgentWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	body := `{"ip":"127.0.0.1"}`
	code, _, err := controller.AddAgent(body)

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

func TestCalledAddAgentWithInValidJsonFormatBody_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidBody := `{"ip"}`

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.AddAgent(invalidBody)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InvalidJSON", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InvalidJSON", err.Error())
	case errors.InvalidJSON:
	}
}

func TestCalledAddAgentWithInvalidBodyNotIncludingIDField_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidBody := `{"key":"value"}`

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.AddAgent(invalidBody)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InvalidJSON", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InvalidJSON", err.Error())
	case errors.InvalidJSON:
	}
}

func TestCalledAddAgentWhenFailedToInsertNewAgentToDB_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().AddAgent(host, port, status).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	body := `{"ip":"127.0.0.1"}`
	code, _, err := controller.AddAgent(body)

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

func TestCalledPingAgentWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, err := controller.PingAgent(agentId, host, "")

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

func TestCalledPingAgentWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, err := controller.PingAgent(agentId, host, "")

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

func TestCalledPingAgentWithInvalidBody_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(agent, nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, err := controller.PingAgent(agentId, host, "")

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "InvalidJSON", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "InvalidJSON", err.Error())
	case errors.InvalidJSON:
	}
}

func TestCalledDeleteAgent_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().DeleteAgent(agentId).Return(nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, err := controller.DeleteAgent(agentId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledDeleteAgentWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, err := controller.DeleteAgent(agentId)

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

func TestCalledDeleteAgentWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().DeleteAgent(agentId).Return(notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, err := controller.DeleteAgent(agentId)

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

func TestCalledGetAgent_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(agent, nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, res, err := controller.GetAgent(agentId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(res, agent) {
		t.Error()
	}
}

func TestCalledGetAgentWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.GetAgent(agentId)

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

func TestCalledGetAgentWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.GetAgent(agentId)

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

func TestCalledGetAgents_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	agents := []map[string]interface{}{agent}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAllAgents().Return(agents, nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, res, err := controller.GetAgents()

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	if !reflect.DeepEqual(res["agents"].([]map[string]interface{}), agents) {
		t.Error()
	}
}

func TestCalledGetAgentsWhenDBConnectionFailed_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(nil, connectionError),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.GetAgents()

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

func TestCalledGetAgentsWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAllAgents().Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj

	code, _, err := controller.GetAgents()

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

	respStr := []string{`{"id":"000000000000000000000000"}`}
	expectedRes := map[string]interface{}{
		"id": "000000000000000000000000",
	}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(agent, nil),
		msgMockObj.EXPECT().DeployApp(address, body).Return(respCode, respStr),
		dbManagerMockObj.EXPECT().AddAppToAgent(agentId, appId).Return(nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, res, err := controller.DeployApp(agentId, body)

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

	code, _, err := controller.DeployApp(agentId, body)

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

func TestCalledDeployAppWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.DeployApp(agentId, body)

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
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(agent, nil),
		msgMockObj.EXPECT().DeployApp(address, body).Return(respCode, invalidRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.DeployApp(agentId, body)

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
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(agent, nil),
		msgMockObj.EXPECT().DeployApp(address, body).Return(respCode, respStr),
		dbManagerMockObj.EXPECT().AddAppToAgent(agentId, appId).Return(notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.DeployApp(agentId, body)

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

func TestCalledGetApps_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	respStr := []string{`{"description":"description"}`}
	expectedRes := map[string]interface{}{
		"description": "description",
	}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(agent, nil),
		msgMockObj.EXPECT().InfoApps(address).Return(respCode, respStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, res, err := controller.GetApps(agentId)

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

	code, _, err := controller.GetApps(agentId)

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

func TestCalledGetAppsWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(agent, nil),
		msgMockObj.EXPECT().InfoApps(address).Return(respCode, invalidRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.GetApps(agentId)

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

func TestCalledGetAppsWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgent(agentId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.GetApps(agentId)

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

func TestCalledGetApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	respStr := []string{`{"description":"description"}`}
	expectedRes := map[string]interface{}{
		"description": "description",
	}

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().InfoApp(address, appId).Return(respCode, respStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, res, err := controller.GetApp(agentId, appId)

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

	code, _, err := controller.GetApp(agentId, appId)

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

func TestCalledGetAppWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().InfoApp(address, appId).Return(respCode, invalidRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.GetApp(agentId, appId)

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

func TestCalledGetAppWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.GetApp(agentId, appId)

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

func TestCalledUpdateAppInfo_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().UpdateAppInfo(address, appId, body).Return(respCode, respStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.UpdateAppInfo(agentId, appId, body)

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

	code, _, err := controller.UpdateAppInfo(agentId, appId, body)

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

func TestCalledUpdateAppInfoWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.UpdateAppInfo(agentId, appId, body)

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

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().UpdateAppInfo(address, appId, body).Return(respCode, invalidRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.UpdateAppInfo(agentId, appId, body)

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

func TestCalledUpdateApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().UpdateApp(address, appId).Return(respCode, respStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.UpdateApp(agentId, appId)

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

	code, _, err := controller.UpdateApp(agentId, appId)

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

func TestCalledUpdateAppWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.UpdateApp(agentId, appId)

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

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().UpdateApp(address, appId).Return(respCode, invalidRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.UpdateApp(agentId, appId)

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

func TestCalledStartApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().StartApp(address, appId).Return(respCode, respStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.StartApp(agentId, appId)

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

	code, _, err := controller.StartApp(agentId, appId)

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

func TestCalledStartAppWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.StartApp(agentId, appId)

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

func TestCalledStartAppWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().StartApp(address, appId).Return(respCode, invalidRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.StartApp(agentId, appId)

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

func TestCalledStopApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().StopApp(address, appId).Return(respCode, respStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.StopApp(agentId, appId)

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

	code, _, err := controller.StopApp(agentId, appId)

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

func TestCalledStopAppWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.StopApp(agentId, appId)

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

func TestCalledStopAppWhenMessengerReturnsInvalidResponse_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().StopApp(address, appId).Return(respCode, invalidRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.StopApp(agentId, appId)

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

func TestCalledDeleteApp_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().DeleteApp(address, appId).Return(respCode, respStr),
		dbManagerMockObj.EXPECT().DeleteAppFromAgent(agentId, appId).Return(nil),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.DeleteApp(agentId, appId)

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

	code, _, err := controller.DeleteApp(agentId, appId)

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

func TestCalledDeleteAppWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(nil, notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.DeleteApp(agentId, appId)

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

func TestCalledDeleteAppWhenMessengerReturnsErrorCode_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().DeleteApp(address, appId).Return(errorRespCode, respStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.DeleteApp(agentId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}
}

func TestCalledDeleteAppWhenMessengerReturnsErrorCodeWithInvalidResponse_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().DeleteApp(address, appId).Return(errorRespCode, invalidRespStr),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.DeleteApp(agentId, appId)

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

func TestCalledDeleteAppWhenFailedToDeleteAppIdFromDB_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConnectionMockObj := dbmocks.NewMockDBConnection(ctrl)
	dbManagerMockObj := dbmocks.NewMockDBManager(ctrl)
	msgMockObj := msgmocks.NewMockMessengerInterface(ctrl)

	gomock.InOrder(
		dbConnectionMockObj.EXPECT().Connect().Return(dbManagerMockObj, nil),
		dbManagerMockObj.EXPECT().GetAgentByAppID(agentId, appId).Return(agent, nil),
		msgMockObj.EXPECT().DeleteApp(address, appId).Return(respCode, nil),
		dbManagerMockObj.EXPECT().DeleteAppFromAgent(agentId, appId).Return(notFoundError),
		dbManagerMockObj.EXPECT().Close(),
	)
	// pass mockObj to a real object.
	dbConnector = dbConnectionMockObj
	httpMessenger = msgMockObj

	code, _, err := controller.DeleteApp(agentId, appId)

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
