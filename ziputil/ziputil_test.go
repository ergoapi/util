package ziputil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressDir(t *testing.T) {
	// 创建临时测试目录
	tempDir, err := os.MkdirTemp("", "ziputil-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 创建测试文件结构
	testFiles := map[string]string{
		"file1.txt":           "Hello, World!",
		"dir1/file2.txt":      "Hello, Dir1!",
		"dir1/dir2/file3.txt": "Hello, Dir2!",
	}

	// 创建文件和目录
	for path, content := range testFiles {
		fullPath := filepath.Join(tempDir, path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		assert.NoError(t, err)

		err = os.WriteFile(fullPath, []byte(content), 0644)
		assert.NoError(t, err)
	}

	// 创建临时zip文件
	zipFile := filepath.Join(tempDir, "test.zip")

	// 测试压缩目录
	err = CompressDir(tempDir, zipFile)
	assert.NoError(t, err)

	// 验证zip文件存在
	_, err = os.Stat(zipFile)
	assert.NoError(t, err)

	// 测试解压缩
	extractDir := filepath.Join(tempDir, "extracted")
	err = os.MkdirAll(extractDir, 0755)
	assert.NoError(t, err)

	files, err := Unzip(zipFile, extractDir)
	assert.NoError(t, err)
	assert.NotEmpty(t, files)

	// 验证解压后的文件内容
	for path, expectedContent := range testFiles {
		fullPath := filepath.Join(extractDir, path)
		content, err := os.ReadFile(fullPath)
		assert.NoError(t, err)
		assert.Equal(t, expectedContent, string(content))
	}
}

func TestCompressFile(t *testing.T) {
	// 创建临时测试文件
	tempDir, err := os.MkdirTemp("", "ziputil-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, World!"
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	assert.NoError(t, err)

	// 创建临时zip文件
	zipFile := filepath.Join(tempDir, "test.zip")

	// 测试压缩文件
	err = CompressFile(testFile, zipFile)
	assert.NoError(t, err)

	// 验证zip文件存在
	_, err = os.Stat(zipFile)
	assert.NoError(t, err)

	// 测试解压缩
	extractDir := filepath.Join(tempDir, "extracted")
	err = os.MkdirAll(extractDir, 0755)
	assert.NoError(t, err)

	files, err := Unzip(zipFile, extractDir)
	assert.NoError(t, err)
	assert.NotEmpty(t, files)

	// 验证解压后的文件内容
	extractedFile := filepath.Join(extractDir, "test.txt")
	content, err := os.ReadFile(extractedFile)
	assert.NoError(t, err)
	assert.Equal(t, testContent, string(content))
}

func TestCompressDirWithExclude(t *testing.T) {
	// 创建临时测试目录
	tempDir, err := os.MkdirTemp("", "ziputil-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 创建测试文件结构
	testFiles := map[string]string{
		"file1.txt":           "Hello, World!",
		"dir1/file2.txt":      "Hello, Dir1!",
		"dir1/dir2/file3.txt": "Hello, Dir2!",
		"exclude/file4.txt":   "This should be excluded",
	}

	// 创建文件和目录
	for path, content := range testFiles {
		fullPath := filepath.Join(tempDir, path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		assert.NoError(t, err)

		err = os.WriteFile(fullPath, []byte(content), 0644)
		assert.NoError(t, err)
	}

	// 创建临时zip文件
	zipFile := filepath.Join(tempDir, "test.zip")

	// 测试压缩目录（排除 exclude 目录）
	err = CompressDir(tempDir, zipFile, "exclude")
	assert.NoError(t, err)

	// 验证zip文件存在
	_, err = os.Stat(zipFile)
	assert.NoError(t, err)

	// 测试解压缩
	extractDir := filepath.Join(tempDir, "extracted")
	err = os.MkdirAll(extractDir, 0755)
	assert.NoError(t, err)

	files, err := Unzip(zipFile, extractDir)
	assert.NoError(t, err)
	assert.NotEmpty(t, files)

	// 验证解压后的文件内容（排除的文件不应该存在）
	for path, expectedContent := range testFiles {
		if path == "exclude/file4.txt" {
			// 验证排除的文件不存在
			fullPath := filepath.Join(extractDir, path)
			_, err := os.Stat(fullPath)
			assert.True(t, os.IsNotExist(err))
			continue
		}

		fullPath := filepath.Join(extractDir, path)
		content, err := os.ReadFile(fullPath)
		assert.NoError(t, err)
		assert.Equal(t, expectedContent, string(content))
	}
}

func TestCompressDeepDir(t *testing.T) {
	// 创建临时测试目录
	tempDir, err := os.MkdirTemp("", "ziputil-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 创建深层目录结构
	deepDir := filepath.Join(tempDir, "deep")
	err = os.MkdirAll(deepDir, 0755)
	assert.NoError(t, err)

	// 创建多级目录和文件
	dirStructure := map[string]string{
		"level1/file1.txt":                             "Level 1 File 1",
		"level1/level2/file2.txt":                      "Level 2 File 2",
		"level1/level2/level3/file3.txt":               "Level 3 File 3",
		"level1/level2/level3/level4/file4.txt":        "Level 4 File 4",
		"level1/level2/level3/level4/level5/file5.txt": "Level 5 File 5",
		"other/file6.txt":                              "Other File 6",
		"other/deep/file7.txt":                         "Other Deep File 7",
	}

	// 创建文件和目录
	for path, content := range dirStructure {
		fullPath := filepath.Join(deepDir, path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		assert.NoError(t, err)

		err = os.WriteFile(fullPath, []byte(content), 0644)
		assert.NoError(t, err)
	}

	// 创建临时zip文件
	zipFile := filepath.Join(tempDir, "deep.zip")

	// 测试压缩目录
	err = CompressDir(deepDir, zipFile)
	assert.NoError(t, err)

	// 验证zip文件存在
	_, err = os.Stat(zipFile)
	assert.NoError(t, err)

	// 测试解压缩
	extractDir := filepath.Join(tempDir, "extracted")
	err = os.MkdirAll(extractDir, 0755)
	assert.NoError(t, err)

	files, err := Unzip(zipFile, extractDir)
	assert.NoError(t, err)
	assert.NotEmpty(t, files)

	// 验证解压后的文件内容
	for path, expectedContent := range dirStructure {
		fullPath := filepath.Join(extractDir, path)
		content, err := os.ReadFile(fullPath)
		assert.NoError(t, err)
		assert.Equal(t, expectedContent, string(content))
	}

	// 验证目录结构
	for path := range dirStructure {
		dirPath := filepath.Dir(path)
		fullDirPath := filepath.Join(extractDir, dirPath)
		_, err := os.Stat(fullDirPath)
		assert.NoError(t, err, "Directory should exist: %s", fullDirPath)
	}
}
