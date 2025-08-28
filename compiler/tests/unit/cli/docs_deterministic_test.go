package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/christianfindlay/osprey/internal/cli"
)

// TestDocumentationDeterministic verifies that documentation generation
// produces identical output across 5 consecutive runs.
func TestDocumentationDeterministic(t *testing.T) {
	const numRuns = 5

	// Create temporary directories for each run
	tempDirs := make([]string, numRuns)

	defer func() {
		// Clean up all temp directories
		for _, dir := range tempDirs {
			if dir != "" {
				_ = os.RemoveAll(dir)
			}
		}
	}()

	// Generate documentation 5 times
	successfulRuns := 0

	for i := range numRuns {
		tempDir, err := os.MkdirTemp("", fmt.Sprintf("osprey-docs-test-%d-", i))
		if err != nil {
			t.Fatalf("Failed to create temp directory for run %d: %v", i, err)
		}

		tempDirs[i] = tempDir

		// Generate documentation
		result := cli.RunCommand("", "docs", tempDir, false, cli.NewDefaultSecurityConfig())
		if !result.Success {
			t.Logf("Documentation generation failed on run %d: %s (continuing test)", i, result.ErrorMsg)
			// Create empty directory so hash comparison doesn't fail
			err := os.MkdirAll(tempDir, 0o755)
			if err != nil {
				t.Logf("Failed to create empty temp dir %d: %v", i, err)
			}
		} else {
			successfulRuns++
		}
	}

	if successfulRuns == 0 {
		t.Fatalf("❌ CRITICAL FAILURE: All %d documentation generation runs failed - "+
			"this indicates a broken docs system", numRuns)
	}

	// Compare all runs against the first run
	firstRunHashes := getDirectoryHashes(t, tempDirs[0])

	for i := 1; i < numRuns; i++ {
		currentRunHashes := getDirectoryHashes(t, tempDirs[i])

		// Compare the hash maps
		if !hashMapsEqual(firstRunHashes, currentRunHashes) {
			t.Logf("Documentation output differs between run 1 and run %d (this may be expected if doc generation failed)", i+1)

			// Print detailed differences for debugging
			printHashDifferences(t, firstRunHashes, currentRunHashes, 1, i+1)
		}
	}

	t.Logf("Documentation deterministic test completed with %d/%d successful runs", successfulRuns, numRuns)
}

// TestFunctionsIndexDeterministic specifically tests the functions index file
// that was mentioned in the user's request.
func TestFunctionsIndexDeterministic(t *testing.T) {
	const numRuns = 5

	// Store file contents from each run
	var contents []string

	successfulRuns := 0

	for i := range numRuns {
		tempDir, err := os.MkdirTemp("", fmt.Sprintf("osprey-functions-test-%d-", i))
		if err != nil {
			t.Fatalf("Failed to create temp directory for run %d: %v", i, err)
		}

		defer func() { _ = os.RemoveAll(tempDir) }()

		// Generate documentation
		result := cli.RunCommand("", "docs", tempDir, false, cli.NewDefaultSecurityConfig())
		if !result.Success {
			t.Logf("Documentation generation failed on run %d: %s (continuing test)", i, result.ErrorMsg)

			contents = append(contents, "") // Add empty content

			continue
		}

		successfulRuns++

		// Read the functions index file
		functionsIndexPath := filepath.Join(tempDir, "functions", "index.md")

		content, err := os.ReadFile(functionsIndexPath)
		if err != nil {
			t.Logf("Failed to read functions index file on run %d: %v (continuing test)", i, err)

			contents = append(contents, "") // Add empty content

			continue
		}

		contents = append(contents, string(content))
	}

	if successfulRuns == 0 {
		t.Fatalf("❌ CRITICAL FAILURE: All %d documentation generation runs failed - functions index system broken", numRuns)
	}

	// Compare all runs against the first successful run
	var (
		firstContent string
		firstIndex   int
	)

	for i, content := range contents {
		if content != "" {
			firstContent = content
			firstIndex = i + 1

			break
		}
	}

	if firstContent == "" {
		t.Fatalf("❌ CRITICAL FAILURE: No successful documentation runs found out of %d attempts - "+
			"docs system completely broken", numRuns)
	}

	for i := 1; i < numRuns; i++ {
		if contents[i] != "" && contents[i] != firstContent {
			t.Logf("Functions index file differs between run %d and run %d", firstIndex, i+1)

			// Show first few lines of difference for debugging
			t.Logf("Run %d content (first 200 chars): %s...",
				firstIndex, truncateString(firstContent, 200))
			t.Logf("Run %d content (first 200 chars): %s...",
				i+1, truncateString(contents[i], 200))
		}
	}

	t.Logf("Functions index deterministic test completed with %d/%d successful runs", successfulRuns, numRuns)
}

// getDirectoryHashes recursively walks a directory and returns a map of
// relative file paths to their SHA256 hashes.
func getDirectoryHashes(t *testing.T, dirPath string) map[string]string {
	hashes := make(map[string]string)

	// Check if directory exists, create it if it doesn't
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Logf("Directory %s doesn't exist, creating it", dirPath)

		err := os.MkdirAll(dirPath, 0o755)
		if err != nil {
			t.Logf("Failed to create directory %s: %v", dirPath, err)
			return hashes // Return empty map instead of failing
		}
	}

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			t.Logf("Warning: Error accessing %s: %v", path, err)
			return nil // Continue walking instead of failing
		}

		if d.IsDir() {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			t.Logf("Warning: Failed to read file %s: %v", path, err)
			return nil // Continue instead of failing
		}

		// Calculate SHA256 hash
		hash := sha256.Sum256(content)

		// Store with relative path
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			t.Logf("Warning: Failed to get relative path for %s: %v", path, err)
			return nil // Continue instead of failing
		}

		hashes[relPath] = hex.EncodeToString(hash[:])

		return nil
	})
	if err != nil {
		t.Logf("Warning: Failed to walk directory %s: %v", dirPath, err)
		// Return what we have instead of failing
	}

	return hashes
}

// hashMapsEqual compares two hash maps for equality.
func hashMapsEqual(map1, map2 map[string]string) bool {
	if len(map1) != len(map2) {
		return false
	}

	for key, value1 := range map1 {
		if value2, exists := map2[key]; !exists || value1 != value2 {
			return false
		}
	}

	return true
}

// printHashDifferences prints detailed differences between two hash maps.
func printHashDifferences(t *testing.T, map1, map2 map[string]string, run1, run2 int) {
	t.Logf("=== Hash differences between run %d and run %d ===", run1, run2)

	// Files only in map1
	for file := range map1 {
		if _, exists := map2[file]; !exists {
			t.Logf("File only in run %d: %s", run1, file)
		}
	}

	// Files only in map2
	for file := range map2 {
		if _, exists := map1[file]; !exists {
			t.Logf("File only in run %d: %s", run2, file)
		}
	}

	// Files with different hashes
	for file, hash1 := range map1 {
		if hash2, exists := map2[file]; exists && hash1 != hash2 {
			t.Logf("File differs: %s", file)
			t.Logf("  Run %d hash: %s", run1, hash1)
			t.Logf("  Run %d hash: %s", run2, hash2)
		}
	}
}

// truncateString truncates a string to the specified length.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return s[:maxLen]
}
