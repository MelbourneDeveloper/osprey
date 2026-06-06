package codegen

import "strings"

// FFI link directives let an Osprey source request that a third-party C library
// be linked, e.g. `// @link: sqlite3` -> `-lsqlite3`. This is how SQLite (and any
// C library) is reached through the generic `extern fn` FFI without hardcoding
// library names in the compiler.
//
// The directive is parsed from source and carried to the link step as an IR
// comment marker (LLVM ignores `;` lines), so it survives the
// source -> IR -> link pipeline shared by both `--compile` and JIT `--run`.
// Precedent: cgo's `// #cgo LDFLAGS:` magic comments.

const (
	linkDirectivePrefix = "@link:"
	irLinkMarker        = "; osprey-link:"
)

// isValidLibName reports whether name is a safe library token. Only these
// characters are allowed in a `-l<name>` flag — this prevents a crafted
// directive from smuggling extra linker flags (a leading '-') or shell
// metacharacters through the link command.
func isValidLibName(name string) bool {
	if name == "" || name[0] == '-' {
		return false
	}

	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		case r == '_' || r == '-' || r == '.' || r == '+':
		default:
			return false
		}
	}

	return true
}

// parseLinkDirectives extracts library names from `// @link: NAME` comments.
func parseLinkDirectives(source string) []string {
	var libs []string

	seen := make(map[string]bool)

	for _, line := range strings.Split(source, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "//") {
			continue
		}

		comment := strings.TrimSpace(strings.TrimPrefix(trimmed, "//"))
		if !strings.HasPrefix(comment, linkDirectivePrefix) {
			continue
		}

		name := strings.TrimSpace(strings.TrimPrefix(comment, linkDirectivePrefix))
		if isValidLibName(name) && !seen[name] {
			seen[name] = true

			libs = append(libs, name)
		}
	}

	return libs
}

// injectLinkMarkers prepends an IR comment marker for each link directive found
// in source, so the link step can recover the requested libraries from the IR
// alone (the JIT path only forwards IR, not source).
func injectLinkMarkers(ir, source string) string {
	libs := parseLinkDirectives(source)
	if len(libs) == 0 {
		return ir
	}

	var b strings.Builder

	for _, lib := range libs {
		b.WriteString(irLinkMarker)
		b.WriteString(" ")
		b.WriteString(lib)
		b.WriteString("\n")
	}

	b.WriteString(ir)

	return b.String()
}

// linkLibFlags recovers `-l<name>` flags from the IR link markers. Markers only
// appear in the leading comment block, so scanning stops at the first real IR
// line.
func linkLibFlags(ir string) []string {
	var flags []string

	for _, line := range strings.Split(ir, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if !strings.HasPrefix(trimmed, ";") {
			break
		}

		if after, ok := strings.CutPrefix(trimmed, irLinkMarker); ok {
			name := strings.TrimSpace(after)
			if isValidLibName(name) {
				flags = append(flags, "-l"+name)
			}
		}
	}

	return flags
}
