package zos

import (
	"bytes"
	"os/exec"
	"strings"
)

// UNAME uname -r
func UNAME() string {
	if NotUnix() {
		return "unknow"
	}
	cmd := exec.Command("uname", "-r")
	var b bytes.Buffer
	cmd.Stdout = &b
	err := cmd.Run()
	if err != nil {
		return "unknow"
	}
	return strings.Trim(b.String(), "\n")
}
