package codegen_test

// Coverage tests for HTTP / WebSocket / channel builtin generators. The
// codegen path only needs the program to type-check and lower; we don't
// run the resulting IR, so the underlying servers/clients are never
// touched. CompileToLLVM uses the permissive security profile by default.

import (
	"testing"

	"github.com/christianfindlay/osprey/internal/codegen"
)

func TestHTTPWebSocketCodegen(t *testing.T) {
	programs := map[string]string{
		"httpRequest_GET": `
fn main() -> int {
  let id = httpCreateClient("http://x", 30)
  httpRequest(id, 0, "/p", "", "")
  0
}`,
		"httpRequest_POST": `
fn main() -> int {
  let id = httpCreateClient("http://x", 30)
  httpRequest(id, 1, "/p", "h", "b")
  0
}`,
		"httpCloseClient": `
fn main() -> int {
  let id = httpCreateClient("http://x", 30)
  httpCloseClient(id)
  0
}`,
		"websocketSend_close": `
fn main() -> int {
  websocketSend(1, "msg")
  websocketClose(1)
  0
}`,
		"websocketServer_full_lifecycle": `
fn main() -> int {
  let s = websocketCreateServer(8080, "127.0.0.1", "/x")
  websocketServerListen(s)
  websocketServerBroadcast(s, "hello")
  websocketStopServer(s)
  websocketKeepAlive()
  0
}`,
	}
	for name, src := range programs {
		t.Run(name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(src)
			if err != nil {
				t.Fatalf("%s failed: %v", name, err)
			}
		})
	}
}

// Channels: fiber channel/send/recv builtins.
func TestChannelCodegen(t *testing.T) {
	programs := map[string]string{
		"channel_via_keyword": `
fn main() -> int {
  let c = Channel<int> { capacity: 1 }
  send(c, 42)
  recv(c)
  0
}`,
	}
	for name, src := range programs {
		t.Run(name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(src)
			if err != nil {
				t.Fatalf("%s failed: %v", name, err)
			}
		})
	}
}

// Yield expression and spawn keyword (separate from fiber_spawn builtin).
func TestFiberKeywordsCodegen(t *testing.T) {
	programs := map[string]string{
		"spawn_keyword": `
fn compute(n: int) -> int = n + 1
fn main() -> int {
  let h = spawn compute(5)
  let r = await(h)
  r
  0
}`,
		"yield_expression": `
fn main() -> int {
  let r = yield 42
  r
  0
}`,
	}
	for name, src := range programs {
		t.Run(name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(src)
			if err != nil {
				t.Fatalf("%s failed: %v", name, err)
			}
		})
	}
}
