// JSON document handles: the jsonParse/jsonGet surface, the node index the
// built-in Json prelude decodes through, and the single-buffer writer it
// encodes through. Every handle is a slot in a lock-guarded table, so a
// stale, forged or double-freed handle is a refusal, never a dereference.
// Implements [BUILTIN-JSON], [BUILTIN-JSON-TYPED].
//
// Path syntax (jsonGet/jsonLength): "a.b[0].c". Keys containing '.' or '['
// are not addressable by path; the typed decoder walks node ids instead.

#ifndef __wasm__
#include <pthread.h>
#endif
#include "io_error.h"
#include "json_internal.h"

#include <errno.h>
#include <math.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

char *osp_float_to_string(double d);

#define MAX_JSON_HANDLES 1024

// One node of a document's breadth-first index. Breadth-first order keeps a
// container's children contiguous, so child `i` of node `n` is `first + i`.
typedef struct {
  const JVal *v;
  size_t first;
  union {
    int64_t i;
    double f;
  } number; // a number's value, read once while indexing
} JNode;

typedef struct {
  JVal *root;
  JNode *nodes;
  size_t count;
} JDoc;

typedef struct {
  JsonBuf buf;
  const char *why;
  bool after_value; // the next value needs a ',' in front of it
} JWriter;

static JDoc *g_json_docs[MAX_JSON_HANDLES];
static JWriter *g_json_writers[MAX_JSON_HANDLES];

// Every fiber shares the handle tables, so each access holds this lock.
// wasm32-wasip1 runs one thread — the fiber runtime is not in its archive —
// and the WASI SDK that CI pins ships no pthread mutex, so a JSON program
// could not link there. The lock compiles away on that target. [WASM-TARGET]
#ifdef __wasm__
static void json_table_lock(void) {}
static void json_table_unlock(void) {}
#else
static pthread_mutex_t g_json_mutex = PTHREAD_MUTEX_INITIALIZER;
static void json_table_lock(void) { pthread_mutex_lock(&g_json_mutex); }
static void json_table_unlock(void) { pthread_mutex_unlock(&g_json_mutex); }
#endif

static bool valid_handle(int64_t handle) {
  return handle >= 1 && handle < MAX_JSON_HANDLES;
}

// Claim the first free slot for `item`; -1 when the table is full. Callers
// hold the lock.
static int64_t claim_slot(void **table, void *item) {
  for (int64_t i = 1; i < MAX_JSON_HANDLES; i++) {
    if (!table[i]) {
      table[i] = item;
      return i;
    }
  }
  return -1;
}

static void doc_free(JDoc *doc) {
  json_tree_free(doc->root);
  free(doc->nodes);
  free(doc);
}

// ---- path navigation -------------------------------------------------------

static const JVal *navigate(const JVal *cur, const char *path) {
  const char *p = path;
  while (cur && *p) {
    if (*p == '.') {
      p++;
      continue;
    }
    if (*p == '[') {
      p++;
      long idx = 0;
      bool any = false;
      while (*p >= '0' && *p <= '9') {
        idx = idx * 10 + (*p - '0');
        p++;
        any = true;
      }
      if (!any || *p != ']') {
        return NULL;
      }
      p++;
      if (cur->type != J_ARR || idx < 0 || (size_t)idx >= cur->count) {
        return NULL;
      }
      cur = cur->items[idx];
    } else {
      const char *start = p;
      while (*p && *p != '.' && *p != '[') {
        p++;
      }
      size_t klen = (size_t)(p - start);
      if (cur->type != J_OBJ) {
        return NULL;
      }
      const JVal *next = NULL;
      for (size_t i = 0; i < cur->nmemb; i++) {
        if (strlen(cur->keys[i]) == klen &&
            strncmp(cur->keys[i], start, klen) == 0) {
          next = cur->vals[i];
          break;
        }
      }
      if (!next) {
        return NULL;
      }
      cur = next;
    }
  }
  return cur;
}


static const JVal *lookup(int64_t handle, const char *path) {
  if (!valid_handle(handle) || !g_json_docs[handle]) {
    return NULL;
  }
  return navigate(g_json_docs[handle]->root, path);
}

// ---- handle API ------------------------------------------------------------

int64_t json_parse(char *s) {
  int64_t status = JSON_ERR_MALFORMED;
  JVal *root = json_parse_tree(s, "jsonParse", &status);
  JDoc *doc = root ? calloc(1, sizeof(JDoc)) : NULL;
  if (!doc) {
    if (root) {
      json_tree_free(root);
      osp_io_error_text("jsonParse", "document", "out of memory");
    }
    return status;
  }
  doc->root = root;
  json_table_lock();
  int64_t handle = claim_slot((void **)g_json_docs, doc);
  json_table_unlock();
  if (handle < 0) {
    doc_free(doc);
    osp_io_error_text("jsonParse", "document", "too many open documents");
    return JSON_ERR_TABLE_FULL;
  }
  return handle;
}

char *json_get(int64_t handle, char *path) {
  json_table_lock();
  const JVal *v = lookup(handle, path ? path : "");
  char *out = NULL;
  if (v) {
    switch (v->type) {
    case J_STR:
    case J_NUM:
      out = strdup(v->str);
      break;
    case J_BOOL:
      out = strdup(v->bval ? "true" : "false");
      break;
    case J_NULL:
      out = strdup("null");
      break;
    default:
      out = NULL; // arrays/objects are not scalars
      break;
    }
  }
  json_table_unlock();
  return out;
}

int64_t json_length(int64_t handle, char *path) {
  json_table_lock();
  const JVal *v = lookup(handle, path ? path : "");
  int64_t len = -1;
  if (v) {
    if (v->type == J_ARR) {
      len = (int64_t)v->count;
    } else if (v->type == J_OBJ) {
      len = (int64_t)v->nmemb;
    }
  }
  json_table_unlock();
  return len;
}

int64_t json_free(int64_t handle) {
  if (!valid_handle(handle)) {
    return -1;
  }
  json_table_lock();
  JDoc *doc = g_json_docs[handle];
  g_json_docs[handle] = NULL;
  json_table_unlock();
  if (!doc) {
    return -1;
  }
  doc_free(doc);
  return 0;
}

// ---- typed decode index ----------------------------------------------------

static size_t child_count(const JVal *v) {
  return v->type == J_ARR ? v->count : v->type == J_OBJ ? v->nmemb : 0;
}

static const JVal *child_at(const JVal *v, size_t i) {
  return v->type == J_ARR ? v->items[i] : v->vals[i];
}

static bool is_integer_text(const char *text) {
  return strpbrk(text, ".eE") == NULL;
}

// Read a number into its node: an integer lexeme must fit int64 and a
// fractional/exponent one must be a finite double, so the Json type holds it
// without rounding. The lexeme already passed the RFC 8259 grammar
// [BUILTIN-JSON-NUMBER].
static const char *read_number(JNode *n) {
  errno = 0;
  if (is_integer_text(n->v->str)) {
    n->number.i = strtoll(n->v->str, NULL, 10);
    return errno == ERANGE ? "integer outside the 64-bit range" : NULL;
  }
  n->number.f = strtod(n->v->str, NULL);
  return isfinite(n->number.f) ? NULL : "number outside the float range";
}

// Lay the tree out breadth-first and check every number on the way, so the
// prelude's walk cannot meet a value it has no variant for.
static const char *index_nodes(JDoc *doc) {
  size_t cap = 16;
  doc->nodes = malloc(cap * sizeof(JNode));
  doc->count = doc->nodes ? 1 : 0;
  if (!doc->nodes) {
    return "out of memory";
  }
  doc->nodes[0].v = doc->root;
  for (size_t at = 0; at < doc->count; at++) {
    const JVal *v = doc->nodes[at].v;
    const char *bad = v->type == J_NUM ? read_number(&doc->nodes[at]) : NULL;
    size_t kids = child_count(v);
    if (bad) {
      return bad;
    }
    if (doc->count + kids > cap) {
      cap = (doc->count + kids) * 2;
      JNode *grown = realloc(doc->nodes, cap * sizeof(JNode));
      if (!grown) {
        return "out of memory";
      }
      doc->nodes = grown;
    }
    doc->nodes[at].first = doc->count;
    for (size_t i = 0; i < kids; i++) {
      doc->nodes[doc->count++].v = child_at(v, i);
    }
  }
  return NULL;
}

// Parse and index in one step, so a document the Json type cannot hold is a
// jsonDecode failure with its reason, never a handle the walk trips over.
int64_t json_decode_open(char *text) {
  int64_t status = JSON_ERR_MALFORMED;
  JVal *root = json_parse_tree(text, "jsonDecode", &status);
  if (!root) {
    return status;
  }
  JDoc *doc = calloc(1, sizeof(JDoc));
  const char *why = !doc ? "out of memory" : (doc->root = root, index_nodes(doc));
  if (why) {
    doc ? doc_free(doc) : json_tree_free(root);
    osp_io_error_text("jsonDecode", "document", why);
    return JSON_ERR_MALFORMED;
  }
  json_table_lock();
  int64_t handle = claim_slot((void **)g_json_docs, doc);
  json_table_unlock();
  if (handle < 0) {
    doc_free(doc);
    osp_io_error_text("jsonDecode", "document", "too many open documents");
    return JSON_ERR_TABLE_FULL;
  }
  return handle;
}

// The indexed node, or NULL for a handle without an index or an id outside
// it. Callers hold the lock.
static const JNode *node_at(int64_t handle, int64_t node) {
  JDoc *doc = valid_handle(handle) ? g_json_docs[handle] : NULL;
  bool indexed = doc && doc->nodes && node >= 0;
  return indexed && (size_t)node < doc->count ? &doc->nodes[node] : NULL;
}

// The kind codes the prelude matches on. Booleans carry their value in the
// kind, and numbers their int/float classification, so one call per node
// decides its variant [BUILTIN-JSON-TYPED].
enum { K_NULL, K_FALSE, K_TRUE, K_INT, K_FLOAT, K_STRING, K_ARRAY, K_OBJECT };

static int64_t kind_of(const JVal *v) {
  switch (v->type) {
  case J_NULL: return K_NULL;
  case J_BOOL: return v->bval ? K_TRUE : K_FALSE;
  case J_NUM: return is_integer_text(v->str) ? K_INT : K_FLOAT;
  case J_STR: return K_STRING;
  case J_ARR: return K_ARRAY;
  default: return K_OBJECT;
  }
}

int64_t json_node_kind(int64_t handle, int64_t node) {
  json_table_lock();
  const JNode *n = node_at(handle, node);
  int64_t kind = n ? kind_of(n->v) : -1;
  json_table_unlock();
  return kind;
}

int64_t json_node_count(int64_t handle, int64_t node) {
  json_table_lock();
  const JNode *n = node_at(handle, node);
  int64_t count = n ? (int64_t)child_count(n->v) : -1;
  json_table_unlock();
  return count;
}

int64_t json_node_child(int64_t handle, int64_t node, int64_t i) {
  json_table_lock();
  const JNode *n = node_at(handle, node);
  bool inside = n && i >= 0 && (size_t)i < child_count(n->v);
  int64_t child = inside ? (int64_t)(n->first + (size_t)i) : -1;
  json_table_unlock();
  return child;
}

char *json_node_key(int64_t handle, int64_t node, int64_t i) {
  json_table_lock();
  const JNode *n = node_at(handle, node);
  bool member = n && n->v->type == J_OBJ && i >= 0 && (size_t)i < n->v->nmemb;
  char *key = member ? strdup(n->v->keys[i]) : NULL;
  json_table_unlock();
  return key;
}

// A string's decoded value.
char *json_node_text(int64_t handle, int64_t node) {
  json_table_lock();
  const JNode *n = node_at(handle, node);
  char *out = n && n->v->type == J_STR ? strdup(n->v->str) : NULL;
  json_table_unlock();
  return out;
}

// An integer node's value; 0 for any other node.
int64_t json_node_int(int64_t handle, int64_t node) {
  json_table_lock();
  const JNode *n = node_at(handle, node);
  int64_t value = n && kind_of(n->v) == K_INT ? n->number.i : 0;
  json_table_unlock();
  return value;
}

// A float node's value; 0.0 for any other node.
double json_node_float(int64_t handle, int64_t node) {
  json_table_lock();
  const JNode *n = node_at(handle, node);
  double value = n && kind_of(n->v) == K_FLOAT ? n->number.f : 0.0;
  json_table_unlock();
  return value;
}

// ---- single-buffer writer --------------------------------------------------
//
// The prelude walks a Json value and appends each token here; the document is
// built in one growing C buffer and handed to Osprey as ONE string when the
// writer closes. Separators are the writer's job: a value that follows a
// value gets its ',', and `[`, `{` and `:` reset that, so the walk is a plain
// traversal. Failure is sticky: the first refusal is kept, later writes are
// no-ops, and close reports it -- so the walk needs no per-token checks.

int64_t json_writer_open(void) {
  JWriter *w = calloc(1, sizeof(JWriter));
  if (!w) {
    osp_io_error_text("jsonEncode", "value", "out of memory");
    return -1;
  }
  json_table_lock();
  int64_t handle = claim_slot((void **)g_json_writers, w);
  json_table_unlock();
  if (handle < 0) {
    free(w);
    osp_io_error_text("jsonEncode", "value", "too many open writers");
  }
  return handle;
}

// The live writer behind `handle`, or NULL. Writers belong to the one walk
// that opened them, so appends run outside the lock.
static JWriter *writer_at(int64_t handle) {
  json_table_lock();
  JWriter *w = valid_handle(handle) ? g_json_writers[handle] : NULL;
  json_table_unlock();
  return w && !w->why ? w : NULL;
}

static int64_t writer_status(JWriter *w, bool ok, const char *why) {
  if (w && !ok && !w->why) {
    w->why = why;
  }
  return w && !w->why ? 0 : -1;
}

// Start a value: its separator, then whether it leaves a value behind (a
// scalar does; an opening bracket does not until it closes).
static bool begin_value(JWriter *w, bool complete) {
  bool ok = !w->after_value || json_buf_append(&w->buf, ",", 1);
  w->after_value = complete;
  return ok;
}

// A structural token or a literal (`null`, `true`, `false`).
int64_t json_write_raw(int64_t handle, char *token) {
  JWriter *w = writer_at(handle);
  if (!w || !token) {
    return writer_status(w, false, "out of memory");
  }
  char lead = token[0];
  bool closes = lead == ']' || lead == '}';
  bool ok = closes || lead == ':' || begin_value(w, lead != '[' && lead != '{');
  w->after_value = closes || (w->after_value && lead != ':');
  ok = ok && json_buf_append(&w->buf, token, strlen(token));
  return writer_status(w, ok, "out of memory");
}

int64_t json_write_string(int64_t handle, char *s) {
  JWriter *w = writer_at(handle);
  bool ok = w && s && begin_value(w, true) && json_quote_into(&w->buf, s);
  return writer_status(w, ok, "string is not valid UTF-8");
}

int64_t json_write_int(int64_t handle, int64_t n) {
  char digits[24];
  int len = snprintf(digits, sizeof(digits), "%lld", (long long)n);
  JWriter *w = writer_at(handle);
  bool ok = w && begin_value(w, true) &&
            json_buf_append(&w->buf, digits, (size_t)len);
  return writer_status(w, ok, "out of memory");
}

// The float spelling toString prints, which keeps a `.0` or exponent so the
// text decodes back to a float. NaN and infinities have no JSON spelling.
int64_t json_write_float(int64_t handle, double f) {
  JWriter *w = writer_at(handle);
  if (!isfinite(f)) {
    return writer_status(w, false, "float is not finite");
  }
  char *text = w ? osp_float_to_string(f) : NULL;
  bool ok = text && begin_value(w, true) &&
            json_buf_append(&w->buf, text, strlen(text));
  free(text);
  return writer_status(w, ok, "out of memory");
}

char *json_writer_close(int64_t handle) {
  json_table_lock();
  JWriter *w = valid_handle(handle) ? g_json_writers[handle] : NULL;
  if (w) {
    g_json_writers[handle] = NULL;
  }
  json_table_unlock();
  const char *why = !w ? "invalid writer handle"
                    : w->why ? w->why
                    : !w->buf.data && !json_buf_append(&w->buf, "", 0)
                        ? "out of memory"
                        : NULL;
  char *out = w && !why ? w->buf.data : NULL;
  if (w && why) {
    free(w->buf.data);
  }
  free(w);
  if (why) {
    osp_io_error_text("jsonEncode", "value", why);
  }
  return out;
}
