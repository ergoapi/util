// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package version

import (
	"encoding/json"
	"fmt"
	"runtime"
)

// Info contains versioning information.
// GitVersion, GitCommit and BuildDate are set via -ldflags at build time.
type Info struct {
	GitVersion string `json:"gitVersion"`
	GitCommit  string `json:"gitCommit"`
	BuildDate  string `json:"buildDate"`
	GoVersion  string `json:"goVersion"`
	Platform   string `json:"platform"`
}

func (info Info) String() string {
	return info.GitVersion
}

// Get returns the overall codebase version.
func Get() Info {
	return Info{
		GitVersion: gitVersion,
		GitCommit:  gitCommit,
		BuildDate:  buildDate,
		GoVersion:  runtime.Version(),
		Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// GetJSONString returns version info as indented JSON.
func GetJSONString() string {
	bs, _ := json.MarshalIndent(Get(), "", "  ")
	return string(bs)
}
