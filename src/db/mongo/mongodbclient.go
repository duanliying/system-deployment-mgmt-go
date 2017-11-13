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

// Package db/mongo implements some functions to use mgo which is MongoDB driver for Go.
// Service Deployment Agent Manager creates two collections.
// The first is used for managing a list of agents and second is used for managing a list of group.
package mongo

import (
	"commons/errors"
	"commons/logger"
	. "db/mongo/wrapper"
	"gopkg.in/mgo.v2/bson"
)

const (
	DB_NAME          = "DeploymentManagerDB"
	AGENT_COLLECTION = "AGENT"
	GROUP_COLLECTION = "GROUP"
)

type (
	Agent struct {
		ID     bson.ObjectId `bson:"_id,omitempty"`
		Host   string
		Port   string
		Apps   []string
		Status string
	}
	Group struct {
		ID      bson.ObjectId `bson:"_id,omitempty"`
		Members []string
	}
)

// convertToMap converts Agent object into a map.
func (agent Agent) convertToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":     agent.ID.Hex(),
		"host":   agent.Host,
		"port":   agent.Port,
		"apps":   agent.Apps,
		"status": agent.Status,
	}
}

// convertToMap converts Group object into a map.
func (group Group) convertToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":      group.ID.Hex(),
		"members": group.Members,
	}
}

// MongoDBManager provides persistence logic for "agent" and "group" collection.
type (
	Builder interface {
		Connect(url string) error
		CreateDB() (*MongoDBManager, error)
	}

	MongoBuilder struct {
		session Session
	}

	MongoDBManager struct {
		mgoSession Session
	}
)

var mgoDial Connection

func init() {
	mgoDial = MongoDial{}
}

// Connect establishes a new session to the database identified by the given url.
// If the connection is unsuccessful, this function returns DBConnectionError object.
func (builder *MongoBuilder) Connect(url string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Create a MongoDB Session.
	session, err := mgoDial.Dial(url)
	if err != nil {
		return errors.DBConnectionError{err.Error()}
	}

	builder.session = session
	return err
}

// CreateDB returns the MongoDBManager object used to interact with databases.
// If the session is nil, this function returns DBOperationError object.
// otherwise, this function returns an error as nil.
func (builder *MongoBuilder) CreateDB() (*MongoDBManager, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	if builder.session == nil {
		return nil, errors.DBOperationError{}
	}

	return &MongoDBManager{
		mgoSession: builder.session,
	}, nil
}

// Close terminates the session.
func (client *MongoDBManager) Close() {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	client.mgoSession.Close()
}

// getCollection returns a value representing the named collection.
func (client *MongoDBManager) getCollection(collectionName string) Collection {
	return client.mgoSession.DB(DB_NAME).C(collectionName)
}

// AddAgent inserts new agent to 'agent' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client *MongoDBManager) AddAgent(host string, port string, status string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	agent := Agent{
		ID:     bson.NewObjectId(),
		Host:   host,
		Port:   port,
		Status: status,
	}

	err := client.getCollection(AGENT_COLLECTION).Insert(agent)
	if err != nil {
		return nil, ConvertMongoError(err)
	}

	result := agent.convertToMap()
	return result, err
}

// UpdateAgentAddress updates ip,port of agent specified by agent_id parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client *MongoDBManager) UpdateAgentAddress(agent_id string, host string, port string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(agent_id) {
		err := errors.InvalidObjectId{agent_id}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(agent_id)}
	update := bson.M{"$set": bson.M{"host": host, "port": port}}
	err := client.getCollection(AGENT_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, "Failed to update address")
	}
	return err
}

// UpdateAgentStatus updates status of agent specified by agent_id parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client *MongoDBManager) UpdateAgentStatus(agent_id string, status string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(agent_id) {
		err := errors.InvalidObjectId{agent_id}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(agent_id)}
	update := bson.M{"$set": bson.M{"status": status}}
	err := client.getCollection(AGENT_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, "Failed to update status")
	}
	return err
}

// GetAgent returns single document specified by agent_id parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client *MongoDBManager) GetAgent(agent_id string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(agent_id) {
		err := errors.InvalidObjectId{agent_id}
		return nil, err
	}

	agent := Agent{}
	query := bson.M{"_id": bson.ObjectIdHex(agent_id)}
	err := client.getCollection(AGENT_COLLECTION).Find(query).One(&agent)
	if err != nil {
		return nil, ConvertMongoError(err, agent_id)
	}

	result := agent.convertToMap()
	return result, err
}

// GetAllAgents returns all documents from 'agent' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client *MongoDBManager) GetAllAgents() ([]map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	agents := []Agent{}
	err := client.getCollection(AGENT_COLLECTION).Find(nil).All(&agents)
	if err != nil {
		return nil, ConvertMongoError(err)
	}

	result := make([]map[string]interface{}, len(agents))
	for i, agent := range agents {
		result[i] = agent.convertToMap()
	}
	return result, err
}

// GetAgentByAppID returns single document specified by agent_id parameter.
// If successful, this function returns an error as nil.
// But if the target agent does not include the given app_id,
// an appropriate error will be returned.
func (client *MongoDBManager) GetAgentByAppID(agent_id string, app_id string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(agent_id) {
		err := errors.InvalidObjectId{agent_id}
		return nil, err
	}

	agent := Agent{}
	query := bson.M{"_id": bson.ObjectIdHex(agent_id), "apps": bson.M{"$in": []string{app_id}}}
	err := client.getCollection(AGENT_COLLECTION).Find(query).One(&agent)
	if err != nil {
		return nil, ConvertMongoError(err, agent_id)
	}

	result := agent.convertToMap()
	return result, err
}

// AddAppToAgent adds the specific app to the target agent.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client *MongoDBManager) AddAppToAgent(agent_id string, app_id string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(agent_id) {
		err := errors.InvalidObjectId{agent_id}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(agent_id)}
	update := bson.M{"$addToSet": bson.M{"apps": app_id}}
	err := client.getCollection(AGENT_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, agent_id)
	}
	return err
}

// DeleteAppFromAgent deletes the specific app from the target agent.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client *MongoDBManager) DeleteAppFromAgent(agent_id string, app_id string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(agent_id) {
		err := errors.InvalidObjectId{agent_id}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(agent_id)}
	update := bson.M{"$pull": bson.M{"apps": app_id}}
	err := client.getCollection(AGENT_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, agent_id)
	}
	return err
}

// DeleteAgent deletes single document from 'agent' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client *MongoDBManager) DeleteAgent(agent_id string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(agent_id) {
		err := errors.InvalidObjectId{agent_id}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(agent_id)}
	err := client.getCollection(AGENT_COLLECTION).Remove(query)
	if err != nil {
		return ConvertMongoError(err, agent_id)
	}
	return err
}

// CreateGroup inserts new Group to 'group' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client *MongoDBManager) CreateGroup() (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	group := Group{
		ID: bson.NewObjectId(),
	}

	err := client.getCollection(GROUP_COLLECTION).Insert(group)
	if err != nil {
		return nil, ConvertMongoError(err)
	}

	result := group.convertToMap()
	return result, err
}

// GetGroup returns single document specified by group_id parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client *MongoDBManager) GetGroup(group_id string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(group_id) {
		err := errors.InvalidObjectId{group_id}
		return nil, err
	}

	group := Group{}
	query := bson.M{"_id": bson.ObjectIdHex(group_id)}
	err := client.getCollection(GROUP_COLLECTION).Find(query).One(&group)
	if err != nil {
		return nil, ConvertMongoError(err, group_id)
	}

	result := group.convertToMap()
	return result, err
}

// GetAllGroups returns all documents from 'group' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client *MongoDBManager) GetAllGroups() ([]map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	groups := []Group{}
	err := client.getCollection(GROUP_COLLECTION).Find(nil).All(&groups)
	if err != nil {
		return nil, ConvertMongoError(err)
	}

	result := make([]map[string]interface{}, len(groups))
	for i, group := range groups {
		result[i] = group.convertToMap()
	}
	return result, err
}

// JoinGroup adds the specific agent to a list of group members.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client *MongoDBManager) JoinGroup(group_id string, agent_id string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(group_id) {
		err := errors.InvalidObjectId{group_id}
		return err
	}
	if !bson.IsObjectIdHex(agent_id) {
		err := errors.InvalidObjectId{agent_id}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(group_id)}
	update := bson.M{"$addToSet": bson.M{"members": agent_id}}
	err := client.getCollection(GROUP_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, group_id)
	}
	return err
}

// LeaveGroup deletes the specific agent from a list of group members.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client *MongoDBManager) LeaveGroup(group_id string, agent_id string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(group_id) {
		err := errors.InvalidObjectId{group_id}
		return err
	}
	if !bson.IsObjectIdHex(agent_id) {
		err := errors.InvalidObjectId{agent_id}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(group_id)}
	update := bson.M{"$pull": bson.M{"members": agent_id}}
	err := client.getCollection(GROUP_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, group_id)
	}
	return err
}

// GetGroupMembers returns all agents who belong to the target group.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client *MongoDBManager) GetGroupMembers(group_id string) ([]map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(group_id) {
		err := errors.InvalidObjectId{group_id}
		return nil, err
	}

	group, err := client.GetGroup(group_id)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(group["members"].([]string)))
	for i, agent_id := range group["members"].([]string) {
		var agent map[string]interface{}
		agent, err = client.GetAgent(agent_id)
		if err != nil {
			return nil, err
		}
		result[i] = agent
	}
	return result, err
}

// GetGroupMembersByAppID returns all agents including the app identified
// by the given appid on the target group.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client *MongoDBManager) GetGroupMembersByAppID(group_id string, app_id string) ([]map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(group_id) {
		err := errors.InvalidObjectId{group_id}
		return nil, err
	}

	group, err := client.GetGroup(group_id)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(group["members"].([]string)))
	for i, agent_id := range group["members"].([]string) {
		var agent map[string]interface{}
		agent, err = client.GetAgentByAppID(agent_id, app_id)
		if err != nil {
			return nil, err
		}
		result[i] = agent
	}
	return result, err
}

// DeleteGroup deletes single document specified by group_id parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client *MongoDBManager) DeleteGroup(group_id string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(group_id) {
		err := errors.InvalidObjectId{group_id}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(group_id)}
	err := client.getCollection(GROUP_COLLECTION).Remove(query)
	if err != nil {
		return ConvertMongoError(err, group_id)
	}
	return err
}
