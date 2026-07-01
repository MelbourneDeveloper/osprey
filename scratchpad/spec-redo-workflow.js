export const meta = {
  name: 'spec-dual-flavor-redo',
  description: 'In docs/specs (TRUE SOURCE): rewrite illegal single-flavor callouts to dual-flavor, and add ML twins after every Default example that has a real syntactic difference',
  phases: [
    { title: 'Callouts' },
    { title: 'Twins' },
    { title: 'Verify' },
  ],
}

const SPECS = 'docs/specs'

// docs/specs conventions the agents MUST follow.
const CONVENTIONS = `
## docs/specs conventions (STRICT — this is the source of truth, copy-spec regenerates website/src/spec from it)
- Files are PascalCase (e.g. 0024-MLFlavorSyntax.md). NO YAML front matter — do not add any.
- Cross-references are REPO-RELATIVE markdown links: [ML Flavor Syntax](0024-MLFlavorSyntax.md), [Language Flavors](0023-LanguageFlavors.md). NEVER write /spec/... or https://github.com/... links.
- Code fences: Default examples use \`\`\`osprey. ML twins use \`\`\`osprey-ml. EBNF uses \`\`\`ebnf. Leave ebnf/text fences alone.
- Keep files under 500 lines. Do NOT duplicate prose.
`

const ML_RULES = `
## Osprey ML flavor (.ospml) — GROUND TRUTH (verified against byte-tested examples/tested/ml/*.ospml)
Translate a Default (.osp) example to its ML (.ospml) twin with these rules:
- Bindings: \`let x = v\` -> \`x = v\` (drop \`let\`). \`mut c = 0\` stays; reassignment \`c = e\` -> \`c := e\`.
- Named fn multi-arg: \`fn f(a, b) = e\` -> \`f (a, b) = e\` (uncurried tuple form). Single arg: \`fn f(x) = e\` -> \`f x = e\`.
- Multi-line body: header line then indented body (offside layout), e.g.  \`classify n =\` then indented \`match n\` then indented arms.
- Lambdas: \`fn(x) => e\` -> \`\\x => e\`. \`fn(x) { e }\` -> \`\\x => e\`. Multi-arg: \`fn(a,b) => e\` -> \`\\(a, b) => e\`.
- Application: \`f(x)\` -> \`f x\`; \`f(x, y)\` -> \`f (x, y)\`; nested \`f(g(x))\` -> \`f (g x)\`.
- match: braces -> offside layout. Scrutinee on the \`match\` line; each \`pattern => body\` arm indented. NO braces, use \`=>\` (never \`->\`). Payload binds by juxtaposition: \`Success value => ...\`, \`Error e => ...\`.
- print: \`print("...")\` -> \`print "..."\`; \`print(x)\` -> \`print x\`. Prefer interpolation \`print "v=\${x}"\`.
- Pipe \`a |> f |> g\` is identical in both flavors.
- Records: \`Foo { a: 1, b: 2 }\` -> offside:  \`Foo\` then indented \`a = 1\` / \`b = 2\`.
- extern: \`extern fn name(p: T) -> R\` -> \`extern name (p : T) -> R\`; zero-arg drops parens: \`extern osprey_ffi_null -> Ptr\`.
- if: \`if (c) { a } else { b }\` -> \`if c then a else b\` (but idiomatic ML often prefers \`match\` on the boolean).
- Type annotations: on a SEPARATE line ABOVE the def (\`f : int -> int\` then \`f x = ...\`), and ONLY when load-bearing (returned curried closure, or empty/ambiguous literal like \`value : Any = 42\`). HINDLEY-MILNER: NEVER write a type the compiler can infer. Strip inferable parameter/return annotations.
- Keep comments (// ...) and inline output comments (// Prints: ...) identical. Do not change values or identifiers.
`

const CALLOUT_TASK = `
${CONVENTIONS}
You are fixing an ILLEGAL single-flavor callout. The current callout frames the ENTIRE chapter as belonging to ONE flavor (e.g. "This chapter IS the Default flavor" or "Every spelling in this chapter is the Default flavor"). This framing is WRONG and must be removed.

Rewrite the \`> **Flavor layer ...**\` blockquote so it:
1. States the SEMANTICS/lowering are shared-core and flavor-blind (one canonical osprey_ast::Program, [FLAVOR-BOUNDARY]).
2. Says the chapter shows BOTH flavors: the Default (.osp) spelling AND, where the surface differs, the ML (.ospml) twin shown inline alongside it (\`\`\`osprey-ml blocks). Do NOT say the chapter "is" one flavor.
3. Keeps the accurate AST-node mapping details already present (they are correct — preserve them).
4. Keeps repo-relative links to [ML Flavor Syntax](0024-MLFlavorSyntax.md) and [Language Flavors](0023-LanguageFlavors.md).
Edit ONLY the blockquote. Do not touch code blocks or other prose in this pass.
`

phase('Callouts')

// The three illegal single-flavor "surface (CST)" callouts that frame the whole chapter as Default.
const ILLEGAL_CALLOUTS = [
  '0002-LexicalStructure.md',
  '0003-Syntax.md',
  '0005-FunctionCalls.md',
]

await parallel(ILLEGAL_CALLOUTS.map(f => () =>
  agent(`${CALLOUT_TASK}\nFile: \`${SPECS}/${f}\`. Read it, rewrite the Flavor-layer blockquote to dual-flavor framing per the rules, and Edit the file.`,
    { label: `callout:${f}`, phase: 'Callouts', agentType: 'general-purpose' })
))

phase('Twins')

// Chapters whose Default examples have a REAL syntactic ML difference worth showing inline.
// (Pure shared-core chapters with no syntax difference, or with 0 osp fences, are excluded.)
const TWIN_CHAPTERS = [
  '0002-LexicalStructure.md',
  '0003-Syntax.md',
  '0004-TypeSystem.md',
  '0005-FunctionCalls.md',
  '0006-StringInterpolation.md',
  '0007-PatternMatching.md',
  '0008-BlockExpressions.md',
  '0009-BooleanOperations.md',
  '0010-LoopConstructsAndFunctionalIterators.md',
  '0011-LightweightFibersAndConcurrency.md',
  '0012-Built-InFunctions.md',
  '0013-ErrorHandling.md',
  '0017-AlgebraicEffects.md',
  '0019-ForeignFunctionInterface.md',
]

const TWIN_SCHEMA = {
  type: 'object',
  additionalProperties: false,
  properties: {
    file: { type: 'string' },
    twinsAdded: { type: 'integer' },
    skipped: { type: 'integer', description: 'Default blocks with NO surface difference, left without a twin' },
    notes: { type: 'string' },
  },
  required: ['file', 'twinsAdded'],
}

const twinResults = await parallel(TWIN_CHAPTERS.map(f => () =>
  agent(
    `${CONVENTIONS}\n${ML_RULES}\n\n` +
    `File: \`${SPECS}/${f}\`. Read it. For EACH \`\`\`osprey (Default) code block that demonstrates SURFACE SYNTAX that differs in ML ` +
    `(bindings, fn defs, lambdas, calls/application, match, records, extern, interpolation with calls inside), ` +
    `insert an ML twin immediately AFTER the closing \`\`\`\` of that osprey block: a blank line, then a \`\`\`osprey-ml block ` +
    `containing the correct ML translation.\n\n` +
    `SKIP (do NOT add a twin to) blocks that are byte-identical in both flavors (e.g. a bare interpolated string \`"Hello \${name}"\`, ` +
    `a pipe-only chain \`a |> f |> g\`, a pure type/EBNF illustration, or an operator table). Adding an identical twin is noise — skip it.\n\n` +
    `Do NOT modify the existing osprey blocks, the EBNF, the callout, headings, or prose. Only INSERT osprey-ml blocks. Use the Edit tool per insertion.\n\n` +
    `Return how many twins you added and how many blocks you skipped.`,
    { label: `twins:${f}`, phase: 'Twins', schema: TWIN_SCHEMA, agentType: 'general-purpose' }
  )
))

phase('Verify')

// Adversarially verify EVERY osprey-ml block that now exists in the twin chapters is valid ML; auto-fix.
const VERIFY_SCHEMA = {
  type: 'object',
  additionalProperties: false,
  properties: {
    file: { type: 'string' },
    checked: { type: 'integer' },
    fixed: { type: 'integer' },
    problems: { type: 'array', items: { type: 'string' } },
  },
  required: ['file', 'checked', 'fixed'],
}

const verifyResults = await parallel(TWIN_CHAPTERS.map(f => () =>
  agent(
    `${CONVENTIONS}\n${ML_RULES}\n\n` +
    `ADVERSARIAL VERIFY. Read \`${SPECS}/${f}\`. Inspect EVERY \`\`\`osprey-ml block. Flag and FIX (with Edit) any block that has:\n` +
    `- a \`fn\` keyword (illegal in ML)\n- a \`let\` keyword (illegal in ML)\n- \`->\` inside a match arm (must be \`=>\`)\n` +
    `- braces { } for a block/match body (must be offside layout)\n- a redundant type annotation the compiler can infer (HM violation)\n` +
    `- \`print(...)\`/\`f(x)\` paren-application where \`print ...\`/\`f x\` is idiomatic\n- semantic drift from the paired Default block\n` +
    `- a record/extern still in Default brace/paren spelling\n` +
    `- an ML twin that is byte-identical to its Default block (should have been skipped — DELETE such a twin block).\n\n` +
    `Return how many osprey-ml blocks you checked and how many you fixed/deleted, with a short problem list.`,
    { label: `verify:${f}`, phase: 'Verify', schema: VERIFY_SCHEMA, agentType: 'general-purpose' }
  )
))

const twins = twinResults.filter(Boolean)
const verifies = verifyResults.filter(Boolean)
const totalTwins = twins.reduce((s, r) => s + (r.twinsAdded || 0), 0)
const totalFixed = verifies.reduce((s, r) => s + (r.fixed || 0), 0)
log(`Rewrote ${ILLEGAL_CALLOUTS.length} illegal callouts. Added ${totalTwins} ML twins across ${twins.length} chapters. Verify fixed ${totalFixed}.`)

return {
  calloutsRewritten: ILLEGAL_CALLOUTS,
  totalTwins,
  totalFixed,
  perChapter: twins.map(r => ({ file: r.file, added: r.twinsAdded, skipped: r.skipped })),
  problems: verifies.flatMap(r => (r.problems || []).map(p => `${r.file}: ${p}`)),
}
