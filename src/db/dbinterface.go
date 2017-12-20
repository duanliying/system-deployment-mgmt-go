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
	"commons/logger"
	"db/mongo"
)

type Command interface {
	// AddAgent insert new Agent.
	AddAgent(host string, port string, status string) (map[string]interface{}, error)

	// UpdateAgentAddress updates ip,port of agent from db related to agent.
	UpdateAgentAddress(agent_id string, host string, port string) error

	// UpdateAgentStatus updates status of agent from db related to agent.
	UpdateAgentStatus(agent_id string, status string) error

	// GetAgent returns single document from db related to agent.
	GetAgent(agent_id string) (map[string]interface{}, error)

	// GetAgentByIP returns single document from db related to agent.
	GetAgentByIP(ip string) (map[string]interface{}, error)

	// GetAllAgents returns all documents from db related to agent.
	GetAllAgents() ([]map[string]interface{}, error)

	// GetAgentByAppID returns single document including specific app.
	GetAgentByAppID(agent_id string, app_id string) (map[string]interface{}, error)

	// AddAppToAgent add specific app to the target agent.
	AddAppToAgent(agent_id string, app_id string) error

	// DeleteAppFromAgent delete specific app from the target agent.
	DeleteAppFromAgent(agent_id string, app_id string) error

	// DeleteAgent delete single document from db related to agent.
	DeleteAgent(agent_id string) error

	// CreateGroup insert new Group.
	CreateGroup() (map[string]interface{}, error)

	// GetGroup returns single document from db related to group.
	GetGroup(group_id string) (map[string]interface{}, error)

	// GetAllGroups returns all documents from db related to group.
	GetAllGroups() ([]map[string]interface{}, error)

	// GetGroupMembers returns all agents who belong to the target group.
	GetGroupMembers(group_id string) ([]map[string]interface{}, error)

	// GetGroupMembersByAppID returns all agents including specific app on the target group.
	GetGroupMembersByAppID(group_id string, app_id string) ([]map[string]interface{}, error)

	// JoinGroup add specific agent to the target group.
	JoinGroup(group_id string, agent_id string) error

	// LeaveGroup delete specific agent from the target group.
	LeaveGroup(group_id string, agent_id string) error

	// DeleteGroup delete single document from db related to group.
	DeleteGroup(group_id string) error
}

type Closer interface {
	// Clean up session.
	Close()
}

type DBManager interface {
	Command
	Closer
}

type (
	DBConnection interface {
		Connect() (DBManager, error)
	}

	DBConnector struct{}
)

var mgoBuilder mongo.Builder

func init() {
	mgoBuilder = &mongo.MongoBuilder{}
}

// Connect establishes a new session to the database identified by the given url.
func (DBConnector) Connect() (DBManager, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// TODO: Should be updated to support different types of databases.
	url := "127.0.0.1:27017"

	err := mgoBuilder.Connect(url)
	if err != nil {
		return nil, err
	}

	dbManager, err := mgoBuilder.CreateDB()
	if err != nil {
		return nil, err
	}

	return dbManager, err
}
