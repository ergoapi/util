// Copyright 2022 ysicing
// Copyright 2022 The envd Authors
// Copyright 2022 mateors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ziputil

import (
	"archive/zip"
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
)

// 缓冲区大小
const bufferSize = 32 * 1024 // 32KB

// symlinkHandler 处理软链接的辅助结构
type symlinkHandler struct {
	processedDirs map[string]bool
	basePath      string
}

// newSymlinkHandler 创建新的软链接处理器
func newSymlinkHandler(basePath string) *symlinkHandler {
	absPath, _ := filepath.Abs(basePath)
	handler := &symlinkHandler{
		processedDirs: make(map[string]bool),
		basePath:      absPath,
	}
	handler.processedDirs[absPath] = true
	return handler
}

// isProcessed 检查目录是否已处理
func (h *symlinkHandler) isProcessed(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	return h.processedDirs[absPath]
}

// markProcessed 标记目录为已处理
func (h *symlinkHandler) markProcessed(path string) {
	if absPath, err := filepath.Abs(path); err == nil {
		h.processedDirs[absPath] = true
	}
}

// MakeZip 压缩目录或文件到zip归档
func MakeZip(inputPath, outputFile string) error {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return errors.Wrap(err, "input path does not exist")
	}

	files, err := fileList(inputPath)
	if err != nil {
		return errors.Wrap(err, "unable to list files")
	}

	return ZipFiles(outputFile, files)
}

// fileList 获取文件列表（包括软链接处理）
func fileList(fileDirectory string) ([]string, error) {
	var files []string
	handler := newSymlinkHandler(fileDirectory)

	err := filepath.Walk(fileDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 处理软链接
		if info.Mode()&os.ModeSymlink != 0 {
			return handleSymlinkInList(path, &files, handler)
		}

		// 处理常规文件和目录
		if !info.IsDir() {
			files = append(files, path)
		} else {
			handler.markProcessed(path)
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to walk directory")
	}

	return files, nil
}

// handleSymlinkInList 处理文件列表中的软链接
func handleSymlinkInList(path string, files *[]string, handler *symlinkHandler) error {
	realPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		// 软链接目标不存在，跳过
		return nil
	}

	realInfo, err := os.Stat(realPath)
	if err != nil {
		// 无法访问目标，跳过
		return nil
	}

	// 如果是文件，添加原路径
	if !realInfo.IsDir() {
		*files = append(*files, path)
		return nil
	}

	// 如果是目录，检查是否已处理
	if handler.isProcessed(realPath) {
		return nil
	}

	handler.markProcessed(realPath)

	// 遍历目录
	return filepath.Walk(realPath, func(subPath string, subInfo os.FileInfo, err error) error {
		if err != nil || subPath == realPath {
			return err
		}

		if subInfo.IsDir() {
			if handler.isProcessed(subPath) {
				return filepath.SkipDir
			}
			handler.markProcessed(subPath)
		} else {
			*files = append(*files, subPath)
		}

		return nil
	})
}

// Unzip 解压缩zip文件到目标目录
func Unzip(src string, dest string) ([]string, error) {
	var filenames []string

	// 确保目标目录存在
	if err := ensureDir(dest, 0755); err != nil {
		return filenames, err
	}

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, errors.Wrap(err, "failed to open zip file")
	}
	defer r.Close()

	for _, f := range r.File {
		if err := extractFile(f, dest, &filenames); err != nil {
			return filenames, err
		}
	}

	return filenames, nil
}

// extractFile 解压单个文件
func extractFile(f *zip.File, dest string, filenames *[]string) error {
	// 构建目标路径
	fpath := filepath.Join(dest, f.Name)

	// 防止ZipSlip攻击
	if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
		return errors.Newf("%s: illegal file path (zip slip attack prevented)", fpath)
	}

	*filenames = append(*filenames, fpath)

	// 获取文件信息
	fileInfo := f.FileInfo()

	// 创建目录
	if fileInfo.IsDir() {
		return ensureDir(fpath, fileInfo.Mode())
	}

	// 确保父目录存在
	if err := ensureDir(filepath.Dir(fpath), 0755); err != nil {
		return err
	}

	// 创建文件
	return extractFileContent(f, fpath)
}

// extractFileContent 提取文件内容
func extractFileContent(f *zip.File, fpath string) error {
	rc, err := f.Open()
	if err != nil {
		return errors.Wrap(err, "failed to open file in zip")
	}
	defer rc.Close()

	// 创建输出文件，保留原始权限
	outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return errors.Wrap(err, "failed to create output file")
	}
	defer outFile.Close()

	// 使用缓冲写入器提高性能
	writer := bufio.NewWriterSize(outFile, bufferSize)

	// 复制内容
	if _, err := io.Copy(writer, rc); err != nil {
		return errors.Wrap(err, "failed to write file content")
	}

	// 刷新缓冲区
	if err := writer.Flush(); err != nil {
		return errors.Wrap(err, "failed to flush buffer")
	}

	// 设置文件权限
	return setFilePermissions(fpath, f.FileInfo())
}

// ZipFiles 压缩多个文件到zip归档
func ZipFiles(filename string, files []string) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "failed to create zip file")
	}
	defer newZipFile.Close()

	// 使用缓冲写入器
	bufferedWriter := bufio.NewWriterSize(newZipFile, bufferSize)
	defer bufferedWriter.Flush()

	zipWriter := zip.NewWriter(bufferedWriter)
	defer zipWriter.Close()

	// 添加文件到zip
	for _, file := range files {
		if err = addFileToZip(zipWriter, file); err != nil {
			return err
		}
	}

	return nil
}

// addFileToZip 添加单个文件到zip
func addFileToZip(zipWriter *zip.Writer, filename string) error {
	// 获取文件信息
	info, err := os.Lstat(filename)
	if err != nil {
		return errors.Wrapf(err, "failed to stat file: %s", filename)
	}

	// 处理软链接
	if info.Mode()&os.ModeSymlink != 0 {
		return addSymlinkToZip(zipWriter, filename, info)
	}

	// 处理普通文件
	return addRegularFileToZip(zipWriter, filename, info)
}

// addSymlinkToZip 添加软链接到zip
func addSymlinkToZip(zipWriter *zip.Writer, filename string, _ os.FileInfo) error {
	// 解析软链接
	realPath, err := filepath.EvalSymlinks(filename)
	if err != nil {
		// 软链接目标不存在，跳过
		return nil
	}

	// 获取真实文件信息
	realInfo, err := os.Stat(realPath)
	if err != nil {
		// 无法访问目标，跳过
		return nil
	}

	// 如果指向目录，跳过
	if realInfo.IsDir() {
		return nil
	}

	// 创建文件头，使用原始文件名
	header, err := zip.FileInfoHeader(realInfo)
	if err != nil {
		return errors.Wrapf(err, "failed to create zip header for: %s", filename)
	}

	header.Name = filename
	header.Method = zip.Deflate

	// 写入文件内容（从真实文件读取）
	return writeFileToZip(zipWriter, header, realPath)
}

// addRegularFileToZip 添加常规文件到zip
func addRegularFileToZip(zipWriter *zip.Writer, filename string, info os.FileInfo) error {
	// 创建文件头
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return errors.Wrapf(err, "failed to create zip header for: %s", filename)
	}

	header.Name = filename
	header.Method = zip.Deflate

	// 写入文件内容
	return writeFileToZip(zipWriter, header, filename)
}

// writeFileToZip 写入文件内容到zip
func writeFileToZip(zipWriter *zip.Writer, header *zip.FileHeader, sourcePath string) error {
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return errors.Wrapf(err, "failed to create zip entry for: %s", header.Name)
	}

	// 打开源文件
	file, err := os.Open(sourcePath)
	if err != nil {
		return errors.Wrapf(err, "failed to open file: %s", sourcePath)
	}
	defer file.Close()

	// 使用缓冲读取器
	reader := bufio.NewReaderSize(file, bufferSize)

	// 复制内容
	_, err = io.Copy(writer, reader)
	if err != nil {
		return errors.Wrapf(err, "failed to write file to zip: %s", sourcePath)
	}

	return nil
}

// CompressDir 压缩目录到zip文件
func CompressDir(src, dst string, excludePaths ...string) error {
	// 检查源目录
	srcInfo, err := os.Stat(src)
	if err != nil {
		return errors.Wrapf(err, "failed to access source directory: %s", src)
	}
	if !srcInfo.IsDir() {
		return errors.Newf("%s is not a directory", src)
	}

	// 创建目标文件
	zipFile, err := os.Create(dst)
	if err != nil {
		return errors.Wrapf(err, "failed to create zip file: %s", dst)
	}
	defer zipFile.Close()

	// 使用缓冲写入器
	bufferedWriter := bufio.NewWriterSize(zipFile, bufferSize)
	defer bufferedWriter.Flush()

	zipWriter := zip.NewWriter(bufferedWriter)
	defer zipWriter.Close()

	// 获取源路径的绝对路径
	srcPath, err := filepath.Abs(filepath.Clean(src))
	if err != nil {
		return errors.Wrapf(err, "failed to get absolute path for: %s", src)
	}

	// 创建排除路径映射
	excludeMap := make(map[string]bool)
	for _, exclude := range excludePaths {
		excludePath := filepath.Join(srcPath, exclude)
		excludeMap[excludePath] = true
	}

	// 创建软链接处理器
	handler := newSymlinkHandler(srcPath)

	// 遍历目录
	return filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查排除路径
		if excludeMap[path] {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 获取相对路径
		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return errors.Wrapf(err, "failed to get relative path for: %s", path)
		}

		// 跳过根目录
		if relPath == "." {
			return nil
		}

		// 处理文件或目录
		return compressEntry(zipWriter, path, relPath, info, handler)
	})
}

// compressEntry 压缩单个条目
func compressEntry(zipWriter *zip.Writer, path, relPath string, info os.FileInfo, handler *symlinkHandler) error {
	// 处理软链接
	if info.Mode()&os.ModeSymlink != 0 {
		return compressSymlink(zipWriter, path, relPath, handler)
	}

	// 创建文件头
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return errors.Wrapf(err, "failed to create zip header for: %s", path)
	}

	// 使用相对路径
	header.Name = filepath.ToSlash(relPath)

	// 处理目录
	if info.IsDir() {
		header.Name += "/"
		header.Method = zip.Store
		handler.markProcessed(path)

		_, err = zipWriter.CreateHeader(header)
		return err
	}

	// 处理文件
	header.Method = zip.Deflate
	return writeFileToZip(zipWriter, header, path)
}

// compressSymlink 压缩软链接
func compressSymlink(zipWriter *zip.Writer, path, relPath string, handler *symlinkHandler) error {
	// 解析软链接
	realPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		// 软链接目标不存在，跳过
		return nil
	}

	// 获取目标信息
	realInfo, err := os.Stat(realPath)
	if err != nil {
		// 无法访问目标，跳过
		return nil
	}

	// 如果是目录
	if realInfo.IsDir() {
		// 检查是否已处理
		if handler.isProcessed(realPath) {
			return nil
		}
		handler.markProcessed(realPath)

		// 遍历目录内容
		return walkSymlinkDir(zipWriter, realPath, relPath, handler)
	}

	// 如果是文件，创建文件条目
	header, err := zip.FileInfoHeader(realInfo)
	if err != nil {
		return errors.Wrapf(err, "failed to create zip header for: %s", path)
	}

	header.Name = filepath.ToSlash(relPath)
	header.Method = zip.Deflate

	return writeFileToZip(zipWriter, header, realPath)
}

// walkSymlinkDir 遍历软链接指向的目录
func walkSymlinkDir(zipWriter *zip.Writer, realPath, baseRelPath string, handler *symlinkHandler) error {
	return filepath.Walk(realPath, func(subPath string, subInfo os.FileInfo, err error) error {
		if err != nil || subPath == realPath {
			return err
		}

		// 如果遇到软链接，检查是否已处理（防止循环）
		if subInfo.Mode()&os.ModeSymlink != 0 {
			linkTarget, err := filepath.EvalSymlinks(subPath)
			if err != nil {
				// 软链接无效，跳过
				return nil
			}

			// 如果软链接指向已处理的目录，跳过
			if handler.isProcessed(linkTarget) {
				return nil
			}

			// 获取链接目标信息
			targetInfo, err := os.Stat(linkTarget)
			if err != nil {
				return nil
			}

			// 如果是目录，跳过（防止无限递归）
			if targetInfo.IsDir() {
				return nil
			}
		}

		// 获取相对路径
		subRelPath, err := filepath.Rel(realPath, subPath)
		if err != nil {
			return errors.Wrapf(err, "failed to get relative path for: %s", subPath)
		}

		// 组合最终路径
		finalRelPath := filepath.Join(baseRelPath, subRelPath)

		// 处理条目
		if subInfo.IsDir() {
			if handler.isProcessed(subPath) {
				return filepath.SkipDir
			}
			handler.markProcessed(subPath)

			// 创建目录条目
			header, err := zip.FileInfoHeader(subInfo)
			if err != nil {
				return errors.Wrapf(err, "failed to create zip header for: %s", subPath)
			}

			header.Name = filepath.ToSlash(finalRelPath) + "/"
			header.Method = zip.Store

			_, err = zipWriter.CreateHeader(header)
			return err
		}

		// 跳过软链接本身（如果指向目录）
		if subInfo.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		// 处理文件
		header, err := zip.FileInfoHeader(subInfo)
		if err != nil {
			return errors.Wrapf(err, "failed to create zip header for: %s", subPath)
		}

		header.Name = filepath.ToSlash(finalRelPath)
		header.Method = zip.Deflate

		return writeFileToZip(zipWriter, header, subPath)
	})
}

// CompressFile 压缩单个文件
func CompressFile(src, dst string) error {
	// 检查文件状态
	info, err := os.Lstat(src)
	if err != nil {
		return errors.Wrapf(err, "failed to access source file: %s", src)
	}

	realPath := src
	realInfo := info

	// 处理软链接
	if info.Mode()&os.ModeSymlink != 0 {
		realPath, err = filepath.EvalSymlinks(src)
		if err != nil {
			return errors.Wrapf(err, "failed to evaluate symlink: %s", src)
		}

		realInfo, err = os.Stat(realPath)
		if err != nil {
			return errors.Wrapf(err, "failed to get real file info: %s", realPath)
		}
	}

	// 检查是否为目录
	if realInfo.IsDir() {
		return errors.Newf("%s is a directory, not a file", src)
	}

	// 创建目标文件
	zipFile, err := os.Create(dst)
	if err != nil {
		return errors.Wrapf(err, "failed to create zip file: %s", dst)
	}
	defer zipFile.Close()

	// 使用缓冲写入器
	bufferedWriter := bufio.NewWriterSize(zipFile, bufferSize)
	defer bufferedWriter.Flush()

	zipWriter := zip.NewWriter(bufferedWriter)
	defer zipWriter.Close()

	// 创建文件头
	header, err := zip.FileInfoHeader(realInfo)
	if err != nil {
		return errors.Wrapf(err, "failed to create zip header for: %s", src)
	}

	// 使用文件基名
	header.Name = filepath.Base(src)
	header.Method = zip.Deflate

	// 写入文件内容
	return writeFileToZip(zipWriter, header, realPath)
}

// ChownR 递归修改文件所有者（导出以保持向后兼容）
// 在Windows上此函数不执行任何操作
func ChownR(path string, uid, gid int) error {
	return chownR(path, uid, gid)
}
