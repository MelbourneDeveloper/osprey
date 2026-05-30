// Package version is the single source of truth for the Osprey compiler version.
//
// Implements [SWR-VERSION-BUILD-STAMPING]: the source-controlled value MUST be
// the placeholder "0.0.0-dev". Real versions are stamped at release time from the
// git tag via the Go linker, NOT by editing this file:
//
//	go build -ldflags "-X github.com/christianfindlay/osprey/internal/version.Version=1.2.3"
package version

// BinaryName is the component id reported by the --version contract.
// Implements [SWR-VERSION-CLI-OUTPUT]: the first stdout line is "<BinaryName> <Version>".
const BinaryName = "osprey"

// Version is the compiler version. Placeholder in source; overridden at release
// build time via -ldflags. See package doc. Must be a package-level var (not a
// const) because the Go linker's -X flag can only stamp variables.
//
//nolint:gochecknoglobals // intentional: stamped via -ldflags -X at release time.
var Version = "0.0.0-dev"

// Line returns the mandatory plain --version output: "<BinaryName> <Version>".
// Implements [SWR-VERSION-CLI-OUTPUT].
func Line() string {
	return BinaryName + " " + Version
}

// JSON returns the --version --json output. Implements [SWR-VERSION-JSON-OUTPUT].
//
// The `language` field required by schemas/version-manifest.schema.json is
// intentionally omitted: Osprey compiles to a native binary, so its
// implementation language carries no deployment-relevant information, and the
// schema's enum has no value for it. See Nimblesite/Shipwright#1.
func JSON() string {
	return `{"manifestVersion":1,"name":"` + BinaryName +
		`","version":"` + Version + `","kind":"cli","product":"` + BinaryName + `"}`
}
