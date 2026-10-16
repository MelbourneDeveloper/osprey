//! `textDocument/implementation` over the effect index: an operation lists
//! the handler arms that answer it, and an effect lists the handlers that
//! discharge it, in the open buffer and then its project siblings. Implements
//! [LSP-IMPLEMENTATIONS-EFFECT-HANDLERS].

use lspkit_vfs::PositionEncoding;
use osprey_syntax::Flavor;

use super::{
    byte_column, effect_index, encoded_column, same_effect, same_operation, target_at,
    trailing_path, Discharge, EffectIndex, OperationId, Site,
};
use crate::features::nth_line;
use crate::model::Location;
use crate::text::{measure, prefix_to, word_at};
use crate::workspace;

/// What an implementation query asks about: one operation, or a whole effect.
enum Target {
    Operation(OperationId),
    Effect(String),
}

/// One indexed file of the project; the open buffer comes first.
struct IndexedFile<'a> {
    index: EffectIndex,
    source: &'a str,
    uri: &'a str,
    flavor: Flavor,
}

/// Every handler arm implementing the operation under the cursor, or every
/// handler discharging the effect under the cursor.
pub(crate) fn implementations(
    source: &str,
    uri: &str,
    line: u32,
    character: u32,
    encoding: PositionEncoding,
    project: &workspace::View,
) -> Vec<Location> {
    let files = indexed_files(source, uri, project);
    implementation_target(&files, line, character, encoding).map_or_else(Vec::new, |target| {
        files
            .iter()
            .flat_map(|file| file.implementations(&target, encoding))
            .collect()
    })
}

/// The open buffer, then every project sibling ([LSP-WORKSPACE]).
fn indexed_files<'a>(
    source: &'a str,
    uri: &'a str,
    project: &'a workspace::View,
) -> Vec<IndexedFile<'a>> {
    let current = IndexedFile {
        index: effect_index(&project.program(uri, source)),
        source,
        uri,
        flavor: project.flavor(uri, source),
    };
    std::iter::once(current)
        .chain(project.siblings.iter().map(|sibling| IndexedFile {
            index: effect_index(&sibling.program),
            source: &sibling.source,
            uri: &sibling.uri,
            flavor: sibling.flavor,
        }))
        .collect()
}

fn implementation_target(
    files: &[IndexedFile<'_>],
    line: u32,
    character: u32,
    encoding: PositionEncoding,
) -> Option<Target> {
    let current = files.first()?;
    target_at(
        &current.index,
        current.source,
        line,
        character,
        encoding,
        current.flavor,
    )
    .map(Target::Operation)
    .or_else(|| effect_at(files, line, character, encoding).map(Target::Effect))
}

/// The effect named under the cursor — at its declaration, a discharge, a
/// `perform` qualifier or anywhere else — when some project file declares or
/// discharges it. A builtin such as `Arith` has no declaration to find.
fn effect_at(
    files: &[IndexedFile<'_>],
    line: u32,
    character: u32,
    encoding: PositionEncoding,
) -> Option<String> {
    let text = nth_line(files.first()?.source, line)?;
    let span = word_at(text, character, encoding)?;
    let mention = trailing_path(prefix_to(text, span.start, encoding)) + &span.word;
    files
        .iter()
        .flat_map(|file| {
            let declared = file.index.effect_symbols.iter().map(|symbol| &symbol.name);
            declared.chain(file.index.discharges.iter().map(|site| &site.effect))
        })
        .any(|effect| same_effect(effect, &mention))
        .then_some(mention)
}

impl IndexedFile<'_> {
    fn implementations(&self, target: &Target, encoding: PositionEncoding) -> Vec<Location> {
        match target {
            Target::Operation(id) => self
                .index
                .handlers
                .iter()
                .filter(|site| same_operation(&site.id, id))
                .filter_map(|site| {
                    site_location(site, self.source, self.uri, encoding, self.flavor)
                })
                .collect(),
            Target::Effect(effect) => self
                .index
                .discharges
                .iter()
                .filter(|discharge| same_effect(&discharge.effect, effect))
                .filter_map(|discharge| self.discharge_location(discharge, encoding))
                .collect(),
        }
    }

    /// The effect's name in the handler head, searched from the `handle`
    /// keyword; a name on a later line lands on the keyword itself.
    fn discharge_location(
        &self,
        discharge: &Discharge,
        encoding: PositionEncoding,
    ) -> Option<Location> {
        let position = discharge.position?;
        let line = position.line.checked_sub(1)?;
        let text = nth_line(self.source, line)?;
        let keyword = byte_column(text, position.column, self.flavor)?;
        let name = text
            .get(keyword..)?
            .find(&discharge.effect)
            .map_or(keyword, |offset| keyword.saturating_add(offset));
        let start = measure(text.get(..name)?, encoding);
        Some(Location {
            uri: self.uri.to_owned(),
            span: (
                line,
                start,
                line,
                start.saturating_add(measure(&discharge.effect, encoding)),
            ),
        })
    }
}

fn site_location(
    site: &Site,
    source: &str,
    uri: &str,
    encoding: PositionEncoding,
    flavor: Flavor,
) -> Option<Location> {
    let position = site.position?;
    let line = position.line.checked_sub(1)?;
    let text = nth_line(source, line)?;
    let start = encoded_column(text, position.column, encoding, flavor)?;
    let end = start.saturating_add(measure(&site.id.operation, encoding));
    Some(Location {
        uri: uri.to_owned(),
        span: (line, start, line, end),
    })
}
