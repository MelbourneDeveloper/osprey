package codegen

// Magic number constants.
const (
	TwoArgs              = 2
	ThreeArgs            = 3
	FourArgs             = 4
	FiveArgs             = 5
	OneArg               = 1
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
	TypeString = "string"
	TypeInt    = "int"
	TypeBool   = "bool"
	TypeAny    = "any"
	TypeResult = "Result"
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
	WebSocketKeepAlive = "webSocketKeepAlive"

	// String utility functions
	LengthFunc    = "length"
	ContainsFunc  = "contains"
	SubstringFunc = "substring"

	// Process and system functions
	SpawnProcessFunc = "spawnProcess"
	WriteFileFunc    = "writeFile"
	ReadFileFunc     = "readFile"

	// JSON functions
	ParseJSONFunc   = "parseJSON"
	ExtractCodeFunc = "extractCode"
)

// HTTP Function names.
const (
	// HTTP Server functions.
	HTTPCreateServerFunc = "httpCreateServer"
	HTTPListenFunc       = "httpListen"
	HTTPStopServerFunc   = "httpStopServer"

	// HTTP Client functions.
	HTTPCreateClientFunc = "httpCreateClient"
	HTTPGetFunc          = "httpGet"
	HTTPPostFunc         = "httpPost"
	HTTPPutFunc          = "httpPut"
	HTTPDeleteFunc       = "httpDelete"
	HTTPRequestFunc      = "httpRequest"
	HTTPCloseClientFunc  = "httpCloseClient"

	// WebSocket functions.
	WebSocketConnectFunc = "websocketConnect"
	WebSocketSendFunc    = "websocketSend"
	WebSocketCloseFunc   = "websocketClose"

	// WebSocket Server functions.
	WebSocketCreateServerFunc    = "websocketCreateServer"
	WebSocketServerListenFunc    = "websocketServerListen"
	WebSocketServerSendFunc      = "websocketServerSend"
	WebSocketServerBroadcastFunc = "websocketServerBroadcast"
	WebSocketStopServerFunc      = "websocketStopServer"
)

// Pattern matching constants.
const (
	UnknownPattern  = "unknown"
	WildcardPattern = "_"
	SuccessPattern  = "Success"
	ErrorPattern    = "Error"
)

// Type name constants for case variations.
const (
	StringTypeName = "String"
	IntTypeName    = "Int"
	BoolTypeName   = "Bool"
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

// Any type validation operation types.
const (
	AnyOpArithmetic       = "arithmetic"
	AnyOpFieldAccess      = "field_access"
	AnyOpFunctionArgument = "function_argument"
	AnyOpVariableAccess   = "variable_access"
)
