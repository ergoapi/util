// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package environ

import "strings"

// Proxy return proxy var
func Proxy() map[string]string {
	environ := map[string]string{}
	if value := envAnyCase("no_proxy"); value != "" {
		environ["no_proxy"] = value
		environ["NO_PROXY"] = value
	}
	if value := envAnyCase("http_proxy"); value != "" {
		environ["http_proxy"] = value
		environ["HTTP_PROXY"] = value
	}
	if value := envAnyCase("https_proxy"); value != "" {
		environ["https_proxy"] = value
		environ["HTTPS_PROXY"] = value
	}
	return environ
}

func envAnyCase(name string) (value string) {
	name = strings.ToUpper(name)
	if value := getenv(name); value != "" {
		return value
	}
	name = strings.ToLower(name)
	if value := getenv(name); value != "" {
		return value
	}
	return
}
