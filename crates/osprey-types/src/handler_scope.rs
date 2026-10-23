//! Handler lifetime clauses the checker types beside the arms: a `finally`
//! cleanup runs outside its own activation and returns `Unit`. Implements
//! [EFFECTS-FINALIZATION] (docs/specs/language/0017-AlgebraicEffects.md).

use crate::check::Checker;
use crate::env::TypeEnv;
use crate::error::TypeError;
use crate::ty::Type;
use osprey_ast::{Expr, Position};

impl Checker {
    /// Type `cleanup` in the handler's enclosing environment. It cannot resume
    /// the dropped continuation, so no enclosing arm's `resume` is visible to
    /// it, and it cannot replace the handler's answer, so it must be `Unit`.
    pub(crate) fn infer_cleanup(
        &mut self,
        cleanup: &Expr,
        position: Option<Position>,
        env: &TypeEnv,
    ) {
        let outer_resume = std::mem::take(&mut self.resume_ctx);
        let ty = self.infer_expr(cleanup, env);
        self.resume_ctx = outer_resume;
        if crate::unify::unify(&mut self.ctx, &Type::unit(), &ty).is_err() {
            let found = self.ctx.prune(&ty);
            self.errors.push(
                TypeError::new(format!(
                    "a handler's `finally` clause must be `Unit`, found `{found}`: \
cleanup runs after the answer is decided and cannot replace it"
                ))
                .with_pos(position),
            );
        }
    }
}
