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

package excmd

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
)

//RunCmd is exec on os ,no return
func RunCmd(name string, arg ...string) error {
	cmd := exec.Command(name, arg[:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

//RunCmdRes is exec on os , return result
func RunCmdRes(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg[:]...)
	cmd.Stdin = os.Stdin
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func downloadCmd(url string) string {
	ishttp := isURL(url)
	var c = ""
	if ishttp {
		param := "--no-check-certificate"
		c = fmt.Sprintf(" wget -c %s %s", param, url)
	}
	return c
}

func isURL(u string) bool {
	if uu, err := url.Parse(u); err == nil && uu != nil && uu.Host != "" {
		return true
	}
	return false
}

// DownloadFile 下载文件
func DownloadFile(url string, location string) {
	dwncmd := downloadCmd(url)
	RunCmd("/bin/sh", "-c", "mkdir -p /tmp/ysicing && cd /tmp/ysicing && "+dwncmd)
	RunCmd("/bin/sh", "-c", "cp -a /tmp/ysicing/* "+location)
	RunCmd("/bin/sh", "-c", "rm -rf /tmp/ysicing")
}

// CheckBin 检查二进制是否存在
func CheckBin(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
