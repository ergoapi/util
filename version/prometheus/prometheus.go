package prometheus

import (
	"github.com/ergoapi/util/version"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	buildInfo := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ergo_build_info",
			Help: "A metric with a constant '1' value labeled by major, minor, git version, git commit, git tree state, build date, Go version, and compiler from which Kubernetes was built, and platform on which it is running.",
		},
		[]string{"major", "minor", "gitVersion", "gitCommit", "gitTreeState", "buildDate", "goVersion", "compiler", "platform", "release"},
	)
	info := version.Get()
	buildInfo.WithLabelValues(info.Major, info.Minor, info.GitVersion, info.GitCommit, info.GitTreeState, info.BuildDate, info.GoVersion, info.Compiler, info.Platform, info.Release).Set(1)

	prometheus.MustRegister(buildInfo)
}
