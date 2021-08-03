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

package zos

import (
	"os"
	"os/user"
	"runtime"
)

// IsMacOS is Mac OS
func IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

// IsLinux is linux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// IsUnix macos or linux
func IsUnix() bool {
	return runtime.GOOS == "linux" || runtime.GOOS == "darwin"
}

// IsContainer 是否是容器
func IsContainer() bool {
	_, err := os.Stat("/.dockerenv")
	return err == nil || os.IsExist(err)
}

// GetUserName 获取当前系统登录用户
func GetUserName() string {
	user, err := user.Current()
	if err != nil {
		return ""
	}
	return user.Username
}

// GetUser 获取当前系统登录用户
func GetUser() *user.User {
	user, err := user.Current()
	if err != nil {
		return nil
	}
	return user
}

func GetHostnames() []string {
	host, err := os.Hostname()
	if err != nil {
		return nil
	}
	return []string{host}
}

func GetHostname() string {
	hosts := GetHostnames()
	if len(hosts) == 0 {
		return "unknow"
	}
	return hosts[0]
}

func GetOS() string {
	return runtime.GOOS
}
