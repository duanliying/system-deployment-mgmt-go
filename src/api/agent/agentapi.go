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

// Package api/agent provides functionality to handle request related to agent.
package agent

import (
	"api/common"
	"commons/errors"
	"commons/logger"
	"commons/results"
	URL "commons/url"
	"manager/agent"
	"net/http"
	"strings"
)

const (
	GET    string = "GET"
	PUT    string = "PUT"
	POST   string = "POST"
	DELETE string = "DELETE"
)

type _SDAMAgentApisHandler struct{}
type _SDAMAgentApis struct{}

var sdamH _SDAMAgentApisHandler
var sdam _SDAMAgentApis
var sdamAgentController agent.AgentInterface

func init() {
	SdamAgentHandle = sdamH
	SdamAgent = sdam
	sdamAgentController = agent.AgentController{}
}

// Handle calls a proper function according to the url and method received from remote device.
func (sdamH _SDAMAgentApisHandler) Handle(w http.ResponseWriter, req *http.Request) {
	url := strings.Replace(req.URL.Path, URL.Base()+URL.Agents(), "", -1)
	split := strings.Split(url, "/")
	switch len(split) {
	case 1:
		if req.Method == GET {
			SdamAgent.agents(w, req)
		} else {
			common.WriteError(w, errors.InvalidMethod{req.Method})
		}

	case 2:
		if "/"+split[1] == URL.Register() {
			if req.Method == POST {
				SdamAgent.agentRegister(w, req)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			if req.Method == GET {
				agentID := split[1]
				SdamAgent.agent(w, req, agentID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		}

	case 3:
		agentID := split[1]
		if "/"+split[2] == URL.Deploy() {
			if req.Method == POST {
				SdamAgent.agentDeployApp(w, req, agentID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else if "/"+split[2] == URL.Unregister() {
			if req.Method == POST {
				SdamAgent.agentUnregister(w, req, agentID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else if "/"+split[2] == URL.Ping() {
			if req.Method == POST {
				SdamAgent.agentPing(w, req, agentID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else if "/"+split[2] == URL.Apps() {
			if req.Method == GET {
				SdamAgent.agentInfoApps(w, req, agentID)
			} else {
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}

	case 4:
		if "/"+split[2] == URL.Apps() {
			agentID, appID := split[1], split[3]
			switch req.Method {
			case GET:
				SdamAgent.agentInfoApp(w, req, agentID, appID)

			case POST:
				SdamAgent.agentUpdateAppInfo(w, req, agentID, appID)

			case DELETE:
				SdamAgent.agentDeleteApp(w, req, agentID, appID)

			default:
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}

	case 5:
		if "/"+split[2] == URL.Apps() {
			agentID, appID := split[1], split[3]
			switch {
			case "/"+split[4] == URL.Start() && req.Method == POST:
				SdamAgent.agentStartApp(w, req, agentID, appID)

			case "/"+split[4] == URL.Stop() && req.Method == POST:
				SdamAgent.agentStopApp(w, req, agentID, appID)

			case "/"+split[4] == URL.Update() && req.Method == POST:
				SdamAgent.agentUpdateApp(w, req, agentID, appID)

			default:
				common.WriteError(w, errors.InvalidMethod{req.Method})
			}
		} else {
			common.WriteError(w, errors.NotFoundURL{})
		}
	}
}

// agentRegister handles requests which is used to register agent to a list of agents.
//
//    paths: '/api/v1/agents/{agentID}/register'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (sdam _SDAMAgentApis) agentRegister(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "[AGENT] Register New Service Deployment Agent")

	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := sdamAgentController.AddAgent(body)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agentUnregister handles requests which is used to unregister agent from a list of agents.
//
//    paths: '/api/v1/agents/{agentID}/unregister'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (sdam _SDAMAgentApis) agentUnregister(w http.ResponseWriter, req *http.Request, agentID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Unregister New Service Deployment Agent")

	result, err := sdamAgentController.DeleteAgent(agentID)
	common.MakeResponse(w, result, nil, err)
}

// agentPing handles requests which is used to check whether an agent is up.
//
//    paths: '/api/v1/agents/{agentID}/ping'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (sdam _SDAMAgentApis) agentPing(w http.ResponseWriter, req *http.Request, agentID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Ping From Service Deployment Agent")

	ip := strings.Split(req.RemoteAddr, ":")[0]

	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, err := sdamAgentController.PingAgent(agentID, ip, body)
	common.MakeResponse(w, result, nil, err)
}

// agents handles requests which is used to get information of agent identified by the given agentID.
//
//    paths: '/api/v1/agents/{agentID}'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (sdam _SDAMAgentApis) agent(w http.ResponseWriter, req *http.Request, agentID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Get Service Deployment Agent")
	result, resp, err := sdamAgentController.GetAgent(agentID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agents handles requests which is used to get information of all agents registered.
//
//    paths: '/api/v1/agents'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (sdam _SDAMAgentApis) agents(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "[AGENT] Get All Service Deployment Agents")
	result, resp, err := sdamAgentController.GetAgents()
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agentDeployApp handles requests which is used to deploy new application to agent
// identified by the given agentID.
//
//    paths: '/api/v1/agents/{agentID}/deploy'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (sdam _SDAMAgentApis) agentDeployApp(w http.ResponseWriter, req *http.Request, agentID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Deploy App")
	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := sdamAgentController.DeployApp(agentID, body)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agentInfoApps handles requests which is used to get information of all applications
// installed on agent identified by the given agentID.
//
//    paths: '/api/v1/agents/{agentID}/apps'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (sdam _SDAMAgentApis) agentInfoApps(w http.ResponseWriter, req *http.Request, agentID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Get Info Apps")
	result, resp, err := sdamAgentController.GetApps(agentID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agentInfoApp handles requests which is used to get information of application
// identified by the given appID.
//
//    paths: '/api/v1/agents/{agentID}/apps/{appID}'
//    method: GET
//    responses: if successful, 200 status code will be returned.
func (sdam _SDAMAgentApis) agentInfoApp(w http.ResponseWriter, req *http.Request, agentID string, appID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Get Info App")
	result, resp, err := sdamAgentController.GetApp(agentID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agentUpdateAppInfo handles requests related to updating the application with given yaml in body.
//
//    paths: '/api/v1/agents/{agentID}/apps/{appID}'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (sdam _SDAMAgentApis) agentUpdateAppInfo(w http.ResponseWriter, req *http.Request, agentID string, appID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Update App Info")
	body, err := common.GetBodyFromReq(req)
	if err != nil {
		common.MakeResponse(w, results.ERROR, nil, err)
		return
	}

	result, resp, err := sdamAgentController.UpdateAppInfo(agentID, appID, body)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agentDeleteApp handles requests related to delete application installed on agent
// identified by the given agentID.
//
//    paths: '/api/v1/agents/{agentID}/apps/{appID}'
//    method: DELETE
//    responses: if successful, 200 status code will be returned.
func (sdam _SDAMAgentApis) agentDeleteApp(w http.ResponseWriter, req *http.Request, agentID string, appID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Delete App")
	result, resp, err := sdamAgentController.DeleteApp(agentID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agentStartApp handles requests related to start application installed on agent
// identified by the given agentID.
//
//    paths: '/api/v1/agents/{agentID}/apps/{appID}/start'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (sdam _SDAMAgentApis) agentStartApp(w http.ResponseWriter, req *http.Request, agentID string, appID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Start App")
	result, resp, err := sdamAgentController.StartApp(agentID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agentStopApp handles requests related to stop application installed on agent
// identified by the given agentID.
//
//    paths: '/api/v1/agents/{agentID}/apps/{appID}/stop'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (sdam _SDAMAgentApis) agentStopApp(w http.ResponseWriter, req *http.Request, agentID string, appID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Stop App")
	result, resp, err := sdamAgentController.StopApp(agentID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}

// agentUpdateApp handles requests related to updating application installed on agent
// identified by the given agentID.
//
//    paths: '/api/v1/agents/{agentID}/apps/{appID}/update'
//    method: POST
//    responses: if successful, 200 status code will be returned.
func (sdam _SDAMAgentApis) agentUpdateApp(w http.ResponseWriter, req *http.Request, agentID string, appID string) {
	logger.Logging(logger.DEBUG, "[AGENT] Update App")
	result, resp, err := sdamAgentController.UpdateApp(agentID, appID)
	common.MakeResponse(w, result, common.ChangeToJson(resp), err)
}
