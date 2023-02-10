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

package file

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/ergoapi/util/exstr"

	recursiveCopy "github.com/otiai10/copy"
)

var (
	// FilePath filepath
	FilePath = "."
)

func init() {
	FilePath, _ = filepath.Abs(".")
}

// RealPath get an absolute path
func RealPath(path string, addSlash ...bool) (realPath string) {
	if !filepath.IsAbs(path) {
		path = FilePath + "/" + path
	}
	realPath, _ = filepath.Abs(path)
	realPath = pathAddSlash(filepath.ToSlash(realPath), addSlash...)

	return
}

// WorkDirPath program directory path
func WorkDirPath(addSlash ...bool) (path string) {
	ePath, err := os.Executable()
	if err != nil {
		ePath = FilePath
	}
	path = pathAddSlash(filepath.Dir(ePath), addSlash...)
	return
}

func pathAddSlash(path string, addSlash ...bool) string {
	if len(addSlash) > 0 && addSlash[0] && !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return path
}

// Rmdir rmdir,support to keep the current directory
func Rmdir(path string, notIncludeSelf ...bool) (ok bool) {
	realPath := RealPath(path)
	err := os.RemoveAll(realPath)
	ok = err == nil
	if ok && len(notIncludeSelf) > 0 && notIncludeSelf[0] {
		_ = os.Mkdir(path, os.ModePerm)
	}
	return
}

// file ./test/dir/xxx.txt if dir ./test/dir not exist, create it
func MkFileFullPathDir(fileName string) error {
	localDir := filepath.Dir(fileName)
	if err := Mkdir(localDir); err != nil {
		return errors.Errorf("failed to create local dir %s: %v", localDir, err)
	}
	return nil
}

func Mkdir(dirName string) error {
	return os.MkdirAll(dirName, 0755)
}

func MkDirs(dirs ...string) error {
	if len(dirs) == 0 {
		return nil
	}
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return errors.Errorf("failed to create %s, %v", dir, err)
		}
	}
	return nil
}

func MkTmpdir(dir string) (string, error) {
	return os.MkdirTemp(dir, ".dtmp-")
}

func MkTmpFile(path string) (*os.File, error) {
	return os.CreateTemp(path, ".ftmp-")
}

// Size file size
func Size(path string) float64 {
	fi, err := os.Stat(path)
	if err == nil {
		bs := float64(fi.Size())
		return bs
	}
	return 0
}

// Size2Str file size
func Size2Str(path string) string {
	fi, err := os.Stat(path)
	if err == nil {
		bs := float64(fi.Size())
		kbs := bs / 1024.0
		mbs := kbs / 1024.0
		if mbs < 1024.0 {
			return fmt.Sprintf("%v M", mbs)
		}
		gbs := mbs / 1024.0
		if gbs < 1024.0 {
			return fmt.Sprintf("%v G", gbs)
		}
		tbs := gbs / 1024.0
		return fmt.Sprintf("%v T", tbs)
	}
	return ""
}

// CheckFileExists check file exist
func CheckFileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// WriteToFile 写文件
func WriteToFile(filePath string, data []byte) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0600)
}

// Writefile 写文件
func Writefile(logpath, msg string, truncate bool) error {
	filemode := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	if truncate {
		filemode = os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	}
	file, err := os.OpenFile(logpath, filemode, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(msg)
	write.Flush()
	return nil
}

// WritefileWithLine 换行
func WritefileWithLine(logpath, msg string) error {
	file, err := os.OpenFile(logpath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(msg + "\n")
	write.Flush()
	return nil
}

// DirIsEmpty 验证目录是否为空
func DirIsEmpty(dir string) bool {
	infos, err := os.ReadDir(dir)
	if len(infos) == 0 || err != nil {
		return true
	}
	return false
}

// SearchFileBody 搜索文件中是否含有指定字符串
func SearchFileBody(filename, searchStr string) bool {
	body, _ := os.ReadFile(filename)
	return strings.Contains(string(body), searchStr)
}

// IsHaveFile 指定目录是否含有文件
// .开头文件除外
func IsHaveFile(path string) bool {
	files, _ := os.ReadDir(path)
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			return true
		}
	}
	return false
}

// SearchFile 搜索指定目录是否有指定文件，指定搜索目录层数，-1为全目录搜索
func SearchFile(pathDir, name string, level int) bool {
	if level == 0 {
		return false
	}
	files, _ := os.ReadDir(pathDir)
	var dirs []fs.DirEntry
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file)
			continue
		}
		if file.Name() == name {
			return true
		}
	}
	if level == 1 {
		return false
	}
	for _, dir := range dirs {
		ok := SearchFile(path.Join(pathDir, dir.Name()), name, level-1)
		if ok {
			return ok
		}
	}
	return false
}

// CheckFileExistsWithSuffix 指定目录是否含有指定后缀的文件
func CheckFileExistsWithSuffix(pathDir, suffix string) bool {
	files, _ := os.ReadDir(pathDir)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), suffix) {
			return true
		}
	}
	return false
}

// RemoveFiles 删除文件
func RemoveFiles(path string) bool {
	if err := os.RemoveAll(path); err != nil {
		return false
	}
	return true
}

func ReadLines(fileName string) ([]string, error) {
	var lines []string
	if !CheckFileExists(fileName) {
		return nil, errors.New("no such file")
	}
	file, err := os.Open(filepath.Clean(fileName))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	br := bufio.NewReader(file)
	for {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		lines = append(lines, string(line))
	}
	return lines, nil
}

// ReadAll read file content
func ReadAll(fileName string) ([]byte, error) {
	// step1：check file exist
	if !CheckFileExists(fileName) {
		return nil, errors.New("no such file")
	}
	// step2：open file
	file, err := os.Open(filepath.Clean(fileName))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// step3：read file content
	content, err := os.ReadFile(filepath.Clean(fileName))
	if err != nil {
		return nil, err
	}

	return content, nil
}

// ReadFileOneLine 读取文件一行
func ReadFileOneLine(fileName string) string {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return ""
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	line, err := buf.ReadString('\n')
	if err != nil {
		return ""
	}
	return line
}

// ReadFile reads a file with a given limit
func ReadFile(path string, limit int64) ([]byte, error) {
	if limit <= 0 {
		return os.ReadFile(path)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	st, err := f.Stat()
	if err != nil {
		return nil, err
	}

	size := st.Size()
	if limit > 0 && size > limit {
		size = limit
	}

	buf := bytes.NewBuffer(nil)
	buf.Grow(int(size))
	_, err = io.Copy(buf, io.LimitReader(f, limit))

	return buf.Bytes(), err
}

// Copy copies a file to a destination path
func Copy(sourcePath string, targetPath string, overwrite bool) error {
	if overwrite {
		return recursiveCopy.Copy(sourcePath, targetPath)
	}

	var err error

	// Convert to absolute path
	sourcePath, err = filepath.Abs(sourcePath)
	if err != nil {
		return err
	}

	// Convert to absolute path
	targetPath, err = filepath.Abs(targetPath)
	if err != nil {
		return err
	}

	return filepath.Walk(sourcePath, func(nextSourcePath string, fileInfo os.FileInfo, err error) error {
		nextTargetPath := filepath.Join(targetPath, strings.TrimPrefix(nextSourcePath, sourcePath))
		if fileInfo == nil {
			return nil
		}

		if !fileInfo.Mode().IsRegular() {
			return nil
		}

		if fileInfo.IsDir() {
			_ = os.MkdirAll(nextTargetPath, os.ModePerm)
			return Copy(nextSourcePath, nextTargetPath, overwrite)
		}

		_, statErr := os.Stat(nextTargetPath)
		if statErr != nil {
			return recursiveCopy.Copy(nextSourcePath, nextTargetPath)
		}

		return nil
	})
}

// IsRecursiveSymlink checks if the provided non-resolved file info
// is a recursive symlink
func IsRecursiveSymlink(f os.FileInfo, symlinkPath string) bool {
	// check if recursive symlink
	if f.Mode()&os.ModeSymlink == os.ModeSymlink {
		resolvedPath, err := filepath.EvalSymlinks(symlinkPath)
		if err != nil || strings.HasPrefix(symlinkPath, filepath.ToSlash(resolvedPath)) {
			return true
		}
	}

	return false
}

// IsDir 是否为目录
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// DirFilesList 获取目录下的文件列表
func DirFilesList(sourcePath string, include, exclude []string) (files []string, err error) {
	// bfs遍历文件夹
	var dirs []string
	dirs = append(dirs, sourcePath)
	for len(dirs) > 0 {
		dirName := dirs[0]
		dirs = dirs[1:]

		fileInfos, err := os.ReadDir(dirName)
		if err != nil {
			return nil, err
		}

		for _, f := range fileInfos {
			fileName := dirName + "/" + f.Name()
			if f.IsDir() { // 目录
				dirs = append(dirs, fileName)
			} else if f.Type().IsRegular() { // 普通文件
				fileName = fileName[len(sourcePath)+1:]
				files = append(files, fileName)
			} else {
				// 可能是软连接等
				continue
			}
		}
	}

	if len(include) > 0 {
		var okfile []string
		for _, i := range include {
			ifiles := matchPattern(files, i, true)
			okfile = append(okfile, ifiles...)
		}
		files = exstr.DuplicateStrElement(okfile)
	}
	if len(exclude) > 0 {
		for _, e := range exclude {
			files = matchPattern(files, e, false)
		}
	}

	return files, nil
}

func matchPattern(strs []string, pattern string, include bool) []string {
	res := make([]string, 0)
	re := regexp.MustCompile(pattern)
	for _, s := range strs {
		match := re.MatchString(s)
		if !include {
			match = !match
		}
		if match {
			res = append(res, s)
		}
	}
	return res
}

// GetTempDir 返回用于临时文件的目录.
func GetTempDir() string {
	return os.TempDir()
}

// CountDirFiles reutrns # of files under a directory.
func CountDirFiles(dirName string) int {
	if !IsDir(dirName) {
		return 0
	}
	var count int
	err := filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		count++
		return nil
	})
	if err != nil {
		return 0
	}
	return count
}
