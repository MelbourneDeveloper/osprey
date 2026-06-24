//! `osprey --docs --docs-dir <dir>`: regenerate the built-in function reference
//! under `<dir>/functions/` straight from the compiler's single source of truth
//! ([`osprey_types::builtin_doc_view`]). Every page's signature, parameter
//! types, and return type come from the real type scheme, so the website docs
//! and the editor hover are guaranteed to show the same thing.
//!
//! The directory is treated as generated output: a page is written for every
//! built-in, an index lists them all, and any stale `*.md` left over from a
//! built-in that no longer exists is pruned.

use osprey_types::{builtin_doc_view, builtin_names, BuiltinDocView};
use std::collections::{HashMap, HashSet};
use std::fs;
use std::path::{Path, PathBuf};
use std::process::ExitCode;

/// Entry point for the `--docs` mode. Reads `--docs-dir <dir>` from `args`.
pub fn run(args: &[String]) -> ExitCode {
    let dir = match docs_dir(args) {
        Some(dir) => PathBuf::from(dir),
        None => {
            eprintln!("usage: osprey --docs --docs-dir <dir>");
            return ExitCode::from(2);
        }
    };
    match generate(&dir) {
        Ok(count) => {
            println!("generated {count} function docs in {}", dir.display());
            ExitCode::SUCCESS
        }
        Err(e) => {
            eprintln!("osprey --docs: {e}");
            ExitCode::FAILURE
        }
    }
}

/// The value following `--docs-dir`, if present.
fn docs_dir(args: &[String]) -> Option<&str> {
    args.iter()
        .position(|a| a == "--docs-dir")
        .and_then(|i| args.get(i + 1))
        .map(String::as_str)
}

/// Write a page per built-in plus the index, then prune stale pages. Returns the
/// number of function pages written.
fn generate(docs_dir: &Path) -> std::io::Result<usize> {
    let functions = docs_dir.join("functions");
    fs::create_dir_all(&functions)?;
    let names = builtin_names();
    let slugs = slug_map(&names);
    let views: Vec<BuiltinDocView> = names.iter().filter_map(|n| builtin_doc_view(n)).collect();
    for view in &views {
        fs::write(
            functions.join(format!("{}.md", slugs[&view.name])),
            page(view),
        )?;
    }
    fs::write(functions.join("index.md"), index(&views, &slugs))?;
    prune(&functions, &slugs)?;
    Ok(views.len())
}

/// Assign each built-in a unique, filesystem-safe page stem. Names normally
/// lowercase to their stem; a capitalized type constructor that would clash with
/// a same-spelled function (`Map` vs `map`) is suffixed `-type`, and any further
/// clash gets a numeric suffix — so distinct built-ins never share a page even
/// on a case-insensitive filesystem.
fn slug_map(names: &[String]) -> HashMap<String, String> {
    let mut counts: HashMap<String, usize> = HashMap::new();
    for n in names {
        *counts.entry(n.to_lowercase()).or_insert(0) += 1;
    }
    let mut used: HashSet<String> = HashSet::new();
    let mut out = HashMap::new();
    for n in names {
        let lower = n.to_lowercase();
        let mut slug = if counts[&lower] == 1 || *n == lower {
            lower.clone()
        } else {
            format!("{lower}-type")
        };
        let mut k = 2;
        while !used.insert(slug.clone()) {
            slug = format!("{lower}-{k}");
            k += 1;
        }
        let _ = out.insert(n.clone(), slug);
    }
    out
}

/// Remove `functions/*.md` pages that no longer correspond to a built-in, so the
/// directory mirrors the live set exactly.
fn prune(functions: &Path, slugs: &HashMap<String, String>) -> std::io::Result<()> {
    let keep: HashSet<&str> = slugs.values().map(String::as_str).collect();
    for entry in fs::read_dir(functions)? {
        let path = entry?.path();
        let stem = path
            .file_stem()
            .and_then(|s| s.to_str())
            .unwrap_or_default();
        let is_md = path.extension().is_some_and(|e| e == "md");
        if is_md && stem != "index" && !keep.contains(stem) {
            println!("pruned stale doc {}", path.display());
            fs::remove_file(&path)?;
        }
    }
    Ok(())
}

/// Render one function reference page in the website's front-matter format.
fn page(v: &BuiltinDocView) -> String {
    let mut out = format!(
        "---\nlayout: page\ntitle: \"{} (Function)\"\ndescription: \"{}\"\n---\n\n",
        v.name,
        yaml(&v.summary)
    );
    out.push_str(&format!(
        "**Signature:** `{}`\n\n**Description:** {}\n\n",
        v.signature, v.summary
    ));
    push_parameters(&mut out, v);
    out.push_str(&format!("**Returns:** {}\n", v.return_type));
    push_example(&mut out, v);
    out
}

fn push_parameters(out: &mut String, v: &BuiltinDocView) {
    if v.params.is_empty() {
        return;
    }
    out.push_str("## Parameters\n\n");
    for p in &v.params {
        out.push_str(&format!("- **{}** ({}): {}\n", p.name, p.ty, p.description));
    }
    out.push('\n');
}

fn push_example(out: &mut String, v: &BuiltinDocView) {
    if !v.example.is_empty() {
        out.push_str(&format!("\n## Example\n\n```osprey\n{}\n```\n", v.example));
    }
}

/// Render the `functions/index.md` listing, one entry per built-in.
fn index(views: &[BuiltinDocView], slugs: &HashMap<String, String>) -> String {
    let mut out = String::from(
        "---\nlayout: page\ntitle: \"Built-in Functions\"\n\
         description: \"Complete reference for all built-in functions in Osprey\"\n---\n\n\
         All built-in functions available in Osprey.\n\n",
    );
    for v in views {
        out.push_str(&format!(
            "## [{}]({}/)\n\n**Signature:** `{}`\n\n{}\n\n",
            v.name, slugs[&v.name], v.signature, v.summary
        ));
    }
    out
}

/// Escape a string for a YAML double-quoted scalar (front-matter `description`).
fn yaml(s: &str) -> String {
    s.replace('\\', "\\\\").replace('"', "\\\"")
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn page_uses_scheme_types_and_omits_empty_sections() {
        let view = builtin_doc_view("sleep").expect("sleep documented");
        let md = page(&view);
        // The signature follows the scheme (`-> Unit`), not the old prose `-> int`.
        assert!(
            md.contains("**Signature:** `sleep(milliseconds: int) -> Unit`"),
            "{md}"
        );
        assert!(
            md.contains("- **milliseconds** (int): Number of milliseconds"),
            "{md}"
        );
        assert!(md.contains("**Returns:** Unit"), "{md}");
        // A zero-arg builtin renders no Parameters section.
        let input = page(&builtin_doc_view("input").expect("input documented"));
        assert!(!input.contains("## Parameters"), "{input}");
    }

    #[test]
    fn docs_dir_reads_the_flag_value() {
        let args = vec!["--docs".into(), "--docs-dir".into(), "out".into()];
        assert_eq!(docs_dir(&args), Some("out"));
        assert_eq!(docs_dir(&["--docs".to_string()]), None);
    }

    #[test]
    fn yaml_escapes_quotes_and_backslashes() {
        assert_eq!(yaml(r#"a "b" \c"#), r#"a \"b\" \\c"#);
    }

    #[test]
    fn slugs_disambiguate_case_collisions_and_stay_unique() {
        // `map` (function) keeps the bare stem; `Map` (constructor) cannot share
        // it on a case-insensitive disk, so it is suffixed.
        let names = vec!["map".to_string(), "Map".to_string(), "sleep".to_string()];
        let slugs = slug_map(&names);
        assert_eq!(slugs["map"], "map");
        assert_eq!(slugs["Map"], "map-type");
        assert_eq!(slugs["sleep"], "sleep");
        let distinct: HashSet<&String> = slugs.values().collect();
        assert_eq!(distinct.len(), slugs.len(), "slugs must be unique");
    }
}
