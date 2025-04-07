package ziputil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 辅助函数：创建测试文件
func createTestFile(t assert.TestingT, path, content string) {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	assert.NoError(t, err)
	err = os.WriteFile(path, []byte(content), 0644)
	assert.NoError(t, err)
}

// 辅助函数：验证文件内容
func assertFileContent(t *testing.T, path, expectedContent string) {
	content, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Equal(t, expectedContent, string(content))
}

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

// TestCompressEmptyDir 测试空目录的压缩
func TestCompressEmptyDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ziputil-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	emptyDir := filepath.Join(tempDir, "empty")
	err = os.MkdirAll(emptyDir, 0755)
	assert.NoError(t, err)

	zipFile := filepath.Join(tempDir, "empty.zip")
	err = CompressDir(emptyDir, zipFile)
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
	assert.Empty(t, files, "Empty directory should not contain any files")
}

// TestCompressSpecialChars 测试包含特殊字符的文件名
func TestCompressSpecialChars(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ziputil-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	specialFiles := map[string]string{
		"file with spaces.txt":     "Content with spaces",
		"file-with-特殊字符.txt":       "Content with special chars",
		"dir with spaces/file.txt": "Content in spaced dir",
		"dir-with-特殊字符/file.txt":   "Content in special dir",
	}

	for path, content := range specialFiles {
		fullPath := filepath.Join(tempDir, path)
		createTestFile(t, fullPath, content)
	}

	zipFile := filepath.Join(tempDir, "special.zip")
	err = CompressDir(tempDir, zipFile)
	assert.NoError(t, err)

	extractDir := filepath.Join(tempDir, "extracted")
	err = os.MkdirAll(extractDir, 0755)
	assert.NoError(t, err)

	files, err := Unzip(zipFile, extractDir)
	assert.NoError(t, err)
	assert.NotEmpty(t, files)

	for path, content := range specialFiles {
		fullPath := filepath.Join(extractDir, path)
		assertFileContent(t, fullPath, content)
	}
}

// TestCompressErrorCases 测试错误场景
func TestCompressErrorCases(t *testing.T) {
	// 测试源目录不存在
	err := CompressDir("/non/existent/path", "test.zip")
	assert.Error(t, err)

	// 测试目标路径无权限
	tempDir, err := os.MkdirTemp("", "ziputil-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 创建一个只读目录
	readOnlyDir := filepath.Join(tempDir, "readonly")
	err = os.MkdirAll(readOnlyDir, 0444)
	assert.NoError(t, err)

	// 尝试在只读目录中创建zip文件
	err = CompressDir(tempDir, filepath.Join(readOnlyDir, "test.zip"))
	assert.Error(t, err)
}

// BenchmarkCompressDir 基准测试：目录压缩性能
func BenchmarkCompressDir(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "ziputil-bench-*")
	assert.NoError(b, err)
	defer os.RemoveAll(tempDir)

	// 创建测试数据
	for i := 0; i < 100; i++ {
		path := filepath.Join(tempDir, "file", "dir", "subdir", "file.txt")
		createTestFile(b, path, "Benchmark test content")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		zipFile := filepath.Join(tempDir, "bench.zip")
		err := CompressDir(tempDir, zipFile)
		assert.NoError(b, err)
		os.Remove(zipFile)
	}
}

// TestCompressWithSymlinks 测试压缩包含软链接的目录
func TestCompressWithSymlinks(t *testing.T) {
	// 创建临时测试目录
	tempDir, err := os.MkdirTemp("", "ziputil-test-symlink-*")
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

	// 创建指向文件的软链接
	fileSymlinkPath := filepath.Join(tempDir, "symlink-to-file1.txt")
	err = os.Symlink(filepath.Join(tempDir, "file1.txt"), fileSymlinkPath)
	assert.NoError(t, err)

	// 创建指向目录的软链接
	dirSymlinkPath := filepath.Join(tempDir, "symlink-to-dir1")
	err = os.Symlink(filepath.Join(tempDir, "dir1"), dirSymlinkPath)
	assert.NoError(t, err)

	// 创建临时zip文件
	zipFile := filepath.Join(tempDir, "symlink-test.zip")

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

	// 验证指向文件的软链接被正确处理
	// 在解压目录中应该有内容为 "Hello, World!" 的文件
	extractedSymlinkFile := filepath.Join(extractDir, "symlink-to-file1.txt")
	content, err := os.ReadFile(extractedSymlinkFile)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", string(content))

	// 验证指向目录的软链接内部的文件也被正确解压
	extractedSymlinkDirFile := filepath.Join(extractDir, "symlink-to-dir1", "file2.txt")
	content, err = os.ReadFile(extractedSymlinkDirFile)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, Dir1!", string(content))

	// 验证所有原始文件也都正确解压
	for path, expectedContent := range testFiles {
		fullPath := filepath.Join(extractDir, path)
		content, err := os.ReadFile(fullPath)
		assert.NoError(t, err)
		assert.Equal(t, expectedContent, string(content))
	}
}

// TestCompressFileWithSymlink 测试压缩软链接文件
func TestCompressFileWithSymlink(t *testing.T) {
	// 创建临时测试目录
	tempDir, err := os.MkdirTemp("", "ziputil-test-symlink-file-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 创建目标文件
	targetFile := filepath.Join(tempDir, "target.txt")
	targetContent := "This is the target file content"
	err = os.WriteFile(targetFile, []byte(targetContent), 0644)
	assert.NoError(t, err)

	// 创建指向文件的软链接
	symlinkFile := filepath.Join(tempDir, "symlink.txt")
	err = os.Symlink(targetFile, symlinkFile)
	assert.NoError(t, err)

	// 创建临时zip文件
	zipFile := filepath.Join(tempDir, "symlink-file.zip")

	// 测试压缩软链接文件
	err = CompressFile(symlinkFile, zipFile)
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

	// 验证解压后的文件内容是原始目标文件的内容
	extractedFile := filepath.Join(extractDir, "symlink.txt")
	content, err := os.ReadFile(extractedFile)
	assert.NoError(t, err)
	assert.Equal(t, targetContent, string(content))
}

// TestCompressWithRecursiveSymlinks 测试包含循环软链接的情况
func TestCompressWithRecursiveSymlinks(t *testing.T) {
	// 创建临时测试目录
	tempDir, err := os.MkdirTemp("", "ziputil-test-recursive-symlink-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 创建测试目录结构
	dir1 := filepath.Join(tempDir, "dir1")
	err = os.MkdirAll(dir1, 0755)
	assert.NoError(t, err)

	// 创建文件
	file1 := filepath.Join(dir1, "file1.txt")
	err = os.WriteFile(file1, []byte("File in dir1"), 0644)
	assert.NoError(t, err)

	// 创建从 dir1 指向 tempDir 的软链接，这会形成循环
	recursiveSymlink := filepath.Join(dir1, "loop")
	err = os.Symlink(tempDir, recursiveSymlink)
	assert.NoError(t, err)

	// 创建临时zip文件
	zipFile := filepath.Join(tempDir, "recursive-symlink.zip")

	// 测试压缩目录，应该可以处理循环软链接而不陷入无限循环
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

	// 验证基本文件被正确解压
	extractedFile := filepath.Join(extractDir, "dir1", "file1.txt")
	content, err := os.ReadFile(extractedFile)
	assert.NoError(t, err)
	assert.Equal(t, "File in dir1", string(content))
}

// TestCompressWithExternalSymlinks 测试压缩包含指向外部文件的软链接
func TestCompressWithExternalSymlinks(t *testing.T) {
	// 创建临时测试目录
	tempDir, err := os.MkdirTemp("", "ziputil-test-ext-symlink-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 创建测试目录结构
	srcDir := filepath.Join(tempDir, "src")
	externalDir := filepath.Join(tempDir, "external")

	err = os.MkdirAll(srcDir, 0755)
	assert.NoError(t, err)

	err = os.MkdirAll(externalDir, 0755)
	assert.NoError(t, err)

	// 在外部目录创建文件
	externalFile := filepath.Join(externalDir, "external.txt")
	externalContent := "This is an external file"
	err = os.WriteFile(externalFile, []byte(externalContent), 0644)
	assert.NoError(t, err)

	// 在源目录中创建普通文件
	normalFile := filepath.Join(srcDir, "normal.txt")
	normalContent := "This is a normal file"
	err = os.WriteFile(normalFile, []byte(normalContent), 0644)
	assert.NoError(t, err)

	// 在源目录中创建指向外部文件的软链接
	symlinkFile := filepath.Join(srcDir, "symlink-to-external.txt")
	err = os.Symlink(externalFile, symlinkFile)
	assert.NoError(t, err)

	// 创建临时zip文件
	zipFile := filepath.Join(tempDir, "external-symlink.zip")

	// 测试压缩目录
	err = CompressDir(srcDir, zipFile)
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

	// 验证普通文件正确解压
	extractedNormalFile := filepath.Join(extractDir, "normal.txt")
	content, err := os.ReadFile(extractedNormalFile)
	assert.NoError(t, err)
	assert.Equal(t, normalContent, string(content))

	// 验证指向外部文件的软链接被正确处理
	// 在解压目录中应该有包含外部文件内容的文件
	extractedSymlinkFile := filepath.Join(extractDir, "symlink-to-external.txt")
	content, err = os.ReadFile(extractedSymlinkFile)
	assert.NoError(t, err)
	assert.Equal(t, externalContent, string(content))

	// 验证是作为普通文件存储的，而不是作为软链接
	fileInfo, err := os.Lstat(extractedSymlinkFile)
	assert.NoError(t, err)
	assert.False(t, fileInfo.Mode()&os.ModeSymlink != 0, "应该是普通文件，而不是软链接")
}
