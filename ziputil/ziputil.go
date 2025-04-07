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
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
)

func ensureDir(dirName string) error {
	err := os.Mkdir(dirName, 0755)
	if err == nil || os.IsExist(err) {
		return nil
	}
	if err = ChownR(dirName, os.Getuid(), os.Getgid()); err != nil {
		return errors.Wrap(err, "unable to chown directory")
	}
	return err
}

func ChownR(path string, uid, gid int) error {
	return filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
		if err == nil {
			err = os.Chown(name, uid, gid)
		}
		return err
	})
}

// MakeZip ...
func MakeZip(inputPath, outputFile string) (bool, error) {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return false, err
	}
	files, err := fileList(inputPath)
	if err != nil {
		return false, errors.Wrap(err, "unable to list files")
	}
	err = ZipFiles(outputFile, files)
	if err != nil {
		return false, err
	}
	return true, nil
}

func fileList(fileDirectory string) ([]string, error) {
	var files []string
	err := filepath.Walk(fileDirectory, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, filepath.Join(path))
		}
		return nil
	})

	return files, err
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {
	var filenames []string

	uid := os.Getuid()
	gid := os.Getgid()

	err := ensureDir(dest)
	if err != nil {
		return filenames, err
	}
	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, errors.Newf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return filenames, errors.Wrap(err, "unable to create directory")
		}

		if err := ChownR(filepath.Dir(fpath), uid, gid); err != nil {
			return filenames, errors.Wrap(err, "unable to chown directory")
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			if err := os.MkdirAll(fpath, 0755); err != nil {
				return filenames, errors.Wrap(err, "unable to create directory")
			}
			if err := ChownR(fpath, uid, gid); err != nil {
				return filenames, errors.Wrap(err, "unable to chown directory")
			}
			continue
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, errors.Wrap(err, "unable to create file")
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}

	for _, f := range filenames {
		if err := ChownR(f, uid, gid); err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
func ZipFiles(filename string, files []string) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = addFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	// header.Name = filepath.Base(filename)
	// header.Name = filepath.FromSlash("test")
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

// CompressDir compresses a directory into a zip file
// src: source directory to compress
// dst: destination zip file path
// excludePaths: paths to exclude from compression (relative to src)
func CompressDir(src, dst string, excludePaths ...string) error {
	// Check if source directory exists
	srcInfo, err := os.Stat(src)
	if err != nil {
		return errors.Wrapf(err, "failed to access source directory: %s", src)
	}
	if !srcInfo.IsDir() {
		return errors.Newf("%s is not a directory", src)
	}

	// Create destination file
	zipFile, err := os.Create(dst)
	if err != nil {
		return errors.Wrapf(err, "failed to create zip file: %s", dst)
	}
	defer zipFile.Close()

	// Create a new zip archive.
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Clean the source path and get the absolute path
	srcPath, err := filepath.Abs(filepath.Clean(src))
	if err != nil {
		return errors.Wrapf(err, "failed to get absolute path for: %s", src)
	}

	// Create a map of excluded paths for faster lookup
	excludeMap := make(map[string]bool)
	for _, exclude := range excludePaths {
		excludePath := filepath.Join(srcPath, exclude)
		excludeMap[excludePath] = true
	}

	// Walk through all files in the source directory
	return filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if path is in exclude list
		if excludeMap[path] {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Get the relative path for the file within the zip
		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return errors.Wrapf(err, "failed to get relative path for: %s", path)
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		// Create zip header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return errors.Wrapf(err, "failed to create zip header for: %s", path)
		}

		// Use relative path and fix the separator to make it consistent
		header.Name = filepath.ToSlash(relPath)

		// Set appropriate flags for directories
		if info.IsDir() {
			header.Name += "/" // Ensure trailing slash for directories
			header.Method = zip.Store
		} else {
			header.Method = zip.Deflate
		}

		// Create writer for the file entry
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return errors.Wrapf(err, "failed to create zip entry for: %s", header.Name)
		}

		// If it's a directory, we're done
		if info.IsDir() {
			return nil
		}

		// If it's a file, open it and copy its content
		file, err := os.Open(path)
		if err != nil {
			return errors.Wrapf(err, "failed to open file: %s", path)
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		if err != nil {
			return errors.Wrapf(err, "failed to write file to zip: %s", path)
		}

		return nil
	})
}

// CompressFile is a convenience function to compress a single file
func CompressFile(src, dst string) error {
	// Check if source file exists
	srcInfo, err := os.Stat(src)
	if err != nil {
		return errors.Wrapf(err, "failed to access source file: %s", src)
	}
	if srcInfo.IsDir() {
		return errors.Newf("%s is a directory, not a file", src)
	}

	// Create destination file
	zipFile, err := os.Create(dst)
	if err != nil {
		return errors.Wrapf(err, "failed to create zip file: %s", dst)
	}
	defer zipFile.Close()

	// Create a new zip archive
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Create zip header
	header, err := zip.FileInfoHeader(srcInfo)
	if err != nil {
		return errors.Wrapf(err, "failed to create zip header for: %s", src)
	}

	// Use the base name of the file for the header
	header.Name = filepath.Base(src)
	header.Method = zip.Deflate

	// Create writer for the file entry
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return errors.Wrapf(err, "failed to create zip entry for: %s", header.Name)
	}

	// Open source file
	file, err := os.Open(src)
	if err != nil {
		return errors.Wrapf(err, "failed to open file: %s", src)
	}
	defer file.Close()

	// Copy file content to zip
	_, err = io.Copy(writer, file)
	if err != nil {
		return errors.Wrapf(err, "failed to write file to zip: %s", src)
	}

	return nil
}
