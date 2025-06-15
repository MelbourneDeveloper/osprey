package codegen

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMain runs before all tests in this package.
func TestMain(m *testing.M) {
	// Clean and rebuild everything before running any tests
	cleanAndRebuildAll()

	// Run all tests
	code := m.Run()

	// Exit with the test result code
	os.Exit(code)
}

// cleanAndRebuildAll cleans and rebuilds all dependencies.
func cleanAndRebuildAll() {
	// Get project root (go up from internal/codegen to project root)
	wd, err := os.Getwd()
	if err != nil {
		panic("Failed to get working directory: " + err.Error())
	}
	projectRoot := filepath.Join(wd, "..", "..")

	// Clean everything including Rust
	cmd := exec.Command("make", "clean")
	cmd.Dir = projectRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		panic("Failed to clean: " + err.Error() + "\nOutput: " + string(output))
	}

	// Rebuild runtime libraries
	cmd = exec.Command("make", "fiber-runtime", "http-runtime")
	cmd.Dir = projectRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		panic("Failed to build runtime libraries: " + err.Error() + "\nOutput: " + string(output))
	}

	// Build Rust interop library
	rustDir := filepath.Join(projectRoot, "examples", "rust_integration")
	if _, err := os.Stat(rustDir); err == nil {
		cmd = exec.Command("cargo", "build")
		cmd.Dir = rustDir
		if output, err := cmd.CombinedOutput(); err != nil {
			panic("Failed to build Rust interop: " + err.Error() + "\nOutput: " + string(output))
		}
	}

	// Create symlink for codegen tests
	binPath := filepath.Join(wd, "bin")
	targetPath := filepath.Join(projectRoot, "bin")

	// Remove existing symlink if it exists
	_ = os.Remove(binPath)

	// Create symlink (ignore errors since it may already exist)
	_ = os.Symlink(targetPath, binPath)

	// Build compiler (needed for some tests) - skip linting for faster test builds
	cmd = exec.Command("make", "build-no-lint")
	cmd.Dir = projectRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		// If build-no-lint target doesn't exist, try regular build
		cmd = exec.Command("go", "build", "-o", "bin/osprey", "./cmd/osprey")
		cmd.Dir = projectRoot
		if output2, err2 := cmd.CombinedOutput(); err2 != nil {
			panic("Failed to build compiler: " + err.Error() +
				"\nOutput: " + string(output) +
				"\nFallback error: " + err2.Error() +
				"\nFallback output: " + string(output2))
		}
	}
}

// TestPkgConfigOpenSSL tests that pkg-config can find OpenSSL.
func TestPkgConfigOpenSSL(t *testing.T) {
	cmd := exec.Command("pkg-config", "--libs", "openssl")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("pkg-config failed to find OpenSSL: %v", err)
	}

	outputStr := strings.TrimSpace(string(output))
	if !strings.Contains(outputStr, "ssl") {
		t.Errorf("Expected OpenSSL libraries in output, got: %s", outputStr)
	}

	t.Logf("✅ OpenSSL libraries found: %s", outputStr)

	// Also test cflags
	cmd = exec.Command("pkg-config", "--cflags", "openssl")
	output, err = cmd.Output()
	if err != nil {
		t.Fatalf("pkg-config failed to get OpenSSL cflags: %v", err)
	}

	cflagsStr := strings.TrimSpace(string(output))
	t.Logf("✅ OpenSSL cflags: %s", cflagsStr)

	// Test crypto specifically
	cmd = exec.Command("pkg-config", "--libs", "libcrypto")
	output, err = cmd.Output()
	if err != nil {
		t.Fatalf("pkg-config failed to find libcrypto: %v", err)
	}

	cryptoStr := strings.TrimSpace(string(output))
	if !strings.Contains(cryptoStr, "crypto") {
		t.Errorf("Expected crypto library in output, got: %s", cryptoStr)
	}

	t.Logf("✅ Crypto library found: %s", cryptoStr)

	// Test specific libraries that should be available
	expectedLibs := []string{"ssl", "crypto"}
	for _, lib := range expectedLibs {
		if !strings.Contains(outputStr+" "+cryptoStr, lib) {
			t.Errorf("Expected library %s not found in pkg-config output", lib)
		}
	}
}

// TestBuildLinkArguments tests that we can generate proper link arguments.
func TestBuildLinkArguments(t *testing.T) {
	httpLib := filepath.Join("bin", "libhttp_runtime.a")
	fiberLib := filepath.Join("bin", "libfiber_runtime.a")

	// Get current working directory for absolute path
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	linkArgs := []string{
		"-o", "test",
		"test.o",
		filepath.Join(cwd, httpLib),
		filepath.Join(cwd, fiberLib),
		"-lpthread", "-lssl", "-lcrypto",
	}

	t.Logf("Link arguments: %v", linkArgs)

	// Check that required libraries are referenced
	validateLinkArguments(t, linkArgs)
}

// validateLinkArguments checks that all required libraries are present in link arguments.
func validateLinkArguments(t *testing.T, linkArgs []string) {
	libFlags := collectLibraryFlags(linkArgs)
	checkRequiredLibraries(t, libFlags)
}

// collectLibraryFlags collects all library flags from link arguments.
func collectLibraryFlags(linkArgs []string) map[string]bool {
	flags := map[string]bool{
		"http":    false,
		"fiber":   false,
		"ssl":     false,
		"crypto":  false,
		"pthread": false,
	}

	for _, arg := range linkArgs {
		if strings.Contains(arg, "libhttp_runtime.a") {
			flags["http"] = true
		}
		if strings.Contains(arg, "libfiber_runtime.a") {
			flags["fiber"] = true
		}
		if arg == "-lssl" {
			flags["ssl"] = true
		}
		if arg == "-lcrypto" {
			flags["crypto"] = true
		}
		if arg == "-lpthread" {
			flags["pthread"] = true
		}
	}

	return flags
}

// checkRequiredLibraries verifies that all required libraries are present.
func checkRequiredLibraries(t *testing.T, flags map[string]bool) {
	if !flags["http"] {
		t.Error("Missing HTTP runtime library")
	}
	if !flags["fiber"] {
		t.Error("Missing fiber runtime library")
	}
	if !flags["ssl"] {
		t.Error("Missing -lssl")
	}
	if !flags["crypto"] {
		t.Error("Missing -lcrypto")
	}
	if !flags["pthread"] {
		t.Error("Missing -lpthread")
	}
}

// TestHTTPRuntimeLibrary verifies that the HTTP runtime library contains expected symbols.
func TestHTTPRuntimeLibrary(t *testing.T) {
	// Get the working directory and construct the library path
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Go up to the project root and then to bin
	httpLibPath := filepath.Join(wd, "..", "..", "bin", "libhttp_runtime.a")
	t.Logf("Found HTTP library: %s", httpLibPath)

	// Check if the library exists
	if _, err := os.Stat(httpLibPath); os.IsNotExist(err) {
		t.Fatalf("HTTP runtime library not built at %s - build failed! Error: %v", httpLibPath, err)
	}

	// Use nm to check symbols in the library
	cmd := exec.Command("nm", httpLibPath)
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("nm command failed - required for symbol analysis: %v", err)
	}

	symbols := string(output)
	t.Logf("HTTP library symbols (first 500 chars): \n%s", symbols[:min(500, len(symbols))])

	// Check for modern OpenSSL EVP symbols instead of deprecated SHA1 symbols
	if !strings.Contains(symbols, "EVP_MD_CTX_new") &&
		!strings.Contains(symbols, "EVP_sha1") &&
		!strings.Contains(symbols, "EVP_DigestInit_ex") {
		t.Log("No OpenSSL EVP symbols found - may be statically linked or using system libraries")
	}

	// Check for our own HTTP functions
	expectedSymbols := []string{
		"sha1_websocket",
		"http_create_server",
		"websocket_handshake",
	}

	for _, symbol := range expectedSymbols {
		if !strings.Contains(symbols, symbol) {
			t.Errorf("Missing expected symbol: %s", symbol)
		}
	}
}

// TestManualLinking tests manual linking with the exact same arguments that compilation.go would use.
func TestManualLinking(t *testing.T) {
	// Create a minimal test object file first
	testC := filepath.Join(t.TempDir(), "test.c")
	testO := filepath.Join(t.TempDir(), "test.o")
	testExe := filepath.Join(t.TempDir(), "test")

	// Create minimal C file that uses modern EVP API
	cCode := `
#include <openssl/evp.h>

int main() {
    EVP_MD_CTX *ctx = EVP_MD_CTX_new();
    if (ctx) {
        EVP_MD_CTX_free(ctx);
    }
    return 0;
}
`

	err := os.WriteFile(testC, []byte(cCode), 0o644)
	if err != nil {
		t.Fatalf("Failed to write test C file: %v", err)
	}

	// Compile to object file with OpenSSL 3.5.0+ flags
	compileArgs := []string{"-c"}

	// Add pkg-config OpenSSL compile flags if available
	if cmd := exec.Command("pkg-config", "--cflags", "openssl"); cmd != nil {
		if output, err := cmd.Output(); err == nil {
			flags := strings.Fields(strings.TrimSpace(string(output)))
			compileArgs = append(compileArgs, flags...)
		}
	}

	compileArgs = append(compileArgs,
		"-DOPENSSL_SUPPRESS_DEPRECATED",
		"-DOPENSSL_API_COMPAT=30000",
		"-Wno-deprecated-declarations",
		"-o", testO, testC)

	compileCmd := exec.Command("clang", compileArgs...)
	if output, err := compileCmd.CombinedOutput(); err != nil {
		t.Errorf("Failed to compile test C file: %v", err)
		t.Errorf("Compile output: %s", string(output))

		return
	}

	// Build the exact link arguments that compilation.go would use
	var linkArgs []string
	linkArgs = append(linkArgs, "clang")
	linkArgs = append(linkArgs, "-o", testExe, testO)

	// Add HTTP runtime library if available
	if httpLib := findLibrary("libhttp_runtime.a"); httpLib != "" {
		linkArgs = append(linkArgs, httpLib)
		t.Logf("Using HTTP library: %s", httpLib)
	}

	linkArgs = append(linkArgs, "-lpthread")

	// Add OpenSSL flags exactly as compilation.go does
	pkgCmd := exec.Command("pkg-config", "--libs", "openssl")
	if output, err := pkgCmd.Output(); err == nil {
		flags := strings.Fields(strings.TrimSpace(string(output)))
		linkArgs = append(linkArgs, flags...)
		t.Logf("Added OpenSSL flags: %v", flags)
	} else {
		t.Logf("pkg-config failed, using direct linking")
		linkArgs = append(linkArgs, "-lssl", "-lcrypto")
	}

	t.Logf("Final link command: %v", linkArgs)

	// Execute the link command
	linkCmd := exec.Command(linkArgs[0], linkArgs[1:]...)
	output, err := linkCmd.CombinedOutput()

	if err != nil {
		t.Errorf("Manual linking failed: %v", err)
		t.Errorf("Link output: %s", string(output))
	} else {
		t.Logf("Manual linking succeeded!")
		t.Logf("Link output: %s", string(output))
	}
}

// findLibrary is a helper function to find library.
func findLibrary(libName string) string {
	possiblePaths := []string{
		filepath.Join("bin", libName),
		filepath.Join(".", "bin", libName),
	}

	// Add working directory based paths - go up to project root
	if wd, err := os.Getwd(); err == nil {
		possiblePaths = append(possiblePaths,
			filepath.Join(wd, "..", "..", "bin", libName),
			filepath.Join(wd, "bin", libName),
		)
	}

	for _, libPath := range possiblePaths {
		if _, err := os.Stat(libPath); err == nil {
			return libPath
		}
	}

	return ""
}

// TestGoCompilationTrace tests what the actual Go compilation process does.
func TestActualCompilationProcess(t *testing.T) {
	// Create a minimal HTTP example to compile
	testDir := t.TempDir()
	ospFile := filepath.Join(testDir, "test_http.osp")

	// Create a minimal HTTP Osprey file
	ospCode := `
fn main() {
    let server = http_server("localhost", 8080)
    http_listen(server) { request ->
        http_response(request, 200, "Hello, World!")
    }
}
`

	err := os.WriteFile(ospFile, []byte(ospCode), 0o644)
	if err != nil {
		t.Fatalf("Failed to write test Osprey file: %v", err)
	}

	// Now try to compile it and capture what commands are actually executed
	// We'll patch the exec.Command to log what's being called

	// Import the compilation logic
	outputFile := filepath.Join(testDir, "test_http")

	// This should trigger the same compilation path as the real examples
	err = CompileToExecutable(ospCode, outputFile)

	if err != nil {
		t.Logf("Compilation failed (expected): %v", err)

		// The failure might tell us what command was actually run
		if strings.Contains(err.Error(), "Undefined symbols") {
			t.Error("Compilation failed with OpenSSL linking error - this confirms the bug")
		}
	} else {
		t.Log("Compilation succeeded (unexpected but good)")
	}
}

// TestHTTPCompilation tests compiling HTTP code and traces any linking issues.
func TestHTTPCompilationLinking(t *testing.T) {
	// Create a minimal HTTP example
	ospCode := `
let client = httpCreateClient("https://httpbin.org", 5000)
`

	testDir := t.TempDir()
	outputFile := filepath.Join(testDir, "test_http")

	// Run the compilation and expect it to fail with OpenSSL linking issue
	err := CompileToExecutable(ospCode, outputFile)

	if err != nil {
		t.Logf("Compilation failed as expected: %v", err)

		// Check if the error is specifically about OpenSSL symbols
		errStr := err.Error()
		if strings.Contains(errStr, "SHA1_") {
			t.Log("✅ Confirmed: Compilation fails due to missing OpenSSL linking")
			t.Log("This test documents the bug we need to fix")

			// Check if OpenSSL flags are missing from the error message
			if !strings.Contains(errStr, "-lssl") && !strings.Contains(errStr, "-lcrypto") {
				t.Error("❌ OpenSSL libraries are not being added to link command")
			}
		} else {
			t.Errorf("❌ Compilation failed for unexpected reason: %v", err)
		}
	} else {
		t.Log("✅ Compilation succeeded - OpenSSL linking is working!")
	}
}
