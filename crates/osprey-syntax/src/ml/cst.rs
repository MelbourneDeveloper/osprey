//! The ML **concrete syntax tree**: a faithful record of the ML surface, with
//! no canonicalisation applied. Currying is still a flat parameter/argument
//! list, pipes are still binary operators, parentheses are still present, and a
//! record literal is still its own node. The CST→AST lowering ([`super::lower`])
//! is the *only* place these are normalised into the canonical
//! [`osprey_ast`] — keeping parse and lower cleanly separated
//! ([FLAVOR-FRONTEND], docs/specs/0023-LanguageFlavors.md).
//!
//! This separation is deliberate: the parser ([`super::parser`]) decides only
//! *what was written*; the lowerer decides *what it means*. Nothing in this
//! module references `osprey_ast`.

use osprey_ast::Position;

/// A top-level item or a statement inside a layout block.
#[derive(Debug, Clone, PartialEq)]
pub(crate) enum MlItem {
    /// `mut? name param* = body`. Zero params ⇒ a value binding; one or more
    /// (including the unit marker) ⇒ a function definition. Currying is not yet
    /// applied — `params` is the flat surface list.
    Binding {
        /// Whether `mut` introduced the binding.
        mutable: bool,
        /// The bound name.
        name: String,
        /// The surface parameter list (empty for a value binding).
        params: Vec<MlParam>,
        /// The right-hand side.
        body: MlExpr,
        /// Source position of the name.
        pos: Position,
    },
    /// `name := value` — mutation of an existing binding.
    Assign {
        /// The mutated name.
        name: String,
        /// The new value.
        value: MlExpr,
        /// Source position of the name.
        pos: Position,
    },
    /// `name : type` — a standalone type signature, paired with the binding of
    /// the same name that follows it. Kept in the CST so the lowerer can apply
    /// concrete parameter/return types (which the type checker and codegen rely
    /// on for curried closures and `Result` auto-unwrap).
    Signature {
        /// The signed name.
        name: String,
        /// The declared type.
        ty: MlType,
    },
    /// `type Name param* =` + an indented layout block of variants
    /// ([FLAVOR-ML-TYPE]). A union/enum lists uppercase constructor variants
    /// (each with an optional indented `field : type` block); a record is the
    /// single-variant form whose lines are lowercase `field : type` — the lowerer
    /// gives that variant the type's own name, matching the Default record shape.
    Type {
        /// The type's name.
        name: String,
        /// Type parameters between the name and `=` (e.g. `T`), in order.
        type_params: Vec<String>,
        /// The declared variants (one per constructor; a record has exactly one).
        variants: Vec<MlVariant>,
        /// Source position of the `type` keyword.
        pos: Position,
    },
    /// `extern name (pname : ptype)* -> rettype` — an external (FFI) function
    /// declaration ([FLAVOR-ML-EXTERN]). Each parameter is a parenthesised
    /// `name : type`; the trailing `-> type` is the return type.
    Extern {
        /// The external symbol name.
        name: String,
        /// The typed parameters, in declaration order.
        params: Vec<MlExternParam>,
        /// The declared return type, if any.
        return_type: Option<MlType>,
        /// Source position of the `extern` keyword.
        pos: Position,
    },
    /// A bare expression evaluated for its effect or trailing value.
    Expr {
        /// The expression.
        value: MlExpr,
        /// Source position.
        pos: Position,
    },
}

/// A parenthesised `name : type` parameter of an `extern` declaration.
#[derive(Debug, Clone, PartialEq)]
pub(crate) struct MlExternParam {
    /// The parameter name.
    pub name: String,
    /// The parameter's declared type.
    pub ty: MlType,
}

/// One variant of a `type` declaration: a constructor name and its payload
/// fields (empty for a bare enum case like `Active`).
#[derive(Debug, Clone, PartialEq)]
pub(crate) struct MlVariant {
    /// The constructor name.
    pub name: String,
    /// The payload fields, in declaration order.
    pub fields: Vec<MlTypeField>,
}

/// A `field : type` line inside a variant's payload block.
#[derive(Debug, Clone, PartialEq)]
pub(crate) struct MlTypeField {
    /// The field name.
    pub name: String,
    /// The field's declared type.
    pub ty: MlType,
}

/// An ML type expression. Arrows are right-associative; application binds
/// tighter (`Handler Db`, `Result int string`) ([FLAVOR-ML-FN]).
#[derive(Debug, Clone, PartialEq)]
pub(crate) enum MlType {
    /// A bare type name (`int`, `string`, `Unit`, a user type).
    Name(String),
    /// Type application `head arg…` (`Handler Db`, `Result int string`).
    App {
        /// The head type name.
        head: String,
        /// The applied argument types.
        args: Vec<MlType>,
    },
    /// `a -> b` (right-associative).
    Arrow {
        /// The argument type.
        from: Box<MlType>,
        /// The result type.
        to: Box<MlType>,
    },
    /// `(a, b, …)` a tupled single argument.
    Tuple(Vec<MlType>),
}

/// A surface parameter pattern in a binding or lambda head.
#[derive(Debug, Clone, PartialEq)]
pub(crate) enum MlParam {
    /// A named parameter.
    Named(String),
    /// The unit marker `()` — a zero-argument function boundary, not a value.
    Unit,
}

/// An ML expression, recorded exactly as written.
#[derive(Debug, Clone, PartialEq)]
pub(crate) enum MlExpr {
    /// Integer literal.
    Int(i64),
    /// Float literal.
    Float(f64),
    /// Boolean literal.
    Bool(bool),
    /// Raw string literal text (quotes/escapes/`${…}` unresolved).
    Str(String),
    /// Identifier or constructor reference.
    Ident(String),
    /// Prefix unary (`-x`, `!x`).
    Unary {
        /// Operator spelling.
        op: String,
        /// The operand.
        operand: Box<MlExpr>,
    },
    /// Binary operator, including the pipe `|>` (the lowerer desugars pipes).
    Binary {
        /// Operator spelling.
        op: String,
        /// Left operand.
        left: Box<MlExpr>,
        /// Right operand.
        right: Box<MlExpr>,
    },
    /// Single-argument application `func arg` (the surface curried form).
    App {
        /// The applied expression.
        func: Box<MlExpr>,
        /// The single argument.
        arg: Box<MlExpr>,
    },
    /// Zero-argument application `func ()`.
    UnitApp {
        /// The applied expression.
        func: Box<MlExpr>,
    },
    /// `target.field` access.
    Field {
        /// The receiver.
        target: Box<MlExpr>,
        /// The field name.
        name: String,
    },
    /// `[ a, b, c ]` list literal (possibly empty).
    List(Vec<MlExpr>),
    /// `target[index]` — a glued postfix index (list/map lookup, returns
    /// `Result`). Only formed when the `[` abuts the target with no space.
    Index {
        /// The indexed expression.
        target: Box<MlExpr>,
        /// The index/key expression.
        index: Box<MlExpr>,
    },
    /// `( inner )` — grouping kept in the CST; the lowerer unwraps it.
    Paren(Box<MlExpr>),
    /// `\param* => body` lambda (curried in the lowerer).
    Lambda {
        /// The surface parameter list.
        params: Vec<MlParam>,
        /// The lambda body.
        body: Box<MlExpr>,
        /// Source position.
        pos: Position,
    },
    /// `match scrutinee` + indented arms.
    Match {
        /// The scrutinee.
        scrutinee: Box<MlExpr>,
        /// The arms.
        arms: Vec<MlArm>,
    },
    /// Constructor record literal `Name` + indented `field = value` lines.
    Record {
        /// Constructor/type name.
        name: String,
        /// Field initialisers.
        fields: Vec<MlField>,
    },
    /// A layout block: leading items and an optional trailing value expression.
    Block {
        /// Statements before the trailing value.
        items: Vec<MlItem>,
        /// The trailing value expression, if any.
        value: Option<Box<MlExpr>>,
    },
    /// `spawn body` — start a fiber whose body (an indented block or inline
    /// expression) runs concurrently ([FLAVOR-ML-SPAWN]).
    Spawn(Box<MlExpr>),
}

/// One `pattern => body` arm of a `match`.
#[derive(Debug, Clone, PartialEq)]
pub(crate) struct MlArm {
    /// The arm pattern.
    pub pattern: MlPattern,
    /// The arm body.
    pub body: MlExpr,
}

/// An ML match pattern.
#[derive(Debug, Clone, PartialEq)]
pub(crate) enum MlPattern {
    /// `_`.
    Wildcard,
    /// An integer literal pattern.
    Int(i64),
    /// A string literal pattern (raw).
    Str(String),
    /// A boolean literal pattern.
    Bool(bool),
    /// `Ctor field*` — a constructor binding zero or more payload fields.
    Ctor {
        /// Constructor name.
        name: String,
        /// Bound field names.
        fields: Vec<String>,
    },
    /// A bare lowercase binding.
    Bind(String),
}

/// A `field = value` initialiser inside a record literal.
#[derive(Debug, Clone, PartialEq)]
pub(crate) struct MlField {
    /// The field name.
    pub name: String,
    /// The field value.
    pub value: MlExpr,
}
