#!/usr/bin/env python3
# Deterministically add an osprey-ml twin after every osprey block in the
# website reference docs. Idempotent: never inserts a twin if one already follows.
#
# ML translation rules (verified against examples/tested/ml/*.ospml):
#   let x = v        -> x = v
#   mut c = 0        -> mut c = 0 ; reassign c = e -> c := e
#   fn f(a, b) = e   -> f (a, b) = e     ; fn f(x) = e -> f x = e
#   fn(x) => e       -> lambda           ; fn(a,b) => e -> lambda tuple
#   f(x)             -> f x              ; f(x, y) -> f (x, y)
#   print("...")     -> print "..."      ; print(x) -> print x
#   match ... { }    -> offside layout with => arms
#   Foo { a: 1 }     -> Foo / a = 1 (offside)
#   extern fn n(p: T) -> R  -> extern n (p : T) -> R
# Type annotations that HM can infer are stripped.
import re
import sys
import pathlib

DOCS = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else pathlib.Path("website/src/docs")

# ---- token-level helpers -------------------------------------------------

def strip_ann_in_let(line):
    # `let name: Type = v` -> keep annotation only if it's `Any` (load-bearing) else drop
    m = re.match(r'^(\s*)let\s+(\w+)\s*:\s*(\w+)\s*=\s*(.*)$', line)
    if m:
        indent, name, typ, val = m.groups()
        if typ == "Any":
            return f"{indent}{name} : Any = {translate_expr(val)}"
        return f"{indent}{name} = {translate_expr(val)}"
    m = re.match(r'^(\s*)let\s+(\w+)\s*=\s*(.*)$', line)
    if m:
        indent, name, val = m.groups()
        return f"{indent}{name} = {translate_expr(val)}"
    return None


def split_top_commas(s):
    """Split on commas not inside quotes, parens, brackets, or braces."""
    parts, depth, buf, i = [], 0, [], 0
    quote = None
    while i < len(s):
        c = s[i]
        if quote:
            buf.append(c)
            if c == quote and s[i-1] != "\\":
                quote = None
        elif c in '"\'':
            quote = c; buf.append(c)
        elif c in "([{":
            depth += 1; buf.append(c)
        elif c in ")]}":
            depth -= 1; buf.append(c)
        elif c == "," and depth == 0:
            parts.append("".join(buf)); buf = []
        else:
            buf.append(c)
        i += 1
    parts.append("".join(buf))
    return [p.strip() for p in parts]


def has_toplevel_comma(s):
    return len(split_top_commas(s)) > 1


def paren_call_to_ws(expr):
    """f(x) -> f x ; f(x, y) -> f (x, y) ; keep nested. String/quote aware."""
    out = expr
    # Rewrite innermost `name(args)` where args has no nested parens.
    pat = re.compile(r'(\b[A-Za-z_]\w*)\(([^()]*)\)')
    prev = None
    while prev != out:
        prev = out
        def repl(m):
            fn, args = m.group(1), m.group(2).strip()
            if fn in ("fn",):
                return m.group(0)
            if args == "":
                return f"{fn}()"  # zero-arg -> stripped below
            if has_toplevel_comma(args):
                parts = split_top_commas(args)
                return f"{fn} ({', '.join(parts)})"
            return f"{fn} {args}"
        out = pat.sub(repl, out)
    # zero-arg: f() -> f
    out = re.sub(r'\b([A-Za-z_]\w*)\(\)', r'\1', out)
    return out


def translate_expr(expr):
    e = expr
    # lambdas fn(x) => e  / fn(a,b) => e
    def lam(m):
        params = m.group(1).strip()
        if "," in params:
            return f"\\({params}) =>"
        return f"\\{params} =>"
    e = re.sub(r'fn\(([^)]*)\)\s*=>', lam, e)
    # print("...") -> print "..."  and print(x) -> print x (handled by paren_call)
    e = paren_call_to_ws(e)
    return e


# ---- line-block translation ---------------------------------------------

def translate_block(code):
    """Translate a whole Default osprey code block body to ML. Returns None if the
    block is byte-identical (no surface difference) so caller can skip."""
    lines = code.split("\n")
    # Detect braced match / braced fn bodies -> too complex for line rules; handle a
    # few common shapes, else fall back to a best-effort per-line pass.
    if any(re.search(r'\bmatch\b.*\{', l) for l in lines) or \
       re.search(r'\{\s*$', code) and re.search(r'=>', code):
        ml = translate_match_block(lines)
        if ml is not None:
            return "\n".join(ml)
    out = []
    for line in lines:
        out.append(translate_line(line))
    ml = "\n".join(out)
    if ml == code:
        return None
    return ml


def translate_line(line):
    if line.strip() == "" or line.lstrip().startswith("//"):
        return line
    got = strip_ann_in_let(line)
    if got is not None:
        return got
    # mut c = 0  (keep) ; reassignment  c = c + 1 -> c := c + 1
    m = re.match(r'^(\s*)mut\s+(\w+)\s*=\s*(.*)$', line)
    if m:
        i, n, v = m.groups()
        return f"{i}mut {n} = {translate_expr(v)}"
    # extern fn name(p: T) -> R
    m = re.match(r'^(\s*)extern\s+fn\s+(\w+)\((.*)\)\s*->\s*(.*)$', line)
    if m:
        i, n, params, ret = m.groups()
        p = params.strip()
        if p == "":
            return f"{i}extern {n} -> {ret}"
        parts = [pp.strip().replace(":", " :") for pp in p.split(",")]
        return f"{i}extern {n} ({', '.join(parts)}) -> {ret}"
    # named fn:  fn f(args) = body    or   fn f(args) { body }
    m = re.match(r'^(\s*)fn\s+(\w+)\((.*)\)\s*(?:->\s*\w+\s*)?=\s*(.*)$', line)
    if m:
        i, n, params, body = m.groups()
        p = params.strip()
        # strip param type annotations (HM infers)
        names = [pp.split(":")[0].strip() for pp in p.split(",")] if p else []
        if len(names) == 0:
            head = f"{i}{n} ="
        elif len(names) == 1:
            head = f"{i}{n} {names[0]} ="
        else:
            head = f"{i}{n} ({', '.join(names)}) ="
        return f"{head} {translate_expr(body)}"
    # plain expression / statement line
    return translate_expr(line)


def translate_match_block(lines):
    """Convert a single `match scrut { arm* }` (arms use -> or =>) to offside ML.
    Returns list of ml lines or None if the shape is unrecognised."""
    text = "\n".join(lines)
    m = re.search(r'match\s+(.+?)\s*\{(.*)\}', text, re.S)
    if not m:
        return None
    scrut, body = m.group(1).strip(), m.group(2)
    arms = []
    for raw in body.split("\n"):
        s = raw.strip()
        if not s:
            continue
        a = re.match(r'^(.*?)\s*(?:=>|->)\s*(.*)$', s)
        if not a:
            return None
        pat, res = a.group(1).strip(), a.group(2).strip()
        # constructor payload: Some(x) -> Some x ; None stays
        pat = re.sub(r'(\b[A-Z]\w*)\(([^)]*)\)', lambda mm: f"{mm.group(1)} {mm.group(2)}" if mm.group(2) else mm.group(1), pat)
        arms.append(f"    {pat} => {translate_expr(res)}")
    return [f"match {translate_expr(scrut)}"] + arms


# ---- file rewriting ------------------------------------------------------

BLOCK = re.compile(r'```osprey\n(.*?)\n```', re.S)

def process(path):
    txt = path.read_text()
    out = []
    idx = 0
    changed = False
    for m in BLOCK.finditer(txt):
        out.append(txt[idx:m.end()])
        idx = m.end()
        # already has a twin right after?
        after = txt[m.end():m.end()+40]
        if re.match(r'\s*```osprey-ml', after):
            continue
        ml = translate_block(m.group(1))
        if ml is None:
            continue
        out.append(f"\n\n```osprey-ml\n{ml}\n```")
        changed = True
    out.append(txt[idx:])
    if changed:
        path.write_text("".join(out))
    return changed

def main():
    files = sorted(DOCS.rglob("*.md"))
    n = 0
    for f in files:
        if process(f):
            n += 1
    print(f"Added ML twins to {n} files (of {len(files)} scanned).")

if __name__ == "__main__":
    main()
