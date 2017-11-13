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

import (
	"commons/url"
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

var doSomething func(method string, urls []string, dataOptional ...string) []httpResponse

func mockSendHttpRequest(method string, urls []string, dataOptional ...string) []httpResponse {
	return doSomething(method, urls, dataOptional...)
}

type httpFunc func(method string, urls []string, dataOptional ...string) []httpResponse

var oldSendHttpRequest httpFunc

type tearDown func(t *testing.T)

func setUp(t *testing.T) tearDown {
	oldSendHttpRequest = sendHttpRequest
	sendHttpRequest = mockSendHttpRequest

	return func(t *testing.T) {
		sendHttpRequest = oldSendHttpRequest
	}
}

const (
	deployApp = iota
	infoApp
	deleteApp
	startApp
	stopApp
	updateApp
	infoApps
	updateAppInfo
)

type sdamMsgr []func([]map[string]interface{}, string) ([]int, []string)

var messenger MessengerInterface

func init() {
	messenger = SdamMsgrImpl{}
}

func TestApiFunctions(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown(t)

	var group_members []map[string]interface{}
	members_cnt := 5

	for i := 0; i < members_cnt; i++ {
		localhost := "localhost"
		port := strconv.Itoa(8080 + i)
		group_members = append(group_members, map[string]interface{}{
			"host": localhost,
			"port": port,
		})
	}

	t.Run("DeployApp", func(t *testing.T) {
		data := "data"
		doSomething = func(method string, urls []string, dataOptional ...string) []httpResponse {
			if method != "POST" {
				t.Error()
			}
			if dataOptional[0] != data {
				t.Error("rrr", dataOptional, "rrr")
			}
			for i := 0; i < len(urls); i++ {
				expectedUrl := "http://" + group_members[i]["host"].(string) +
					":" + group_members[i]["port"].(string) + "/api/v1/deploy"
				if expectedUrl != urls[i] {
					t.Error()
				}
			}
			var respList []httpResponse
			for i := 0; i < len(urls); i++ {
				respList = append(respList, httpResponse{index: i, resp: nil, err: ""})
			}
			return respList
		}
		messenger.DeployApp(group_members, data)
	})

	t.Run("InfoApp", func(t *testing.T) {
		appId := "appId"
		doSomething = func(method string, urls []string, dataOptional ...string) []httpResponse {
			if method != "GET" {
				t.Error()
			}
			for i := 0; i < len(urls); i++ {
				expectedUrl := "http://" + group_members[i]["host"].(string) +
					":" + group_members[i]["port"].(string) + "/api/v1/apps/" + appId
				if expectedUrl != urls[i] {
					t.Error()
				}
			}
			var respList []httpResponse
			for i := 0; i < len(urls); i++ {
				respList = append(respList, httpResponse{index: i, resp: nil, err: ""})
			}
			return respList
		}
		messenger.InfoApp(group_members, appId)
	})

	t.Run("DeleteApp", func(t *testing.T) {
		appId := "appId"
		doSomething = func(method string, urls []string, dataOptional ...string) []httpResponse {
			if method != "DELETE" {
				t.Error()
			}
			for i := 0; i < len(urls); i++ {
				expectedUrl := "http://" + group_members[i]["host"].(string) +
					":" + group_members[i]["port"].(string) + "/api/v1/apps/" + appId
				if expectedUrl != urls[i] {
					t.Error()
				}
			}
			var respList []httpResponse
			for i := 0; i < len(urls); i++ {
				respList = append(respList, httpResponse{index: i, resp: nil, err: ""})
			}
			return respList
		}
		messenger.DeleteApp(group_members, appId)
	})

	t.Run("StartApp", func(t *testing.T) {
		appId := "appId"
		doSomething = func(method string, urls []string, dataOptional ...string) []httpResponse {
			if method != "POST" {
				t.Error()
			}
			for i := 0; i < len(urls); i++ {
				expectedUrl := "http://" + group_members[i]["host"].(string) +
					":" + group_members[i]["port"].(string) + "/api/v1/apps/" + appId +
					"/start"
				if expectedUrl != urls[i] {
					t.Error(expectedUrl, urls[i])
				}
			}
			var respList []httpResponse
			for i := 0; i < len(urls); i++ {
				respList = append(respList, httpResponse{index: i, resp: nil, err: ""})
			}
			return respList
		}
		messenger.StartApp(group_members, appId)
	})

	t.Run("StopApp", func(t *testing.T) {
		appId := "appId"
		doSomething = func(method string, urls []string, dataOptional ...string) []httpResponse {
			if method != "POST" {
				t.Error()
			}
			for i := 0; i < len(urls); i++ {
				expectedUrl := "http://" + group_members[i]["host"].(string) +
					":" + group_members[i]["port"].(string) + "/api/v1/apps/" + appId +
					"/stop"
				if expectedUrl != urls[i] {
					t.Error()
				}
			}
			var respList []httpResponse
			for i := 0; i < len(urls); i++ {
				respList = append(respList, httpResponse{index: i, resp: nil, err: ""})
			}
			return respList
		}
		messenger.StopApp(group_members, appId)
	})

	t.Run("UpdateApp", func(t *testing.T) {
		appId := "appId"
		doSomething = func(method string, urls []string, dataOptional ...string) []httpResponse {
			if method != "POST" {
				t.Error()
			}
			for i := 0; i < len(urls); i++ {
				expectedUrl := "http://" + group_members[i]["host"].(string) +
					":" + group_members[i]["port"].(string) + "/api/v1/apps/" + appId +
					"/update"
				if expectedUrl != urls[i] {
					t.Error()
				}
			}
			var respList []httpResponse
			for i := 0; i < len(urls); i++ {
				respList = append(respList, httpResponse{index: i, resp: nil, err: ""})
			}
			return respList
		}
		messenger.UpdateApp(group_members, appId)
	})

	t.Run("InfoApps", func(t *testing.T) {
		doSomething = func(method string, urls []string, dataOptional ...string) []httpResponse {
			if method != "GET" {
				t.Error()
			}
			for i := 0; i < len(urls); i++ {
				expectedUrl := "http://" + group_members[i]["host"].(string) +
					":" + group_members[i]["port"].(string) + "/api/v1/apps"
				if expectedUrl != urls[i] {
					t.Error()
				}
			}
			var respList []httpResponse
			for i := 0; i < len(urls); i++ {
				respList = append(respList, httpResponse{index: i, resp: nil, err: ""})
			}
			return respList
		}
		messenger.InfoApps(group_members)
	})

	t.Run("UpdateAppInfo", func(t *testing.T) {
		appId := "appId"
		data := "data"
		doSomething = func(method string, urls []string, dataOptional ...string) []httpResponse {
			if method != "POST" {
				t.Error()
			}
			if dataOptional[0] != data {
				t.Error()
			}
			for i := 0; i < len(urls); i++ {
				expectedUrl := "http://" + group_members[i]["host"].(string) +
					":" + group_members[i]["port"].(string) + "/api/v1/apps/" + appId
				if expectedUrl != urls[i] {
					t.Error()
				}
			}
			var respList []httpResponse
			for i := 0; i < len(urls); i++ {
				respList = append(respList, httpResponse{index: i, resp: nil, err: ""})
			}
			return respList
		}
		messenger.UpdateAppInfo(group_members, appId, data)
	})
}

func TestLen(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown(t)

	t.Run("TestLen", func(t *testing.T) {
		var test_input sortRespSlice
		expected := 5
		for i := 0; i < expected; i++ {
			test_input = append(test_input, httpResponse{i, nil, ""})
		}
		if expected != test_input.Len() {
			t.Error()
		}
	})
}

func TestLess(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown(t)

	t.Run("TestLess", func(t *testing.T) {
		var test_input sortRespSlice
		test_input = append(test_input, httpResponse{0, nil, ""})
		test_input = append(test_input, httpResponse{1, nil, ""})

		if test_input.Less(0, 1) != true {
			t.Error()
		}
		if test_input.Less(1, 0) != false {
			t.Error()
		}
		if test_input.Less(1, 1) != false {
			t.Error()
		}
	})
}

func TestSwap(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown(t)

	t.Run("TestSwap", func(t *testing.T) {
		var expectedSwapedList [2]httpResponse
		var test_input sortRespSlice
		test_input = append(test_input, httpResponse{0, nil, ""})
		test_input = append(test_input, httpResponse{1, nil, ""})

		expectedSwapedList[0] = test_input[1]
		expectedSwapedList[1] = test_input[0]

		test_input.Swap(0, 1)
		if test_input[0] != expectedSwapedList[0] {
			t.Error()
		}
		if test_input[1] != expectedSwapedList[1] {
			t.Error()
		}
	})
}

func TestChangeToReturnValue(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown(t)

	testResp := http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("Some data")),
		StatusCode: 200,
	}

	var respList []httpResponse

	respList = append(respList, httpResponse{index: 0, resp: nil, err: "errorMsg"})
	respList = append(respList, httpResponse{index: 1, resp: nil, err: "errorMsg"})
	respList = append(respList, httpResponse{index: 2, resp: &testResp, err: ""})

	expectedRespCode := [3]int{500, 500, 200}
	expectedRespBody := [3]string{`{"message":"errorMsg"}`, `{"message":"errorMsg"}`, "Some data"}

	t.Run("TestChangeToReturnValue", func(t *testing.T) {
		respCode, respBody := changeToReturnValue(respList)
		for i := 0; i < len(respList); i++ {
			if expectedRespCode[i] != respCode[i] || expectedRespBody[i] != respBody[i] {
				t.Error()
			}
		}
	})
}

func TestSetUrlList(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown(t)
	var group_members []map[string]interface{}
	var expectedUrlList []string

	appId := "appId"
	for i := 0; i < 3; i++ {
		localhost := "localhost"
		port := strconv.Itoa(8080 + i)

		group_members = append(group_members, map[string]interface{}{
			"host": localhost,
			"port": port,
		})
		expectedUrlList = append(expectedUrlList, "http://"+localhost+
			":"+port+"/api/v1"+"/"+appId+url.Start())
	}

	resultUrlList := setUrlList(group_members, "/", appId, url.Start())
	t.Run("TestSetUrlList", func(t *testing.T) {
		for i, _ := range group_members {
			if expectedUrlList[i] != resultUrlList[i] {
				t.Error()
			}
		}
	})
}

type _MockHttp struct{}

var doWrapperReturn func(req *http.Request) (*http.Response, error)

func (mockHttp _MockHttp) DoWrapper(req *http.Request) (*http.Response, error) {
	return doWrapperReturn(req)
}

func setUpHttpRequester() func() {
	var mockHttp _MockHttp
	httpInterface = mockHttp
	return func() {
		httpInterface = useHttp
	}
}

func TestHttpRequester(t *testing.T) {
	tearDown := setUpHttpRequester()
	defer tearDown()

	testURLs := []string{
		"http://0.0.0.0:8080",
		"http://0.0.0.1:8080",
		"http://0.0.0.2:8080",
	}

	doWrapperReturn = func(req *http.Request) (*http.Response, error) {
		switch req.URL.Scheme + "://" + req.URL.Host {
		case testURLs[2]:
			return &http.Response{}, errors.New("Error")
		default:
			return &http.Response{}, nil
		}
	}

	result := httpRequester("GET", testURLs)
	for i, val := range result {
		switch {
		case i != 2 && val.err != "":
			t.Error()
		case i == 2 && val.err != "Error":
			t.Error()
		}
	}
}

func TestHttpRequesterinBody(t *testing.T) {
	tearDown := setUpHttpRequester()
	defer tearDown()

	testURLs := []string{
		"http://0.0.0.0:8080",
		"http://0.0.0.1:8080",
		"http://0.0.0.2:8080",
		"",
	}

	doWrapperReturn = func(req *http.Request) (*http.Response, error) {
		switch req.URL.Scheme + "://" + req.URL.Host {
		case testURLs[2]:
			return &http.Response{}, errors.New("Error")
		default:
			return &http.Response{}, nil
		}
	}
	result := httpRequester("GET", testURLs, "testData")
	for i, val := range result {
		switch {
		case i != 2 && val.err != "":
			t.Error()
		case i == 2 && val.err != "Error":
			t.Error()
		}
	}
}
