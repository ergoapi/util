package zos

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// https://github.com/lima-vm/lima/blob/19e79df9673c5122fbe3e3418b6297c6296ec8eb/pkg/localpathutil/localpathutil.go

// Expand expands a path like "~", "~/", "~/foo".
// Paths like "~foo/bar" are unsupported.
//
// FIXME: is there an existing library for this?
func HomeExpand(orig string) (string, error) {
	s := orig
	if s == "" {
		return "", errors.New("empty path")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(s, "~") {
		if s == "~" || strings.HasPrefix(s, "~/") {
			s = strings.Replace(s, "~", homeDir, 1)
		} else {
			// Paths like "~foo/bar" are unsupported.
			return "", fmt.Errorf("unexpandable path %q", orig)
		}
	}
	return filepath.Abs(s)
}
