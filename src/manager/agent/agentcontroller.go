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

// Package agent provides an interfaces to deploy, update, start, stop, delete
// an application to target edge device.
package agent

import (
	"commons/errors"
	"commons/logger"
	"commons/results"
	"db"
	"encoding/json"
	"messenger"
	"strconv"
	"time"
)

const (
	AGENTS                      = "agents"       // used to indicate a list of agents.
	ID                          = "id"           // used to indicate an agent id.
	HOST                        = "host"         // used to indicate an agent address.
	PORT                        = "port"         // used to indicate an agent port.
	DEFAULT_SDA_PORT            = "48098"        // default service deployment agent port.
	STATUS_CONNECTED            = "connected"    // used to update agent status with connected.
	STATUS_DISCONNECTED         = "disconnected" // used to update agent status with disconnected.
	INTERVAL                    = "interval"     // a period between two healthcheck message.
	MAXIMUM_NETWORK_LATENCY_SEC = 3              // the term used to indicate any kind of delay that happens in data communication over a network.
	TIME_UNIT                   = time.Minute    // the minute is a unit of time for healthcheck.
)

type AgentController struct{}

var dbConnector db.DBConnection
var httpMessenger messenger.MessengerInterface
var timers map[string]chan bool

func init() {
	dbConnector = db.DBConnector{}
	httpMessenger = messenger.SdamMsgrImpl{}

	timers = make(map[string]chan bool)
}

// AddAgent inserts a new agent with ip which is passed in call to function.
// If successful, a unique id that is created automatically will be returned.
// otherwise, an appropriate error will be returned.
func (AgentController) AddAgent(body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// If body is not empty, try to get agent id from body.
	// This code will be used to update the information of agent without changing id.
	bodyMap, err := convertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Check whether 'ip' is included.
	_, exists := bodyMap["ip"]
	if !exists {
		return results.ERROR, nil, errors.InvalidJSON{"ip field is required"}
	}

	// Get agent with given ip.
	agent, err := db.GetAgentByIP(bodyMap["ip"].(string))
	if err == nil {
		// Agent with that ip address already exists in the database.
		res := make(map[string]interface{})
		res[ID] = agent[ID]
		return results.OK, res, err
	}

	// Add new agent to database with given ip, port, status.
	agent, err = db.AddAgent(bodyMap["ip"].(string), DEFAULT_SDA_PORT, STATUS_CONNECTED)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	res := make(map[string]interface{})
	res[ID] = agent[ID]
	return results.OK, res, err
}

// PingAgent starts timer with received interval.
// If agent does not send next healthcheck message in interval time,
// change the status of device from connected to disconnected.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (AgentController) PingAgent(agentId string, ip string, body string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}
	defer db.Close()

	// Get agent specified by agentId parameter.
	_, err = db.GetAgent(agentId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	bodyMap, err := convertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	_, exists := timers[agentId]
	if !exists {
		logger.Logging(logger.DEBUG, "first ping request is received from agent")
	} else {
		if timers[agentId] != nil {
			// If ping request is received in interval time, send signal to stop timer.
			timers[agentId] <- true
			logger.Logging(logger.DEBUG, "ping request is received in interval time")
		} else {
			logger.Logging(logger.DEBUG, "ping request is received after interval time-out")
			err = db.UpdateAgentStatus(agentId, STATUS_CONNECTED)
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
			}
		}
	}

	// Start timer with received interval time.
	interval, err := strconv.Atoi(bodyMap[INTERVAL].(string))
	timeDurationMin := time.Duration(interval+MAXIMUM_NETWORK_LATENCY_SEC) * TIME_UNIT
	timer := time.NewTimer(timeDurationMin)
	go func() {
		quit := make(chan bool)
		timers[agentId] = quit

		select {
		// Block until timer finishes.
		case <-timer.C:
			logger.Logging(logger.ERROR, "ping request is not received in interval time")

			// Connect to the database.
			db, err := dbConnector.Connect()
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
			}
			defer db.Close()

			// Status is updated with 'disconnected'.
			err = db.UpdateAgentStatus(agentId, STATUS_DISCONNECTED)
			if err != nil {
				logger.Logging(logger.ERROR, err.Error())
			}

		case <-quit:
			timer.Stop()
			return
		}

		timers[agentId] = nil
		close(quit)
	}()

	return results.OK, err
}

// DeleteAgent deletes the agent with a primary key matching the agentId argument.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (AgentController) DeleteAgent(agentId string) (int, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}
	defer db.Close()

	// Get agent specified by agentId parameter.
	agent, err := db.GetAgent(agentId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	// Send request to unregister a specific agent.
	address := getAgentAddress(agent)
	codes, _ := httpMessenger.Unregister(address)

	result := codes[0]
	if !isSuccessCode(result) {
		return results.ERROR, err
	}

	if timers[agentId] != nil {
		timers[agentId] <- true
	}
	delete(timers, agentId)

	// Delete agent specified by agentId parameter.
	err = db.DeleteAgent(agentId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, err
	}

	return results.OK, err
}

// GetAgent returns the agent with a primary key matching the agentId argument.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (AgentController) GetAgent(agentId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// Get agent specified by agentId parameter.
	agent, err := db.GetAgent(agentId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return results.OK, agent, err
}

// GetAgents returns all agents in databases as an array.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (AgentController) GetAgents() (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// Get all agents stored in the database.
	agents, err := db.GetAllAgents()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	res := make(map[string]interface{})
	res[AGENTS] = agents

	return results.OK, res, err
}

// DeployApp request an deployment of edge services to an agent specified by agentId parameter.
// If response code represents success, add an app id to a list of installed app and returns it.
// Otherwise, an appropriate error will be returned.
func (AgentController) DeployApp(agentId string, body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// Get agent specified by agentId parameter.
	agent, err := db.GetAgent(agentId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Request an deployment of edge services to a specific agent.
	address := getAgentAddress(agent)
	codes, respStr := httpMessenger.DeployApp(address, body)

	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// if response code represents success, insert the installed appId into db.
	result := codes[0]
	if isSuccessCode(result) {
		err = db.AddAppToAgent(agentId, respMap[ID].(string))
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}
	}

	return result, respMap, err
}

// GetApps request a list of applications that is deployed to an agent
// specified by agentId parameter.
// If response code represents success, returns a list of applications.
// Otherwise, an appropriate error will be returned.
func (AgentController) GetApps(agentId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// Get agent specified by agentId parameter.
	agent, err := db.GetAgent(agentId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Request list of applications that is deployed to agent.
	address := getAgentAddress(agent)
	codes, respStr := httpMessenger.InfoApps(address)

	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
}

// GetApp gets the application's information of the agent specified by agentId parameter.
// If response code represents success, returns information of application.
// Otherwise, an appropriate error will be returned.
func (AgentController) GetApp(agentId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// Get agent including app specified by appId parameter.
	agent, err := db.GetAgentByAppID(agentId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Request get target application's information
	address := getAgentAddress(agent)
	codes, respStr := httpMessenger.InfoApp(address, appId)

	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
}

// UpdateApp request to update an application specified by appId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (AgentController) UpdateAppInfo(agentId string, appId string, body string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// Get agent including app specified by appId parameter.
	agent, err := db.GetAgentByAppID(agentId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Request update target application's information.
	address := getAgentAddress(agent)
	codes, respStr := httpMessenger.UpdateAppInfo(address, appId, body)

	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
}

// DeleteApp request to delete an application specified by appId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (AgentController) DeleteApp(agentId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// Get agent including app specified by appId parameter.
	agent, err := db.GetAgentByAppID(agentId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Request delete target application
	address := getAgentAddress(agent)
	codes, respStr := httpMessenger.DeleteApp(address, appId)

	result := codes[0]
	if !isSuccessCode(result) {
		respMap, err := convertRespToMap(respStr)
		if err != nil {
			logger.Logging(logger.ERROR, err.Error())
			return results.ERROR, nil, err
		}
		return result, respMap, err
	}

	// if response code represents success, delete the appId from db.
	err = db.DeleteAppFromAgent(agentId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, nil, err
}

// UpdateAppInfo request to update all of images which is included an application
// specified by appId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (AgentController) UpdateApp(agentId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// Get agent including app specified by appId parameter.
	agent, err := db.GetAgentByAppID(agentId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Request checking and updating all of images which is included target.
	address := getAgentAddress(agent)
	codes, respStr := httpMessenger.UpdateApp(address, appId)

	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
}

// StartApp request to start an application specified by appId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (AgentController) StartApp(agentId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// Get agent including app specified by appId parameter.
	agent, err := db.GetAgentByAppID(agentId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Request start target application.
	address := getAgentAddress(agent)
	codes, respStr := httpMessenger.StartApp(address, appId)

	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
}

// StopApp request to stop an application specified by appId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (AgentController) StopApp(agentId string, appId string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Connect to the database.
	db, err := dbConnector.Connect()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}
	defer db.Close()

	// Get agent including app specified by appId parameter.
	agent, err := db.GetAgentByAppID(agentId, appId)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	// Request stop target application.
	address := getAgentAddress(agent)
	codes, respStr := httpMessenger.StopApp(address, appId)

	result := codes[0]
	respMap, err := convertRespToMap(respStr)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	return result, respMap, err
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

// getAgentAddress returns an address as an array.
func getAgentAddress(agent map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 1)
	result[0] = map[string]interface{}{
		"host": agent["host"],
		"port": agent["port"],
	}
	return result
}

// convertRespToMap converts a response in the form of JSON data into a map.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func convertRespToMap(respStr []string) (map[string]interface{}, error) {
	resp, err := convertJsonToMap(respStr[0])
	if err != nil {
		logger.Logging(logger.ERROR, "Failed to convert response from string to map")
		return nil, errors.InternalServerError{"Json Converting Failed"}
	}
	return resp, err
}

// isSuccessCode returns true in case of success and false otherwise.
func isSuccessCode(code int) bool {
	if code >= 200 && code <= 299 {
		return true
	}
	return false
}
