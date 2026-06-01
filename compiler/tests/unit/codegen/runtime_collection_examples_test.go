package codegen_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/christianfindlay/osprey/internal/codegen"
)

func TestTestedExamplesCompileInProcess(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to locate test file")
	}
	testDir := filepath.Dir(currentFile)
	examplesRoot := filepath.Clean(filepath.Join(testDir, "../../../examples/tested"))

	var examples []string
	err := filepath.WalkDir(examplesRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(path) != ".osp" {
			return nil
		}
		examples = append(examples, path)
		return nil
	})
	if err != nil {
		t.Fatalf("find tested examples: %v", err)
	}

	for _, example := range examples {
		name, err := filepath.Rel(examplesRoot, example)
		if err != nil {
			t.Fatalf("name example: %v", err)
		}
		t.Run(name, func(t *testing.T) {
			source, err := os.ReadFile(example)
			if err != nil {
				t.Fatalf("read example: %v", err)
			}
			_, err = codegen.CompileToLLVM(string(source))
			if err != nil {
				t.Fatalf("compile example in process: %v", err)
			}
		})
	}
}
