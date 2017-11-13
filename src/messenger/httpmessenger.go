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

// Package messenger provides abstracted interfaces for HTTP messages,
// including requests and responses.
package messenger

import (
	"bytes"
	"commons/logger"
	"commons/url"
	"net/http"
	"sort"
	"sync"
)

func init() {
	sendHttpRequest = httpRequester
	httpInterface = useHttp
}

var sendHttpRequest func(method string, urls []string, dataOptional ...string) []httpResponse

// A httpResponse represents an HTTP response received from remote device.
type httpResponse struct {
	index int
	resp  *http.Response
	err   string
}
type sortRespSlice []httpResponse

type SdamMsgrImpl struct{}

// DeployApp make a url using /api/v1/deploy and send a HTTP(POST) request.
func (SdamMsgrImpl) DeployApp(members []map[string]interface{}, data string) (respCode []int, respBody []string) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	urls := setUrlList(members, url.Deploy())
	respList := sendHttpRequest("POST", urls, data)
	return changeToReturnValue(respList)
}

// InfoApp make a url using /api/v1/apps/{appId} and send a HTTP(GET) request.
func (SdamMsgrImpl) InfoApp(members []map[string]interface{}, appId string) (respCode []int, respBody []string) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	urls := setUrlList(members, url.Apps(), "/", appId)
	respList := sendHttpRequest("GET", urls)
	return changeToReturnValue(respList)
}

// DeleteApp make a url using /api/v1/apps/{appId} and send a HTTP(DELETE) request.
func (SdamMsgrImpl) DeleteApp(members []map[string]interface{}, appId string) (respCode []int, respBody []string) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	urls := setUrlList(members, url.Apps(), "/", appId)
	respList := sendHttpRequest("DELETE", urls)
	return changeToReturnValue(respList)
}

// StartApp make a url using /api/v1/apps/{appId}/start and send a HTTP(POST) request.
func (SdamMsgrImpl) StartApp(members []map[string]interface{}, appId string) (respCode []int, respBody []string) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	urls := setUrlList(members, url.Apps(), "/", appId, url.Start())
	respList := sendHttpRequest("POST", urls)
	return changeToReturnValue(respList)
}

// StopApp make a url using /api/v1/apps/{appId}/stop and send a HTTP(POST) request.
func (SdamMsgrImpl) StopApp(members []map[string]interface{}, appId string) (respCode []int, respBody []string) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	urls := setUrlList(members, url.Apps(), "/", appId, url.Stop())
	respList := sendHttpRequest("POST", urls)
	return changeToReturnValue(respList)
}

// UpdateApp make a url using /api/v1/apps/{appId}/update and send a HTTP(POST) request.
func (SdamMsgrImpl) UpdateApp(members []map[string]interface{}, appId string) (respCode []int, respBody []string) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	urls := setUrlList(members, url.Apps(), "/", appId, url.Update())
	respList := sendHttpRequest("POST", urls)
	return changeToReturnValue(respList)
}

// InfoApps make a url using /api/v1/apps and send a HTTP(GET) request.
func (SdamMsgrImpl) InfoApps(member []map[string]interface{}) (respCode []int, respBody []string) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	urls := setUrlList(member, url.Apps())
	respList := sendHttpRequest("GET", urls)
	return changeToReturnValue(respList)
}

// UpdateAppInfo make a url using /api/v1/apps/{appId} and send a HTTP(POST) request.
func (SdamMsgrImpl) UpdateAppInfo(member []map[string]interface{}, appId string, data string) (respCode []int, respBody []string) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	urls := setUrlList(member, url.Apps(), "/", appId)
	respList := sendHttpRequest("POST", urls, data)
	return changeToReturnValue(respList)
}

// Len returns length of httpResponse.
func (arr sortRespSlice) Len() int {
	return len(arr)
}

// Less returns whether the its first argument compares less than the second.
func (arr sortRespSlice) Less(i, j int) bool {
	return arr[i].index < arr[j].index
}

// Swap exchange its first argument with the second.
func (arr sortRespSlice) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

// _HTTPInterface is an interface including the function to execute a single HTTP transaction.
type _HTTPInterface interface {
	DoWrapper(req *http.Request) (*http.Response, error)
}
type _UseHttp struct{}

var httpInterface _HTTPInterface
var useHttp _UseHttp

// DoWrapper calls Do function of http.DefaultClient to send an HTTP request.
func (useHttp _UseHttp) DoWrapper(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

// httpRequester make a new request given a method, url, and optional body.
// and send a request to target device.
// A list of httpResponse structure will be returned by this function.
func httpRequester(method string, urls []string, dataOptional ...string) []httpResponse {
	var wg sync.WaitGroup
	wg.Add(len(urls))

	respChannel := make(chan httpResponse, len(urls))
	for i := range urls {
		go func(idx int) {
			logger.Logging(logger.DEBUG, "sending http request:", urls[idx])

			var err error
			var req *http.Request
			var resp httpResponse

			resp.index = idx

			switch len(dataOptional) {
			case 0:
				req, err = http.NewRequest(method, urls[idx], bytes.NewBuffer(nil))
			case 1:
				req, err = http.NewRequest(method, urls[idx], bytes.NewBuffer([]byte(dataOptional[0])))
			}

			if err != nil {
				resp.resp = nil
				resp.err = err.Error()
				respChannel <- resp
			} else {
				resp.resp, err = httpInterface.DoWrapper(req)
				if err != nil {
					resp.err = err.Error()
				} else {
					resp.err = ""
				}
				respChannel <- resp
			}
			defer wg.Done()
		}(i)
	}
	wg.Wait()

	var respList []httpResponse
	for range urls {
		respList = append(respList, <-respChannel)
	}
	sort.Sort(sortRespSlice(respList))
	return respList
}

// changeToReturnValue parses a response code and body from httpResponse structure.
func changeToReturnValue(respList []httpResponse) (respCode []int, respBody []string) {
	var buf bytes.Buffer

	for i := 0; i < len(respList); i++ {
		buf.Reset()
		if respList[i].resp == nil {
			message := `{"message":"` + respList[i].err + `"}`
			respBody = append(respBody, message)
			respCode = append(respCode, 500)
		} else {
			buf.ReadFrom(respList[i].resp.Body)
			respBody = append(respBody, buf.String())
			respCode = append(respCode, respList[i].resp.StatusCode)
		}
	}
	return respCode, respBody
}

// setUrlList make a list of urls that can be used to send a http request.
func setUrlList(members []map[string]interface{}, api_parts ...string) (urls []string) {
	var httpTag string = "http://"
	var full_url bytes.Buffer

	for i := range members {
		full_url.Reset()
		full_url.WriteString(httpTag + members[i]["host"].(string) +
			":" + members[i]["port"].(string) +
			url.Base())
		for _, api_part := range api_parts {
			full_url.WriteString(api_part)
		}
		urls = append(urls, full_url.String())
	}
	return urls
}
