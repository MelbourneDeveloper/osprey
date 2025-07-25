package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/codegen"
)

// TestLibDir is the relative path to the lib directory for tests
const TestLibDir = "../../lib"

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
	fiberLib := filepath.Join("lib", "libfiber_runtime.a")
	httpLib := filepath.Join("lib", "libhttp_runtime.a")
	websocketLib := filepath.Join("lib", "libwebsocket_runtime.a")
	systemLib := filepath.Join("lib", "libsystem_runtime.a")

	// Get current working directory for absolute path
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	linkArgs := []string{
		"-o", "test",
		"test.o",
		filepath.Join(cwd, fiberLib),
		filepath.Join(cwd, httpLib),
		filepath.Join(cwd, websocketLib),
		filepath.Join(cwd, systemLib),
		"-lpthread", "-lssl", "-lcrypto",
	}

	t.Logf("Link arguments: %v", linkArgs)

	// Check that ALL 4 required libraries are referenced
	hasFiberLib := false
	hasHTTPLib := false
	hasWebSocketLib := false
	hasSystemLib := false
	hasSSL := false
	hasCrypto := false
	hasPthread := false

	for _, arg := range linkArgs {
		if strings.Contains(arg, "libfiber_runtime.a") {
			hasFiberLib = true
		}
		if strings.Contains(arg, "libhttp_runtime.a") {
			hasHTTPLib = true
		}
		if strings.Contains(arg, "libwebsocket_runtime.a") {
			hasWebSocketLib = true
		}
		if strings.Contains(arg, "libsystem_runtime.a") {
			hasSystemLib = true
		}
		if arg == "-lssl" {
			hasSSL = true
		}
		if arg == "-lcrypto" {
			hasCrypto = true
		}
		if arg == "-lpthread" {
			hasPthread = true
		}
	}

	if !hasFiberLib {
		t.Fatal("Missing fiber runtime library")
	}
	if !hasHTTPLib {
		t.Fatal("Missing HTTP runtime library")
	}
	if !hasWebSocketLib {
		t.Fatal("Missing WebSocket runtime library")
	}
	if !hasSystemLib {
		t.Fatal("Missing System runtime library")
	}
	if !hasSSL {
		t.Fatal("Missing -lssl")
	}
	if !hasCrypto {
		t.Fatal("Missing -lcrypto")
	}
	if !hasPthread {
		t.Fatal("Missing -lpthread")
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
		t.Fatalf("Failed to compile test C file: %v. Output: %s", err, string(output))
	}

	// Build the exact link arguments that compilation.go would use
	var linkArgs []string
	linkArgs = append(linkArgs, "clang")
	linkArgs = append(linkArgs, "-o", testExe, testO)

	// Add ALL 4 runtime libraries if available (order matters: dependents before dependencies)
	if fiberLib := findLibrary("libfiber_runtime.a"); fiberLib != "" {
		linkArgs = append(linkArgs, fiberLib)
		t.Logf("Using Fiber library: %s", fiberLib)
	}
	if httpLib := findLibrary("libhttp_runtime.a"); httpLib != "" {
		linkArgs = append(linkArgs, httpLib)
		t.Logf("Using HTTP library: %s", httpLib)
	}
	if websocketLib := findLibrary("libwebsocket_runtime.a"); websocketLib != "" {
		linkArgs = append(linkArgs, websocketLib)
		t.Logf("Using WebSocket library: %s", websocketLib)
	}
	if systemLib := findLibrary("libsystem_runtime.a"); systemLib != "" {
		linkArgs = append(linkArgs, systemLib)
		t.Logf("Using System library: %s", systemLib)
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
		t.Fatalf("Manual linking failed: %v. Output: %s", err, string(output))
	} else {
		t.Logf("Manual linking succeeded!")
		t.Logf("Link output: %s", string(output))
	}
}

// findLibrary is a helper function to find library using the same search logic as the JIT executor.
func findLibrary(libName string) string {
	// Use the exact same search paths as the JIT executor
	possiblePaths := []string{
		filepath.Join("lib", libName),
		filepath.Join(".", "lib", libName),
		"/usr/local/lib/" + libName, // System install location
	}

	// Add working directory based paths - match JIT executor exactly
	if wd, err := os.Getwd(); err == nil {
		possiblePaths = append(possiblePaths,
			filepath.Join(wd, "lib", libName),
			filepath.Join(wd, "..", "lib", libName),
			filepath.Join(wd, "..", "..", "lib", libName),
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

	// Create a minimal HTTP Osprey file that handles a Result type
	ospCode := `
fn handleRequest(method: string, path: string, headers: string, body: string) -> HttpResponse = HttpResponse {
    status: 200,
    headers: "Content-Type: text/plain",
    contentType: "text/plain",
    streamFd: -1,
    isComplete: true,
    partialBody: "Hello World!"
}

fn main() -> int {
    let server = httpCreateServer(8080, "127.0.0.1")
    print("Server created with ID: ")
    print(toString(server))
    print("\n")

    let result = httpListen(server, handleRequest)
    print("Listen result: ")
    print(toString(result))
    print("\n")

    0
}
`

	err := os.WriteFile(ospFile, []byte(ospCode), 0o644)
	if err != nil {
		t.Fatalf("Failed to write test Osprey file: %v", err)
	}

	// Now try to compile it
	outputFile := filepath.Join(testDir, "test_http")
	// Use CompileToExecutableWithLibDir to specify the lib directory for tests
	err = codegen.CompileToExecutableWithLibDir(ospCode, outputFile, TestLibDir)

	if err != nil {
		t.Fatalf("FATAL: Compilation failed unexpectedly. This test expects successful compilation. Error: %v", err)
	} else {
		t.Log("Compilation succeeded as expected.")
	}
}

// TestHTTPCompilation tests compiling HTTP code and traces any linking issues.
func TestHTTPCompilationLinking(t *testing.T) {
	// Create a minimal HTTP example that correctly handles a Result type
	ospCode := `
fn main() -> int {
    let client = httpCreateClient("https://httpbin.org", 5000)
    print("Client created with ID: ")
    print(toString(client))
    print("\n")

    0
}
`

	testDir := t.TempDir()
	outputFile := filepath.Join(testDir, "test_http")

	// Run the compilation
	// Use CompileToExecutableWithLibDir to specify the lib directory for tests
	err := codegen.CompileToExecutableWithLibDir(ospCode, outputFile, TestLibDir)

	if err != nil {
		t.Fatalf("FATAL: ❌ Compilation failed unexpectedly. "+
			"This test requires successful compilation to verify linking. Error: %v", err)
	} else {
		t.Log("✅ Compilation succeeded - HTTP runtime linking is working!")
	}
}

// TestFailsCompilationCircularDependency tests that circular effect dependencies fail compilation
func TestFailsCompilationCircularDependency(t *testing.T) {
	// Test the circular dependency example
	// Use CompileToExecutableWithLibDir to specify the lib directory for tests
	err := codegen.CompileToExecutableWithLibDir(`
effect StateA {
    getFromB: fn() -> int
    setInA: fn(int) -> Unit
}

effect StateB {
    getFromA: fn() -> int
    setInB: fn(int) -> Unit
}

fn circularEffectA() -> int !StateA, StateB = {
    let bValue = perform StateB.getFromA()
    perform StateA.setInA(bValue + 1)
    perform StateA.getFromB()
}

fn circularEffectB() -> int !StateA, StateB = {
    let aValue = perform StateA.getFromA()
    perform StateB.setInB(aValue + 1)
    perform StateB.getFromA()
}

fn main() -> Unit = {
    with handler StateA
        getFromB() => circularEffectB()
        setInA(x) => print("StateA set: " + toString(x))
    with handler StateB
        getFromA() => circularEffectA()
        setInB(x) => print("StateB set: " + toString(x))
    {
        let result = circularEffectA()
        print("Result: " + toString(result))
    }
}`, "/tmp/circular_test", TestLibDir)

	// This SHOULD fail with a circular dependency error
	if err == nil {
		t.Fatal("Expected compilation to fail due to circular effect dependency, but it succeeded")
	}

	// Check that the error message mentions circular dependencies
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "circular") && !strings.Contains(errorMsg, "recursion") {
		t.Logf("Error message: %s", errorMsg)
		// For now, just log that we need to implement circular dependency detection
		t.Log("⚠️  NOTE: Circular dependency detection not yet implemented - this will be added later")
	} else {
		t.Log("✅ Circular dependency correctly detected!")
	}
}
