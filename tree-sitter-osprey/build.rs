fn main() {
    let src_dir = std::path::Path::new("src");
    let mut build = cc::Build::new();
    build.include(src_dir);
    build.flag_if_supported("-Wno-unused-parameter");
    build.flag_if_supported("-Wno-unused-but-set-variable");
    build.file(src_dir.join("parser.c"));
    println!("cargo:rerun-if-changed=src/parser.c");
    build.compile("tree-sitter-osprey");
}
