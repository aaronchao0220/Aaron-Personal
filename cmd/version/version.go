// Package version contains build information. The build information is populated at build-time.
package version

var (
	// Version is the service version
	Version string

	// Revision is the SHA-1 of the git revision
	Revision string

	// Branch is the name of the git branch
	Branch string

	// BuildTime is the build's timestamp
	BuildTime string
)
