export const meta = {
  name: 'add-ml-twins-to-docs',
  description: 'Add an Osprey ML (.ospml) twin code block after every Default osprey example in the 133 /docs reference pages',
  phases: [
    { title: 'Transform' },
    { title: 'Verify' },
  ],
}

const DOCS_ROOT = 'website/src/docs'

// Every doc page that currently has exactly one ```osprey Example block.
const FILES = [
  "functions/abs.md","functions/await.md","functions/awaitprocess.md","functions/byteat.md","functions/bytelength.md","functions/channel.md","functions/cleanupprocess.md","functions/codepointat.md","functions/codepointwidth.md","functions/contains.md","functions/deletefile.md","functions/drop.md","functions/endswith.md","functions/fiber_yield.md","functions/fiberdone.md","functions/filter.md","functions/fold.md","functions/foreach.md","functions/foreachlist.md","functions/fromcodepoint.md","functions/httpcloseclient.md","functions/httpcreateclient.md","functions/httpcreateserver.md","functions/httpdelete.md","functions/httpget.md","functions/httpgetresponse.md","functions/httplisten.md","functions/httppost.md","functions/httpput.md","functions/httpresponsebody.md","functions/httpresponsefree.md","functions/httpresponseheader.md","functions/httpresponsestatus.md","functions/httpstopserver.md","functions/indexof.md","functions/input.md","functions/intdiv.md","functions/isempty.md","functions/join.md","functions/jsonfree.md","functions/jsonget.md","functions/jsonlength.md","functions/jsonparse.md","functions/length.md","functions/lines.md","functions/list.md","functions/listappend.md","functions/listconcat.md","functions/listcontains.md","functions/listget.md","functions/listlength.md","functions/listprepend.md","functions/listreverse.md","functions/map-type.md","functions/map.md","functions/mapcontains.md","functions/mapget.md","functions/mapkeys.md","functions/maplength.md","functions/mapmerge.md","functions/mapremove.md","functions/mapset.md","functions/mapvalues.md","functions/not.md","functions/padend.md","functions/padstart.md","functions/parsefloat.md","functions/parseint.md","functions/print.md","functions/random.md","functions/randombelow.md","functions/range.md","functions/readfile.md","functions/recv.md","functions/repeat.md","functions/replace.md","functions/reverse.md","functions/send.md","functions/sleep.md","functions/spawnprocess.md","functions/split.md","functions/startswith.md","functions/substring.md","functions/take.md","functions/termclear.md","functions/termcols.md","functions/termhidecursor.md","functions/termmovecursor.md","functions/termrawmode.md","functions/termreadkey.md","functions/termrows.md","functions/termshowcursor.md","functions/tolowercase.md","functions/tostring.md","functions/touppercase.md","functions/trim.md","functions/trimend.md","functions/trimstart.md","functions/websocketclose.md","functions/websocketconnect.md","functions/websocketcreateserver.md","functions/websocketkeepalive.md","functions/websocketsend.md","functions/websocketserverbroadcast.md","functions/websocketserverlisten.md","functions/words.md","functions/writefile.md","functions/yield.md","keywords/false.md","keywords/fn.md","keywords/import.md","keywords/let.md","keywords/match.md","keywords/true.md","keywords/type.md","operators/divide.md","operators/equal.md","operators/greater-equal.md","operators/greater-than.md","operators/less-equal.md","operators/less-than.md","operators/minus.md","operators/modulo.md","operators/multiply.md","operators/not-equal.md","operators/pipe-operator.md","operators/plus.md","types/any.md","types/bool.md","types/httpresponse.md","types/int.md","types/processhandle.md","types/string.md"
]

const ML_RULES = `
You are translating an Osprey **Default** (\`.osp\`) code example into its **ML** (\`.ospml\`) twin.

## Osprey ML flavor — GROUND TRUTH (verified against byte-tested examples/tested/ml/*.ospml)

Translation rules (Default → ML):
- **Bindings**: \`let x = v\`  →  \`x = v\`   (drop the \`let\` keyword entirely).
- **Mutable**: \`mut c = 0\` stays \`mut c = 0\`; reassignment \`c = c + 1\`  →  \`c := c + 1\` (walrus).
- **Named function, multi-arg**: \`fn f(a, b) = e\`  →  \`f (a, b) = e\`  (uncurried tuple form).
- **Named function, single arg**: \`fn f(x) = e\`  →  \`f x = e\`.
- **Multi-line function body**: put the body on indented lines under the header, e.g.
    \`\`\`
    classify n =
        match n
            0 => "zero"
            _ => "many"
    \`\`\`
- **Lambdas**: \`fn(x) => e\`  or \`fn(x) { e }\`  →  \`\\x => e\`. Multi-arg lambda \`fn(a, b) => e\` → \`\\(a, b) => e\`.
- **Application**: \`f(x)\`  →  \`f x\`. \`f(x, y)\`  →  \`f (x, y)\`. Nested \`f(g(x))\`  →  \`f (g x)\`.
- **Match** (braces → offside layout): scrutinee on the \`match\` line, each \`pattern => body\` arm indented under it. NO braces, NO \`->\` (ML uses \`=>\`). Constructor patterns bind payload by juxtaposition: \`Success value => ...\`, \`Error e => ...\`.
- **print**: \`print("...")\`  →  \`print "..."\`. \`print(x)\`  →  \`print x\`. Prefer interpolation: \`print "got \${value}"\`.
- **Pipe**: \`a |> f |> g\` is identical in both flavors.
- **Type annotations**: put on a SEPARATE line ABOVE the definition: \`f : int -> int\` then \`f x = ...\`. **Only include an annotation when it is genuinely load-bearing** (a returned curried closure's types, or an empty/ambiguous literal with no inferable context, e.g. \`value : Any = 42\` where \`Any\` is intentional). **Hindley-Milner: NEVER write a type the compiler can infer.** Do NOT annotate parameters or return types that are inferable.
- **Records / construction**: \`Foo { a: 1, b: 2 }\`  →  offside layout:
    \`\`\`
    Foo
        a = 1
        b = 2
    \`\`\`
- **extern**: \`extern fn name(p: T) -> R\`  →  \`extern name (p : T) -> R\`; a zero-arg extern drops parens: \`extern osprey_ffi_null -> Ptr\`.
- **if**: \`if (c) { a } else { b }\` → \`if c then a else b\`. But idiomatic ML usually prefers \`match\` on a boolean. Keep it simple and faithful to what the Default block shows.

CRITICAL:
- Keep comments (\`// ...\`) as-is; they are identical in both flavors.
- Preserve the exact inline output comments (e.g. \`// Prints: 2, 4, 6\`).
- Do NOT change the meaning, the values, or the identifiers.
- The ML twin must be VALID ML — no braces around blocks, no \`fn\` keyword, no \`let\` keyword, no \`->\` in match arms.
- Do NOT add redundant type annotations. Strip any the Default block had that ML can infer.
`

const OUT_SCHEMA = {
  type: 'object',
  additionalProperties: false,
  properties: {
    file: { type: 'string' },
    ml: { type: 'string', description: 'The ML twin code (contents ONLY, no ```osprey-ml fences)' },
    notes: { type: 'string' },
  },
  required: ['file', 'ml'],
}

const VERDICT_SCHEMA = {
  type: 'object',
  additionalProperties: false,
  properties: {
    file: { type: 'string' },
    valid: { type: 'boolean' },
    problems: { type: 'array', items: { type: 'string' } },
    correctedMl: { type: 'string', description: 'If invalid, the corrected ML twin contents (no fences). Empty string if valid.' },
  },
  required: ['file', 'valid'],
}

phase('Transform')

const results = await pipeline(
  FILES,
  // Stage 1: read the Default example, produce the ML twin, WRITE it into the file.
  async (file) => {
    const path = `${DOCS_ROOT}/${file}`
    const out = await agent(
      `${ML_RULES}\n\n` +
      `Read the file \`${path}\`. It contains exactly one Osprey Default example fenced as \`\`\`osprey ... \`\`\` (usually under a \`## Example\` heading).\n\n` +
      `1. Extract the Default example's code.\n` +
      `2. Translate it to a correct, idiomatic ML (.ospml) twin using the rules above.\n` +
      `3. Edit the file: immediately AFTER the closing \`\`\`\` of the existing \`\`\`osprey block, insert a blank line and then a new fenced block:\n` +
      "```osprey-ml\n<your ML twin here>\n```\n" +
      `Use the Edit tool. Do NOT touch the existing osprey block, the front matter, or any prose. Do NOT add any heading between the two blocks.\n\n` +
      `Return the ML twin contents (without fences) in the \`ml\` field.`,
      { label: `twin:${file}`, phase: 'Transform', schema: OUT_SCHEMA, agentType: 'general-purpose' }
    )
    return { file, out }
  },
  // Stage 2: adversarially verify the ML twin is valid ML; auto-fix if not.
  async (prev) => {
    if (!prev || !prev.out) return prev
    const { file } = prev
    const path = `${DOCS_ROOT}/${file}`
    const v = await agent(
      `${ML_RULES}\n\n` +
      `ADVERSARIAL VERIFY. Read \`${path}\`. It should now contain a \`\`\`osprey block (Default) followed by a \`\`\`osprey-ml block (ML twin).\n\n` +
      `Check the ML twin RUTHLESSLY for these defects (default to valid=false if ANY are present):\n` +
      `- Any \`fn\` keyword (illegal in ML).\n` +
      `- Any \`let\` keyword (illegal in ML).\n` +
      `- Any \`->\` inside a match arm (ML uses \`=>\`).\n` +
      `- Braces \`{ }\` used for a block body or match body (ML uses offside layout).\n` +
      `- A redundant type annotation the compiler could infer (HM violation). Parameter/return annotations that mirror the Default block but are inferable MUST be stripped.\n` +
      `- \`print(...)\` with parens where \`print ...\` is idiomatic, or \`f(x)\` application where \`f x\` is meant.\n` +
      `- Semantic drift from the Default block (different values, identifiers, or output).\n` +
      `- A record/extern still in Default brace/paren spelling.\n\n` +
      `If the twin is valid ML, return valid=true. If NOT, set valid=false, list the problems, put the CORRECTED ML twin contents (no fences) in correctedMl, AND apply the correction with the Edit tool (replace the bad osprey-ml block body with the corrected one).`,
      { label: `verify:${file}`, phase: 'Verify', schema: VERDICT_SCHEMA, agentType: 'general-purpose' }
    )
    return { file, out: prev.out, verdict: v }
  }
)

const done = results.filter(Boolean)
const bad = done.filter(r => r.verdict && r.verdict.valid === false)
log(`Transformed ${done.length}/${FILES.length} files. ${bad.length} needed correction (auto-fixed).`)

return {
  transformed: done.length,
  total: FILES.length,
  corrected: bad.map(r => ({ file: r.file, problems: r.verdict.problems })),
}
