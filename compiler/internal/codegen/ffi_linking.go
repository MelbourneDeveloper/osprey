package codegen

import "strings"

// FFI link directives let an Osprey source request that a third-party C library
// be linked, e.g. `// @link: sqlite3` -> `-lsqlite3`, and that an extra search
// directory be added, e.g. `// @linkdir: /opt/homebrew/opt/libpq/lib` -> `-L...`.
// This is how SQLite, libpq (and any C library) are reached through the generic
// `extern fn` FFI without hardcoding library names in the compiler.
//
// Directives are parsed from source and carried to the link step as IR comment
// markers (LLVM ignores `;` lines), so they survive the source -> IR -> link
// pipeline shared by both `--compile` and JIT `--run`.
// Precedent: cgo's `// #cgo LDFLAGS:` magic comments.

const (
	linkDirectivePrefix    = "@link:"
	linkDirDirectivePrefix = "@linkdir:"
	irLinkMarker           = "; osprey-link:"
	irLinkDirMarker        = "; osprey-linkdir:"
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

// isValidLibPath is isValidLibName plus the path separator, for `-L<path>`.
func isValidLibPath(path string) bool {
	if path == "" || path[0] == '-' {
		return false
	}

	for _, r := range path {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		case r == '_' || r == '-' || r == '.' || r == '+' || r == '/':
		default:
			return false
		}
	}

	return true
}

// directiveValues returns the trimmed argument of every `// <prefix>VALUE`
// comment in text, de-duplicated, that passes the validator.
func directiveValues(text, prefix string, valid func(string) bool) []string {
	var out []string

	seen := make(map[string]bool)

	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "//") {
			continue
		}

		comment := strings.TrimSpace(strings.TrimPrefix(trimmed, "//"))

		value, ok := strings.CutPrefix(comment, prefix)
		if !ok {
			continue
		}

		value = strings.TrimSpace(value)
		if valid(value) && !seen[value] {
			seen[value] = true

			out = append(out, value)
		}
	}

	return out
}

// injectLinkMarkers prepends an IR comment marker for each link / linkdir
// directive in source, so the link step can recover them from the IR alone (the
// JIT path only forwards IR, not source).
func injectLinkMarkers(ir, source string) string {
	libs := directiveValues(source, linkDirectivePrefix, isValidLibName)
	dirs := directiveValues(source, linkDirDirectivePrefix, isValidLibPath)

	if len(libs) == 0 && len(dirs) == 0 {
		return ir
	}

	var b strings.Builder

	for _, dir := range dirs {
		b.WriteString(irLinkDirMarker + " " + dir + "\n")
	}

	for _, lib := range libs {
		b.WriteString(irLinkMarker + " " + lib + "\n")
	}

	b.WriteString(ir)

	return b.String()
}

// linkLibFlags recovers linker flags from the IR markers: `-L<dir>` search
// directories first, then `-l<name>` libraries. Markers only appear in the
// leading comment block, so scanning stops at the first real IR line.
func linkLibFlags(ir string) []string {
	var dirFlags []string

	var libFlags []string

	for _, line := range strings.Split(ir, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if !strings.HasPrefix(trimmed, ";") {
			break
		}

		if after, ok := strings.CutPrefix(trimmed, irLinkDirMarker); ok {
			dir := strings.TrimSpace(after)
			if isValidLibPath(dir) {
				dirFlags = append(dirFlags, "-L"+dir)
			}

			continue
		}

		if after, ok := strings.CutPrefix(trimmed, irLinkMarker); ok {
			name := strings.TrimSpace(after)
			if isValidLibName(name) {
				libFlags = append(libFlags, "-l"+name)
			}
		}
	}

	return append(dirFlags, libFlags...)
}
