#include "map_runtime_internal.h"

#include <stdint.h>
#include <stdlib.h>
#include <string.h>

/*
 * HAMT node algebra: hashing, key equality, assoc, lookup, remove.
 * Public API and iterator live in map_runtime.c. Split per CLAUDE.md
 * 500-LoC limit.
 *
 * Implements [TYPE-MAP] node-level invariants from
 * compiler/spec/0004-TypeSystem.md.
 */

static uint32_t hash_int(int64_t key) {
  uint64_t k = (uint64_t)key;
  k = (k ^ (k >> 30)) * (uint64_t)0xbf58476d1ce4e5b5ULL;
  k = (k ^ (k >> 27)) * (uint64_t)0x94d049bb133111ebULL;
  k = k ^ (k >> 31);
  return (uint32_t)(k & 0xffffffffu);
}

static uint32_t hash_bool(int64_t key) { return (key != 0) ? 1u : 0u; }

static uint32_t hash_string(int64_t key) {
  const char *s = (const char *)(uintptr_t)key;
  if (s == NULL) {
    return 0u;
  }
  uint32_t h = 0x811c9dc5u;
  while (*s != '\0') {
    h ^= (uint32_t)(unsigned char)*s;
    h *= 0x01000193u;
    s++;
  }
  return h;
}

uint32_t osprey_map_hash_key(OspreyKeyType kt, int64_t key) {
  switch (kt) {
  case OSPREY_KEY_INT:
    return hash_int(key);
  case OSPREY_KEY_BOOL:
    return hash_bool(key);
  case OSPREY_KEY_STRING:
    return hash_string(key);
  default:
    return 0u;
  }
}

int osprey_map_keys_equal(OspreyKeyType kt, int64_t a, int64_t b) {
  if (kt == OSPREY_KEY_STRING) {
    const char *sa = (const char *)(uintptr_t)a;
    const char *sb = (const char *)(uintptr_t)b;
    if (sa == sb) {
      return 1;
    }
    if (sa == NULL || sb == NULL) {
      return 0;
    }
    return (strcmp(sa, sb) == 0) ? 1 : 0;
  }
  return (a == b) ? 1 : 0;
}

static uint32_t bit_for(uint32_t hash, int32_t shift) {
  return 1u << ((hash >> shift) & (uint32_t)OSPREY_LIST_MASK);
}

static uint32_t index_for(uint32_t bitmap, uint32_t bit) {
  return (uint32_t)__builtin_popcount(bitmap & (bit - 1u));
}

static OspreyMapNode *make_leaf(uint32_t h, int64_t key, int64_t value) {
  OspreyMapNode *n = (OspreyMapNode *)calloc(1, sizeof(OspreyMapNode));
  n->kind = NODE_LEAF;
  n->hash = h;
  n->leaf_key = key;
  n->leaf_value = value;
  return n;
}

static OspreyMapNode *make_internal(uint32_t bitmap, uint32_t count,
                                    OspreyMapNode **children) {
  OspreyMapNode *n = (OspreyMapNode *)calloc(1, sizeof(OspreyMapNode));
  n->kind = NODE_INTERNAL;
  n->bitmap = bitmap;
  n->count = count;
  n->children = children;
  return n;
}

static OspreyMapNode *make_collision(uint32_t h, uint32_t count, int64_t *keys,
                                     int64_t *values) {
  OspreyMapNode *n = (OspreyMapNode *)calloc(1, sizeof(OspreyMapNode));
  n->kind = NODE_COLLISION;
  n->hash = h;
  n->count = count;
  n->coll_keys = keys;
  n->coll_values = values;
  return n;
}

static OspreyMapNode **clone_children(OspreyMapNode **src, uint32_t count) {
  OspreyMapNode **out =
      (OspreyMapNode **)calloc((size_t)count + 1, sizeof(OspreyMapNode *));
  if (src != NULL) {
    memcpy(out, src, (size_t)count * sizeof(OspreyMapNode *));
  }
  return out;
}

static OspreyMapNode *merge_leaves(OspreyMapNode *a, OspreyMapNode *b,
                                   int32_t shift, OspreyKeyType kt) {
  if (shift >= MAP_HASH_BITS) {
    if (a->kind == NODE_COLLISION && b->kind == NODE_LEAF) {
      int64_t *ks = (int64_t *)calloc((size_t)a->count + 1, sizeof(int64_t));
      int64_t *vs = (int64_t *)calloc((size_t)a->count + 1, sizeof(int64_t));
      memcpy(ks, a->coll_keys, (size_t)a->count * sizeof(int64_t));
      memcpy(vs, a->coll_values, (size_t)a->count * sizeof(int64_t));
      ks[a->count] = b->leaf_key;
      vs[a->count] = b->leaf_value;
      return make_collision(a->hash, a->count + 1u, ks, vs);
    }
    int64_t *ks = (int64_t *)calloc(2, sizeof(int64_t));
    int64_t *vs = (int64_t *)calloc(2, sizeof(int64_t));
    ks[0] = a->leaf_key;
    vs[0] = a->leaf_value;
    ks[1] = b->leaf_key;
    vs[1] = b->leaf_value;
    return make_collision(a->hash, 2u, ks, vs);
  }
  (void)kt;
  uint32_t bit_a = bit_for(a->hash, shift);
  uint32_t bit_b = bit_for(b->hash, shift);
  if (bit_a == bit_b) {
    OspreyMapNode **kids = (OspreyMapNode **)calloc(1, sizeof(OspreyMapNode *));
    kids[0] = merge_leaves(a, b, shift + OSPREY_MAP_BITS, kt);
    return make_internal(bit_a, 1u, kids);
  }
  OspreyMapNode **kids = (OspreyMapNode **)calloc(2, sizeof(OspreyMapNode *));
  uint32_t idx_a = index_for(bit_a | bit_b, bit_a);
  uint32_t idx_b = index_for(bit_a | bit_b, bit_b);
  kids[idx_a] = a;
  kids[idx_b] = b;
  return make_internal(bit_a | bit_b, 2u, kids);
}

OspreyMapNode *osprey_map_node_assoc(OspreyMapNode *node, int32_t shift,
                                     uint32_t hash, int64_t key, int64_t value,
                                     OspreyKeyType kt, int *grew) {
  if (node == NULL) {
    *grew = 1;
    return make_leaf(hash, key, value);
  }
  if (node->kind == NODE_LEAF) {
    if (node->hash == hash && osprey_map_keys_equal(kt, node->leaf_key, key)) {
      *grew = 0;
      return make_leaf(hash, key, value);
    }
    OspreyMapNode *new_leaf = make_leaf(hash, key, value);
    *grew = 1;
    return merge_leaves(node, new_leaf, shift, kt);
  }
  if (node->kind == NODE_COLLISION) {
    if (node->hash == hash) {
      for (uint32_t i = 0; i < node->count; i++) {
        if (osprey_map_keys_equal(kt, node->coll_keys[i], key)) {
          int64_t *ks =
              (int64_t *)calloc((size_t)node->count, sizeof(int64_t));
          int64_t *vs =
              (int64_t *)calloc((size_t)node->count, sizeof(int64_t));
          memcpy(ks, node->coll_keys, (size_t)node->count * sizeof(int64_t));
          memcpy(vs, node->coll_values, (size_t)node->count * sizeof(int64_t));
          vs[i] = value;
          *grew = 0;
          return make_collision(hash, node->count, ks, vs);
        }
      }
      int64_t *ks =
          (int64_t *)calloc((size_t)node->count + 1u, sizeof(int64_t));
      int64_t *vs =
          (int64_t *)calloc((size_t)node->count + 1u, sizeof(int64_t));
      memcpy(ks, node->coll_keys, (size_t)node->count * sizeof(int64_t));
      memcpy(vs, node->coll_values, (size_t)node->count * sizeof(int64_t));
      ks[node->count] = key;
      vs[node->count] = value;
      *grew = 1;
      return make_collision(hash, node->count + 1u, ks, vs);
    }
    OspreyMapNode *new_leaf = make_leaf(hash, key, value);
    *grew = 1;
    return merge_leaves(node, new_leaf, shift, kt);
  }
  uint32_t bit = bit_for(hash, shift);
  uint32_t idx = index_for(node->bitmap, bit);
  if ((node->bitmap & bit) != 0u) {
    OspreyMapNode *child = node->children[idx];
    OspreyMapNode *new_child = osprey_map_node_assoc(
        child, shift + OSPREY_MAP_BITS, hash, key, value, kt, grew);
    OspreyMapNode **new_kids = clone_children(node->children, node->count);
    new_kids[idx] = new_child;
    return make_internal(node->bitmap, node->count, new_kids);
  }
  *grew = 1;
  OspreyMapNode **new_kids = (OspreyMapNode **)calloc(
      (size_t)node->count + 1u, sizeof(OspreyMapNode *));
  memcpy(new_kids, node->children, (size_t)idx * sizeof(OspreyMapNode *));
  new_kids[idx] = make_leaf(hash, key, value);
  memcpy(new_kids + idx + 1, node->children + idx,
         (size_t)(node->count - idx) * sizeof(OspreyMapNode *));
  return make_internal(node->bitmap | bit, node->count + 1u, new_kids);
}

int osprey_map_node_lookup(OspreyMapNode *node, int32_t shift, uint32_t hash,
                           int64_t key, OspreyKeyType kt, int64_t *out) {
  while (node != NULL) {
    if (node->kind == NODE_LEAF) {
      if (node->hash == hash &&
          osprey_map_keys_equal(kt, node->leaf_key, key)) {
        *out = node->leaf_value;
        return 1;
      }
      return 0;
    }
    if (node->kind == NODE_COLLISION) {
      if (node->hash != hash) {
        return 0;
      }
      for (uint32_t i = 0; i < node->count; i++) {
        if (osprey_map_keys_equal(kt, node->coll_keys[i], key)) {
          *out = node->coll_values[i];
          return 1;
        }
      }
      return 0;
    }
    uint32_t bit = bit_for(hash, shift);
    if ((node->bitmap & bit) == 0u) {
      return 0;
    }
    node = node->children[index_for(node->bitmap, bit)];
    shift += OSPREY_MAP_BITS;
  }
  return 0;
}

OspreyMapNode *osprey_map_node_remove(OspreyMapNode *node, int32_t shift,
                                      uint32_t hash, int64_t key,
                                      OspreyKeyType kt, int *shrunk) {
  if (node == NULL) {
    *shrunk = 0;
    return NULL;
  }
  if (node->kind == NODE_LEAF) {
    if (node->hash == hash && osprey_map_keys_equal(kt, node->leaf_key, key)) {
      *shrunk = 1;
      return NULL;
    }
    *shrunk = 0;
    return node;
  }
  if (node->kind == NODE_COLLISION) {
    if (node->hash != hash) {
      *shrunk = 0;
      return node;
    }
    for (uint32_t i = 0; i < node->count; i++) {
      if (osprey_map_keys_equal(kt, node->coll_keys[i], key)) {
        *shrunk = 1;
        if (node->count == 2u) {
          uint32_t other = (i == 0u) ? 1u : 0u;
          return make_leaf(hash, node->coll_keys[other],
                           node->coll_values[other]);
        }
        uint32_t new_count = node->count - 1u;
        int64_t *ks = (int64_t *)calloc((size_t)new_count, sizeof(int64_t));
        int64_t *vs = (int64_t *)calloc((size_t)new_count, sizeof(int64_t));
        memcpy(ks, node->coll_keys, (size_t)i * sizeof(int64_t));
        memcpy(vs, node->coll_values, (size_t)i * sizeof(int64_t));
        memcpy(ks + i, node->coll_keys + i + 1,
               (size_t)(node->count - i - 1u) * sizeof(int64_t));
        memcpy(vs + i, node->coll_values + i + 1,
               (size_t)(node->count - i - 1u) * sizeof(int64_t));
        return make_collision(hash, new_count, ks, vs);
      }
    }
    *shrunk = 0;
    return node;
  }
  uint32_t bit = bit_for(hash, shift);
  if ((node->bitmap & bit) == 0u) {
    *shrunk = 0;
    return node;
  }
  uint32_t idx = index_for(node->bitmap, bit);
  OspreyMapNode *child = node->children[idx];
  OspreyMapNode *new_child = osprey_map_node_remove(
      child, shift + OSPREY_MAP_BITS, hash, key, kt, shrunk);
  if (!*shrunk) {
    return node;
  }
  if (new_child == NULL) {
    if (node->count == 1u) {
      return NULL;
    }
    uint32_t new_count = node->count - 1u;
    OspreyMapNode **new_kids =
        (OspreyMapNode **)calloc((size_t)new_count, sizeof(OspreyMapNode *));
    memcpy(new_kids, node->children, (size_t)idx * sizeof(OspreyMapNode *));
    memcpy(new_kids + idx, node->children + idx + 1,
           (size_t)(node->count - idx - 1u) * sizeof(OspreyMapNode *));
    return make_internal(node->bitmap & ~bit, new_count, new_kids);
  }
  OspreyMapNode **new_kids = clone_children(node->children, node->count);
  new_kids[idx] = new_child;
  return make_internal(node->bitmap, node->count, new_kids);
}
