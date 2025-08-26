// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package exhttp

type ClientOptionFunc func(*Client) error

// WithDevMode 开启调试模式, 默认情况下会读取DEBUG变量
func WithDevMode() ClientOptionFunc {
	return func(c *Client) error {
		return c.setDebug()
	}
}

// WithDumpAll 开启请求/响应的详细日志
func WithDumpAll() ClientOptionFunc {
	return func(c *Client) error {
		return c.setDumpAll()
	}
}

// WithoutProxy 禁用代理, 默认情况下会读取HTTP_PROXY/HTTPS_PROXY/http_proxy/https_proxy变量
func WithoutProxy() ClientOptionFunc {
	return func(c *Client) error {
		return c.setDisableProxy()
	}
}

// WithUserAgent 设置请求的User-Agent
func WithUserAgent(ua string) ClientOptionFunc {
	return func(c *Client) error {
		if ua == "" {
			ua = "github.com/ergoapi/util/exhttp"
		}
		return c.setReqUserAgent(ua)
	}
}
