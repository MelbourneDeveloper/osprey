//! Built-in documentation data (language & collections). Generated companion to
//! `builtins.rs`: every entry's prose pairs with the type scheme of the
//! same name. Edit prose here; edit types in `builtins.rs`. The parity
//! test in `builtin_docs.rs` guarantees the two stay in lockstep.
//!
//! Param order and count MUST match the builtin's real arity.

use crate::builtin_docs::{BuiltinDoc, ParamDoc};

/// `core` built-in documentation. Prose only — types come from the
/// authoritative scheme in `builtins.rs`, joined by name.
pub(crate) static CORE: &[BuiltinDoc] = &[
    BuiltinDoc {
        name: "print",
        summary: "Prints a value to the console. Automatically converts the value to a string representation.",
        params: &[ParamDoc { name: "value", description: "The value to print" }],
        example: "print(\"Hello, World!\")  // Prints: Hello, World!\nprint(42)             // Prints: 42\nprint(true)           // Prints: true",
    },
    BuiltinDoc {
        name: "input",
        summary: "Reads a string from the user's input.",
        params: &[],
        example: "let userInput = input()\nprint(userInput)",
    },
    BuiltinDoc {
        name: "toString",
        summary: "Converts a value to its string representation.",
        params: &[ParamDoc { name: "value", description: "The value to convert to string" }],
        example: "let str = toString(42)\nprint(str)  // Prints: 42",
    },
    BuiltinDoc {
        name: "length",
        summary: "Returns the byte length of a string. Total — never fails.",
        params: &[ParamDoc { name: "s", description: "The string to measure" }],
        example: "let len = length(\"hello\")  // 5",
    },
    BuiltinDoc {
        name: "sleep",
        summary: "Pauses execution for the specified number of milliseconds.",
        params: &[ParamDoc { name: "milliseconds", description: "Number of milliseconds to sleep" }],
        example: "sleep(1000)  // Sleep for 1 second\nprint(\"Awake!\")",
    },
    BuiltinDoc {
        name: "range",
        summary: "Creates an iterator that generates numbers from start to end (exclusive).",
        params: &[ParamDoc { name: "start", description: "The starting number (inclusive)" }, ParamDoc { name: "end", description: "The ending number (exclusive)" }],
        example: "forEach(range(0, 5), fn(x) { print(x) })  // Prints: 0, 1, 2, 3, 4",
    },
    BuiltinDoc {
        name: "abs",
        summary: "Returns the absolute value of an integer.",
        params: &[ParamDoc { name: "value", description: "The integer whose magnitude to take" }],
        example: "let d = abs(0 - 5)  // 5",
    },
    BuiltinDoc {
        name: "not",
        summary: "Returns the logical negation of a boolean.",
        params: &[ParamDoc { name: "value", description: "The boolean to negate" }],
        example: "let off = not(true)  // false",
    },
];

/// `strings` built-in documentation. Prose only — types come from the
/// authoritative scheme in `builtins.rs`, joined by name.
pub(crate) static STRINGS: &[BuiltinDoc] = &[
    BuiltinDoc {
        name: "contains",
        summary: "True if needle appears anywhere in s. Empty needle returns true.",
        params: &[ParamDoc { name: "s", description: "The string to search in" }, ParamDoc { name: "needle", description: "The substring to search for" }],
        example: "let found = contains(\"hello world\", \"world\")  // true",
    },
    BuiltinDoc {
        name: "startsWith",
        summary: "True if s begins with prefix.",
        params: &[ParamDoc { name: "s", description: "The string to test" }, ParamDoc { name: "prefix", description: "The prefix to look for" }],
        example: "startsWith(\"GET /api\", \"GET \")  // true",
    },
    BuiltinDoc {
        name: "endsWith",
        summary: "True if s ends with suffix.",
        params: &[ParamDoc { name: "s", description: "The string to test" }, ParamDoc { name: "suffix", description: "The suffix to look for" }],
        example: "endsWith(\"image.png\", \".png\")  // true",
    },
    BuiltinDoc {
        name: "indexOf",
        summary: "Returns byte-index of first occurrence of needle, or Error(NotFound).",
        params: &[ParamDoc { name: "s", description: "The string to search in" }, ParamDoc { name: "needle", description: "The substring to locate" }],
        example: "match indexOf(\"foo=bar\", \"=\") { Success { value } => print(value) ... }",
    },
    BuiltinDoc {
        name: "split",
        summary: "Splits s on separator. Error(InvalidArgument) on empty separator.",
        params: &[ParamDoc { name: "s", description: "The string to split" }, ParamDoc { name: "separator", description: "Non-empty separator" }],
        example: "split(\"a,b,c\", \",\")  // Success { value: [\"a\",\"b\",\"c\"] }",
    },
    BuiltinDoc {
        name: "join",
        summary: "Concatenates parts with separator between each pair.",
        params: &[ParamDoc { name: "parts", description: "Strings to join" }, ParamDoc { name: "separator", description: "Separator string" }],
        example: "join([\"a\",\"b\",\"c\"], \"-\")  // \"a-b-c\"",
    },
    BuiltinDoc {
        name: "parseInt",
        summary: "Strict base-10 signed-int parser. No whitespace tolerance.",
        params: &[ParamDoc { name: "s", description: "The string to parse" }],
        example: "parseInt(\"42\")  // Success { value: 42 }",
    },
    BuiltinDoc {
        name: "lines",
        summary: "Splits on '\\n'. A trailing newline does not produce an empty entry.",
        params: &[ParamDoc { name: "s", description: "The string to split" }],
        example: "lines(\"a\\\nb\\\nc\")  // [\"a\",\"b\",\"c\"]",
    },
    BuiltinDoc {
        name: "words",
        summary: "Splits on runs of whitespace; empty results dropped.",
        params: &[ParamDoc { name: "s", description: "The string to split" }],
        example: "words(\"a  b\\\\tc\")  // [\"a\",\"b\",\"c\"]",
    },
    BuiltinDoc {
        name: "replace",
        summary: "Replaces every occurrence of needle. Error(InvalidArgument) on empty needle.",
        params: &[ParamDoc { name: "s", description: "The source string" }, ParamDoc { name: "needle", description: "The substring to find" }, ParamDoc { name: "replacement", description: "The replacement string" }],
        example: "replace(\"a-b-c\", \"-\", \"_\")  // Success { value: \"a_b_c\" }",
    },
    BuiltinDoc {
        name: "repeat",
        summary: "Concatenates s with itself n times. Error(InvalidArgument) on negative n.",
        params: &[ParamDoc { name: "s", description: "The string to repeat" }, ParamDoc { name: "n", description: "Repeat count, must be >= 0" }],
        example: "repeat(\"ab\", 3)  // Success { value: \"ababab\" }",
    },
    BuiltinDoc {
        name: "substring",
        summary: "Extracts s[start, end). Returns Error(IndexOutOfRange) if start<0, end>len, or start>end.",
        params: &[ParamDoc { name: "s", description: "The source string" }, ParamDoc { name: "start", description: "Starting index (inclusive)" }, ParamDoc { name: "end", description: "Ending index (exclusive)" }],
        example: "substring(\"hello\", 1, 4)  // Success { value: \"ell\" }",
    },
    BuiltinDoc {
        name: "take",
        summary: "Returns at most the first n bytes of s. Clamps; never fails.",
        params: &[ParamDoc { name: "s", description: "The source string" }, ParamDoc { name: "n", description: "How many bytes to take" }],
        example: "take(\"hello\", 3)  // \"hel\"",
    },
    BuiltinDoc {
        name: "drop",
        summary: "Returns s without its first n bytes. Clamps; never fails.",
        params: &[ParamDoc { name: "s", description: "The source string" }, ParamDoc { name: "n", description: "How many bytes to drop" }],
        example: "drop(\"hello\", 3)  // \"lo\"",
    },
    BuiltinDoc {
        name: "isEmpty",
        summary: "True if string has zero length.",
        params: &[ParamDoc { name: "s", description: "The string to test" }],
        example: "let blank = isEmpty(\"\")  // true",
    },
    BuiltinDoc {
        name: "parseFloat",
        summary: "Strict base-10 floating-point parser. No whitespace tolerance.",
        params: &[ParamDoc { name: "s", description: "The string to parse" }],
        example: "parseFloat(\"3.14\")  // Success { value: 3.14 }",
    },
    BuiltinDoc {
        name: "padStart",
        summary: "Pads s on the left with copies of fill to reach targetLength bytes.",
        params: &[ParamDoc { name: "s", description: "The string to pad" }, ParamDoc { name: "targetLength", description: "Desired total length" }, ParamDoc { name: "fill", description: "Padding string (non-empty)" }],
        example: "padStart(\"7\", 3, \"0\")  // Success { value: \"007\" }",
    },
    BuiltinDoc {
        name: "padEnd",
        summary: "Pads s on the right with copies of fill to reach targetLength bytes.",
        params: &[ParamDoc { name: "s", description: "The string to pad" }, ParamDoc { name: "targetLength", description: "Desired total length" }, ParamDoc { name: "fill", description: "Padding string (non-empty)" }],
        example: "padEnd(\"7\", 3, \".\")  // Success { value: \"7..\" }",
    },
    BuiltinDoc {
        name: "byteLength",
        summary: "Returns the number of bytes in the string's UTF-8 encoding.",
        params: &[ParamDoc { name: "text", description: "The string to measure" }],
        example: "let n = byteLength(\"héllo\")  // 6",
    },
    BuiltinDoc {
        name: "byteAt",
        summary: "Returns the byte at the given index (0-255), or an error if the index is out of range.",
        params: &[ParamDoc { name: "text", description: "The string to read from" }, ParamDoc { name: "index", description: "Zero-based byte offset" }],
        example: "match byteAt(\"hi\", 0) {\n  Success { value } => print(\"byte: ${value}\")\n  Error { message } => print(message)\n}",
    },
    BuiltinDoc {
        name: "codePointAt",
        summary: "Returns the Unicode code point that begins at the given byte index. Fails on an invalid index or malformed UTF-8.",
        params: &[ParamDoc { name: "text", description: "The string to read from" }, ParamDoc { name: "index", description: "Byte offset where the code point starts" }],
        example: "match codePointAt(\"héllo\", 1) {\n  Success { value } => print(\"U+${value}\")\n  Error { message } => print(message)\n}",
    },
    BuiltinDoc {
        name: "codePointWidth",
        summary: "Returns how many bytes the given Unicode code point occupies in UTF-8 (1-4).",
        params: &[ParamDoc { name: "codePoint", description: "The Unicode scalar value" }],
        example: "match codePointWidth(233) {\n  Success { value } => print(\"${value} bytes\")\n  Error { message } => print(message)\n}",
    },
    BuiltinDoc {
        name: "fromCodePoint",
        summary: "Returns the single-character string for a Unicode code point, or an error if it is not a valid scalar value.",
        params: &[ParamDoc { name: "codePoint", description: "The Unicode scalar value to encode" }],
        example: "match fromCodePoint(233) {\n  Success { value } => print(value)  // é\n  Error { message } => print(message)\n}",
    },
    BuiltinDoc {
        name: "toUpperCase",
        summary: "ASCII-aware uppercase. Unicode simple case mapping is a future addition.",
        params: &[ParamDoc { name: "s", description: "The string to transform" }],
        example: "toUpperCase(\"hello\")  // \"HELLO\"",
    },
    BuiltinDoc {
        name: "toLowerCase",
        summary: "ASCII-aware lowercase.",
        params: &[ParamDoc { name: "s", description: "The string to transform" }],
        example: "toLowerCase(\"HELLO\")  // \"hello\"",
    },
    BuiltinDoc {
        name: "trim",
        summary: "Removes leading and trailing whitespace.",
        params: &[ParamDoc { name: "s", description: "The string to trim" }],
        example: "trim(\"  hi  \")  // \"hi\"",
    },
    BuiltinDoc {
        name: "trimStart",
        summary: "Removes leading whitespace.",
        params: &[ParamDoc { name: "s", description: "The string to trim" }],
        example: "trimStart(\"  hi  \")  // \"hi  \"",
    },
    BuiltinDoc {
        name: "trimEnd",
        summary: "Removes trailing whitespace.",
        params: &[ParamDoc { name: "s", description: "The string to trim" }],
        example: "trimEnd(\"  hi  \")  // \"  hi\"",
    },
    BuiltinDoc {
        name: "reverse",
        summary: "Reverses byte order. Grapheme-cluster reversal is future work.",
        params: &[ParamDoc { name: "s", description: "The string to reverse" }],
        example: "reverse(\"abc\")  // \"cba\"",
    },
];

/// `functional` built-in documentation. Prose only — types come from the
/// authoritative scheme in `builtins.rs`, joined by name.
pub(crate) static FUNCTIONAL: &[BuiltinDoc] = &[
    BuiltinDoc {
        name: "forEach",
        summary: "Applies a function to each element in an iterator.",
        params: &[ParamDoc { name: "iterator", description: "The iterator to process" }, ParamDoc { name: "function", description: "The function to apply to each element" }],
        example: "forEach(range(1, 4), fn(x) { print(x * 2) })  // Prints: 2, 4, 6",
    },
    BuiltinDoc {
        name: "map",
        summary: "Transforms each element in an iterator using a function, returning a new iterator.",
        params: &[ParamDoc { name: "iterator", description: "The iterator to transform" }, ParamDoc { name: "fn", description: "The transformation function" }],
        example: "let doubled = map(range(1, 4), fn(x) { x * 2 })\nforEach(doubled, print)  // Prints: 2, 4, 6",
    },
    BuiltinDoc {
        name: "filter",
        summary: "Filters elements in an iterator based on a predicate function.",
        params: &[ParamDoc { name: "iterator", description: "The iterator to filter" }, ParamDoc { name: "predicate", description: "The predicate function that returns true for elements to keep" }],
        example: "let evens = filter(range(1, 6), fn(x) { x % 2 == 0 })\nforEach(evens, print)  // Prints: 2, 4",
    },
    BuiltinDoc {
        name: "fold",
        summary: "Reduces an iterator to a single value by repeatedly applying a function.",
        params: &[ParamDoc { name: "iterator", description: "The iterator to reduce" }, ParamDoc { name: "initial", description: "The initial value for the accumulator" }, ParamDoc { name: "fn", description: "The reduction function that takes (accumulator, current) and returns new accumulator" }],
        example: "range(1, 5) |> fold(0, add)  // sum: 0+1+2+3+4 = 10",
    },
];

/// `lists` built-in documentation. Prose only — types come from the
/// authoritative scheme in `builtins.rs`, joined by name.
pub(crate) static LISTS: &[BuiltinDoc] = &[
    BuiltinDoc {
        name: "List",
        summary: "Creates a new empty list.",
        params: &[],
        example: "let myList = List()\nprint(\"Created empty list\")",
    },
    BuiltinDoc {
        name: "listAppend",
        summary: "Returns a new list with value at the end. O(log32 n) amortised.",
        params: &[ParamDoc { name: "list", description: "The list" }, ParamDoc { name: "value", description: "Value to append" }],
        example: "listAppend([1, 2], 3)  // [1, 2, 3]",
    },
    BuiltinDoc {
        name: "listPrepend",
        summary: "Returns a new list with value at the front. O(n).",
        params: &[ParamDoc { name: "list", description: "The list" }, ParamDoc { name: "value", description: "Value to prepend" }],
        example: "listPrepend([2, 3], 1)  // [1, 2, 3]",
    },
    BuiltinDoc {
        name: "listConcat",
        summary: "Returns left ++ right. Same as left + right.",
        params: &[ParamDoc { name: "left", description: "Left operand" }, ParamDoc { name: "right", description: "Right operand" }],
        example: "listConcat([1, 2], [3, 4])  // [1, 2, 3, 4]",
    },
    BuiltinDoc {
        name: "listReverse",
        summary: "Returns a new list in reverse order.",
        params: &[ParamDoc { name: "list", description: "The list" }],
        example: "listReverse([1, 2, 3])  // [3, 2, 1]",
    },
    BuiltinDoc {
        name: "listLength",
        summary: "Returns the number of elements in a list. O(1).",
        params: &[ParamDoc { name: "list", description: "The list" }],
        example: "listLength([1, 2, 3])  // 3",
    },
    BuiltinDoc {
        name: "listGet",
        summary: "Returns the element at the given index, or an error if the index is out of range.",
        params: &[ParamDoc { name: "list", description: "The list to read from" }, ParamDoc { name: "index", description: "Zero-based element index" }],
        example: "match listGet(myList, 0) {\n  Success { value } => print(value)\n  Error { message } => print(message)\n}",
    },
    BuiltinDoc {
        name: "listContains",
        summary: "True iff some element equals value. O(n).",
        params: &[ParamDoc { name: "list", description: "The list" }, ParamDoc { name: "value", description: "Value to find" }],
        example: "listContains([1, 2, 3], 2)  // true",
    },
    BuiltinDoc {
        name: "forEachList",
        summary: "Apply function to every element of list. Phase 7 of collections plan.",
        params: &[ParamDoc { name: "list", description: "The list" }, ParamDoc { name: "function", description: "Function applied per element" }],
        example: "forEachList(xs, print)",
    },
];

/// `maps` built-in documentation. Prose only — types come from the
/// authoritative scheme in `builtins.rs`, joined by name.
pub(crate) static MAPS: &[BuiltinDoc] = &[
    BuiltinDoc {
        name: "Map",
        summary: "Creates a new, empty persistent map.",
        params: &[],
        example: "let m = Map()",
    },
    BuiltinDoc {
        name: "mapSet",
        summary: "Returns a new map with key bound to value (replaces prior binding).",
        params: &[ParamDoc { name: "map", description: "The map" }, ParamDoc { name: "key", description: "Key" }, ParamDoc { name: "value", description: "Value" }],
        example: "mapSet({\"a\": 1}, \"b\", 2)  // {\"a\": 1, \"b\": 2}",
    },
    BuiltinDoc {
        name: "mapGet",
        summary: "Returns the value associated with the key, or an error if the key is absent.",
        params: &[ParamDoc { name: "map", description: "The map to look up in" }, ParamDoc { name: "key", description: "The key to find" }],
        example: "match mapGet(scores, \"alice\") {\n  Success { value } => print(value)\n  Error { message } => print(message)\n}",
    },
    BuiltinDoc {
        name: "mapRemove",
        summary: "Returns a new map without key. No-op if key is absent.",
        params: &[ParamDoc { name: "map", description: "The map" }, ParamDoc { name: "key", description: "Key" }],
        example: "mapRemove({\"a\": 1, \"b\": 2}, \"a\")  // {\"b\": 2}",
    },
    BuiltinDoc {
        name: "mapMerge",
        summary: "Right-biased union. Same as left + right.",
        params: &[ParamDoc { name: "left", description: "Left" }, ParamDoc { name: "right", description: "Right" }],
        example: "mapMerge({\"a\": 1}, {\"b\": 2})  // {\"a\": 1, \"b\": 2}",
    },
    BuiltinDoc {
        name: "mapContains",
        summary: "True iff key is present in map.",
        params: &[ParamDoc { name: "map", description: "The map" }, ParamDoc { name: "key", description: "Key to find" }],
        example: "mapContains({\"a\": 1}, \"a\")  // true",
    },
    BuiltinDoc {
        name: "mapLength",
        summary: "Returns the number of entries in a map. O(1).",
        params: &[ParamDoc { name: "map", description: "The map" }],
        example: "mapLength({\"a\": 1, \"b\": 2})  // 2",
    },
    BuiltinDoc {
        name: "mapKeys",
        summary: "All keys of the map as a list. Order unspecified.",
        params: &[ParamDoc { name: "map", description: "The map" }],
        example: "mapKeys(m)  // List<K>",
    },
    BuiltinDoc {
        name: "mapValues",
        summary: "All values of the map as a list. Order matches mapKeys.",
        params: &[ParamDoc { name: "map", description: "The map" }],
        example: "mapValues(m)  // List<V>",
    },
];
