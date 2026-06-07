//! Osprey abstract syntax tree — a Rust port of `compiler/internal/ast/ast.go`.
//!
//! Go used two marker interfaces (`Statement`, `Expression`) plus ~50 structs;
//! Rust models that as two enums (`Stmt`, `Expr`) with struct-like variants, which
//! is both more compact and exhaustively matchable (the type checker and codegen
//! port get compiler-enforced totality for free). Field names track the Go structs.

/// 1-based line, 0-based column — mirrors `ast.Position`.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Default)]
pub struct Position {
    pub line: u32,
    pub column: u32,
}

/// A parsed program: the sequence of top-level statements (`ast.Program`).
#[derive(Debug, Clone, PartialEq, Default)]
pub struct Program {
    pub statements: Vec<Stmt>,
}

/// A type expression — `Result<Int, Error>`, `[String]`, `fn(Int) -> Bool`
/// (`ast.TypeExpression`).
#[derive(Debug, Clone, PartialEq)]
pub struct TypeExpr {
    pub name: String,
    pub generic_params: Vec<TypeExpr>,
    pub is_array: bool,
    pub array_element: Option<Box<TypeExpr>>,
    pub is_function: bool,
    pub parameter_types: Vec<TypeExpr>,
    pub return_type: Option<Box<TypeExpr>>,
    pub position: Option<Position>,
}

impl TypeExpr {
    /// A bare named type like `Int` or `Ptr`.
    pub fn named(name: impl Into<String>) -> Self {
        TypeExpr {
            name: name.into(),
            generic_params: Vec::new(),
            is_array: false,
            array_element: None,
            is_function: false,
            parameter_types: Vec::new(),
            return_type: None,
            position: None,
        }
    }
}

/// A function parameter with an optional type annotation (`ast.Parameter`).
#[derive(Debug, Clone, PartialEq)]
pub struct Parameter {
    pub name: String,
    pub ty: Option<TypeExpr>,
}

/// An `extern fn` parameter — type annotation required (`ast.ExternParameter`).
#[derive(Debug, Clone, PartialEq)]
pub struct ExternParameter {
    pub name: String,
    pub ty: TypeExpr,
}

/// A variant of a union type (`ast.TypeVariant`).
#[derive(Debug, Clone, PartialEq)]
pub struct TypeVariant {
    pub name: String,
    pub fields: Vec<TypeField>,
}

/// A field within a record/variant, with an optional `where` constraint
/// (`ast.TypeField`).
#[derive(Debug, Clone, PartialEq)]
pub struct TypeField {
    pub name: String,
    pub ty: String,
    pub constraint: Option<Box<Expr>>,
}

/// An operation declared inside an `effect` block (`ast.EffectOperation`).
#[derive(Debug, Clone, PartialEq)]
pub struct EffectOperation {
    pub name: String,
    pub ty: String,
    pub parameters: Vec<Parameter>,
    pub return_type: String,
}

/// A statement (`ast.Statement` marker interface).
#[derive(Debug, Clone, PartialEq)]
pub enum Stmt {
    Import {
        module: Vec<String>,
    },
    Let {
        name: String,
        mutable: bool,
        ty: Option<TypeExpr>,
        value: Expr,
        position: Option<Position>,
    },
    Assignment {
        name: String,
        value: Expr,
        position: Option<Position>,
    },
    Function {
        name: String,
        parameters: Vec<Parameter>,
        return_type: Option<TypeExpr>,
        effects: Vec<String>,
        body: Expr,
        position: Option<Position>,
    },
    Extern {
        name: String,
        parameters: Vec<ExternParameter>,
        return_type: Option<TypeExpr>,
    },
    Type {
        name: String,
        type_params: Vec<String>,
        variants: Vec<TypeVariant>,
        validation_func: Option<String>,
    },
    Effect {
        name: String,
        operations: Vec<EffectOperation>,
    },
    Module {
        name: String,
        body: Vec<Stmt>,
    },
    Expr(Expr),
}

/// A named argument `name: value` (`ast.NamedArgument`).
#[derive(Debug, Clone, PartialEq)]
pub struct NamedArgument {
    pub name: String,
    pub value: Expr,
}

/// A part of an interpolated string — literal text or an embedded expression
/// (`ast.InterpolatedPart`).
#[derive(Debug, Clone, PartialEq)]
pub enum InterpolatedPart {
    Text(String),
    Expr(Expr),
}

/// A match arm `pattern => body` (`ast.MatchArm`).
#[derive(Debug, Clone, PartialEq)]
pub struct MatchArm {
    pub pattern: Pattern,
    pub body: Expr,
}

/// A pattern in a match/select arm (`ast.Pattern`).
#[derive(Debug, Clone, PartialEq)]
pub enum Pattern {
    Wildcard,
    Literal(Box<Expr>),
    /// `Ctor` / `Ctor { a, b }` / `Ctor(p, ...)` constructor & destructuring forms.
    Constructor {
        name: String,
        fields: Vec<String>,
        sub_patterns: Vec<Pattern>,
    },
    /// `value: Int` type-annotated binding.
    TypeAnnotated {
        name: String,
        ty: TypeExpr,
    },
    /// `{ name, age }` anonymous structural.
    Structural {
        fields: Vec<String>,
    },
    /// A bare identifier capture.
    Binding(String),
}

/// A field assignment `name: value` in an object/type constructor
/// (`ast.FieldAssignment`-style).
#[derive(Debug, Clone, PartialEq)]
pub struct FieldAssignment {
    pub name: String,
    pub value: Expr,
}

/// A map entry `key: value` (`ast.MapEntry`).
#[derive(Debug, Clone, PartialEq)]
pub struct MapEntry {
    pub key: Expr,
    pub value: Expr,
}

/// An expression (`ast.Expression` marker interface). Boxing breaks the recursive
/// cycle; positions are attached where the Go AST carried them.
#[derive(Debug, Clone, PartialEq)]
pub enum Expr {
    Integer(i64),
    Float(f64),
    Str(String),
    Bool(bool),
    InterpolatedStr(Vec<InterpolatedPart>),
    Identifier(String),
    List(Vec<Expr>),
    Map(Vec<MapEntry>),
    Object(Vec<FieldAssignment>),
    Binary {
        op: String,
        left: Box<Expr>,
        right: Box<Expr>,
    },
    Unary {
        op: String,
        operand: Box<Expr>,
    },
    /// `f(args)` — positional or named (UFCS dispatch resolved later).
    Call {
        function: Box<Expr>,
        arguments: Vec<Expr>,
        named_arguments: Vec<NamedArgument>,
    },
    /// `a |> b` pipe.
    Pipe {
        left: Box<Expr>,
        right: Box<Expr>,
    },
    /// `obj.field` field access.
    FieldAccess {
        target: Box<Expr>,
        field: String,
    },
    /// `obj.method(args)` method call.
    MethodCall {
        target: Box<Expr>,
        method: String,
        arguments: Vec<Expr>,
        named_arguments: Vec<NamedArgument>,
    },
    /// `a[i]` index access (returns `Result`).
    Index {
        target: Box<Expr>,
        index: Box<Expr>,
    },
    Lambda {
        parameters: Vec<Parameter>,
        return_type: Option<TypeExpr>,
        body: Box<Expr>,
    },
    Match {
        value: Box<Expr>,
        arms: Vec<MatchArm>,
    },
    Block {
        statements: Vec<Stmt>,
        /// The trailing value-expression, if any.
        value: Option<Box<Expr>>,
    },
    /// `Type<T> { field: value }` type constructor.
    TypeConstructor {
        name: String,
        type_args: Vec<TypeExpr>,
        fields: Vec<FieldAssignment>,
    },
    /// `record { field: newValue }` non-destructive update.
    Update {
        record: String,
        fields: Vec<FieldAssignment>,
    },
    Spawn(Box<Expr>),
    Yield(Option<Box<Expr>>),
    Await(Box<Expr>),
    Send {
        channel: Box<Expr>,
        value: Box<Expr>,
    },
    Recv(Box<Expr>),
    Select {
        arms: Vec<MatchArm>,
    },
    /// `perform Effect.operation(args)`.
    Perform {
        effect: String,
        operation: String,
        arguments: Vec<Expr>,
        named_arguments: Vec<NamedArgument>,
    },
    /// `handle Effect op params => body ... in body`.
    Handler {
        effect: String,
        arms: Vec<HandlerArm>,
        body: Box<Expr>,
    },
}

/// One arm of a `handle ... in` expression (`ast` handler arm).
#[derive(Debug, Clone, PartialEq)]
pub struct HandlerArm {
    pub operation: String,
    pub params: Vec<String>,
    pub body: Expr,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn builds_a_small_program() {
        let p = Program {
            statements: vec![Stmt::Let {
                name: "x".into(),
                mutable: false,
                ty: None,
                value: Expr::Integer(1),
                position: None,
            }],
        };
        assert_eq!(p.statements.len(), 1);
        match &p.statements[0] {
            Stmt::Let { name, value, .. } => {
                assert_eq!(name, "x");
                assert_eq!(*value, Expr::Integer(1));
            }
            _ => panic!("expected let"),
        }
    }

    #[test]
    fn named_type_helper() {
        let t = TypeExpr::named("Ptr");
        assert_eq!(t.name, "Ptr");
        assert!(!t.is_function);
    }
}
