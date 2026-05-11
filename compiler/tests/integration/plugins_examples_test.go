package integration

// DO NOT EVER SKIP TESTS!!!!

import "testing"

// TestPluginsExamples exercises the language plugin system end-to-end: it compiles and runs
// every .osp file under examples/tested/plugins, where each example uses `fn <plugin> ...`
// declarations. The plugin subprocess is invoked at compile time (validating the language
// body), and the runtime output is matched byte-for-byte against the corresponding
// .expectedoutput file — proving the plugin's response actually drives codegen.
func TestPluginsExamples(t *testing.T) {
	checkLLVMTools(t)

	runTestExamplesRecursive(t, "../../examples/tested/plugins", getExpectedOutputs())
}
