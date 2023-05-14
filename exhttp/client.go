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
	"github.com/imroc/req/v3"
)

type Client struct {
	*req.Client
}

func GetClientRequest(opts ...ClientOptionFunc) (*req.Request, error) {
	c, err := GetClient(opts...)
	if err != nil {
		return nil, err
	}
	return c.R(), nil
}

func GetClient(opts ...ClientOptionFunc) (*Client, error) {
	c := &Client{req.C().SetLogger(nil)}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Client) setDebug() error {
	c.EnableDebugLog()
	return nil
}

func (c *Client) setDumpAll() error {
	c.EnableDumpAll()
	return nil
}

func (c *Client) setDisableProxy() error {
	c.SetProxy(nil)
	return nil
}

func (c *Client) setReqUserAgent(ua string) error {
	c.SetUserAgent(ua)
	return nil
}
