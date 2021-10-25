//  Copyright (c) 2020. The EFF Team Authors.
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

package environ

import (
	"os"
	"strings"
)

// function used to expand environment variables.
var getenv = os.Getenv

// GetEnv 获取环境变量
func GetEnv(envstr string, fallback ...string) string {
	e := getenv(envstr)
	if e == "" && len(fallback) > 0 {
		e = fallback[0]
	}
	return e
}

// Environ 类似 os.Environ, 返回key-value map[string]string.
func Environ() map[string]string {
	envList := os.Environ()
	envMap := make(map[string]string, len(envList))

	for _, str := range envList {
		nodes := strings.SplitN(str, "=", 2)

		if len(nodes) < 2 {
			envMap[nodes[0]] = ""
		} else {
			envMap[nodes[0]] = nodes[1]
		}
	}
	return envMap
}

// Expand is a helper function to expand the PATH parameter in
// the pipeline environment.
func Expand(env map[string]string) map[string]string {
	c := map[string]string{}
	for k, v := range env {
		c[k] = v
	}
	if path := c["PATH"]; path != "" {
		c["PATH"] = os.Expand(path, getenv)
	}
	return c
}
