//! The **ML flavor** frontend: a layout-based, curry-by-default source surface
//! that lowers to the same canonical [`osprey_ast::Program`] as the Default
//! frontend, after which the two are indistinguishable ([FLAVOR-BOUNDARY]).
//!
//! Surface reference: `docs/specs/0024-MLFlavorSyntax.md`. Boundary and
//! lowering contract: `docs/specs/0023-LanguageFlavors.md`. Build sequence:
//! `docs/plans/0013-ml-flavor-frontend.md`.
//!
//! The layout-sensitive grammar (external INDENT/DEDENT/NEWLINE scanner, plan
//! phase 2) and the `MlLowerer` (phase 3) are not yet wired. Until they land,
//! [`parse_ml`] fails loudly rather than silently misparsing ML source as
//! Default — there are no silent placeholders.

use crate::{Flavor, Parsed, SyntaxError};
use osprey_ast::{Position, Program};

/// Parse ML-flavor source into the canonical [`Program`].
///
/// Phase-1 seam: returns a single explanatory syntax error so selecting the ML
/// flavor today is a clear, actionable failure, not a corrupt parse. Phases 2–3
/// replace this body with the real layout parse + lowering.
pub(crate) fn parse_ml(_source: &str) -> Parsed {
    Parsed {
        program: Program {
            statements: Vec::new(),
        },
        errors: vec![SyntaxError {
            message: "the ML flavor frontend is not yet implemented \
                      (docs/plans/0013-ml-flavor-frontend.md, phases 2–3); \
                      use the Default flavor or omit --flavor"
                .to_owned(),
            position: Position { line: 1, column: 0 },
        }],
        flavor: Flavor::Ml,
    }
}
