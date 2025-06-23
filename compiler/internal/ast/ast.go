// Package ast provides Abstract Syntax Tree definitions for the Osprey language.
package ast

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
}

func (l *LetDeclaration) isStatement() {}

// FunctionDeclaration represents a function declaration.
type FunctionDeclaration struct {
	Name       string
	Parameters []Parameter     // Updated to support type annotations
	ReturnType *TypeExpression // Optional return type annotation
	Body       Expression
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
	Name       string
	TypeParams []string // Generic type parameters
	Variants   []TypeVariant
}

func (t *TypeDeclaration) isStatement() {}

// TypeVariant represents a variant in a union type.
type TypeVariant struct {
	Name   string
	Fields []TypeField
}

// TypeField represents a field in a type variant.
type TypeField struct {
	Name       string
	Type       string
	Constraint *FunctionCallExpression // Optional WHERE constraint
}

// FunctionCallExpression represents a function call in constraints.
type FunctionCallExpression struct {
	Function  string
	Arguments []Expression
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
}

func (i *IntegerLiteral) isExpression() {}

// StringLiteral represents a string literal.
type StringLiteral struct {
	Value string
}

func (s *StringLiteral) isExpression() {}

// BooleanLiteral represents a boolean literal.
type BooleanLiteral struct {
	Value bool
}

func (b *BooleanLiteral) isExpression() {}

// InterpolatedStringLiteral represents an interpolated string.
type InterpolatedStringLiteral struct {
	Parts []InterpolatedPart
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
}

func (i *Identifier) isExpression() {}

// BinaryExpression represents a binary expression.
type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (b *BinaryExpression) isExpression() {}

// UnaryExpression represents a unary expression.
type UnaryExpression struct {
	Operator string
	Operand  Expression
}

func (u *UnaryExpression) isExpression() {}

// CallExpression represents a function call (with optional parentheses).
type CallExpression struct {
	Function       Expression
	Arguments      []Expression
	HasParentheses bool // Whether call uses parentheses
	NamedArguments []NamedArgument
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
}

func (m *MethodCallExpression) isExpression() {}

// LambdaExpression represents anonymous functions.
type LambdaExpression struct {
	Parameters []Parameter     // Updated to support type annotations
	ReturnType *TypeExpression // Optional return type annotation
	Body       Expression
}

func (l *LambdaExpression) isExpression() {}

// MatchExpression represents pattern matching.
type MatchExpression struct {
	Expression Expression
	Arms       []MatchArm
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
}

func (r *ResultExpression) isExpression() {}

// FieldAccessExpression represents field access like obj.field.
type FieldAccessExpression struct {
	Object    Expression
	FieldName string
}

func (f *FieldAccessExpression) isExpression() {}

// ModuleAccessExpression represents module member access.
type ModuleAccessExpression struct {
	ModuleName     string
	MemberName     string
	Arguments      []Expression    // For function calls
	NamedArguments []NamedArgument // For named arguments
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
}

func (s *SpawnExpression) isExpression() {}

// AwaitExpression represents awaiting a fiber result.
type AwaitExpression struct {
	Expression Expression
}

func (a *AwaitExpression) isExpression() {}

// YieldExpression represents yielding in a fiber.
type YieldExpression struct {
	Value Expression // Optional value to yield
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
}

func (l *ListLiteral) isExpression() {}

// ListAccessExpression represents safe list access like list[0] -> Result<T, IndexError>.
type ListAccessExpression struct {
	List  Expression
	Index Expression
}

func (l *ListAccessExpression) isExpression() {}

// UpdateExpression represents non-destructive updates.
type UpdateExpression struct {
	Target Expression
	Fields map[string]Expression
}

func (u *UpdateExpression) isExpression() {}
