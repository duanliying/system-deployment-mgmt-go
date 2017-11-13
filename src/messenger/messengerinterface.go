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
package messenger

type MessengerInterface interface {
	DeployApp(members []map[string]interface{}, data string) (respCode []int, respBody []string)
	InfoApp(members []map[string]interface{}, appId string) (respCode []int, respBody []string)
	DeleteApp(members []map[string]interface{}, appId string) (respCode []int, respBody []string)
	StartApp(members []map[string]interface{}, appId string) (respCode []int, respBody []string)
	StopApp(members []map[string]interface{}, appId string) (respCode []int, respBody []string)
	UpdateApp(members []map[string]interface{}, appId string) (respCode []int, respBody []string)
	InfoApps(member []map[string]interface{}) (respCode []int, respBody []string)
	UpdateAppInfo(member []map[string]interface{}, appId string, data string) (respCode []int, respBody []string)
}		