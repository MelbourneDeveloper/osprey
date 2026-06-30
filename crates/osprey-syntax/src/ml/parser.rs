//! The ML-flavor parser: a hand-written **recursive-descent** parser with a
//! **Pratt / precedence-climbing** expression core, run over the layout-resolved
//! token stream from [`super::lexer`]. It produces the ML **concrete syntax
//! tree** ([`super::cst`]) and nothing else — every canonicalisation (currying,
//! pipe desugaring, record/block normalisation, string interpolation) is the
//! lowerer's job ([`super::lower`]). This keeps a clean parse/lower seam: the
//! parser decides *what was written*, the lowerer decides *what it means*
//! ([FLAVOR-FRONTEND], docs/specs/0023-LanguageFlavors.md).
//!
//! ## Design, and the authorities it follows
//!
//! The expression grammar is parsed by binding powers in one driving loop
//! ([`Parser::expr`]) rather than one routine per precedence level. This is
//! Pratt's *top-down operator precedence*; precedence climbing is the same
//! algorithm phrased with explicit minimum-binding-power, so the two names
//! describe one technique. The statement grammar is straight predictive
//! recursive descent. Layout (`Indent`/`Dedent`/`Newline`) is the offside rule,
//! resolved in the lexer and consumed here as ordinary tokens.
//!
//! References (verified 2026-06-30):
//! - V. R. Pratt, "Top Down Operator Precedence", POPL 1973, pp. 41–51.
//!   DOI <https://doi.org/10.1145/512927.512931>. The origin of binding-power
//!   expression parsing used by [`Parser::expr`].
//! - T. Norvell, "Parsing Expressions by Recursive Descent", Memorial Univ.,
//!   1999. <https://www.engr.mun.ca/~theo/Misc/exp_parsing.htm>. Establishes
//!   precedence climbing (origin: M. Richards / K. Clarke) and that it "is a
//!   special case of … Pratt parsing".
//! - A. V. Aho, M. S. Lam, R. Sethi, J. D. Ullman, *Compilers: Principles,
//!   Techniques, and Tools*, 2nd ed., 2006, ISBN 978-0-321-48681-3, ch. 4 §4.4
//!   (recursive-descent / predictive parsing) and §4.1.3–4.1.4 (error recovery:
//!   panic-mode, used by [`Parser::recover`]).
//! - P. J. Landin, "The Next 700 Programming Languages", CACM 9(3), 1966,
//!   pp. 157–166. DOI <https://doi.org/10.1145/365230.365257>. Origin of the
//!   offside rule the layout lexer implements ([FLAVOR-ML-LAYOUT]).
//! - *Haskell 2010 Report*, ch. 10 §10.3 "Layout".
//!   <https://www.haskell.org/onlinereport/haskell2010/haskellch10.html>. A
//!   concrete authoritative spec of layout-driven token insertion.

use super::cst::{
    MlArm, MlExpr, MlExternParam, MlField, MlItem, MlParam, MlPattern, MlType, MlTypeField,
    MlVariant,
};
use super::lexer::lex;
use super::token::{TokKind, Token};
use crate::SyntaxError;
use osprey_ast::Position;

/// Parse ML-flavor `source` into the ML CST plus any syntax errors. Best-effort:
/// errors never abort the parse ([FLAVOR-LOWER-CONTRACT]).
pub(crate) fn parse(source: &str) -> (Vec<MlItem>, Vec<SyntaxError>) {
    let (tokens, mut errors) = lex(source);
    let items = {
        let mut parser = Parser {
            toks: &tokens,
            i: 0,
            errors: &mut errors,
        };
        parser.program()
    };
    (items, errors)
}

/// Binding powers, mirroring the Default grammar's precedence table so equal
/// programs in either flavor produce the same canonical AST (higher binds
/// tighter): or < and < compare < add < mul < pipe. Application (whitespace)
/// and prefix unary bind tighter still and are handled structurally.
fn infix_bp(op: &str) -> Option<u8> {
    let bp = match op {
        "||" => 2,
        "&&" => 3,
        "==" | "!=" | "<" | ">" | "<=" | ">=" => 4,
        "+" | "-" => 5,
        "*" | "/" | "%" => 6,
        "|>" => 8,
        _ => return None,
    };
    Some(bp)
}

/// Recursive-descent + Pratt parser over the layout-resolved token slice.
struct Parser<'t> {
    toks: &'t [Token],
    i: usize,
    errors: &'t mut Vec<SyntaxError>,
}

impl Parser<'_> {
    fn peek(&self) -> &TokKind {
        self.toks.get(self.i).map_or(&TokKind::Eof, |t| &t.kind)
    }

    fn peek_at(&self, ahead: usize) -> &TokKind {
        self.toks
            .get(self.i + ahead)
            .map_or(&TokKind::Eof, |t| &t.kind)
    }

    fn pos(&self) -> Position {
        self.toks.get(self.i).map_or(Position::default(), |t| t.pos)
    }

    /// Consume the current token, discarding it (callers peek first when they
    /// need its payload).
    fn advance(&mut self) {
        if self.i < self.toks.len() {
            self.i += 1;
        }
    }

    fn eat(&mut self, kind: &TokKind) -> bool {
        if self.peek() == kind {
            self.i += 1;
            true
        } else {
            false
        }
    }

    fn error(&mut self, message: impl Into<String>) {
        let position = self.pos();
        self.errors.push(SyntaxError {
            message: message.into(),
            position,
        });
    }

    /// Panic-mode recovery (Dragon Book §4.1.4): drop tokens up to the next
    /// statement separator so one bad line cannot derail the rest.
    fn recover(&mut self) {
        while !matches!(
            self.peek(),
            TokKind::Newline | TokKind::Dedent | TokKind::Eof
        ) {
            self.i += 1;
        }
    }

    fn skip_separators(&mut self) {
        while matches!(self.peek(), TokKind::Newline) {
            self.i += 1;
        }
    }

    fn at_block_end(&self) -> bool {
        matches!(self.peek(), TokKind::Dedent | TokKind::Eof)
    }

    // --- statements -------------------------------------------------------

    fn program(&mut self) -> Vec<MlItem> {
        let mut out = Vec::new();
        loop {
            self.skip_separators();
            if matches!(self.peek(), TokKind::Eof) {
                break;
            }
            match self.item() {
                Some(item) => out.push(item),
                None => self.recover(),
            }
        }
        out
    }

    /// Parse one item, or `None` for a skipped signature line or a recoverable
    /// error.
    fn item(&mut self) -> Option<MlItem> {
        match self.peek() {
            TokKind::KwMut => self.mut_binding(),
            TokKind::KwType => self.type_decl(),
            TokKind::KwExtern => self.extern_decl(),
            TokKind::Reserved(word) => {
                let word = word.clone();
                self.error(format!(
                    "ML construct '{word}' is not yet supported (effects/handlers \
                     are plan 0013 phase 0); use the Default flavor for now"
                ));
                None
            }
            TokKind::Ident(_) => self.ident_item(),
            _ => Some(self.expr_item()),
        }
    }

    /// `mut name = body` → a mutable binding.
    fn mut_binding(&mut self) -> Option<MlItem> {
        let pos = self.pos();
        self.advance(); // `mut`
        let name = self.ident()?;
        let _ = self.expect_eq();
        let body = self.body_after_eq();
        Some(MlItem::Binding {
            mutable: true,
            name,
            params: Vec::new(),
            body,
            pos,
        })
    }

    /// `type Name param* =` + an indented block of variants ([FLAVOR-ML-TYPE]).
    /// A union/enum lists uppercase constructor lines (each with an optional
    /// nested `field : type` block); a record is the single-variant form whose
    /// first block line is a lowercase `field : type`, in which case the lone
    /// variant takes the type's own name (matching the Default record shape).
    fn type_decl(&mut self) -> Option<MlItem> {
        let pos = self.pos();
        self.advance(); // `type`
        let name = self.ident()?;
        let type_params = self.type_params();
        let _ = self.expect_eq();
        let variants = self.type_body(&name);
        Some(MlItem::Type {
            name,
            type_params,
            variants,
            pos,
        })
    }

    /// Bare type-parameter names between the type name and `=` (e.g. `T` in
    /// `type Box T = …`), in order.
    fn type_params(&mut self) -> Vec<String> {
        let mut out = Vec::new();
        while let TokKind::Ident(name) = self.peek() {
            out.push(name.clone());
            self.advance();
        }
        out
    }

    /// The indented body of a `type`. If the first non-blank line is a lowercase
    /// `field : type`, the whole block is one record variant named after the
    /// type; otherwise each uppercase line is a union/enum constructor variant.
    fn type_body(&mut self, type_name: &str) -> Vec<MlVariant> {
        if !self.eat(&TokKind::Indent) {
            return Vec::new();
        }
        self.skip_separators();
        let variants = if self.at_record_field() {
            let fields = self.type_fields();
            vec![MlVariant {
                name: type_name.to_owned(),
                fields,
            }]
        } else {
            self.union_variants()
        };
        let _ = self.eat(&TokKind::Dedent);
        variants
    }

    /// Whether the current block line is a record field `name : type` (a
    /// lowercase identifier directly followed by `:`), versus a constructor line.
    fn at_record_field(&self) -> bool {
        matches!(self.peek(), TokKind::Ident(name) if !is_constructor(name))
            && matches!(self.peek_at(1), TokKind::Colon)
    }

    /// The uppercase constructor variants of a union/enum, each optionally
    /// followed by an indented `field : type` payload block.
    fn union_variants(&mut self) -> Vec<MlVariant> {
        let mut variants = Vec::new();
        while !self.at_block_end() {
            self.skip_separators();
            if self.at_block_end() {
                break;
            }
            let before = self.i;
            match self.ident() {
                Some(name) => {
                    let fields = if matches!(self.peek(), TokKind::Indent) {
                        self.advance(); // `Indent`
                        let fields = self.type_fields();
                        let _ = self.eat(&TokKind::Dedent);
                        fields
                    } else {
                        Vec::new()
                    };
                    variants.push(MlVariant { name, fields });
                }
                None => self.recover(),
            }
            if self.i == before {
                self.recover();
            }
        }
        variants
    }

    /// A run of `field : type` lines (a variant payload or a record body).
    fn type_fields(&mut self) -> Vec<MlTypeField> {
        let mut fields = Vec::new();
        while !self.at_block_end() {
            self.skip_separators();
            if self.at_block_end() {
                break;
            }
            let before = self.i;
            match self.ident() {
                Some(name) => {
                    if !self.eat(&TokKind::Colon) {
                        self.error("expected ':' in type field");
                    }
                    let ty = self.ty();
                    fields.push(MlTypeField { name, ty });
                }
                None => self.recover(),
            }
            if self.i == before {
                self.recover();
            }
        }
        fields
    }

    /// `extern name (pname : ptype)* -> rettype` — an external (FFI) function
    /// declaration ([FLAVOR-ML-EXTERN]). Each parameter is a parenthesised
    /// `name : type`; an optional trailing `-> type` gives the return type.
    fn extern_decl(&mut self) -> Option<MlItem> {
        let pos = self.pos();
        self.advance(); // `extern`
        let name = self.ident()?;
        let mut params = Vec::new();
        while matches!(self.peek(), TokKind::LParen) {
            if let Some(param) = self.extern_param() {
                params.push(param);
            }
        }
        let return_type = if self.eat(&TokKind::Arrow) {
            Some(self.ty())
        } else {
            None
        };
        Some(MlItem::Extern {
            name,
            params,
            return_type,
            pos,
        })
    }

    /// One `( name : type )` parameter of an `extern` declaration.
    fn extern_param(&mut self) -> Option<MlExternParam> {
        self.advance(); // `(`
        let name = self.ident()?;
        if !self.eat(&TokKind::Colon) {
            self.error("expected ':' in extern parameter");
        }
        let ty = self.ty();
        if !self.eat(&TokKind::RParen) {
            self.error("expected ')'");
        }
        Some(MlExternParam { name, ty })
    }

    /// Dispatch an identifier-led item: signature (skipped), assignment,
    /// binding/function, or a bare expression.
    fn ident_item(&mut self) -> Option<MlItem> {
        match self.peek_at(1) {
            TokKind::Colon => self.signature(),
            TokKind::ColonEq => self.assignment(),
            _ if self.is_binding_head() => self.binding(),
            _ => Some(self.expr_item()),
        }
    }

    /// `name := value` → an assignment.
    fn assignment(&mut self) -> Option<MlItem> {
        let pos = self.pos();
        let name = self.ident()?;
        self.advance(); // `:=`
        let value = self.body_after_eq();
        Some(MlItem::Assign { name, value, pos })
    }

    /// `name : type` → a type signature for the binding that follows.
    fn signature(&mut self) -> Option<MlItem> {
        let name = self.ident()?;
        self.advance(); // `:`
        let ty = self.ty();
        Some(MlItem::Signature { name, ty })
    }

    /// A type: arrows are right-associative (`a -> b -> c` = `a -> (b -> c)`).
    fn ty(&mut self) -> MlType {
        let from = self.ty_app();
        if self.eat(&TokKind::Arrow) {
            return MlType::Arrow {
                from: Box::new(from),
                to: Box::new(self.ty()),
            };
        }
        from
    }

    /// Type application `head arg…` — a head name applied to atom types.
    fn ty_app(&mut self) -> MlType {
        let head = self.ty_atom();
        let mut args = Vec::new();
        while self.starts_ty_atom() {
            args.push(self.ty_atom());
        }
        match head {
            MlType::Name(head) if !args.is_empty() => MlType::App { head, args },
            head => head,
        }
    }

    fn starts_ty_atom(&self) -> bool {
        matches!(self.peek(), TokKind::Ident(_) | TokKind::LParen)
    }

    /// A type atom: a name, or a parenthesised group / tuple.
    fn ty_atom(&mut self) -> MlType {
        match self.peek().clone() {
            TokKind::Ident(name) => {
                self.advance();
                MlType::Name(name)
            }
            TokKind::LParen => self.ty_paren(),
            other => {
                self.error(format!("unexpected token {other:?} in type"));
                MlType::Name("Unit".to_owned())
            }
        }
    }

    /// `( t )` grouping or `( t, t, … )` a tupled argument.
    fn ty_paren(&mut self) -> MlType {
        self.advance(); // `(`
        let mut parts = vec![self.ty()];
        while self.eat(&TokKind::Comma) {
            parts.push(self.ty());
        }
        let _ = self.eat(&TokKind::RParen);
        if parts.len() == 1 {
            parts
                .into_iter()
                .next()
                .unwrap_or(MlType::Name("Unit".to_owned()))
        } else {
            MlType::Tuple(parts)
        }
    }

    /// `name param* = body` → a binding (value when `param*` is empty, function
    /// otherwise). Currying is applied later, in the lowerer.
    fn binding(&mut self) -> Option<MlItem> {
        let pos = self.pos();
        let name = self.ident()?;
        let params = self.params();
        let _ = self.expect_eq();
        let body = self.body_after_eq();
        Some(MlItem::Binding {
            mutable: false,
            name,
            params,
            body,
            pos,
        })
    }

    fn expr_item(&mut self) -> MlItem {
        let pos = self.pos();
        let value = self.expr(0);
        MlItem::Expr { value, pos }
    }

    /// Collect zero or more surface parameter patterns up to the `=`/`=>`.
    fn params(&mut self) -> Vec<MlParam> {
        let mut out = Vec::new();
        loop {
            match self.peek() {
                TokKind::Ident(name) => {
                    let name = name.clone();
                    self.advance();
                    out.push(MlParam::Named(name));
                }
                TokKind::LParen => match self.paren_param() {
                    Some(name) => out.push(MlParam::Named(name)),
                    None => out.push(MlParam::Unit),
                },
                _ => break,
            }
        }
        out
    }

    /// A parenthesised parameter: `()` (the unit marker → `None`) or `(name)`.
    fn paren_param(&mut self) -> Option<String> {
        self.advance(); // `(`
        let name = match self.peek() {
            TokKind::Ident(name) => {
                let name = name.clone();
                self.advance();
                Some(name)
            }
            _ => None,
        };
        let _ = self.eat(&TokKind::RParen);
        name
    }

    /// Lookahead (non-consuming): does the run from the current identifier end
    /// in `=` on this logical line (`Ident paramAtom* =`)?
    fn is_binding_head(&self) -> bool {
        let mut j = self.i + 1; // past the leading identifier
        loop {
            match self.toks.get(j).map(|t| &t.kind) {
                Some(TokKind::Ident(_)) => j += 1,
                Some(TokKind::LParen) => {
                    j += 1;
                    while !matches!(
                        self.toks.get(j).map(|t| &t.kind),
                        Some(TokKind::RParen | TokKind::Eof) | None
                    ) {
                        j += 1;
                    }
                    j += 1; // past `)`
                }
                Some(TokKind::Eq) => return true,
                _ => return false,
            }
        }
    }

    // --- expressions (Pratt) ---------------------------------------------

    /// Parse an expression whose operators bind at least as tightly as `min_bp`
    /// — the driving loop of Pratt / precedence climbing (Pratt 1973; Norvell).
    fn expr(&mut self, min_bp: u8) -> MlExpr {
        let mut left = self.unary();
        while let TokKind::Op(op) = self.peek() {
            let op = op.clone();
            let Some(bp) = infix_bp(&op) else { break };
            if bp < min_bp {
                break;
            }
            self.advance();
            let right = self.expr(bp + 1);
            left = MlExpr::Binary {
                op,
                left: Box::new(left),
                right: Box::new(right),
            };
        }
        left
    }

    /// A prefix unary (`-x`, `!x`) or an application.
    fn unary(&mut self) -> MlExpr {
        if let TokKind::Op(op) = self.peek() {
            if op == "-" || op == "!" {
                let op = op.clone();
                self.advance();
                let operand = self.unary();
                return MlExpr::Unary {
                    op,
                    operand: Box::new(operand),
                };
            }
        }
        self.application()
    }

    /// Whitespace application `f a b`, left-associative, recorded as nested
    /// single-argument [`MlExpr::App`] ([FLAVOR-ML-CALL]).
    fn application(&mut self) -> MlExpr {
        let mut func = self.postfix();
        // `Ctor(field = v, …)` is an inline record literal, not application: an
        // uppercase constructor immediately followed by `(ident = …`. It lowers
        // to the same `MlExpr::Record` the layout form does ([FLAVOR-ML-RECORD]).
        if let MlExpr::Ident(name) = &func {
            if is_constructor(name) && self.at_inline_record() {
                let name = name.clone();
                func = self.inline_record(name);
            }
        }
        // `f ()` is a zero-argument application, not application to unit.
        if matches!(self.peek(), TokKind::LParen) && matches!(self.peek_at(1), TokKind::RParen) {
            self.advance();
            self.advance();
            func = MlExpr::UnitApp {
                func: Box::new(func),
            };
        }
        while self.starts_atom() {
            let arg = self.postfix();
            func = MlExpr::App {
                func: Box::new(func),
                arg: Box::new(arg),
            };
        }
        func
    }

    /// Postfix `.field` access and glued `[index]` chained onto an atom. A `[`
    /// only indexes when it abuts the target (`xs[0]`); a spaced `[` is a list
    /// literal argument, left for [`Self::application`] ([FLAVOR-ML-INDEX]).
    fn postfix(&mut self) -> MlExpr {
        let mut target = self.atom();
        loop {
            if self.eat(&TokKind::Dot) {
                if let Some(name) = self.ident() {
                    target = MlExpr::Field {
                        target: Box::new(target),
                        name,
                    };
                }
            } else if matches!(self.peek(), TokKind::LBracket) && self.glued() {
                target = self.index(target);
            } else {
                return target;
            }
        }
    }

    /// `target[index]` — consume a glued bracket index.
    fn index(&mut self, target: MlExpr) -> MlExpr {
        self.advance(); // `[`
        let index = self.expr(0);
        if !self.eat(&TokKind::RBracket) {
            self.error("expected ']'");
        }
        MlExpr::Index {
            target: Box::new(target),
            index: Box::new(index),
        }
    }

    /// Whether the current token abuts the previous one with no whitespace.
    fn glued(&self) -> bool {
        self.toks.get(self.i).is_some_and(|t| t.glued)
    }

    /// Whether the next token can begin an argument atom.
    fn starts_atom(&self) -> bool {
        matches!(
            self.peek(),
            TokKind::Int(_)
                | TokKind::Float(_)
                | TokKind::Str(_)
                | TokKind::Ident(_)
                | TokKind::KwTrue
                | TokKind::KwFalse
                | TokKind::LParen
                | TokKind::LBracket
        )
    }

    fn atom(&mut self) -> MlExpr {
        match self.peek().clone() {
            TokKind::Int(n) => {
                self.advance();
                MlExpr::Int(n)
            }
            TokKind::Float(f) => {
                self.advance();
                MlExpr::Float(f)
            }
            TokKind::KwTrue => {
                self.advance();
                MlExpr::Bool(true)
            }
            TokKind::KwFalse => {
                self.advance();
                MlExpr::Bool(false)
            }
            TokKind::Str(raw) => {
                self.advance();
                MlExpr::Str(raw)
            }
            TokKind::KwMatch => self.match_expr(),
            TokKind::KwSpawn => self.spawn_expr(),
            TokKind::Backslash => self.lambda(),
            TokKind::LParen => self.paren(),
            TokKind::LBracket => self.list(),
            TokKind::Ident(name) => {
                self.advance();
                self.ident_atom(name)
            }
            other => {
                self.error(format!("unexpected token {other:?} in expression"));
                self.advance();
                MlExpr::Bool(false)
            }
        }
    }

    /// An identifier atom: a bare reference, or — for an uppercase constructor
    /// directly followed by an indented `field = value` block — a record
    /// literal ([FLAVOR-ML-RECORD]).
    fn ident_atom(&mut self, name: String) -> MlExpr {
        if is_constructor(&name) && matches!(self.peek(), TokKind::Indent) {
            let fields = self.record_fields();
            MlExpr::Record { name, fields }
        } else {
            MlExpr::Ident(name)
        }
    }

    /// `[ a, b, c ]` list literal (possibly empty). Layout is suppressed inside
    /// brackets, so elements may span lines ([FLAVOR-ML-LIST]).
    fn list(&mut self) -> MlExpr {
        self.advance(); // `[`
        let mut items = Vec::new();
        if !matches!(self.peek(), TokKind::RBracket) {
            items.push(self.expr(0));
            while self.eat(&TokKind::Comma) {
                if matches!(self.peek(), TokKind::RBracket) {
                    break; // tolerate a trailing comma
                }
                items.push(self.expr(0));
            }
        }
        if !self.eat(&TokKind::RBracket) {
            self.error("expected ']'");
        }
        MlExpr::List(items)
    }

    /// `( expr )` grouping, kept as an [`MlExpr::Paren`] node.
    fn paren(&mut self) -> MlExpr {
        self.advance(); // `(`
        let inner = self.expr(0);
        if !self.eat(&TokKind::RParen) {
            self.error("expected ')'");
        }
        MlExpr::Paren(Box::new(inner))
    }

    /// `\param* => body` lambda.
    fn lambda(&mut self) -> MlExpr {
        let pos = self.pos();
        self.advance(); // `\`
        let params = self.params();
        if !self.eat(&TokKind::FatArrow) {
            self.error("expected '=>' in lambda");
        }
        let body = self.body_after_eq();
        MlExpr::Lambda {
            params,
            body: Box::new(body),
            pos,
        }
    }

    /// `spawn body` — start a fiber. The body is an indented layout block or an
    /// inline expression, parsed exactly like a `=`/`=>` body ([FLAVOR-ML-SPAWN]).
    fn spawn_expr(&mut self) -> MlExpr {
        self.advance(); // `spawn`
        MlExpr::Spawn(Box::new(self.body_after_eq()))
    }

    /// `match scrutinee` + indented `pattern => body` arms.
    fn match_expr(&mut self) -> MlExpr {
        self.advance(); // `match`
        let scrutinee = self.expr(0);
        let mut arms = Vec::new();
        if self.eat(&TokKind::Indent) {
            while !self.at_block_end() {
                self.skip_separators();
                if self.at_block_end() {
                    break;
                }
                arms.push(self.match_arm());
            }
            let _ = self.eat(&TokKind::Dedent);
        }
        MlExpr::Match {
            scrutinee: Box::new(scrutinee),
            arms,
        }
    }

    fn match_arm(&mut self) -> MlArm {
        let pattern = self.pattern();
        if !self.eat(&TokKind::FatArrow) {
            self.error("expected '=>' in match arm");
        }
        let body = self.body_after_eq();
        MlArm { pattern, body }
    }

    /// A match pattern: `_`, a literal, `Ctor field…`, or a bare binding.
    fn pattern(&mut self) -> MlPattern {
        match self.peek().clone() {
            TokKind::Int(n) => {
                self.advance();
                MlPattern::Int(n)
            }
            TokKind::Str(raw) => {
                self.advance();
                MlPattern::Str(raw)
            }
            TokKind::KwTrue => {
                self.advance();
                MlPattern::Bool(true)
            }
            TokKind::KwFalse => {
                self.advance();
                MlPattern::Bool(false)
            }
            TokKind::Ident(name) => {
                self.advance();
                self.ident_pattern(name)
            }
            other => {
                self.error(format!("unexpected token {other:?} in pattern"));
                MlPattern::Wildcard
            }
        }
    }

    /// `_` → wildcard; `Ctor a b` → constructor binding payload fields; a bare
    /// lowercase name → a binding ([FLAVOR-ML-MATCH]).
    fn ident_pattern(&mut self, name: String) -> MlPattern {
        if name == "_" {
            return MlPattern::Wildcard;
        }
        if is_constructor(&name) {
            let mut fields = Vec::new();
            while let TokKind::Ident(field) = self.peek() {
                fields.push(field.clone());
                self.advance();
            }
            return MlPattern::Ctor { name, fields };
        }
        MlPattern::Bind(name)
    }

    /// The indented `field = value` lines of a layout record literal.
    fn record_fields(&mut self) -> Vec<MlField> {
        let mut fields = Vec::new();
        let _ = self.eat(&TokKind::Indent);
        while !self.at_block_end() {
            self.skip_separators();
            if self.at_block_end() {
                break;
            }
            match self.parse_record_field() {
                Some(field) => fields.push(field),
                None => self.recover(),
            }
        }
        let _ = self.eat(&TokKind::Dedent);
        fields
    }

    /// `( field = expr ( , field = expr )* )` — an inline record literal in
    /// expression/argument position ([FLAVOR-ML-RECORD]). Layout is suppressed
    /// inside parens, so the fields are a simple comma list; it lowers to the
    /// same [`MlExpr::Record`] the layout form produces.
    fn inline_record(&mut self, name: String) -> MlExpr {
        self.advance(); // `(`
        let mut fields = Vec::new();
        if !matches!(self.peek(), TokKind::RParen) {
            loop {
                match self.parse_record_field() {
                    Some(field) => fields.push(field),
                    None => self.recover(),
                }
                if !self.eat(&TokKind::Comma) {
                    break;
                }
                if matches!(self.peek(), TokKind::RParen) {
                    break; // tolerate a trailing comma
                }
            }
        }
        if !self.eat(&TokKind::RParen) {
            self.error("expected ')'");
        }
        MlExpr::Record { name, fields }
    }

    /// One `field = value` initialiser, shared by the layout and inline record
    /// forms so neither duplicates the field-parsing rule.
    fn parse_record_field(&mut self) -> Option<MlField> {
        let name = self.ident()?;
        let _ = self.expect_eq();
        let value = self.body_after_eq();
        Some(MlField { name, value })
    }

    /// Whether the current `(` opens an inline record literal — its first two
    /// tokens are `Ident` then `=`. Used to disambiguate `Ctor(field = v)` (a
    /// record) from `Ctor (expr)` (application) and `Ctor ()` (unit application).
    fn at_inline_record(&self) -> bool {
        matches!(self.peek(), TokKind::LParen)
            && matches!(self.peek_at(1), TokKind::Ident(_))
            && matches!(self.peek_at(2), TokKind::Eq)
    }

    // --- bodies and helpers ----------------------------------------------

    /// The body after `=`/`=>`: an inline expression, or an indented layout
    /// block whose trailing expression is its value ([FLAVOR-ML-BLOCK]).
    fn body_after_eq(&mut self) -> MlExpr {
        if !matches!(self.peek(), TokKind::Indent) {
            return self.expr(0);
        }
        self.advance(); // `Indent`
        let (items, value) = self.block_items();
        let _ = self.eat(&TokKind::Dedent);
        MlExpr::Block { items, value }
    }

    /// The items (and optional trailing value) of an indented block.
    fn block_items(&mut self) -> (Vec<MlItem>, Option<Box<MlExpr>>) {
        let mut items = Vec::new();
        let mut value = None;
        while !self.at_block_end() {
            self.skip_separators();
            if self.at_block_end() {
                break;
            }
            let before = self.i;
            value = self.block_line(&mut items);
            // Forward-progress guard ([FLAVOR-LOWER-CONTRACT]): a `block_line`
            // whose `item()` errored without consuming a token — a reserved word
            // (`do`/`effect`/…) or a malformed line inside the block — would
            // otherwise spin this loop forever. Recover past the offending token,
            // exactly as the top-level `program()` loop does, so any input
            // terminates.
            if self.i == before {
                self.recover();
            }
        }
        (items, value)
    }

    /// Parse one block line. A trailing bare expression with nothing after it is
    /// the block value; anything else is appended as an item.
    fn block_line(&mut self, items: &mut Vec<MlItem>) -> Option<Box<MlExpr>> {
        match self.item() {
            Some(MlItem::Expr { value, .. }) if self.at_block_end() => Some(Box::new(value)),
            Some(item) => {
                items.push(item);
                None
            }
            None => None,
        }
    }

    fn ident(&mut self) -> Option<String> {
        if let TokKind::Ident(name) = self.peek() {
            let name = name.clone();
            self.advance();
            Some(name)
        } else {
            self.error("expected an identifier");
            None
        }
    }

    fn expect_eq(&mut self) -> bool {
        if self.eat(&TokKind::Eq) {
            true
        } else {
            self.error("expected '='");
            false
        }
    }
}

/// An uppercase initial marks a constructor/type name; lowercase marks a value
/// binding or variable, mirroring the Default flavor's lexical convention.
fn is_constructor(name: &str) -> bool {
    name.chars().next().is_some_and(char::is_uppercase)
}
