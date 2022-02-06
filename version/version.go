package version

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

type Info struct {
	Major        string `json:"major"`
	Minor        string `json:"minor"`
	GitVersion   string `json:"gitVersion"`
	GitBranch    string `json:"gitBranch"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

func (info Info) String() string {
	return info.GitVersion
}

func Get() Info {
	// These variables typically come from -ldflags settings and in
	// their absence fallback to the settings in pkg/version/base.go
	return Info{
		Major:        gitMajor,
		Minor:        gitMinor,
		GitVersion:   gitVersion,
		GitBranch:    gitBranch,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func GetJsonString() string {
	bs, _ := json.MarshalIndent(Get(), "", "  ")
	return string(bs)
}

func shortDate(dateStr string) string {
	var buf strings.Builder
	for _, c := range dateStr {
		if c >= '0' && c <= '9' {
			buf.WriteRune(c)
		}
	}
	dateStr = buf.String()
	if strings.HasPrefix(dateStr, "20") {
		dateStr = dateStr[2:]
	}
	if len(dateStr) > 8 {
		dateStr = dateStr[:8]
	}
	return dateStr
}

func GetShortString() string {
	v := Get()
	return fmt.Sprintf("%s(%s%s)", v.GitBranch, v.GitCommit, shortDate(v.BuildDate))
}
