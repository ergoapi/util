//go:build windows
// +build windows

package ziputil

import (
	"os"

	"github.com/cockroachdb/errors"
)

// chownR 在Windows上不执行任何操作
func chownR(path string, uid, gid int) error {
	// Windows不支持Unix风格的文件所有权
	return nil
}

// getFileOwnership 在Windows上返回默认值
func getFileOwnership() (uid, gid int) {
	return -1, -1
}

// setFilePermissions 在Windows上设置文件权限
func setFilePermissions(path string, info os.FileInfo) error {
	// Windows权限模型不同，只设置基本权限
	return os.Chmod(path, info.Mode())
}

// ensureDir 确保目录存在（Windows版本）
func ensureDir(dirName string, perm os.FileMode) error {
	if perm == 0 {
		perm = 0755
	}
	
	err := os.MkdirAll(dirName, perm)
	if err != nil && !os.IsExist(err) {
		return errors.Wrap(err, "failed to create directory")
	}
	
	return nil
}