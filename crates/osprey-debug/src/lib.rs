//! Generic debugger support primitives.
//!
//! This crate intentionally avoids Osprey parser, type-checker, codegen, and
//! editor APIs. It holds small debugger concepts that are candidates to move to
//! `lspkit` once the shape proves useful across languages.

use std::path::{Path, PathBuf};

/// Source file identity used by debug-info producers and editor debug launches.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct DebugSource {
    /// Basename of the source file.
    pub filename: String,
    /// Directory containing the source file.
    pub directory: String,
}

impl DebugSource {
    /// Build a source identity from a source path.
    #[must_use]
    pub fn from_path(path: &str) -> Self {
        let path = Path::new(path);
        let filename = path
            .file_name()
            .and_then(|s| s.to_str())
            .unwrap_or("input.osp")
            .to_string();
        let directory = path
            .parent()
            .and_then(|p| p.to_str())
            .unwrap_or(".")
            .to_string();
        DebugSource {
            filename,
            directory,
        }
    }

    /// The full source path represented by this identity.
    #[must_use]
    pub fn path(&self) -> PathBuf {
        Path::new(&self.directory).join(&self.filename)
    }
}

/// Debug-build switches shared by native compiler front-ends.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct DebugBuild {
    /// Whether source-level debug information is requested.
    pub enabled: bool,
}

impl DebugBuild {
    /// A non-debug build.
    pub const OFF: DebugBuild = DebugBuild { enabled: false };

    /// A source-level debug build.
    pub const ON: DebugBuild = DebugBuild { enabled: true };

    /// The optimizer flag to use for this build.
    #[must_use]
    pub fn opt_flag(self, release_default: String, debug_override: Option<String>) -> String {
        if self.enabled {
            return debug_override.unwrap_or_else(|| "-O0".to_string());
        }
        release_default
    }

    /// Extra C/LLVM driver flags for native debug builds.
    #[must_use]
    pub fn native_driver_flags(self) -> Vec<String> {
        if self.enabled {
            vec!["-g".to_string(), "-fno-omit-frame-pointer".to_string()]
        } else {
            Vec::new()
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn debug_source_splits_file_and_directory() {
        let src = DebugSource::from_path("/tmp/example.osp");
        assert_eq!(src.filename, "example.osp");
        assert_eq!(src.directory, "/tmp");
        assert_eq!(src.path(), PathBuf::from("/tmp/example.osp"));
    }

    #[test]
    fn debug_build_selects_flags() {
        assert_eq!(
            DebugBuild::OFF.opt_flag("-O2".to_string(), Some("-O0".to_string())),
            "-O2"
        );
        assert_eq!(DebugBuild::ON.opt_flag("-O2".to_string(), None), "-O0");
        assert_eq!(
            DebugBuild::ON.opt_flag("-O2".to_string(), Some("-Og".to_string())),
            "-Og"
        );
        assert!(DebugBuild::OFF.native_driver_flags().is_empty());
        assert!(DebugBuild::ON
            .native_driver_flags()
            .iter()
            .any(|f| f == "-g"));
    }
}
