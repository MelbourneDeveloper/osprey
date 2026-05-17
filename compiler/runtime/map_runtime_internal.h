#ifndef OSPREY_MAP_RUNTIME_INTERNAL_H
#define OSPREY_MAP_RUNTIME_INTERNAL_H

#include "collection_runtime.h"
#include <stdint.h>

/*
 * Internal types shared between map_runtime.c (public API + iter + builder)
 * and map_runtime_hamt.c (HAMT node algebra: hash, assoc, lookup, remove).
 *
 * Not part of the public ABI — included only by the two map_runtime
 * translation units to keep each file under 500 LoC per CLAUDE.md.
 */

#define MAP_HASH_BITS 32

typedef enum {
  NODE_INTERNAL = 0,
  NODE_LEAF = 1,
  NODE_COLLISION = 2
} OspreyMapNodeKind;

typedef struct OspreyMapNode {
  OspreyMapNodeKind kind;
  uint32_t bitmap;
  uint32_t count;
  uint32_t hash;
  struct OspreyMapNode **children;
  int64_t leaf_key;
  int64_t leaf_value;
  int64_t *coll_keys;
  int64_t *coll_values;
} OspreyMapNode;

struct OspreyMap {
  OspreyKeyType key_type;
  int64_t length;
  OspreyMapNode *root;
};

struct OspreyMapBuilder {
  OspreyKeyType key_type;
  int64_t length;
  OspreyMapNode *root;
};

struct OspreyMapIter {
  OspreyMap *map;
  OspreyMapNode *stack_nodes[8];
  uint32_t stack_slots[8];
  int32_t stack_depth;
  uint32_t coll_index;
};

/* Hashing / equality. */
uint32_t osprey_map_hash_key(OspreyKeyType kt, int64_t key);
int osprey_map_keys_equal(OspreyKeyType kt, int64_t a, int64_t b);

/* Node algebra. The grew / shrunk out-params report whether the operation
   changed cardinality. */
OspreyMapNode *osprey_map_node_assoc(OspreyMapNode *node, int32_t shift,
                                     uint32_t hash, int64_t key, int64_t value,
                                     OspreyKeyType kt, int *grew);
int osprey_map_node_lookup(OspreyMapNode *node, int32_t shift, uint32_t hash,
                           int64_t key, OspreyKeyType kt, int64_t *out);
OspreyMapNode *osprey_map_node_remove(OspreyMapNode *node, int32_t shift,
                                      uint32_t hash, int64_t key,
                                      OspreyKeyType kt, int *shrunk);

#endif /* OSPREY_MAP_RUNTIME_INTERNAL_H */
