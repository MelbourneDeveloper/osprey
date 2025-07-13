package codegen

// Magic number constants.
const (
	TwoArgs              = 2
	ThreeArgs            = 3
	FourArgs             = 4
	FiveArgs             = 5
	OneArg               = 1
	HTTPMethodGet        = 0
	HTTPMethodPost       = 1
	HTTPMethodPut        = 2
	HTTPMethodDelete     = 3
	BufferSize1KB        = 1024
	BufferSize64Bytes    = 64
	FilePermissions      = 0o644
	FilePermissionsLess  = 0o644 // Less secure permissions for temp files
	DirPermissions       = 0o750 // More secure permissions
	ArrayIndexZero       = 0
	ArrayIndexOne        = 1
	StringTerminatorSize = 2  // For adding "\x00"
	MinArgs              = 2  // Minimum command line arguments
	ExpressionOffset     = 2  // Offset for expression parsing
	DefaultPlaceholder   = 42 // Default placeholder value for LLVM constants
)

// String constants.
const (
	FormatStringInt    = "%ld\x00"
	FormatStringString = "%s"
	TrueString         = "true\x00"
	FalseString        = "false\x00"
	StringTerminator   = "\x00"
	PercentEscape      = "%%"
)

// Type names.
const (
	TypeString       = "string"
	TypeInt          = "int"
	TypeBool         = "bool"
	TypeAny          = "any"
	TypeUnit         = "Unit"
	TypeResult       = "Result"
	TypeHTTPResponse = "HttpResponse"
	TypeFunction     = "Function"
)

// Function names.
const (
	ToStringFunc       = "toString"
	PrintFunc          = "print"
	InputFunc          = "input"
	RangeFunc          = "range"
	ForEachFunc        = "forEach"
	MapFunc            = "map"
	FilterFunc         = "filter"
	FoldFunc           = "fold"
	MainFunctionName   = "main"
	WebSocketKeepAlive = "websocket_keep_alive"

	// String utility functions
	LengthFunc    = "length"
	ContainsFunc  = "contains"
	SubstringFunc = "substring"

	// Process and system functions
	SpawnProcessFunc = "spawnProcess"

	WriteFileFunc  = "writeFile"
	ReadFileFunc   = "readFile"
	DeleteFileFunc = "deleteFile"
	SleepFunc      = "sleep"

	// Process management functions
	AwaitProcessFunc   = "awaitProcess"
	CleanupProcessFunc = "cleanupProcess"
)

// HTTP Function names.
const (
	// HTTP Server functions.
	HTTPCreateServerFunc = "http_create_server"
	HTTPListenFunc       = "http_listen"
	HTTPStopServerFunc   = "http_stop_server"

	// HTTP Client functions.
	HTTPCreateClientFunc = "http_create_client"
	HTTPGetFunc          = "http_get"
	HTTPPostFunc         = "http_post"
	HTTPPutFunc          = "http_put"
	HTTPDeleteFunc       = "http_delete"
	HTTPRequestFunc      = "http_request"
	HTTPCloseClientFunc  = "http_close_client"

	// WebSocket functions.
	WebSocketConnectFunc = "websocket_connect"
	WebSocketSendFunc    = "websocket_send"
	WebSocketCloseFunc   = "websocket_close"

	// WebSocket Server functions.
	WebSocketCreateServerFunc    = "websocket_create_server"
	WebSocketServerListenFunc    = "websocket_server_listen"
	WebSocketServerSendFunc      = "websocket_server_send"
	WebSocketServerBroadcastFunc = "websocket_server_broadcast"
	WebSocketStopServerFunc      = "websocket_stop_server"
)

// Pattern matching constants.
const (
	UnknownPattern  = "unknown"
	WildcardPattern = "_"
	SuccessPattern  = "Success"
	ErrorPattern    = "Error"
)

// HTTP Error codes (matching C runtime).
const (
	HTTPSuccess         = 0
	HTTPErrorBind       = -1
	HTTPErrorConnection = -2
	HTTPErrorTimeout    = -3
	HTTPErrorParse      = -4
	HTTPErrorServer     = -5
)

// Constants for argument counts
const (
	// Only keep truly unique constants that aren't function argument counts
	MaxFunctionParams = 10
	DefaultBufferSize = 1024
)

// NOTE: Function argument counts have been moved to the unified built-in function registry
// (builtin_registry.go). Use len(GlobalBuiltInRegistry.GetFunction(name).ParameterTypes) instead.
