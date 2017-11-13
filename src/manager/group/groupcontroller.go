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

// Package group provides an interfaces to deploy, update, start, stop, delete
// an application to target edge group.
// and also provides operations to manage edge device group (e.g., create, join, leave, delete...).
package group

import (
	"commons/errors"
	"commons/logger"
	"commons/results"
	"db"
	"encoding/json"
	"messenger"
)

const (
	AGENTS        = "agents"      // used to indicate a list of agents.
	GROUPS        = "groups"      // used to indicate a list of groups.
	MEMBERS       = "members"     // used to indicate a list of members.
	APPS          = "apps"        // used to indicate a list of apps.
	ID            = "id"          // used to indicate an id.
	RESPONSE_CODE = "code"        // used to indicate a code.
	ERROR_MESSAGE = "message"     // used to indicate a message.
	RESPONSES     = "responses"   // used to indicate a list of responses.
	DESCRIPTION   = "description" // used to indicate a description.
)

type GroupController struct{}

var dbConnector db.DBConnection
var httpMessenger messenger.MessengerInterface

func init() {
	dbConnector = db.DBConnector{}
	httpMessenger = messenger.SdamMsgrImpl{}
}

// CreateGroup inserts a new group to databases.
// This function returns a unique id in case of success and an error otherwise.
func (GroupController) CreateGroup() (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	group, err := db.CreateGroup()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return results.OK, group, err
}

// GetGroup returns the information of the group specified by groupId parameter.
// If response code represents success, returns information about the group.
// Otherwise, an appropriate error will be returned.
func (GroupController) GetGroup(groupId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	group, err := db.GetGroup(groupId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return results.OK, group, err
}

// GetGroups returns a list of groups that is created on databases.
// If response code represents success, returns a list of groups.
// Otherwise, an appropriate error will be returned.
func (GroupController) GetGroups() (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	groups, err := db.GetAllGroups()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	res := make(map[string]interface{})
	res[GROUPS] = groups

	return results.OK, res, err
}

// JoinGroup adds the agent to a list of members.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (GroupController) JoinGroup(groupId string, body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	bodyMap, err := convertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Check whether 'agents' is included.
	_, exists := bodyMap[AGENTS]
	if !exists {
		return results.ERROR, nil, errors.InvalidJSON{"agents field is required"}
	}

	for _, agentId := range bodyMap[AGENTS].([]interface{}) {
		err = db.JoinGroup(groupId, agentId.(string))
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}
	}

	return results.OK, nil, err
}

// LeaveGroup removes the agent from a list of members.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (GroupController) LeaveGroup(groupId string, body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	bodyMap, err := convertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Check whether 'agents' is included.
	_, exists := bodyMap[AGENTS]
	if !exists {
		return results.ERROR, nil, errors.InvalidJSON{"agents field is required"}
	}

	for _, agentId := range bodyMap[AGENTS].([]interface{}) {
		err = db.LeaveGroup(groupId, agentId.(string))
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}
	}

	return results.OK, nil, err
}

// DeleteGroup deletes the group with a primary key matching the groupId argument.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (GroupController) DeleteGroup(groupId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	err = db.DeleteGroup(groupId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return results.OK, nil, err
}

// DeployApp request an deployment of edge services to a group specified by groupId parameter.
// If response code represents success, add an app id to a list of installed app and returns it.
// Otherwise, an appropriate error will be returned.
func (GroupController) DeployApp(groupId string, body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// Get group members from the database.
	members, err := db.GetGroupMembers(groupId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Request an deployment of edge services to a specific group.
	address := getMemberAddress(members)
	codes, respStr := httpMessenger.DeployApp(address, body)
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// if response code represents success, insert the installed appId into db.
	installedAppId := ""
	for i, agent := range members {
		if isSuccessCode(codes[i]) {
			err = db.AddAppToAgent(agent[ID].(string), respMap[i][ID].(string))
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return results.ERROR, nil, err
			}
			installedAppId = respMap[i][ID].(string)
		}
	}

	result := decideResultCode(codes)
	if result != results.OK {
		// Make separate responses to represent partial failure case.
		resp := make(map[string]interface{})
		resp[RESPONSES] = makeSeparateResponses(members, codes, respMap)
		if installedAppId != "" {
			resp[ID] = installedAppId
		}
		return result, resp, err
	}

	resp := make(map[string]interface{})
	resp[ID] = installedAppId

	return result, resp, err
}

// GetApps request a list of applications that is deployed to a group
// specified by groupId parameter.
// If response code represents success, returns a list of applications.
// Otherwise, an appropriate error will be returned.
func (GroupController) GetApps(groupId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// Get group members from the database.
	members, err := db.GetGroupMembers(groupId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	contains := func(list []map[string]interface{}, appId string) map[string]interface{} {
		for _, item := range list {
			if item[ID] == appId {
				return item
			}
		}
		return nil
	}

	respValue := make([]map[string]interface{}, 0)
	for _, agent := range members {
		for _, appId := range agent[APPS].([]string) {
			item := contains(respValue, appId)
			if item != nil {
				item[MEMBERS] = append(item[MEMBERS].([]string), agent[ID].(string))
			} else {
				item = map[string]interface{}{
					ID:      appId,
					MEMBERS: []string{agent[ID].(string)},
				}
				respValue = append(respValue, item)
			}
		}
	}

	res := make(map[string]interface{})
	res[APPS] = respValue

	return results.OK, res, err
}

// GetApp gets the application's information of the group specified by groupId parameter.
// If response code represents success, returns information of application.
// Otherwise, an appropriate error will be returned.
func (GroupController) GetApp(groupId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// Get group members including app specified by appId parameter.
	members, err := db.GetGroupMembersByAppID(groupId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Request get target application's information.
	address := getMemberAddress(members)
	codes, respStr := httpMessenger.InfoApp(address, appId)
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	result := decideResultCode(codes)
	if result != results.OK {
		// Make separate responses to represent partial failure case.
		resp := make(map[string]interface{})
		resp[RESPONSES] = makeSeparateResponses(members, codes, respMap)

		for i, _ := range members {
			if isSuccessCode(codes[i]) {
				respValue := resp[RESPONSES].([]map[string]interface{})
				for key, value := range respMap[i] {
					respValue[i][key] = value
				}
			}
		}
		return result, resp, err
	}

	resp := make(map[string]interface{})
	respValue := make([]map[string]interface{}, len(members))
	resp[RESPONSES] = respValue

	for i, agent := range members {
		respValue[i] = make(map[string]interface{})
		respValue[i][ID] = agent[ID].(string)
		for key, value := range respMap[i] {
			respValue[i][key] = value
		}
	}

	return result, resp, err
}

// UpdateApp request to update an application specified by appId parameter
// to all members of the group.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (GroupController) UpdateAppInfo(groupId string, appId string, body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// Get group members including app specified by appId parameter.
	members, err := db.GetGroupMembersByAppID(groupId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Request update target application's information.
	address := getMemberAddress(members)
	codes, respStr := httpMessenger.UpdateAppInfo(address, appId, body)
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	result := decideResultCode(codes)
	if result != results.OK {
		// Make separate responses to represent partial failure case.
		resp := make(map[string]interface{})
		resp[RESPONSES] = makeSeparateResponses(members, codes, respMap)
		return result, resp, err
	}

	return result, nil, err
}

// DeleteApp request to delete an application specified by appId parameter
// to all members of the group.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (GroupController) DeleteApp(groupId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// Get group members including app specified by appId parameter.
	members, err := db.GetGroupMembersByAppID(groupId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Request delete target application.
	address := getMemberAddress(members)
	codes, respStr := httpMessenger.DeleteApp(address, appId)
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// if response code represents success, delete the appId from db.
	for i, agent := range members {
		if isSuccessCode(codes[i]) {
			err = db.DeleteAppFromAgent(agent[ID].(string), appId)
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
				return results.ERROR, nil, err
			}
		}
	}

	result := decideResultCode(codes)
	if result != results.OK {
		// Make separate responses to represent partial failure case.
		resp := make(map[string]interface{})
		resp[RESPONSES] = makeSeparateResponses(members, codes, respMap)
		return result, resp, err
	}

	return result, nil, err
}

// UpdateAppInfo request to update all of images which is included an application
// specified by appId parameter to all members of the group.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (GroupController) UpdateApp(groupId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// Get group members including app specified by appId parameter.
	members, err := db.GetGroupMembersByAppID(groupId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Request checking and updating all of images which is included target.
	address := getMemberAddress(members)
	codes, respStr := httpMessenger.UpdateApp(address, appId)
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	result := decideResultCode(codes)
	if result != results.OK {
		// Make separate responses to represent partial failure case.
		resp := make(map[string]interface{})
		resp[RESPONSES] = makeSeparateResponses(members, codes, respMap)
		return result, resp, err
	}

	return result, nil, err
}

// StartApp request to start an application specified by appId parameter
// to all members of the group.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (GroupController) StartApp(groupId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// Get group members including app specified by appId parameter.
	members, err := db.GetGroupMembersByAppID(groupId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Request start target application.
	address := getMemberAddress(members)
	codes, respStr := httpMessenger.StartApp(address, appId)
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	result := decideResultCode(codes)
	if result != results.OK {
		// Make separate responses to represent partial failure case.
		resp := make(map[string]interface{})
		resp[RESPONSES] = makeSeparateResponses(members, codes, respMap)
		return result, resp, err
	}

	return result, nil, err
}

// StopApp request to stop an application specified by appId parameter
// to all members of the group.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (GroupController) StopApp(groupId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// Get group members including app specified by appId parameter.
	members, err := db.GetGroupMembersByAppID(groupId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Request stop target application.
	address := getMemberAddress(members)
	codes, respStr := httpMessenger.StopApp(address, appId)
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	result := decideResultCode(codes)
	if result != results.OK {
		// Make separate responses to represent partial failure case.
		resp := make(map[string]interface{})
		resp[RESPONSES] = makeSeparateResponses(members, codes, respMap)
		return result, resp, err
	}

	return result, nil, err
}

// convertJsonToMap converts JSON data into a map.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func convertJsonToMap(jsonStr string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, errors.InvalidJSON{"Unmarshalling Failed"}
	}
	return result, err
}

// getAgentAddress returns an member's address as an array.
func getMemberAddress(members []map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, len(members))
	for i, agent := range members {
		result[i] = map[string]interface{}{
			"host": agent["host"],
			"port": agent["port"],
		}
	}
	return result
}

// convertRespToMap converts a response in the form of JSON data into a map.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func convertRespToMap(respStr []string) ([]map[string]interface{}, error) {
	respMap := make([]map[string]interface{}, len(respStr))
	for i, v := range respStr {
		resp, err := convertJsonToMap(v)
		if err != nil {
			logger.Logging(logger.ERROR, "Failed to convert response from string to map")
			return nil, errors.InternalServerError{"Json Converting Failed"}
		}
		respMap[i] = resp
	}

	return respMap, nil
}

// isSuccessCode returns true in case of success and false otherwise.
func isSuccessCode(code int) bool {
	if code >= 200 && code <= 299 {
		return true
	}
	return false
}

// decideResultCode returns a result of group operations.
// OK: Returned when all members of the group send a success response.
// MULTI_STATUS: Partial success for multiple requests. Some requests succeeded
//               but at least one failed.
// ERROR: Returned when all members of the gorup send an error response.
func decideResultCode(codes []int) int {
	successCounts := 0
	for _, code := range codes {
		if isSuccessCode(code) {
			successCounts++
		}
	}

	result := results.OK
	switch successCounts {
	case len(codes):
		result = results.OK
	case 0:
		result = results.ERROR
	default:
		result = results.MULTI_STATUS
	}
	return result
}

// makeSeparateResponses used to make a separate response
// when the group operations is a partial success.
func makeSeparateResponses(members []map[string]interface{}, codes []int,
	respMap []map[string]interface{}) []map[string]interface{} {

	respValue := make([]map[string]interface{}, len(members))

	for i, agent := range members {
		respValue[i] = make(map[string]interface{})
		respValue[i][ID] = agent[ID].(string)
		respValue[i][RESPONSE_CODE] = codes[i]

		if !isSuccessCode(codes[i]) {
			respValue[i][ERROR_MESSAGE] = respMap[i][ERROR_MESSAGE]
		}
	}

	return respValue
}
