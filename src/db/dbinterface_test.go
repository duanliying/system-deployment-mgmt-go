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
package db

import (
	"commons/errors"
	"db/mongo"
	"db/mongo/mocks"
	gomock "github.com/golang/mock/gomock"
	"testing"
)

func TestCalledConnectWhenConnectReturnError_ExpectRetrunError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	url := "127.0.0.1:27017"
	dummyError := errors.DBConnectionError{}

	builderMockObj := mocks.NewMockBuilder(mockCtrl)

	gomock.InOrder(
		builderMockObj.EXPECT().Connect(url).Return(dummyError),
	)
	mgoBuilder = builderMockObj

	dbConnector := DBConnector{}
	_, err := dbConnector.Connect()

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "UnknownError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "UnknownError", err.Error())
	case errors.DBConnectionError:
	}
}

func TestCalledConnectWhenCreateDBReturnError_ExpectRetrunError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	url := "127.0.0.1:27017"
	dummyError := errors.DBConnectionError{}

	builderMockObj := mocks.NewMockBuilder(mockCtrl)

	gomock.InOrder(
		builderMockObj.EXPECT().Connect(url).Return(nil),
		builderMockObj.EXPECT().CreateDB().Return(nil, dummyError),
	)
	mgoBuilder = builderMockObj

	dbConnector := DBConnector{}
	_, err := dbConnector.Connect()

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "UnknownError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "UnknownError", err.Error())
	case errors.DBConnectionError:
	}
}

func TestCalled_Connect_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	url := "127.0.0.1:27017"
	dbManager := mongo.MongoDBManager{}

	builderMockObj := mocks.NewMockBuilder(mockCtrl)

	gomock.InOrder(
		builderMockObj.EXPECT().Connect(url).Return(nil),
		builderMockObj.EXPECT().CreateDB().Return(&dbManager, nil),
	)
	mgoBuilder = builderMockObj

	dbConnector := DBConnector{}
	_, err := dbConnector.Connect()

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}
