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
	"os/exec"
	"os/user"
	"runtime"
	"strings"

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
	h := HostInfo()
	return h.Container.EqualBool(true)
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

func GetHostname() string {
	h := HostInfo()
	return h.Hostname
}

func GetOS() string {
	h := HostInfo()
	return h.OS
}

func GetDistro() string {
	h := HostInfo()
	return h.Distro
}

func GetDistroVersion() string {
	h := HostInfo()
	return h.DistroVersion
}

func GetDistroCodeName() string {
	h := HostInfo()
	return h.DistroCodeName
}

func GetOSVersion() string {
	h := HostInfo()
	return h.OSVersion
}

// GetHomeDir 获取home目录
func GetHomeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		return "/root"
	}
	if home == "/" || home == "" {
		return "/tmp"
	}
	return home
}

// ExpandPath will parse `~` as user home dir path.
func ExpandPath(path string) string {
	path, _ = homedir.Expand(path)
	return path
}

// IsDebian debian
func IsDebian() bool {
	h := HostInfo()
	return h.Distro == "debian"
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
