package version

import "github.com/blang/semver/v4"

func LT(v1, v2 string) bool {
	vv1, _ := semver.Make(v1)
	vv2, _ := semver.Make(v2)
	return vv2.LT(vv1)
}

func Parse(v string) (semver.Version, error) {
	return semver.Make(v)
}
