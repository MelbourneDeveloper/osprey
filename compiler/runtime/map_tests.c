#include "collection_runtime.h"

#include <assert.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/* Test helpers ----------------------------------------------------- */

static OspreyMap *make_int_map(int64_t start, int64_t end) {
  OspreyMap *m = osprey_map_empty(OSPREY_KEY_INT);
  for (int64_t i = start; i < end; i++) {
    m = osprey_map_set(m, i, i * 11);
  }
  return m;
}

/* Tests ------------------------------------------------------------ */

static void test_empty(void) {
  OspreyMap *m = osprey_map_empty(OSPREY_KEY_INT);
  assert(osprey_map_length(m) == 0);
  assert(osprey_map_contains(m, 0) == 0);
  assert(osprey_map_contains(m, 42) == 0);
  /* set on empty creates 1-entry map */
  OspreyMap *one = osprey_map_set(m, 1, 100);
  assert(osprey_map_length(one) == 1);
  assert(osprey_map_get(one, 1) == 100);
  /* original empty untouched */
  assert(osprey_map_length(m) == 0);
  printf("  map_empty OK\n");
}

static void test_int_keys_persistent(void) {
  OspreyMap *m = make_int_map(0, 200);
  assert(osprey_map_length(m) == 200);
  for (int64_t i = 0; i < 200; i++) {
    assert(osprey_map_contains(m, i) == 1);
    assert(osprey_map_get(m, i) == i * 11);
  }
  assert(osprey_map_contains(m, 200) == 0);
  assert(osprey_map_contains(m, -1) == 0);
  printf("  map_int_keys (size=200) OK\n");
}

static void test_overwrite_persistent(void) {
  OspreyMap *m = make_int_map(0, 50);
  OspreyMap *m2 = osprey_map_set(m, 25, -99);
  assert(osprey_map_length(m2) == 50);
  assert(osprey_map_get(m2, 25) == -99);
  /* original keeps 25 */
  assert(osprey_map_get(m, 25) == 275);
  /* spot-check neighbours unchanged */
  assert(osprey_map_get(m2, 24) == 264);
  assert(osprey_map_get(m2, 26) == 286);
  printf("  map_overwrite_persistent OK\n");
}

static void test_string_keys(void) {
  OspreyMap *m = osprey_map_empty(OSPREY_KEY_STRING);
  const char *keys[] = {"alice", "bob", "charlie", "dave", "eve",
                         "frank", "grace", "heidi"};
  for (int64_t i = 0; i < 8; i++) {
    m = osprey_map_set(m, (int64_t)(uintptr_t)keys[i], i * 100);
  }
  assert(osprey_map_length(m) == 8);
  /* Look up using independent buffers (forces value-equality, not
     pointer-equality). */
  for (int64_t i = 0; i < 8; i++) {
    char buf[16];
    strncpy(buf, keys[i], sizeof(buf) - 1);
    buf[sizeof(buf) - 1] = '\0';
    assert(osprey_map_contains(m, (int64_t)(uintptr_t)buf) == 1);
    assert(osprey_map_get(m, (int64_t)(uintptr_t)buf) == i * 100);
  }
  /* Negative lookup */
  assert(osprey_map_contains(m, (int64_t)(uintptr_t)"missing") == 0);
  printf("  map_string_keys (size=8, value-equality) OK\n");
}

static void test_bool_keys(void) {
  OspreyMap *m = osprey_map_empty(OSPREY_KEY_BOOL);
  m = osprey_map_set(m, 0, 100);
  m = osprey_map_set(m, 1, 200);
  assert(osprey_map_length(m) == 2);
  assert(osprey_map_get(m, 0) == 100);
  assert(osprey_map_get(m, 1) == 200);
  /* Overwrite */
  m = osprey_map_set(m, 1, 999);
  assert(osprey_map_length(m) == 2);
  assert(osprey_map_get(m, 1) == 999);
  printf("  map_bool_keys OK\n");
}

static void test_string_keys_with_prefix(void) {
  /* Strings sharing prefixes shouldn't collide. */
  OspreyMap *m = osprey_map_empty(OSPREY_KEY_STRING);
  m = osprey_map_set(m, (int64_t)(uintptr_t)"a", 1);
  m = osprey_map_set(m, (int64_t)(uintptr_t)"ab", 2);
  m = osprey_map_set(m, (int64_t)(uintptr_t)"abc", 3);
  m = osprey_map_set(m, (int64_t)(uintptr_t)"abcd", 4);
  assert(osprey_map_length(m) == 4);
  assert(osprey_map_get(m, (int64_t)(uintptr_t)"a") == 1);
  assert(osprey_map_get(m, (int64_t)(uintptr_t)"ab") == 2);
  assert(osprey_map_get(m, (int64_t)(uintptr_t)"abc") == 3);
  assert(osprey_map_get(m, (int64_t)(uintptr_t)"abcd") == 4);
  printf("  map_string_keys_with_prefix OK\n");
}

static void test_remove_basic(void) {
  OspreyMap *m = make_int_map(0, 50);
  OspreyMap *m2 = osprey_map_remove(m, 25);
  assert(osprey_map_length(m2) == 49);
  assert(osprey_map_contains(m2, 25) == 0);
  /* original keeps it */
  assert(osprey_map_contains(m, 25) == 1);
  /* Removing absent is no-op */
  OspreyMap *m3 = osprey_map_remove(m, 999);
  assert(osprey_map_length(m3) == 50);
  /* Remove every entry */
  OspreyMap *cleared = m;
  for (int64_t i = 0; i < 50; i++) {
    cleared = osprey_map_remove(cleared, i);
  }
  assert(osprey_map_length(cleared) == 0);
  /* m unchanged */
  assert(osprey_map_length(m) == 50);
  printf("  map_remove_basic + clear-all OK\n");
}

static void test_remove_then_set_round_trip(void) {
  OspreyMap *m = make_int_map(0, 20);
  OspreyMap *removed = osprey_map_remove(m, 10);
  OspreyMap *re_added = osprey_map_set(removed, 10, 999);
  assert(osprey_map_length(re_added) == 20);
  assert(osprey_map_get(re_added, 10) == 999);
  /* Other entries intact */
  for (int64_t i = 0; i < 20; i++) {
    if (i == 10) continue;
    assert(osprey_map_get(re_added, i) == i * 11);
  }
  printf("  map_remove_then_set_round_trip OK\n");
}

static void test_merge_right_biased(void) {
  OspreyMap *a = osprey_map_empty(OSPREY_KEY_INT);
  OspreyMap *b = osprey_map_empty(OSPREY_KEY_INT);
  for (int64_t i = 0; i < 10; i++) a = osprey_map_set(a, i, 100 + i);
  for (int64_t i = 5; i < 15; i++) b = osprey_map_set(b, i, 200 + i);
  OspreyMap *m = osprey_map_merge(a, b);
  assert(osprey_map_length(m) == 15);
  /* a-only keys keep a's values (0..4) */
  for (int64_t i = 0; i < 5; i++) assert(osprey_map_get(m, i) == 100 + i);
  /* overlap (5..9): b wins */
  for (int64_t i = 5; i < 10; i++) assert(osprey_map_get(m, i) == 200 + i);
  /* b-only keys (10..14) keep b's values */
  for (int64_t i = 10; i < 15; i++) assert(osprey_map_get(m, i) == 200 + i);
  printf("  map_merge_right_biased OK\n");
}

static void test_merge_empty_edges(void) {
  OspreyMap *e = osprey_map_empty(OSPREY_KEY_INT);
  OspreyMap *a = make_int_map(0, 10);
  /* empty + a == a (semantically) */
  OspreyMap *ea = osprey_map_merge(e, a);
  assert(osprey_map_length(ea) == 10);
  for (int64_t i = 0; i < 10; i++) assert(osprey_map_get(ea, i) == i * 11);
  /* a + empty == a */
  OspreyMap *ae = osprey_map_merge(a, e);
  assert(osprey_map_length(ae) == 10);
  /* empty + empty == empty */
  OspreyMap *ee = osprey_map_merge(e, e);
  assert(osprey_map_length(ee) == 0);
  printf("  map_merge_empty_edges OK\n");
}

static void test_merge_with_self(void) {
  OspreyMap *m = make_int_map(0, 30);
  OspreyMap *self = osprey_map_merge(m, m);
  assert(osprey_map_length(self) == 30);
  for (int64_t i = 0; i < 30; i++) assert(osprey_map_get(self, i) == i * 11);
  printf("  map_merge_with_self OK\n");
}

static void test_iter_completeness(void) {
  OspreyMap *m = make_int_map(0, 250);
  OspreyMapIter *it = osprey_map_iter_new(m);
  int64_t count = 0;
  int64_t sum_k = 0;
  int64_t sum_v = 0;
  int64_t k, v;
  while (osprey_map_iter_next(it, &k, &v)) {
    count++;
    sum_k += k;
    sum_v += v;
    /* Invariant: v == k * 11 */
    assert(v == k * 11);
  }
  assert(count == 250);
  /* Gauss: 0..249 sum = 249*250/2 = 31125; v sum = 11x */
  assert(sum_k == 31125);
  assert(sum_v == 11 * 31125);
  free(it);
  printf("  map_iter_completeness (count + sum invariants) OK\n");
}

static void test_iter_empty(void) {
  OspreyMap *e = osprey_map_empty(OSPREY_KEY_INT);
  OspreyMapIter *it = osprey_map_iter_new(e);
  int64_t k, v;
  assert(osprey_map_iter_next(it, &k, &v) == 0);
  free(it);
  printf("  map_iter_empty OK\n");
}

static void test_builder_vs_set_equivalence(void) {
  /* Build a map two ways and verify they have the same content. */
  OspreyMap *via_set = osprey_map_empty(OSPREY_KEY_INT);
  for (int64_t i = 0; i < 100; i++) via_set = osprey_map_set(via_set, i, i + 1);

  OspreyMapBuilder *b = osprey_map_builder_new(OSPREY_KEY_INT);
  for (int64_t i = 0; i < 100; i++) osprey_map_builder_put(b, i, i + 1);
  OspreyMap *via_builder = osprey_map_builder_seal(b);

  assert(osprey_map_length(via_set) == 100);
  assert(osprey_map_length(via_builder) == 100);
  for (int64_t i = 0; i < 100; i++) {
    assert(osprey_map_get(via_set, i) == i + 1);
    assert(osprey_map_get(via_builder, i) == i + 1);
  }
  printf("  map_builder_vs_set_equivalence OK\n");
}

static void test_builder_overwrite(void) {
  /* Putting the same key twice should overwrite. */
  OspreyMapBuilder *b = osprey_map_builder_new(OSPREY_KEY_INT);
  osprey_map_builder_put(b, 1, 100);
  osprey_map_builder_put(b, 2, 200);
  osprey_map_builder_put(b, 1, 999); /* overwrite */
  OspreyMap *m = osprey_map_builder_seal(b);
  assert(osprey_map_length(m) == 2);
  assert(osprey_map_get(m, 1) == 999);
  assert(osprey_map_get(m, 2) == 200);
  printf("  map_builder_overwrite OK\n");
}

static void test_stress_5000(void) {
  /* 5000 entries forces deep HAMT nodes. */
  OspreyMap *m = make_int_map(0, 5000);
  assert(osprey_map_length(m) == 5000);
  /* Random-ish spot checks across the range */
  int64_t spots[] = {0, 1, 31, 32, 100, 1023, 1024, 2500, 4999};
  for (size_t i = 0; i < sizeof(spots) / sizeof(spots[0]); i++) {
    assert(osprey_map_contains(m, spots[i]) == 1);
    assert(osprey_map_get(m, spots[i]) == spots[i] * 11);
  }
  printf("  map_stress_5000 OK\n");
}

static void test_remove_during_iteration_irrelevant(void) {
  /* Iterator snapshots: removing from the source after creating an iter
     must not change what the iter sees (because the source is immutable —
     map_remove returns a NEW map). */
  OspreyMap *m = make_int_map(0, 20);
  OspreyMapIter *it = osprey_map_iter_new(m);
  /* This produces a new map; m is unchanged. */
  OspreyMap *m2 = osprey_map_remove(m, 0);
  (void)m2;
  int64_t count = 0;
  int64_t k, v;
  while (osprey_map_iter_next(it, &k, &v)) count++;
  assert(count == 20);
  free(it);
  printf("  map_iter_immutable_snapshot OK\n");
}

void run_map_tests(void) {
  printf("Map:\n");
  test_empty();
  test_int_keys_persistent();
  test_overwrite_persistent();
  test_string_keys();
  test_bool_keys();
  test_string_keys_with_prefix();
  test_remove_basic();
  test_remove_then_set_round_trip();
  test_merge_right_biased();
  test_merge_empty_edges();
  test_merge_with_self();
  test_iter_completeness();
  test_iter_empty();
  test_builder_vs_set_equivalence();
  test_builder_overwrite();
  test_stress_5000();
  test_remove_during_iteration_irrelevant();
}
