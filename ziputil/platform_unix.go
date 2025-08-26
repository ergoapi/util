// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

//go:build !windows
// +build !windows

package ziputil

import (
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
)

// chownR 递归修改文件所有者（Unix系统）
func chownR(path string, uid, gid int) error {
	return filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
		if err == nil {
			err = os.Chown(name, uid, gid)
		}
		return err
	})
}

// getFileOwnership 获取文件所有者信息
func getFileOwnership() (uid, gid int) {
	return os.Getuid(), os.Getgid()
}

// setFilePermissions 设置文件权限
func setFilePermissions(path string, info os.FileInfo) error {
	if err := os.Chmod(path, info.Mode()); err != nil {
		return errors.Wrap(err, "failed to set file permissions")
	}

	// Unix系统上尝试设置所有者
	uid, gid := getFileOwnership()
	return os.Chown(path, uid, gid)
}

// ensureDir 确保目录存在并设置正确的权限
func ensureDir(dirName string, perm os.FileMode) error {
	if perm == 0 {
		perm = 0755
	}

	err := os.MkdirAll(dirName, perm)
	if err != nil && !os.IsExist(err) {
		return errors.Wrap(err, "failed to create directory")
	}

	// Unix系统上设置所有者
	uid, gid := getFileOwnership()
	// 忽略chown错误，因为某些文件系统可能不支持
	_ = chownR(dirName, uid, gid)

	return nil
}
