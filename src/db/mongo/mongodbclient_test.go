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
package mongo

import (
	errors "commons/errors"
	mgomocks "db/mongo/wrapper/mocks"
	"github.com/golang/mock/gomock"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"testing"
)

const (
	validUrl        = "192.168.0.1:27017"
	dbName          = "DeploymentManagerDB"
	collectionName  = "AGENT"
	status          = "connected"
	appId           = "000000000000000000000000"
	agentId         = "000000000000000000000001"
	groupId         = "000000000000000000000002"
	invalidObjectId = ""
)

var (
	dummySession       = mgomocks.MockSession{}
	connectionError    = errors.DBConnectionError{}
	invalidObjectError = errors.InvalidObjectId{invalidObjectId}
	notFoundError      = errors.NotFound{}
)

func TestCalledConnectWithEmptyURL_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	invalidUrl := ""

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(invalidUrl).Return(&dummySession, connectionError),
	)
	mgoDial = connectionMockObj

	builder := MongoBuilder{}
	err := builder.Connect(invalidUrl)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "UnknownError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "UnknownError", err.Error())
	case errors.DBConnectionError:
	}
}

func TestCalledConnectWithValidURL_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(&dummySession, nil),
	)
	mgoDial = connectionMockObj

	builder := MongoBuilder{}
	err := builder.Connect(validUrl)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledCreateDBWithInvalidSession_ExpectErrorReturn(t *testing.T) {
	builder := MongoBuilder{}
	_, err := builder.CreateDB()

	switch err.(type) {
	default:
		t.Error()
	case errors.DBOperationError:
	}
}

func TestCalledCreateDBWithValidSession_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(&dummySession, nil),
	)
	mgoDial = connectionMockObj

	builder := MongoBuilder{}
	_ = builder.Connect(validUrl)

	_, err := builder.CreateDB()

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledClose_ExpectSessionClosed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sessionMockObj := mgomocks.NewMockSession(mockCtrl)

	dbManager := MongoDBManager{
		mgoSession: sessionMockObj,
	}

	gomock.InOrder(
		sessionMockObj.EXPECT().Close(),
	)

	dbManager.Close()
}

func TestCalled_GetCollcetion_ExpectToCCalled(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	collectionMockObj := mgomocks.NewMockCollection(mockCtrl)
	dbMockObj := mgomocks.NewMockDatabase(mockCtrl)

	dbManager := MongoDBManager{
		mgoSession: sessionMockObj,
	}

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(dbName).Return(dbMockObj),
		dbMockObj.EXPECT().C(collectionName).Return(collectionMockObj),
	)

	dbManager.getCollection(collectionName)
}

func TestCalledAddAgent_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	host := "192.168.0.1"
	port := "8080"

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Insert(gomock.Any()).Return(nil),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	_, err := dbManager.AddAgent(host, port, status)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledAddAgentWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	host := "192.168.0.1"
	port := "8080"

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Insert(gomock.Any()).Return(mgo.ErrNotFound),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	_, err := dbManager.AddAgent(host, port, status)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledUpdateAgentAddress_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}
	update := bson.M{"$set": bson.M{"host": "192.168.0.1", "port": "48098"}}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Update(query, update).Return(nil),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	err := dbManager.UpdateAgentAddress(agentId, "192.168.0.1", "48098")

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledUpdateAgentAddressWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManager := MongoDBManager{}
	err := dbManager.UpdateAgentAddress(invalidObjectId, "192.168.0.1", "48098")

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledUpdateAgentAddressWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}
	update := bson.M{"$set": bson.M{"host": "192.168.0.1", "port": "48098"}}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Update(query, update).Return(mgo.ErrNotFound),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	err := dbManager.UpdateAgentAddress(agentId, "192.168.0.1", "48098")

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledUpdateAgentStatus_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}
	update := bson.M{"$set": bson.M{"status": "connected"}}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Update(query, update).Return(nil),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	err := dbManager.UpdateAgentStatus(agentId, "connected")

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledUpdateAgentStatusWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManager := MongoDBManager{}
	err := dbManager.UpdateAgentStatus(invalidObjectId, "connected")

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledUpdateAgentStatusWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}
	update := bson.M{"$set": bson.M{"status": "connected"}}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Update(query, update).Return(mgo.ErrNotFound),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	err := dbManager.UpdateAgentStatus(agentId, "connected")

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

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}
	arg := Agent{ID: bson.ObjectIdHex(agentId), Host: "192.168.0.1", Port: "8888", Apps: []string{}, Status: status}
	expectedRes := map[string]interface{}{
		"id":     agentId,
		"host":   "192.168.0.1",
		"port":   "8888",
		"apps":   []string{},
		"status": status,
	}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)
	queryMockObj := mgomocks.NewMockQuery(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(query).Return(queryMockObj),
		queryMockObj.EXPECT().One(gomock.Any()).SetArg(0, arg).Return(nil),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	res, err := dbManager.GetAgent(agentId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledGetAgentWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManager := MongoDBManager{}
	_, err := dbManager.GetAgent(invalidObjectId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledGetAgentWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)
	queryMockObj := mgomocks.NewMockQuery(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(query).Return(queryMockObj),
		queryMockObj.EXPECT().One(gomock.Any()).Return(mgo.ErrNotFound),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	_, err := dbManager.GetAgent(agentId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetAllAgents_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	args := []Agent{{ID: bson.ObjectIdHex(agentId), Host: "192.168.0.1", Port: "8888", Apps: []string{}, Status: status}}
	expectedRes := []map[string]interface{}{{
		"id":     agentId,
		"host":   "192.168.0.1",
		"port":   "8888",
		"apps":   []string{},
		"status": status,
	}}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)
	queryMockObj := mgomocks.NewMockQuery(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(nil).Return(queryMockObj),
		queryMockObj.EXPECT().All(gomock.Any()).SetArg(0, args).Return(nil),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	res, err := dbManager.GetAllAgents()

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledGetAllAgentsWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)
	queryMockObj := mgomocks.NewMockQuery(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(nil).Return(queryMockObj),
		queryMockObj.EXPECT().All(gomock.Any()).Return(mgo.ErrNotFound),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	_, err := dbManager.GetAllAgents()

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetAgentByAppID_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId), "apps": bson.M{"$in": []string{appId}}}
	arg := Agent{ID: bson.ObjectIdHex(agentId), Host: "192.168.0.1", Port: "8888", Apps: []string{}, Status: status}
	expectedRes := map[string]interface{}{
		"id":     agentId,
		"host":   "192.168.0.1",
		"port":   "8888",
		"apps":   []string{},
		"status": status,
	}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)
	queryMockObj := mgomocks.NewMockQuery(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(query).Return(queryMockObj),
		queryMockObj.EXPECT().One(gomock.Any()).SetArg(0, arg).Return(nil),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	res, err := dbManager.GetAgentByAppID(agentId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledGetAgentByAppIDWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId), "apps": bson.M{"$in": []string{appId}}}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)
	queryMockObj := mgomocks.NewMockQuery(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(query).Return(queryMockObj),
		queryMockObj.EXPECT().One(gomock.Any()).Return(mgo.ErrNotFound),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	_, err := dbManager.GetAgentByAppID(agentId, appId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetAgentByAppIDWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManager := MongoDBManager{}
	_, err := dbManager.GetAgentByAppID(invalidObjectId, appId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledAddAppToAgent_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}
	update := bson.M{"$addToSet": bson.M{"apps": appId}}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Update(query, update).Return(nil),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	err := dbManager.AddAppToAgent(agentId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledAddAppToAgentWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManager := MongoDBManager{}
	err := dbManager.AddAppToAgent(invalidObjectId, appId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledAddAppToAgentWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}
	update := bson.M{"$addToSet": bson.M{"apps": appId}}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Update(query, update).Return(mgo.ErrNotFound),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	err := dbManager.AddAppToAgent(agentId, appId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledDeleteAppFromAgent_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}
	update := bson.M{"$pull": bson.M{"apps": appId}}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Update(query, update).Return(nil),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	err := dbManager.DeleteAppFromAgent(agentId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledDeleteAppFromAgentWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManager := MongoDBManager{}
	err := dbManager.DeleteAppFromAgent(invalidObjectId, appId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledDeleteAppFromAgentWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}
	update := bson.M{"$pull": bson.M{"apps": appId}}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Update(query, update).Return(mgo.ErrNotFound),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	err := dbManager.DeleteAppFromAgent(agentId, appId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledDeleteAgent_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Remove(query).Return(nil),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	err := dbManager.DeleteAgent(agentId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledDeleteAgentWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManager := MongoDBManager{}
	err := dbManager.DeleteAgent(invalidObjectId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledDeleteAgentWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Remove(query).Return(mgo.ErrNotFound),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	err := dbManager.DeleteAgent(agentId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledCreateGroup_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Insert(gomock.Any()).Return(nil),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	_, err := dbManager.CreateGroup()

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledCreateGroupWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Insert(gomock.Any()).Return(mgo.ErrNotFound),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	_, err := dbManager.CreateGroup()

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

	query := bson.M{"_id": bson.ObjectIdHex(groupId)}
	arg := Group{ID: bson.ObjectIdHex(groupId), Members: []string{}}
	expectedRes := map[string]interface{}{
		"id":      groupId,
		"members": []string{},
	}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)
	queryMockObj := mgomocks.NewMockQuery(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(query).Return(queryMockObj),
		queryMockObj.EXPECT().One(gomock.Any()).SetArg(0, arg).Return(nil),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	res, err := dbManager.GetGroup(groupId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledGetGroupWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManager := MongoDBManager{}
	_, err := dbManager.GetGroup(invalidObjectId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledGetGroupWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(groupId)}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)
	queryMockObj := mgomocks.NewMockQuery(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(query).Return(queryMockObj),
		queryMockObj.EXPECT().One(gomock.Any()).Return(mgo.ErrNotFound),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	_, err := dbManager.GetGroup(groupId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetAllGroups_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	args := []Group{{ID: bson.ObjectIdHex(groupId), Members: []string{}}}
	expectedRes := []map[string]interface{}{{
		"id":      groupId,
		"members": []string{},
	}}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)
	queryMockObj := mgomocks.NewMockQuery(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(nil).Return(queryMockObj),
		queryMockObj.EXPECT().All(gomock.Any()).SetArg(0, args).Return(nil),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	res, err := dbManager.GetAllGroups()

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledGetAllGroupsWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)
	queryMockObj := mgomocks.NewMockQuery(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(nil).Return(queryMockObj),
		queryMockObj.EXPECT().All(gomock.Any()).Return(mgo.ErrNotFound),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	_, err := dbManager.GetAllGroups()

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

	query := bson.M{"_id": bson.ObjectIdHex(groupId)}
	update := bson.M{"$addToSet": bson.M{"members": agentId}}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Update(query, update).Return(nil),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	err := dbManager.JoinGroup(groupId, agentId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledJoinGroupWithInvalidObjectIdAboutGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManager := MongoDBManager{}
	err := dbManager.JoinGroup(invalidObjectId, agentId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledJoinGroupWithInvalidObjectIdAboutAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManager := MongoDBManager{}
	err := dbManager.JoinGroup(groupId, invalidObjectId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledJoinGroupWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(groupId)}
	update := bson.M{"$addToSet": bson.M{"members": agentId}}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Update(query, update).Return(mgo.ErrNotFound),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	err := dbManager.JoinGroup(groupId, agentId)

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

	query := bson.M{"_id": bson.ObjectIdHex(groupId)}
	update := bson.M{"$pull": bson.M{"members": agentId}}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Update(query, update).Return(nil),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	err := dbManager.LeaveGroup(groupId, agentId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledLeaveGroupWithInvalidObjectIdAboutGroup_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManager := MongoDBManager{}
	err := dbManager.LeaveGroup(invalidObjectId, agentId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledLeaveGroupWithInvalidObjectIdAboutAgent_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManager := MongoDBManager{}
	err := dbManager.LeaveGroup(groupId, invalidObjectId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledLeaveGroupWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(groupId)}
	update := bson.M{"$pull": bson.M{"members": agentId}}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Update(query, update).Return(mgo.ErrNotFound),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	err := dbManager.LeaveGroup(groupId, agentId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetGroupMembers_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupQuery := bson.M{"_id": bson.ObjectIdHex(groupId)}
	agentQuery := bson.M{"_id": bson.ObjectIdHex(agentId)}
	groupArg := Group{ID: bson.ObjectIdHex(groupId), Members: []string{agentId}}
	agentArg := Agent{ID: bson.ObjectIdHex(agentId), Host: "192.168.0.1", Port: "8888", Apps: []string{}, Status: status}
	expectedRes := []map[string]interface{}{{
		"id":     agentId,
		"host":   "192.168.0.1",
		"port":   "8888",
		"apps":   []string{},
		"status": status,
	}}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)
	queryMockObj := mgomocks.NewMockQuery(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(groupQuery).Return(queryMockObj),
		queryMockObj.EXPECT().One(gomock.Any()).SetArg(0, groupArg).Return(nil),

		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(agentQuery).Return(queryMockObj),
		queryMockObj.EXPECT().One(gomock.Any()).SetArg(0, agentArg).Return(nil),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	res, err := dbManager.GetGroupMembers(groupId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledGetGroupMembersWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManager := MongoDBManager{}
	_, err := dbManager.GetGroupMembers(invalidObjectId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledGetGroupMembersByAppID_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupQuery := bson.M{"_id": bson.ObjectIdHex(groupId)}
	agentQuery := bson.M{"_id": bson.ObjectIdHex(agentId), "apps": bson.M{"$in": []string{appId}}}
	groupArg := Group{ID: bson.ObjectIdHex(groupId), Members: []string{agentId}}
	agentArg := Agent{ID: bson.ObjectIdHex(agentId), Host: "192.168.0.1", Port: "8888", Apps: []string{appId}, Status: status}
	expectedRes := []map[string]interface{}{{
		"id":     agentId,
		"host":   "192.168.0.1",
		"port":   "8888",
		"apps":   []string{appId},
		"status": status,
	}}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)
	queryMockObj := mgomocks.NewMockQuery(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(groupQuery).Return(queryMockObj),
		queryMockObj.EXPECT().One(gomock.Any()).SetArg(0, groupArg).Return(nil),

		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(agentQuery).Return(queryMockObj),
		queryMockObj.EXPECT().One(gomock.Any()).SetArg(0, agentArg).Return(nil),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	res, err := dbManager.GetGroupMembersByAppID(groupId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledGetGroupMembersByAppIDWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManager := MongoDBManager{}
	_, err := dbManager.GetGroupMembersByAppID(invalidObjectId, appId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "invalidObjectError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "invalidObjectError", err.Error())
	case errors.InvalidObjectId:
	}
}

func TestCalledDeleteGroup_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(groupId)}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Remove(query).Return(nil),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	err := dbManager.DeleteGroup(groupId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledDeleteGroupWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbManager := MongoDBManager{}
	err := dbManager.DeleteGroup(invalidObjectId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "invalidObjectError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "invalidObjectError", err.Error())
	case errors.InvalidObjectId:
	}
}

func TestCalledDeleteGroupWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(groupId)}

	sessionMockObj := mgomocks.NewMockSession(ctrl)
	dbMockObj := mgomocks.NewMockDatabase(ctrl)
	collectionMockObj := mgomocks.NewMockCollection(ctrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Remove(query).Return(mgo.ErrNotFound),
	)

	dbManager := MongoDBManager{mgoSession: sessionMockObj}
	err := dbManager.DeleteGroup(groupId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}
