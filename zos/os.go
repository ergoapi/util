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
	"os/exec"
	"os/user"
	"runtime"
	"strings"

	"github.com/acobaugh/osrelease"
	"github.com/ergoapi/util/file"
	"github.com/mitchellh/go-homedir"
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

// NotUnix  not macos and not linux
func NotUnix() bool {
	return runtime.GOOS != "linux" && runtime.GOOS != "darwin"
}

// IsContainer 是否是容器
func IsContainer() bool {
	return file.CheckFileExists("/.dockerenv")
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

// GetHomeDir 获取home目录
func GetHomeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		return "/root"
	}
	return home
}

// ExpandPath will parse `~` as user home dir path.
func ExpandPath(path string) string {
	path, _ = homedir.Expand(path)
	return path
}

// OSRelease get os release
func OSRelease() (map[string]string, error) {
	return osrelease.Read()
}

// IsDebian debian
func IsDebian() bool {
	os, err := osrelease.Read()
	if err != nil {
		return false
	}
	i, exist := os["ID"]
	if exist && (i == "debian" || i == "ubuntu") {
		return true
	}
	return false
}

func IsWsl() bool {
	// Return false if meet error
	cmd := exec.Command("cat", "/proc/version")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(output)), "microsoft")
}

func GetWindowsConfigHome() (string, error) {
	userCmd := exec.Command("wslvar", "LOCALAPPDATA")
	userOutput, err := userCmd.Output()
	if err != nil {
		return "", err
	}

	cmd := exec.Command("wslpath", string(userOutput))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.Trim(string(output), "\n"), nil
}
