// Package version provides the Shipwright-compliant version contract for the
// osprey binary. See https://github.com/nimblesite/Shipwright docs/specs/binary-version-contract.md
// (SWR-VERSION-*).
package version

import (
	"encoding/json"
	"fmt"
	"io"
)

// ComponentID is the Shipwright component id printed as the first stdout token
// for `osprey --version`. It must match the `id` field for the osprey component
// in every shipwright.json that references this binary.
const ComponentID = "osprey"

// Kind is the Shipwright component kind for this binary.
const Kind = "cli"

// Language is the implementation language reported in JSON version output.
const Language = "go"

// SourcePlaceholder is the valid semantic-version placeholder used in source
// control per SWR-VERSION-BUILD-STAMPING. Release builds override Version via
// -ldflags "-X github.com/christianfindlay/osprey/internal/version.Version=X.Y.Z".
const SourcePlaceholder = "0.0.0-dev"

// Version is the stamped semantic version. It defaults to the source
// placeholder and is overridden at link time during release builds.
var Version = SourcePlaceholder //nolint:gochecknoglobals // ldflags target

// PrintPlain writes the SWR-VERSION plain text contract line to w.
//
// The first line MUST be exactly "<ComponentID> <Version>" per
// `docs/specs/binary-version-contract.md` § "Required --version Behavior".
func PrintPlain(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%s %s\n", ComponentID, Version)
	return err
}

// Manifest is the JSON shape mandated by
// `schemas/version-manifest.schema.json` (SWR-VERSION JSON output).
type Manifest struct {
	ManifestVersion int    `json:"manifestVersion"`
	Name            string `json:"name"`
	Version         string `json:"version"`
	Kind            string `json:"kind"`
	Language        string `json:"language"`
}

// PrintJSON writes the SWR-VERSION JSON contract to w.
func PrintJSON(w io.Writer) error {
	manifest := Manifest{
		ManifestVersion: 1,
		Name:            ComponentID,
		Version:         Version,
		Kind:            Kind,
		Language:        Language,
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(manifest)
}
