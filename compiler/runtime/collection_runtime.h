#ifndef OSPREY_COLLECTION_RUNTIME_H
#define OSPREY_COLLECTION_RUNTIME_H

#include <stdint.h>

/*
 * Shared types and ABI for List<T> and Map<K,V>.
 *
 * Spec: docs/specs/0004-TypeSystem.md#collection-types
 * Plan: docs/plans/collections.md
 *
 * Implements [TYPE-LIST], [TYPE-MAP], [TYPE-MAP-LOOKUP], [TYPE-MAP-OPS].
 *
 * Every element is stored as an int64_t. Pointers (strings, nested
 * collections, records) are cast to int64_t at storage time. Codegen on the
 * Go side is responsible for boxing/unboxing.
 *
 * Memory model: leak-semantic, matching the existing Osprey runtime
 * (malloc, never free). Structural sharing is achieved via path-copying;
 * old versions remain valid because no node is ever freed.
 */

#define OSPREY_LIST_BITS 5
#define OSPREY_LIST_BRANCH 32
#define OSPREY_LIST_MASK 31

#define OSPREY_MAP_BITS 5
#define OSPREY_MAP_BRANCH 32

/* Key-type tags for Map. Codegen passes one of these at map creation. */
typedef enum {
  OSPREY_KEY_INT = 0,
  OSPREY_KEY_STRING = 1,
  OSPREY_KEY_BOOL = 2
} OspreyKeyType;

/* Opaque handles. */
typedef struct OspreyList OspreyList;
typedef struct OspreyListBuilder OspreyListBuilder;
typedef struct OspreyListIter OspreyListIter;

typedef struct OspreyMap OspreyMap;
typedef struct OspreyMapBuilder OspreyMapBuilder;
typedef struct OspreyMapIter OspreyMapIter;

/* ============ List API ============ */

/*
 * Pointer parameters are not declared `const` even though the runtime never
 * mutates inputs — returning a (possibly aliased) input from `concat` /
 * `drop` etc. would require dropping const, which `-Wcast-qual` rejects.
 * Treat every collection pointer as read-only by contract.
 */

OspreyList *osprey_list_empty(void);
int64_t osprey_list_length(OspreyList *l);
/* 1 if 0 <= i < length, else 0. */
int osprey_list_in_bounds(OspreyList *l, int64_t i);
/* Caller must ensure osprey_list_in_bounds(l, i) before calling. */
int64_t osprey_list_get(OspreyList *l, int64_t i);
OspreyList *osprey_list_set(OspreyList *l, int64_t i, int64_t v);
OspreyList *osprey_list_append(OspreyList *l, int64_t v);
OspreyList *osprey_list_prepend(OspreyList *l, int64_t v);
OspreyList *osprey_list_concat(OspreyList *a, OspreyList *b);
OspreyList *osprey_list_drop(OspreyList *l, int64_t n);
OspreyList *osprey_list_reverse(OspreyList *l);

OspreyListBuilder *osprey_list_builder_new(void);
void osprey_list_builder_push(OspreyListBuilder *b, int64_t v);
OspreyList *osprey_list_builder_seal(OspreyListBuilder *b);

OspreyListIter *osprey_list_iter_new(OspreyList *l);
/* Returns 1 if a value was produced (in *out), 0 if exhausted. */
int osprey_list_iter_next(OspreyListIter *it, int64_t *out);

/* ============ Map API ============ */

OspreyMap *osprey_map_empty(OspreyKeyType key_type);
int64_t osprey_map_length(OspreyMap *m);
/* 1 if present, 0 if absent. */
int osprey_map_contains(OspreyMap *m, int64_t key);
/* Caller must ensure osprey_map_contains(m, k) before calling. */
int64_t osprey_map_get(OspreyMap *m, int64_t key);
OspreyMap *osprey_map_set(OspreyMap *m, int64_t key, int64_t value);
OspreyMap *osprey_map_remove(OspreyMap *m, int64_t key);
/* Right-biased: keys in b override keys in a. */
OspreyMap *osprey_map_merge(OspreyMap *a, OspreyMap *b);

OspreyMapBuilder *osprey_map_builder_new(OspreyKeyType key_type);
void osprey_map_builder_put(OspreyMapBuilder *b, int64_t key, int64_t value);
OspreyMap *osprey_map_builder_seal(OspreyMapBuilder *b);

OspreyMapIter *osprey_map_iter_new(OspreyMap *m);
/* Returns 1 if a (key, value) was produced, 0 if exhausted. */
int osprey_map_iter_next(OspreyMapIter *it, int64_t *out_key, int64_t *out_value);

#endif /* OSPREY_COLLECTION_RUNTIME_H */
