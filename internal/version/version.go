// Package version enables setting build-time version using ldflags.
package version

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	// ProjectName is the canonical project name set by ldl flags
	ProjectName = ""
	// ProjectURL is the canonical project url  set by ldl flags
	ProjectURL = ""
	// Version specifies Semantic versioning increment (MAJOR.MINOR.PATCH).
	Version = "v0.0.0"
	// GitCommit specifies the git commit sha, set by the compiler.
	GitCommit = ""
	// BuildMeta specifies release type (dev,rc1,beta,etc)
	BuildMeta = ""
	// Timestamp specifies the build time, set by the compiler as UNIX timestamp.
	Timestamp = ""

	runtimeVersion = runtime.Version()
)

func BuildTime() string {
	tm, err := strconv.Atoi(Timestamp)
	if err != nil {
		return ""
	}
	return time.Unix(int64(tm), 0).UTC().Format(time.RFC3339)
}

// FullVersion returns a version string.
func FullVersion() string {
	var sb strings.Builder
	sb.Grow(len(Version) + len(GitCommit) + len(BuildMeta) + len("-") + len("+"))
	sb.WriteString(Version)
	if BuildMeta != "" {
		sb.WriteString("-" + BuildMeta)
	}
	if GitCommit != "" {
		sb.WriteString("+" + GitCommit)
	}
	return sb.String()
}

// UserAgent returns a user-agent string as specified in RFC 2616:14.43
// https://tools.ietf.org/html/rfc2616
func UserAgent() string {
	return fmt.Sprintf("%s/%s (+%s; %s; %s)", ProjectName, Version, ProjectURL, GitCommit, runtimeVersion)
}
