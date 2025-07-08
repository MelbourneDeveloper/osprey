// Package ast provides Abstract Syntax Tree definitions for the Osprey language.
package ast

// Position represents a position in the source code.
type Position struct {
	Line   int // 1-based line number
	Column int // 0-based column number
}

// Statement represents a statement in the AST.
type Statement interface {
	isStatement()
}

// Expression represents an expression in the AST.
type Expression interface {
	isExpression()
}

// TypeExpression represents a type in the AST.
type TypeExpression struct {
	Name          string
	GenericParams []TypeExpression // For generic types like Result<Int, Error>
	IsArray       bool
	ArrayElement  *TypeExpression

	// Function type support
	IsFunction     bool             // true if this is a function type
	ParameterTypes []TypeExpression // Parameter types for function types
	ReturnType     *TypeExpression  // Return type for function types

	// Position information
	Position *Position
}

// Parameter represents a function parameter with optional type annotation.
type Parameter struct {
	Name string
	Type *TypeExpression // Optional type annotation
}

// ExternParameter represents a parameter in an extern function declaration.
type ExternParameter struct {
	Name string
	Type TypeExpression // Required type annotation for extern functions
}

// Program represents the root of the AST.
type Program struct {
	Statements []Statement
}

// ImportStatement represents an import declaration.
type ImportStatement struct {
	Module []string
}

func (i *ImportStatement) isStatement() {}

// LetDeclaration represents a variable declaration.
type LetDeclaration struct {
	Name    string
	Mutable bool
	Type    *TypeExpression // Optional type annotation
	Value   Expression

	// Position information
	Position *Position
}

func (l *LetDeclaration) isStatement() {}

// AssignmentStatement represents a mutable variable assignment.
type AssignmentStatement struct {
	Name  string
	Value Expression

	// Position information
	Position *Position
}

func (a *AssignmentStatement) isStatement() {}

// FunctionDeclaration represents a function declaration.
type FunctionDeclaration struct {
	Name       string
	Parameters []Parameter     // Updated to support type annotations
	ReturnType *TypeExpression // Optional return type annotation
	Effects    []string        // CRITICAL: Effect declarations like !Logger, !IO
	Body       Expression

	// Position information
	Position *Position
}

func (f *FunctionDeclaration) isStatement() {}

// ExternDeclaration represents an external function declaration.
type ExternDeclaration struct {
	Name       string
	Parameters []ExternParameter // Required type annotations
	ReturnType *TypeExpression   // Optional return type annotation
}

func (e *ExternDeclaration) isStatement() {}

// TypeDeclaration represents a type declaration with union types.
type TypeDeclaration struct {
	Name           string
	TypeParams     []string // Generic type parameters
	Variants       []TypeVariant
	ValidationFunc *string // Optional WHERE validation function name
}

func (t *TypeDeclaration) isStatement() {}

// TypeVariant represents a variant in a union type.
type TypeVariant struct {
	Name   string
	Fields []TypeField
}

// TypeField represents a field in a type variant.
type TypeField struct {
	Name string
	Type string
}

// FunctionCallExpression represents a function call in constraints.
type FunctionCallExpression struct {
	Function  string
	Arguments []Expression

	// Position information
	Position *Position
}

func (f *FunctionCallExpression) isExpression() {}

// ExpressionStatement represents an expression used as a statement.
type ExpressionStatement struct {
	Expression Expression
}

func (e *ExpressionStatement) isStatement() {}

// Expressions

// IntegerLiteral represents an integer literal.
type IntegerLiteral struct {
	Value int64

	// Position information
	Position *Position
}

func (i *IntegerLiteral) isExpression() {}

// StringLiteral represents a string literal.
type StringLiteral struct {
	Value string

	// Position information
	Position *Position
}

func (s *StringLiteral) isExpression() {}

// BooleanLiteral represents a boolean literal.
type BooleanLiteral struct {
	Value bool

	// Position information
	Position *Position
}

func (b *BooleanLiteral) isExpression() {}

// InterpolatedStringLiteral represents an interpolated string.
type InterpolatedStringLiteral struct {
	Parts []InterpolatedPart

	// Position information
	Position *Position
}

func (i *InterpolatedStringLiteral) isExpression() {}

// InterpolatedPart represents a part of an interpolated string.
type InterpolatedPart struct {
	IsExpression bool
	Text         string
	Expression   Expression
}

// Identifier represents an identifier.
type Identifier struct {
	Name string

	// Position information
	Position *Position
}

func (i *Identifier) isExpression() {}

// BinaryExpression represents a binary expression.
type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression

	// Position information
	Position *Position
}

func (b *BinaryExpression) isExpression() {}

// UnaryExpression represents a unary expression.
type UnaryExpression struct {
	Operator string
	Operand  Expression

	// Position information
	Position *Position
}

func (u *UnaryExpression) isExpression() {}

// CallExpression represents a function call (with optional parentheses).
type CallExpression struct {
	Function       Expression
	Arguments      []Expression
	HasParentheses bool // Whether call uses parentheses
	NamedArguments []NamedArgument

	// Position information
	Position *Position
}

func (c *CallExpression) isExpression() {}

// NamedArgument represents a named argument in a function call.
type NamedArgument struct {
	Name  string
	Value Expression
}

// MethodCallExpression represents method chaining like obj.method().
type MethodCallExpression struct {
	Object         Expression
	MethodName     string
	Arguments      []Expression
	HasParentheses bool // Whether call uses parentheses
	NamedArguments []NamedArgument

	// Position information
	Position *Position
}

func (m *MethodCallExpression) isExpression() {}

// LambdaExpression represents anonymous functions.
type LambdaExpression struct {
	Parameters []Parameter     // Updated to support type annotations
	ReturnType *TypeExpression // Optional return type annotation
	Body       Expression

	// Position information
	Position *Position
}

func (l *LambdaExpression) isExpression() {}

// MatchExpression represents pattern matching.
type MatchExpression struct {
	Expression Expression
	Arms       []MatchArm

	// Position information
	Position *Position
}

func (m *MatchExpression) isExpression() {}

// MatchArm represents a single arm in a match expression.
type MatchArm struct {
	Pattern    Pattern
	Expression Expression
}

// Pattern represents a pattern in pattern matching.
type Pattern struct {
	Constructor string
	Variable    string
	Fields      []string  // For field destructuring like Ok { value }
	Nested      []Pattern // For nested destructuring patterns
	IsWildcard  bool      // For _ patterns
}

// ResultExpression represents Result<T, E> construction.
type ResultExpression struct {
	Success   bool       // true for Success, false for Error
	Value     Expression // Value for Success or Error
	ErrorType string     // Type of error (e.g., "DivisionByZero")

	// Position information
	Position *Position
}

func (r *ResultExpression) isExpression() {}

// FieldAccessExpression represents field access like obj.field.
type FieldAccessExpression struct {
	Object    Expression
	FieldName string

	// Position information
	Position *Position
}

func (f *FieldAccessExpression) isExpression() {}

// ModuleAccessExpression represents module member access.
type ModuleAccessExpression struct {
	ModuleName     string
	MemberName     string
	Arguments      []Expression    // For function calls
	NamedArguments []NamedArgument // For named arguments

	// Position information
	Position *Position
}

func (m *ModuleAccessExpression) isExpression() {}

// ModuleDeclaration represents a module declaration.
type ModuleDeclaration struct {
	Name       string
	Statements []Statement
}

func (m *ModuleDeclaration) isStatement() {}

// SpawnExpression represents spawning a fiber.
type SpawnExpression struct {
	Expression Expression

	// Position information
	Position *Position
}

func (s *SpawnExpression) isExpression() {}

// AwaitExpression represents awaiting a fiber result.
type AwaitExpression struct {
	Expression Expression

	// Position information
	Position *Position
}

func (a *AwaitExpression) isExpression() {}

// YieldExpression represents yielding in a fiber.
type YieldExpression struct {
	Value Expression // Optional value to yield

	// Position information
	Position *Position
}

func (y *YieldExpression) isExpression() {}

// ChannelExpression represents channel creation.
type ChannelExpression struct {
	ElementType TypeExpression
	Capacity    Expression // Optional capacity expression
}

func (c *ChannelExpression) isExpression() {}

// ChannelSendExpression represents sending to a channel.
type ChannelSendExpression struct {
	Channel Expression
	Value   Expression
}

func (c *ChannelSendExpression) isExpression() {}

// ChannelRecvExpression represents receiving from a channel.
type ChannelRecvExpression struct {
	Channel Expression
}

func (c *ChannelRecvExpression) isExpression() {}

// SelectExpression represents select statements for channel operations.
type SelectExpression struct {
	Arms []SelectArm
}

func (s *SelectExpression) isExpression() {}

// ChannelCreateExpression represents creating a channel.
type ChannelCreateExpression struct {
	Capacity Expression
}

func (c *ChannelCreateExpression) isExpression() {}

// TypeConstructorExpression represents type construction.
type TypeConstructorExpression struct {
	TypeName string
	Fields   map[string]Expression

	// Position information
	Position *Position
}

func (t *TypeConstructorExpression) isExpression() {}

// SelectArm represents an arm in a select expression.
type SelectArm struct {
	Pattern    Pattern
	Operation  Expression // The channel operation (send/recv)
	Expression Expression // The result expression
}

// BlockExpression represents a block with statements and optional return expression.
type BlockExpression struct {
	Statements []Statement
	Expression Expression // Optional return expression
}

func (b *BlockExpression) isExpression() {}

// ListLiteral represents a list literal like [1, 2, 3].
type ListLiteral struct {
	Elements []Expression

	// Position information
	Position *Position
}

func (l *ListLiteral) isExpression() {}

// ObjectLiteral represents an object literal like { field: value }.
type ObjectLiteral struct {
	Fields map[string]Expression

	// Position information
	Position *Position
}

func (o *ObjectLiteral) isExpression() {}

// ListAccessExpression represents safe list access like list[0] -> Result<T, IndexError>.
type ListAccessExpression struct {
	List  Expression
	Index Expression

	// Position information
	Position *Position
}

func (l *ListAccessExpression) isExpression() {}

// UpdateExpression represents non-destructive updates.
type UpdateExpression struct {
	Target Expression
	Fields map[string]Expression

	// Position information
	Position *Position
}

func (u *UpdateExpression) isExpression() {}
