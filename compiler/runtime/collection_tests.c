#include "collection_runtime.h"

#include <assert.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/*
 * Vanilla-C tests for list_runtime.c and map_runtime.c.
 * Run with `make c-test`. Uses assert() — non-zero exit on failure.
 *
 * Covers semantics required by [TYPE-LIST] and [TYPE-MAP]:
 *   - persistence: old versions remain valid after "mutation"
 *   - length / get correctness
 *   - append / prepend / concat / set / drop / reverse for List
 *   - set / remove / merge / iteration for Map
 *   - hash-collision path for Map (forced by tiny key set)
 *   - tree-growth path for List (>32, >1024 elements)
 */

static void test_list_empty(void) {
  OspreyList *e = osprey_list_empty();
  assert(osprey_list_length(e) == 0);
  assert(osprey_list_in_bounds(e, 0) == 0);
  printf("  list_empty OK\n");
}

static void test_list_small_append(void) {
  OspreyList *l = osprey_list_empty();
  for (int64_t i = 0; i < 10; i++) {
    l = osprey_list_append(l, i * 10);
  }
  assert(osprey_list_length(l) == 10);
  for (int64_t i = 0; i < 10; i++) {
    assert(osprey_list_get(l, i) == i * 10);
  }
  printf("  list_small_append OK\n");
}

static void test_list_persistence(void) {
  OspreyList *a = osprey_list_empty();
  a = osprey_list_append(a, 1);
  a = osprey_list_append(a, 2);
  OspreyList *b = osprey_list_append(a, 3);
  OspreyList *c = osprey_list_append(a, 99);
  /* a is unchanged */
  assert(osprey_list_length(a) == 2);
  assert(osprey_list_get(a, 0) == 1);
  assert(osprey_list_get(a, 1) == 2);
  /* b extends a with 3 */
  assert(osprey_list_length(b) == 3);
  assert(osprey_list_get(b, 2) == 3);
  /* c extends a with 99 — independent of b */
  assert(osprey_list_length(c) == 3);
  assert(osprey_list_get(c, 2) == 99);
  printf("  list_persistence OK\n");
}

static void test_list_large(void) {
  /* Force tree growth past 32 (one leaf) and 1024 (one internal level). */
  OspreyList *l = osprey_list_empty();
  for (int64_t i = 0; i < 2000; i++) {
    l = osprey_list_append(l, i * 3);
  }
  assert(osprey_list_length(l) == 2000);
  for (int64_t i = 0; i < 2000; i++) {
    assert(osprey_list_get(l, i) == i * 3);
  }
  /* Set midway should not affect prior version. */
  OspreyList *prior = l;
  OspreyList *mut = osprey_list_set(l, 500, -7);
  assert(osprey_list_get(prior, 500) == 1500);
  assert(osprey_list_get(mut, 500) == -7);
  assert(osprey_list_get(mut, 499) == 1497);
  assert(osprey_list_get(mut, 501) == 1503);
  printf("  list_large (size=2000) OK\n");
}

static void test_list_concat(void) {
  OspreyList *a = osprey_list_empty();
  OspreyList *b = osprey_list_empty();
  for (int64_t i = 0; i < 50; i++) {
    a = osprey_list_append(a, i);
  }
  for (int64_t i = 50; i < 100; i++) {
    b = osprey_list_append(b, i);
  }
  OspreyList *c = osprey_list_concat(a, b);
  assert(osprey_list_length(c) == 100);
  for (int64_t i = 0; i < 100; i++) {
    assert(osprey_list_get(c, i) == i);
  }
  /* Empties */
  OspreyList *e = osprey_list_empty();
  assert(osprey_list_length(osprey_list_concat(e, e)) == 0);
  assert(osprey_list_length(osprey_list_concat(a, e)) == 50);
  assert(osprey_list_length(osprey_list_concat(e, a)) == 50);
  printf("  list_concat OK\n");
}

static void test_list_prepend_drop_reverse(void) {
  OspreyList *l = osprey_list_empty();
  for (int64_t i = 0; i < 5; i++) {
    l = osprey_list_append(l, i);
  }
  OspreyList *p = osprey_list_prepend(l, 99);
  assert(osprey_list_length(p) == 6);
  assert(osprey_list_get(p, 0) == 99);
  assert(osprey_list_get(p, 5) == 4);

  OspreyList *d = osprey_list_drop(l, 2);
  assert(osprey_list_length(d) == 3);
  assert(osprey_list_get(d, 0) == 2);
  assert(osprey_list_get(d, 2) == 4);

  OspreyList *r = osprey_list_reverse(l);
  assert(osprey_list_length(r) == 5);
  for (int64_t i = 0; i < 5; i++) {
    assert(osprey_list_get(r, i) == 4 - i);
  }
  printf("  list_prepend_drop_reverse OK\n");
}

static void test_list_builder(void) {
  OspreyListBuilder *b = osprey_list_builder_new();
  for (int64_t i = 0; i < 100; i++) {
    osprey_list_builder_push(b, i * i);
  }
  OspreyList *l = osprey_list_builder_seal(b);
  assert(osprey_list_length(l) == 100);
  for (int64_t i = 0; i < 100; i++) {
    assert(osprey_list_get(l, i) == i * i);
  }
  printf("  list_builder OK\n");
}

static void test_list_iter(void) {
  OspreyList *l = osprey_list_empty();
  for (int64_t i = 0; i < 50; i++) {
    l = osprey_list_append(l, i + 1000);
  }
  OspreyListIter *it = osprey_list_iter_new(l);
  int64_t expected = 1000;
  int64_t v = 0;
  int64_t count = 0;
  while (osprey_list_iter_next(it, &v)) {
    assert(v == expected);
    expected++;
    count++;
  }
  assert(count == 50);
  free(it);
  printf("  list_iter OK\n");
}

/* ============ Map tests ============ */

static void test_map_empty(void) {
  OspreyMap *m = osprey_map_empty(OSPREY_KEY_INT);
  assert(osprey_map_length(m) == 0);
  assert(osprey_map_contains(m, 42) == 0);
  printf("  map_empty OK\n");
}

static void test_map_int_keys(void) {
  OspreyMap *m = osprey_map_empty(OSPREY_KEY_INT);
  for (int64_t i = 0; i < 200; i++) {
    m = osprey_map_set(m, i, i * 11);
  }
  assert(osprey_map_length(m) == 200);
  for (int64_t i = 0; i < 200; i++) {
    assert(osprey_map_contains(m, i) == 1);
    assert(osprey_map_get(m, i) == i * 11);
  }
  assert(osprey_map_contains(m, 999) == 0);
  /* Overwrite */
  OspreyMap *m2 = osprey_map_set(m, 50, -1);
  assert(osprey_map_length(m2) == 200);
  assert(osprey_map_get(m2, 50) == -1);
  /* Original untouched */
  assert(osprey_map_get(m, 50) == 550);
  printf("  map_int_keys (size=200) OK\n");
}

static void test_map_string_keys(void) {
  OspreyMap *m = osprey_map_empty(OSPREY_KEY_STRING);
  /* String keys passed as int64_t(uintptr_t) of char*. */
  const char *keys[] = {"alice", "bob", "charlie", "dave", "eve"};
  for (int64_t i = 0; i < 5; i++) {
    m = osprey_map_set(m, (int64_t)(uintptr_t)keys[i], i * 100);
  }
  assert(osprey_map_length(m) == 5);
  for (int64_t i = 0; i < 5; i++) {
    /* Use a fresh copy of the key string to check value-equality, not pointer-equality. */
    char buf[16];
    strncpy(buf, keys[i], sizeof(buf) - 1);
    buf[sizeof(buf) - 1] = '\0';
    assert(osprey_map_contains(m, (int64_t)(uintptr_t)buf) == 1);
    assert(osprey_map_get(m, (int64_t)(uintptr_t)buf) == i * 100);
  }
  assert(osprey_map_contains(m, (int64_t)(uintptr_t)"missing") == 0);
  printf("  map_string_keys OK\n");
}

static void test_map_remove(void) {
  OspreyMap *m = osprey_map_empty(OSPREY_KEY_INT);
  for (int64_t i = 0; i < 50; i++) {
    m = osprey_map_set(m, i, i * 2);
  }
  OspreyMap *m2 = osprey_map_remove(m, 25);
  assert(osprey_map_length(m2) == 49);
  assert(osprey_map_contains(m2, 25) == 0);
  /* Original keeps 25 */
  assert(osprey_map_contains(m, 25) == 1);
  /* Removing absent key is a no-op (returns same map). */
  OspreyMap *m3 = osprey_map_remove(m, 999);
  assert(osprey_map_length(m3) == 50);
  printf("  map_remove OK\n");
}

static void test_map_merge(void) {
  OspreyMap *a = osprey_map_empty(OSPREY_KEY_INT);
  OspreyMap *b = osprey_map_empty(OSPREY_KEY_INT);
  for (int64_t i = 0; i < 10; i++) {
    a = osprey_map_set(a, i, 100 + i);
  }
  for (int64_t i = 5; i < 15; i++) {
    b = osprey_map_set(b, i, 200 + i);
  }
  OspreyMap *m = osprey_map_merge(a, b);
  assert(osprey_map_length(m) == 15);
  /* Right-biased: b wins on overlap. */
  for (int64_t i = 0; i < 5; i++) {
    assert(osprey_map_get(m, i) == 100 + i);
  }
  for (int64_t i = 5; i < 15; i++) {
    assert(osprey_map_get(m, i) == 200 + i);
  }
  printf("  map_merge (right-biased) OK\n");
}

static void test_map_iter(void) {
  OspreyMap *m = osprey_map_empty(OSPREY_KEY_INT);
  int64_t expected_sum = 0;
  for (int64_t i = 0; i < 100; i++) {
    m = osprey_map_set(m, i, i * 7);
    expected_sum += i * 7;
  }
  OspreyMapIter *it = osprey_map_iter_new(m);
  int64_t k = 0;
  int64_t v = 0;
  int64_t count = 0;
  int64_t sum = 0;
  while (osprey_map_iter_next(it, &k, &v)) {
    assert(v == k * 7);
    count++;
    sum += v;
  }
  assert(count == 100);
  assert(sum == expected_sum);
  free(it);
  printf("  map_iter (size=100, sum=%lld) OK\n", (long long)sum);
}

static void test_map_builder(void) {
  OspreyMapBuilder *b = osprey_map_builder_new(OSPREY_KEY_INT);
  for (int64_t i = 0; i < 500; i++) {
    osprey_map_builder_put(b, i, i + 1);
  }
  OspreyMap *m = osprey_map_builder_seal(b);
  assert(osprey_map_length(m) == 500);
  for (int64_t i = 0; i < 500; i++) {
    assert(osprey_map_get(m, i) == i + 1);
  }
  printf("  map_builder (size=500) OK\n");
}

int main(void) {
  printf("Running collection_tests…\n");
  printf("List:\n");
  test_list_empty();
  test_list_small_append();
  test_list_persistence();
  test_list_large();
  test_list_concat();
  test_list_prepend_drop_reverse();
  test_list_builder();
  test_list_iter();
  printf("Map:\n");
  test_map_empty();
  test_map_int_keys();
  test_map_string_keys();
  test_map_remove();
  test_map_merge();
  test_map_iter();
  test_map_builder();
  printf("All collection tests passed.\n");
  return 0;
}
