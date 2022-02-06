package version

var (
	// TODO: Deprecate gitMajor and gitMinor, use only gitVersion
	// instead. First step in deprecation, keep the fields but make
	// them irrelevant. (Next we'll take it out, which may muck with
	// scripts consuming the kubectl version output - but most of
	// these should be looking at gitVersion already anyways.)
	gitMajor string = "1" // major version, always numeric
	gitMinor string = "6" // minor version, numeric possibly followed by "+"

	// semantic version, derived by build scripts (see
	// https://github.com/kubernetes/kubernetes/blob/master/docs/design/versioning.md
	// for a detailed discussion of this field)
	//
	// TODO: This field is still called "gitVersion" for legacy
	// reasons. For prerelease versions, the build metadata on the
	// semantic version is a git hash, but the version itself is no
	// longer the direct output of "git describe", but a slight
	// translation to be semver compliant.
	gitVersion   string = "v1.6.2+$Format:%h$"
	gitCommit    string = "$Format:%H$"    // sha1 from git, output of $(git rev-parse HEAD)
	gitTreeState string = "not a git tree" // state of git tree, either "clean" or "dirty"
	gitBranch    string = "unknown"

	buildDate string = "1970-01-01T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
)
