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
	"fmt"
	"os"

	"github.com/ergoapi/util/zos"
	"github.com/ergoapi/util/ztime"
)

//CheckFileExists check file exist
func CheckFileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// Writefile 写文件
func Writefile(logpath, msg string, name ...string) (err error) {
	svcname := "unknown"
	if len(name) != 0 {
		svcname = name[0]
	}
	prepath := "/var/log/"
	if zos.IsMacOS() {
		prepath = "/tmp"
	}
	logpath = fmt.Sprintf("%v/%v/%v", prepath, svcname, logpath)
	file, err := os.OpenFile(logpath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(fmt.Sprintf("%v - %v\n", ztime.NowFormat(), msg))
	write.Flush()
	return nil
}
