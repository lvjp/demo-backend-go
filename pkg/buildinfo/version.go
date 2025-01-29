package buildinfo

import (
	"fmt"
	"runtime/debug"
)

type BuildInfo struct {
	Revision     string
	RevisionTime string
	Modified     bool

	GoVersion string
	GoOS      string
	GoArch    string
}

var buildInfo BuildInfo

func init() {
	buildInfo.setDefaults()

	raw, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	buildInfo.Load(raw)
}

func (bi *BuildInfo) setDefaults() {
	bi.Revision = "-"
	bi.RevisionTime = "-"
	bi.Modified = false
	bi.GoVersion = "-"
	bi.GoOS = "-"
	bi.GoArch = "-"
}

func (bi *BuildInfo) Load(raw *debug.BuildInfo) {
	bi.GoVersion = raw.GoVersion

	for _, s := range raw.Settings {
		switch s.Key {
		case "vcs.revision":
			bi.Revision = s.Value
		case "vcs.time":
			bi.RevisionTime = s.Value
		case "vcs.modified":
			bi.Modified = s.Value == "true"
		case "GOOS":
			bi.GoOS = s.Value
		case "GOARCH":
			bi.GoArch = s.Value
		}
	}
}

func (bi BuildInfo) VersionString() string {
	ret := fmt.Sprintf(
		"%s %s %s %s/%s",
		bi.Revision,
		bi.RevisionTime,
		bi.GoVersion,
		bi.GoOS,
		bi.GoArch,
	)

	if bi.Modified {
		ret += " (modified)"
	}

	return ret
}

func Get() BuildInfo {
	return buildInfo
}
