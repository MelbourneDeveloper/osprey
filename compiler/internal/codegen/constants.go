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
	ResultFieldCount     = 2  // Number of fields in Result type struct (value, discriminant)
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

// GOOSDarwin is the runtime.GOOS value for macOS, extracted as a const
// to satisfy goconst (three+ call sites use it).
const GOOSDarwin = "darwin"

// GOOSWindows is the runtime.GOOS value for Windows. [WINDOWS-PORT]
const GOOSWindows = "windows"

// Type names.
const (
	TypeString       = "string"
	TypeInt          = "int"
	TypeFloat        = "float"
	TypeBool         = "bool"
	TypeAny          = "any"
	TypeUnit         = "Unit"
	TypeResult       = "Result"
	TypeMathError    = "MathError"
	TypeHTTPResponse = "HttpResponse"
	TypeFunction     = "Function"
	TypeFiber        = "Fiber"
	TypeChannel      = "Channel"
	TypeList         = "List"
	TypeMap          = "Map"
)

// Type argument counts.
const (
	TwoTypeArgs = 2
)

// Size constants.
const (
	PointerPairSize = 16 // Size of two pointers (key, value) in bytes
)

// Function names.
const (
	ToStringFunc     = "toString"
	PrintFunc        = "print"
	InputFunc        = "input"
	RangeFunc        = "range"
	ForEachFunc      = "forEach"
	MapFunc          = "map"
	FilterFunc       = "filter"
	FoldFunc         = "fold"
	MainFunctionName = "main"

	LengthFunc    = "length"
	ContainsFunc  = "contains"
	SubstringFunc = "substring"
	ParseIntFunc  = "parseInt"
	JoinFunc      = "join"

	// IsEmptyFunc is the Osprey name for the isEmpty(s) builtin. Implements [BUILTIN-STRING-INSPECTION].
	IsEmptyFunc = "isEmpty"

	// StartsWithFunc is the Osprey name for startsWith(s, prefix). Implements [BUILTIN-STRING-SEARCH].
	StartsWithFunc = "startsWith"
	// EndsWithFunc is the Osprey name for endsWith(s, suffix). Implements [BUILTIN-STRING-SEARCH].
	EndsWithFunc = "endsWith"
	// IndexOfFunc is the Osprey name for indexOf(s, needle). Implements [BUILTIN-STRING-SEARCH].
	IndexOfFunc = "indexOf"

	// TakeFunc is the Osprey name for take(s, n). Implements [BUILTIN-STRING-SUBSTRINGS].
	TakeFunc = "take"
	// DropFunc is the Osprey name for drop(s, n). Implements [BUILTIN-STRING-SUBSTRINGS].
	DropFunc = "drop"

	// ToUpperCaseFunc is the Osprey name for toUpperCase(s). Implements [BUILTIN-STRING-TRANSFORM].
	ToUpperCaseFunc = "toUpperCase"
	// ToLowerCaseFunc is the Osprey name for toLowerCase(s). Implements [BUILTIN-STRING-TRANSFORM].
	ToLowerCaseFunc = "toLowerCase"
	// TrimFunc is the Osprey name for trim(s). Implements [BUILTIN-STRING-TRANSFORM].
	TrimFunc = "trim"
	// TrimStartFunc is the Osprey name for trimStart(s). Implements [BUILTIN-STRING-TRANSFORM].
	TrimStartFunc = "trimStart"
	// TrimEndFunc is the Osprey name for trimEnd(s). Implements [BUILTIN-STRING-TRANSFORM].
	TrimEndFunc = "trimEnd"
	// ReverseFunc is the Osprey name for reverse(s). Implements [BUILTIN-STRING-TRANSFORM].
	ReverseFunc = "reverse"
	// ReplaceFunc is the Osprey name for replace(s, needle, replacement). Implements [BUILTIN-STRING-TRANSFORM].
	ReplaceFunc = "replace"
	// RepeatFunc is the Osprey name for repeat(s, n). Implements [BUILTIN-STRING-TRANSFORM].
	RepeatFunc = "repeat"
	// PadStartFunc is the Osprey name for padStart(s, target, fill). Implements [BUILTIN-STRING-TRANSFORM].
	PadStartFunc = "padStart"
	// PadEndFunc is the Osprey name for padEnd(s, target, fill). Implements [BUILTIN-STRING-TRANSFORM].
	PadEndFunc = "padEnd"

	// ParseFloatFunc is the Osprey name for parseFloat(s). Implements [BUILTIN-STRING-PARSING].
	ParseFloatFunc = "parseFloat"

	// SplitFunc is the Osprey name for split(s, sep). Implements [BUILTIN-STRING-LIST].
	SplitFunc = "split"
	// LinesFunc is the Osprey name for lines(s). Implements [BUILTIN-STRING-LIST].
	LinesFunc = "lines"
	// WordsFunc is the Osprey name for words(s). Implements [BUILTIN-STRING-LIST].
	WordsFunc = "words"

	SpawnProcessFunc = "spawnProcess"

	WriteFileFunc  = "writeFile"
	ReadFileFunc   = "readFile"
	DeleteFileFunc = "deleteFile"
	SleepFunc      = "sleep"

	AwaitProcessFunc   = "awaitProcess"
	CleanupProcessFunc = "cleanupProcess"
)

// Osprey HTTP Function names.
// These are the names of the functions at the Osprey level.
const (
	HTTPCreateServerOsprey = "httpCreateServer"
	HTTPListenOsprey       = "httpListen"
	HTTPStopServerOsprey   = "httpStopServer"

	HTTPCreateClientOsprey = "httpCreateClient"
	HTTPGetOsprey          = "httpGet"
	HTTPPostOsprey         = "httpPost"
	HTTPPutOsprey          = "httpPut"
	HTTPDeleteOsprey       = "httpDelete"
	HTTPRequestOsprey      = "httpRequest"
	HTTPCloseClientOsprey  = "httpCloseClient"

	WebSocketConnectOsprey = "websocketConnect"
	WebSocketSendOsprey    = "websocketSend"
	WebSocketCloseOsprey   = "websocketClose"

	WebSocketCreateServerOsprey    = "websocketCreateServer"
	WebSocketServerListenOsprey    = "websocketServerListen"
	WebSocketServerSendOsprey      = "websocketServerSend"
	WebSocketServerBroadcastOsprey = "websocketServerBroadcast"
	WebSocketStopServerOsprey      = "websocketStopServer"
	WebSocketKeepAliveOsprey       = "websocketKeepAlive"
)

// C Runtime HTTP Function names.
// These are the names of the functions at the C level.
// There are NOT the Osprey function names
const (
	HTTPCreateServerFunc = "http_create_server"
	HTTPListenFunc       = "http_listen"
	HTTPStopServerFunc   = "http_stop_server"

	HTTPCreateClientFunc = "http_create_client"
	HTTPGetFunc          = "http_get"
	HTTPPostFunc         = "http_post"
	HTTPPutFunc          = "http_put"
	HTTPDeleteFunc       = "http_delete"
	HTTPRequestFunc      = "http_request"
	HTTPCloseClientFunc  = "http_close_client"

	WebSocketConnectFunc = "websocket_connect"
	WebSocketSendFunc    = "websocket_send"
	WebSocketCloseFunc   = "websocket_close"

	WebSocketCreateServerFunc    = "websocket_create_server"
	WebSocketServerListenFunc    = "websocket_server_listen"
	WebSocketServerSendFunc      = "websocket_server_send"
	WebSocketServerBroadcastFunc = "websocket_server_broadcast"
	WebSocketStopServerFunc      = "websocket_stop_server"
	WebSocketKeepAliveFunc       = "websocket_keep_alive"
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

// NOTE: Function argument counts have been moved to the unified built-in function registry
// (builtin_registry.go). Use len(GlobalBuiltInRegistry.GetFunction(name).ParameterTypes) instead.
