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

// PluginFunctionDeclaration represents a language plugin function declaration.
type PluginFunctionDeclaration struct {
	PluginName       string          // e.g., "sql", "graphql", "template"
	FunctionName     string          // Function name
	Parameters       []Parameter     // Function parameters
	ReturnType       *TypeExpression // Return type specification
	IsTypeGeneration bool            // true for "as Type", false for "-> Type"
	PluginContent    string          // Raw plugin language content (SQL, GraphQL, etc.)
}

func (p *PluginFunctionDeclaration) isStatement() {}

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
	Function  Expression
	Arguments []Expression
	// For named arguments (multi-parameter functions)
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
	Object     Expression
	MethodName string
	Arguments  []Expression
	// For named arguments
	NamedArguments []NamedArgument
}

func (m *MethodCallExpression) isExpression() {}

// LambdaExpression represents anonymous functions.
type LambdaExpression struct {
	Parameters []Parameter     // Updated to support type annotations
	ReturnType *TypeExpression // Optional return type
	Body       Expression
}

func (l *LambdaExpression) isExpression() {}

// MatchExpression represents pattern matching.
type MatchExpression struct {
	Expression Expression
	Arms       []MatchArm
}

func (m *MatchExpression) isExpression() {}

// MatchArm represents an arm in a match expression.
type MatchArm struct {
	Pattern    Pattern
	Expression Expression
}

// Pattern represents a pattern in match expressions.
type Pattern struct {
	Constructor string
	Variable    string
	Fields      []string  // For field destructuring like Ok { value }
	Nested      []Pattern // For nested destructuring patterns
	IsWildcard  bool      // For _ patterns
}

// ResultExpression represents a result type for arithmetic operations.
type ResultExpression struct {
	IsSuccess bool
	Value     Expression // The actual value for Success, error message for Err
	ErrorType string     // Type of error (e.g., "DivisionByZero")
}

func (r *ResultExpression) isExpression() {}

// FieldAccessExpression represents field access like obj.field.
type FieldAccessExpression struct {
	Object    Expression
	FieldName string
}

func (f *FieldAccessExpression) isExpression() {}

// ModuleAccessExpression represents accessing a member of a module like Module.function.
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

// YieldExpression represents yielding control to the scheduler.
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

// ChannelSendExpression represents sending data through a channel.
type ChannelSendExpression struct {
	Channel Expression
	Value   Expression
}

func (c *ChannelSendExpression) isExpression() {}

// ChannelRecvExpression represents receiving data from a channel.
type ChannelRecvExpression struct {
	Channel Expression
}

func (c *ChannelRecvExpression) isExpression() {}

// SelectExpression represents a select statement for channel operations.
type SelectExpression struct {
	Arms []SelectArm
}

// ChannelCreateExpression represents creating a channel with capacity.
type ChannelCreateExpression struct {
	Capacity Expression
}

// TypeConstructorExpression represents generic type construction like TypeName { field: value }.
type TypeConstructorExpression struct {
	TypeName string
	Fields   map[string]Expression
}

// SelectArm represents a single arm in a select expression.
type SelectArm struct {
	Pattern    Pattern
	Operation  Expression // The channel operation (send/recv)
	Expression Expression // The result expression
}

func (s *SelectExpression) isExpression() {}

func (c *ChannelCreateExpression) isExpression() {}

func (t *TypeConstructorExpression) isExpression() {}

// BlockExpression represents a block of statements followed by an optional expression.
type BlockExpression struct {
	Statements []Statement
	Expression Expression // Optional return expression
}

func (b *BlockExpression) isExpression() {}
