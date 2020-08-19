// Copyright 2020 Huawei Technologies Co.,Ltd.
//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package sdkerr

import (
	"bytes"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type CredentialsTypeError struct {
	ErrorMessage string
}

func NewCredentialsTypeError(msg string) *CredentialsTypeError {
	c := &CredentialsTypeError{
		ErrorMessage: msg,
	}
	return c
}

func (c *CredentialsTypeError) Error() string {
	return fmt.Sprintf("{\"ErrorMessage\": \"%s\"}", c.ErrorMessage)
}

type ConnectionError struct {
	ErrorMessage string
}

func NewConnectionError(msg string) *ConnectionError {
	c := &ConnectionError{
		ErrorMessage: msg,
	}
	return c
}

func (c *ConnectionError) Error() string {
	return fmt.Sprintf("{\"ErrorMessage\": \"%s\"}", c.ErrorMessage)
}

type RequestTimeoutError struct {
	ErrorMessage string
}

func NewRequestTimeoutError(msg string) *RequestTimeoutError {
	rt := &RequestTimeoutError{
		ErrorMessage: msg,
	}
	return rt
}

func (rt *RequestTimeoutError) Error() string {
	return fmt.Sprintf("{\"ErrorMessage\": \"%s\"}", rt.ErrorMessage)
}

type ServiceResponseError struct {
	StatusCode   int    `json:"status_code"`
	RequestId    string `json:"request_id"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func NewServiceResponseError(resp *http.Response) *ServiceResponseError {
	sr := &ServiceResponseError{
		StatusCode: resp.StatusCode,
		RequestId:  resp.Header.Get("X-Request-Id"),
	}

	dataBuf := make(map[string]string)
	data, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		if err := jsoniter.Unmarshal(data, &dataBuf); err != nil {
			sr.ErrorMessage = string(data)
		}
	}
	if err := resp.Body.Close(); err == nil {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	}

	if sr.ErrorCode == "" && sr.ErrorMessage == "" {
		sr.ErrorCode = dataBuf["error_code"]
		sr.ErrorMessage = dataBuf["error_msg"]
	}

	if sr.ErrorCode == "" && sr.ErrorMessage == "" {
		sr.ErrorCode = dataBuf["code"]
		sr.ErrorMessage = dataBuf["message"]
	}

	if sr.ErrorCode == "" && sr.ErrorMessage == "" {
		for _, value := range dataBuf {
			buf := make(map[string]string)
			err := jsoniter.Unmarshal([]byte(value), &buf)
			if err == nil && buf["code"] != "" {
				sr.ErrorCode = buf["code"]
			}
			if err == nil && buf["message"] != "" {
				sr.ErrorMessage = buf["message"]
			}
		}
	}

	return sr
}

func (sr ServiceResponseError) Error() string {
	data, err := json.Marshal(sr)
	if err != nil {
		return fmt.Sprintf("{\"ErrorMessage\": \"%s\",\"ErrorCode\": \"%s\"}", sr.ErrorMessage, sr.ErrorCode)
	}
	return fmt.Sprintf(string(data))
}