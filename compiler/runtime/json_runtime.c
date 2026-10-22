// JSON runtime - the strict RFC 8259 recursive-descent parser every JSON
// entry point shares: the handle API (jsonParse) and the typed Json decoder
// both read documents through json_parse_tree, so they accept and reject
// exactly the same texts. String literals decode in json_text.c; handles,
// the decode index and the writer live in json_codec.c. Implements
// [BUILTIN-JSON], [BUILTIN-JSON-LIMITS].

#include "io_error.h"
#include "json_internal.h"

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static JVal *jval_new(JType type) {
  JVal *v = calloc(1, sizeof(JVal));
  if (v) {
    v->type = type;
  }
  return v;
}

void json_tree_free(JVal *v) {
  if (!v) {
    return;
  }
  switch (v->type) {
  case J_STR:
  case J_NUM:
    free(v->str);
    break;
  case J_ARR:
    for (size_t i = 0; i < v->count; i++) {
      json_tree_free(v->items[i]);
    }
    free(v->items);
    break;
  case J_OBJ:
    for (size_t i = 0; i < v->nmemb; i++) {
      free(v->keys[i]);
      json_tree_free(v->vals[i]);
    }
    free(v->keys);
    free(v->vals);
    break;
  default:
    break;
  }
  free(v);
}

// ---- parser ----------------------------------------------------------------

void json_reject(Cursor *c, const char *why) {
  if (c->ok) {
    c->ok = false;
    c->why = why;
    c->at = c->p;
  }
}

static void skip_ws(Cursor *c) {
  while (*c->p == ' ' || *c->p == '\t' || *c->p == '\n' || *c->p == '\r') {
    c->p++;
  }
}

static JVal *parse_value(Cursor *c);

static JVal *parse_string(Cursor *c) {
  char *s = json_parse_string_raw(c);
  if (!s) {
    return NULL;
  }
  JVal *v = jval_new(J_STR);
  if (!v) {
    free(s);
    json_reject(c, "out of memory");
    return NULL;
  }
  v->str = s;
  return v;
}

// One or more digits; false when there were none.
static bool scan_digits(Cursor *c) {
  const char *start = c->p;
  while (*c->p >= '0' && *c->p <= '9') {
    c->p++;
  }
  return c->p != start;
}

// RFC 8259's number grammar exactly: an optional minus, an integer part whose
// leading zero stands alone, an optional fraction of at least one digit, and
// an optional exponent of at least one digit.
//
// The scanner this replaces took any run of [0-9.eE+-], so `-`, `1.`, `1e`,
// `1+2`, `--1` and `1e+` all parsed -- and since a number is STORED as its
// source text, every one of them came straight back out of jsonGet as a
// "number" no other reader would ever produce [BUILTIN-JSON-NUMBER].
static bool scan_number(Cursor *c) {
  if (*c->p == '-') {
    c->p++;
  }
  if (*c->p == '0') {
    c->p++; // a leading zero is the entire integer part
  } else if (!scan_digits(c)) {
    return false;
  }
  if (*c->p == '.') {
    c->p++;
    if (!scan_digits(c)) {
      return false;
    }
  }
  if (*c->p != 'e' && *c->p != 'E') {
    return true;
  }
  c->p++;
  if (*c->p == '+' || *c->p == '-') {
    c->p++;
  }
  return scan_digits(c);
}

static JVal *parse_number(Cursor *c) {
  const char *start = c->p;
  if (!scan_number(c)) {
    c->p = start; // name the number, not the byte the scan stopped on
    json_reject(c, "invalid number");
    return NULL;
  }
  size_t len = (size_t)(c->p - start);
  JVal *v = jval_new(J_NUM);
  if (!v) {
    json_reject(c, "out of memory");
    return NULL;
  }
  v->str = malloc(len + 1);
  if (!v->str) {
    json_tree_free(v);
    json_reject(c, "out of memory");
    return NULL;
  }
  memcpy(v->str, start, len);
  v->str[len] = '\0';
  return v;
}

static JVal *parse_array(Cursor *c) {
  c->p++; // consume '['
  JVal *v = jval_new(J_ARR);
  if (!v) {
    json_reject(c, "out of memory");
    return NULL;
  }
  skip_ws(c);
  if (*c->p == ']') {
    c->p++;
    return v;
  }
  for (;;) {
    JVal *item = parse_value(c);
    if (!c->ok) {
      json_tree_free(v);
      return NULL;
    }
    JVal **ni = realloc(v->items, (v->count + 1) * sizeof(JVal *));
    if (!ni) {
      json_tree_free(item);
      json_tree_free(v);
      json_reject(c, "out of memory");
      return NULL;
    }
    v->items = ni;
    v->items[v->count++] = item;
    skip_ws(c);
    if (*c->p == ',') {
      c->p++;
      skip_ws(c);
      continue;
    }
    if (*c->p == ']') {
      c->p++;
      return v;
    }
    json_tree_free(v);
    json_reject(c, "expected ',' or ']' in array");
    return NULL;
  }
}

static int compare_keys(const void *left, const void *right) {
  return strcmp(*(char *const *)left, *(char *const *)right);
}

// Two members with one name make every lookup a guess about which one the
// author meant -- first-wins here, last-wins in most other readers -- and
// hostile input exploits exactly that disagreement. A document that repeats a
// name is rejected whole. Sorting a copy keeps the check O(n log n), so a huge
// object cannot turn it into a denial of service [BUILTIN-JSON-LIMITS].
static JVal *unique_members(Cursor *c, JVal *v) {
  if (v->nmemb < 2) {
    return v;
  }
  char **sorted = malloc(v->nmemb * sizeof(char *));
  if (!sorted) {
    json_tree_free(v);
    json_reject(c, "out of memory");
    return NULL;
  }
  memcpy(sorted, v->keys, v->nmemb * sizeof(char *));
  qsort(sorted, v->nmemb, sizeof(char *), compare_keys);
  bool repeated = false;
  for (size_t i = 1; i < v->nmemb && !repeated; i++) {
    repeated = strcmp(sorted[i - 1], sorted[i]) == 0;
  }
  free(sorted);
  if (repeated) {
    json_tree_free(v);
    json_reject(c, "duplicate member name");
    return NULL;
  }
  return v;
}

static JVal *parse_object(Cursor *c) {
  c->p++; // consume '{'
  JVal *v = jval_new(J_OBJ);
  if (!v) {
    json_reject(c, "out of memory");
    return NULL;
  }
  skip_ws(c);
  if (*c->p == '}') {
    c->p++;
    return v;
  }
  for (;;) {
    skip_ws(c);
    char *key = json_parse_string_raw(c);
    if (!key) {
      json_tree_free(v);
      json_reject(c, "expected a member name");
      return NULL;
    }
    skip_ws(c);
    if (*c->p != ':') {
      free(key);
      json_tree_free(v);
      json_reject(c, "expected ':' after a member name");
      return NULL;
    }
    c->p++;
    JVal *val = parse_value(c);
    if (!c->ok) {
      free(key);
      json_tree_free(v);
      return NULL;
    }
    // Each block is published into `v` the instant it exists. realloc has
    // already freed the old one, so holding a new block in a local while the
    // NEXT allocation runs leaves `v` pointing at freed memory — and the
    // cleanup below then walks it. That was a use-after-free reached by
    // nothing worse than a second realloc saying no.
    char **nk = realloc(v->keys, (v->nmemb + 1) * sizeof(char *));
    if (nk != NULL) {
      v->keys = nk;
    }
    JVal **nv = realloc(v->vals, (v->nmemb + 1) * sizeof(JVal *));
    if (nv != NULL) {
      v->vals = nv;
    }
    if (nk == NULL || nv == NULL) {
      free(key);
      json_tree_free(val);
      json_tree_free(v);
      json_reject(c, "out of memory");
      return NULL;
    }
    v->keys[v->nmemb] = key;
    v->vals[v->nmemb] = val;
    v->nmemb++;
    skip_ws(c);
    if (*c->p == ',') {
      c->p++;
      continue;
    }
    if (*c->p == '}') {
      c->p++;
      return unique_members(c, v);
    }
    json_tree_free(v);
    json_reject(c, "expected ',' or '}' in object");
    return NULL;
  }
}

#define JSON_STRINGIFY(x) #x
#define JSON_DEPTH_TEXT(x) "nesting deeper than " JSON_STRINGIFY(x) " levels"

// Every container costs one parser frame, and the typed decoder another, so
// depth is bounded before the frame is pushed [BUILTIN-JSON-LIMITS].
static JVal *parse_nested(Cursor *c, JVal *(*container)(Cursor *)) {
  if (c->depth >= JSON_MAX_DEPTH) {
    json_reject(c, JSON_DEPTH_TEXT(JSON_MAX_DEPTH));
    return NULL;
  }
  c->depth++;
  JVal *v = container(c);
  c->depth--;
  return v;
}

static JVal *parse_value(Cursor *c) {
  skip_ws(c);
  switch (*c->p) {
  case '"':
    return parse_string(c);
  case '{':
    return parse_nested(c, parse_object);
  case '[':
    return parse_nested(c, parse_array);
  case 't':
    if (strncmp(c->p, "true", 4) == 0) {
      c->p += 4;
      JVal *v = jval_new(J_BOOL);
      if (v) {
        v->bval = true;
      } else {
        json_reject(c, "out of memory");
      }
      return v;
    }
    break;
  case 'f':
    if (strncmp(c->p, "false", 5) == 0) {
      c->p += 5;
      JVal *v = jval_new(J_BOOL);
      if (v) {
        v->bval = false;
      } else {
        json_reject(c, "out of memory");
      }
      return v;
    }
    break;
  case 'n':
    if (strncmp(c->p, "null", 4) == 0) {
      c->p += 4;
      JVal *v = jval_new(J_NULL);
      if (!v) {
        json_reject(c, "out of memory");
      }
      return v;
    }
    break;
  default:
    if (*c->p == '-' || (*c->p >= '0' && *c->p <= '9')) {
      return parse_number(c);
    }
    break;
  }
  json_reject(c, *c->p ? "unexpected character" : "unexpected end of text");
  return NULL;
}

// One complete value and nothing after it but whitespace. The reason names
// the byte where the text stopped being JSON, so `Error { message }` says why
// -- "jsonParse: byte 7: expected ',' or ']' in array" -- instead of the
// placeholder word "Error" [BUILTIN-JSON-STRING].
JVal *json_parse_tree(const char *text, const char *op, int64_t *status) {
  if (!text) {
    *status = JSON_ERR_NULL;
    osp_io_error_text(op, "document", "no text");
    return NULL;
  }
  Cursor c = {.p = text, .base = text, .ok = true};
  JVal *root = parse_value(&c);
  skip_ws(&c);
  bool trailing = c.ok && *c.p != '\0';
  if (trailing) {
    json_reject(&c, "trailing text after the value");
  }
  if (c.ok) {
    return root;
  }
  json_tree_free(root);
  *status = trailing ? JSON_ERR_TRAILING : JSON_ERR_MALFORMED;
  char where[32];
  snprintf(where, sizeof(where), "byte %zu", (size_t)(c.at - c.base));
  osp_io_error_text(op, where, c.why);
  return NULL;
}
