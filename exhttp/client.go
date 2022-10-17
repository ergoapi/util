//  Copyright (c) 2021. The EFF Team Authors.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  See the License for the specific language governing permissions and
//  limitations under the License.

package exhttp

import (
	"net/http"
	"time"
)

type ExHttpClient struct {
	*http.Client
}

var httpClient *ExHttpClient

func init() {
	httpClient = &ExHttpClient{http.DefaultClient}
	httpClient.Timeout = time.Second * 60
	httpClient.Transport = &http.Transport{
		TLSHandshakeTimeout:   time.Second * 5,
		IdleConnTimeout:       time.Second * 10,
		ResponseHeaderTimeout: time.Second * 10,
		ExpectContinueTimeout: time.Second * 20,
		Proxy:                 http.ProxyFromEnvironment,
	}
}

func GetHttpClient() *ExHttpClient {
	c := *httpClient
	return &c
}
