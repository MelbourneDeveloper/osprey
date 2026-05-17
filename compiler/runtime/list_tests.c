#include "collection_runtime.h"

#include <assert.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

/* Test helpers ----------------------------------------------------- */

static OspreyList *make_range_list(int64_t start, int64_t end) {
  OspreyListBuilder *b = osprey_list_builder_new();
  for (int64_t i = start; i < end; i++) {
    osprey_list_builder_push(b, i);
  }
  return osprey_list_builder_seal(b);
}

static void assert_list_equals_range(OspreyList *l, int64_t start, int64_t end) {
  assert(osprey_list_length(l) == end - start);
  for (int64_t i = 0; i < end - start; i++) {
    int64_t expected = start + i;
    int64_t got = osprey_list_get(l, i);
    if (got != expected) {
      fprintf(stderr, "expected %lld at index %lld, got %lld\n",
              (long long)expected, (long long)i, (long long)got);
      assert(0);
    }
  }
}

/* Tests ------------------------------------------------------------ */

static void test_empty(void) {
  OspreyList *e = osprey_list_empty();
  assert(osprey_list_length(e) == 0);
  assert(osprey_list_in_bounds(e, 0) == 0);
  assert(osprey_list_in_bounds(e, -1) == 0);
  /* empty() is idempotent (returns same singleton or equivalent). */
  assert(osprey_list_length(osprey_list_empty()) == 0);
  printf("  list_empty OK\n");
}

static void test_append_small(void) {
  OspreyList *l = osprey_list_empty();
  for (int64_t i = 0; i < 10; i++) {
    l = osprey_list_append(l, i * 10);
  }
  assert(osprey_list_length(l) == 10);
  for (int64_t i = 0; i < 10; i++) {
    assert(osprey_list_get(l, i) == i * 10);
  }
  printf("  list_append_small OK\n");
}

static void test_persistence_chain(void) {
  /* 5 generations: each `v[k]` extends `v[k-1]` with k*1000. Verify every
     generation keeps its independent contents. */
  OspreyList *v[5];
  v[0] = osprey_list_empty();
  for (int g = 1; g < 5; g++) {
    v[g] = osprey_list_append(v[g - 1], (int64_t)g * 1000);
  }
  for (int g = 0; g < 5; g++) {
    assert(osprey_list_length(v[g]) == (int64_t)g);
    for (int64_t i = 0; i < (int64_t)g; i++) {
      assert(osprey_list_get(v[g], i) == (int64_t)(i + 1) * 1000);
    }
  }
  printf("  list_persistence_chain (5 gens) OK\n");
}

static void test_branching_persistence(void) {
  /* Build two divergent branches from a common ancestor; both must remain
     valid simultaneously. */
  OspreyList *a = make_range_list(0, 5);
  OspreyList *branch_x = osprey_list_append(a, 99);
  OspreyList *branch_y = osprey_list_append(a, 77);
  OspreyList *branch_y2 = osprey_list_append(branch_y, 88);

  assert(osprey_list_length(a) == 5);
  assert(osprey_list_length(branch_x) == 6);
  assert(osprey_list_length(branch_y) == 6);
  assert(osprey_list_length(branch_y2) == 7);

  assert(osprey_list_get(branch_x, 5) == 99);
  assert(osprey_list_get(branch_y, 5) == 77);
  assert(osprey_list_get(branch_y2, 5) == 77);
  assert(osprey_list_get(branch_y2, 6) == 88);

  /* a unchanged */
  for (int64_t i = 0; i < 5; i++) {
    assert(osprey_list_get(a, i) == i);
  }
  printf("  list_branching_persistence OK\n");
}

static void test_tree_growth_boundaries(void) {
  /* Verify operations at the exact boundaries where the trie depth grows:
     n = 32 (tail full), n = 33 (first push to tree), n = 1024 (level 1
     boundary), n = 1025 (level 2 begins). */
  int64_t boundaries[] = {31, 32, 33, 1023, 1024, 1025, 2000};
  for (size_t bi = 0; bi < sizeof(boundaries) / sizeof(boundaries[0]); bi++) {
    int64_t n = boundaries[bi];
    OspreyList *l = make_range_list(0, n);
    assert(osprey_list_length(l) == n);
    /* spot-check several indices: 0, mid, last */
    assert(osprey_list_get(l, 0) == 0);
    if (n > 1) {
      assert(osprey_list_get(l, n / 2) == n / 2);
      assert(osprey_list_get(l, n - 1) == n - 1);
    }
  }
  printf("  list_tree_growth_boundaries OK\n");
}

static void test_set_deep(void) {
  /* Set element midway through a deep-tree list, verify prior version
     unaffected and new version has only that one slot changed. */
  OspreyList *base = make_range_list(0, 2000);
  OspreyList *mutated = osprey_list_set(base, 500, -1);
  assert(osprey_list_get(base, 500) == 500);
  assert(osprey_list_get(mutated, 500) == -1);
  /* surroundings unchanged */
  assert(osprey_list_get(mutated, 499) == 499);
  assert(osprey_list_get(mutated, 501) == 501);
  /* far-from-modification index untouched */
  assert(osprey_list_get(mutated, 1999) == 1999);
  assert(osprey_list_get(mutated, 0) == 0);
  printf("  list_set_deep OK\n");
}

static void test_set_at_extremes(void) {
  OspreyList *l = make_range_list(10, 50);
  OspreyList *first = osprey_list_set(l, 0, -7);
  OspreyList *last = osprey_list_set(l, osprey_list_length(l) - 1, -9);
  assert(osprey_list_get(first, 0) == -7);
  assert(osprey_list_get(first, 1) == 11);
  assert(osprey_list_get(last, 39) == -9);
  assert(osprey_list_get(last, 38) == 48);
  /* original list unaffected */
  assert(osprey_list_get(l, 0) == 10);
  assert(osprey_list_get(l, 39) == 49);
  printf("  list_set_at_extremes OK\n");
}

static void test_concat_variations(void) {
  OspreyList *e = osprey_list_empty();
  OspreyList *a = make_range_list(0, 5);
  OspreyList *b = make_range_list(5, 10);

  /* empty + empty */
  assert(osprey_list_length(osprey_list_concat(e, e)) == 0);
  /* empty + a */
  OspreyList *ea = osprey_list_concat(e, a);
  assert_list_equals_range(ea, 0, 5);
  /* a + empty */
  OspreyList *ae = osprey_list_concat(a, e);
  assert_list_equals_range(ae, 0, 5);
  /* a + b */
  OspreyList *ab = osprey_list_concat(a, b);
  assert_list_equals_range(ab, 0, 10);
  /* (a + b) + a */
  OspreyList *aba = osprey_list_concat(ab, a);
  assert(osprey_list_length(aba) == 15);
  for (int64_t i = 0; i < 10; i++) assert(osprey_list_get(aba, i) == i);
  for (int64_t i = 0; i < 5; i++)  assert(osprey_list_get(aba, 10 + i) == i);
  printf("  list_concat_variations OK\n");
}

static void test_concat_crosses_tail_boundary(void) {
  /* Both sides have non-full tails. After concat, tail of left becomes
     part of the tree of result; verify all elements present. */
  OspreyList *left = make_range_list(0, 50);   /* tail_count = 18 */
  OspreyList *right = make_range_list(50, 73); /* tail_count = 23 */
  OspreyList *r = osprey_list_concat(left, right);
  assert_list_equals_range(r, 0, 73);
  printf("  list_concat_crosses_tail_boundary OK\n");
}

static void test_prepend(void) {
  OspreyList *l = make_range_list(1, 5);
  OspreyList *p = osprey_list_prepend(l, 0);
  assert_list_equals_range(p, 0, 5);
  /* prepend to empty */
  OspreyList *just_one = osprey_list_prepend(osprey_list_empty(), 42);
  assert(osprey_list_length(just_one) == 1);
  assert(osprey_list_get(just_one, 0) == 42);
  /* original list unaffected */
  assert(osprey_list_length(l) == 4);
  assert(osprey_list_get(l, 0) == 1);
  printf("  list_prepend OK\n");
}

static void test_drop(void) {
  OspreyList *l = make_range_list(0, 10);
  /* drop 0 = identity */
  assert(osprey_list_length(osprey_list_drop(l, 0)) == 10);
  /* drop 5 */
  OspreyList *d5 = osprey_list_drop(l, 5);
  assert_list_equals_range(d5, 5, 10);
  /* drop length = empty */
  assert(osprey_list_length(osprey_list_drop(l, 10)) == 0);
  /* drop > length = empty */
  assert(osprey_list_length(osprey_list_drop(l, 100)) == 0);
  /* drop negative = identity */
  assert(osprey_list_length(osprey_list_drop(l, -1)) == 10);
  printf("  list_drop OK\n");
}

static void test_reverse(void) {
  /* reverse(reverse(l)) == l for any l */
  OspreyList *l = make_range_list(1, 50);
  OspreyList *r = osprey_list_reverse(l);
  OspreyList *rr = osprey_list_reverse(r);
  assert(osprey_list_length(r) == osprey_list_length(l));
  for (int64_t i = 0; i < 49; i++) {
    assert(osprey_list_get(r, i) == 49 - i);
    assert(osprey_list_get(rr, i) == osprey_list_get(l, i));
  }
  /* reverse(empty) == empty */
  assert(osprey_list_length(osprey_list_reverse(osprey_list_empty())) == 0);
  /* reverse of single-element list = itself */
  OspreyList *one = osprey_list_append(osprey_list_empty(), 99);
  OspreyList *one_r = osprey_list_reverse(one);
  assert(osprey_list_length(one_r) == 1);
  assert(osprey_list_get(one_r, 0) == 99);
  printf("  list_reverse (involution) OK\n");
}

static void test_builder_matches_incremental(void) {
  /* Builder and incremental append must produce equal lists. */
  OspreyList *incr = osprey_list_empty();
  for (int64_t i = 0; i < 300; i++) incr = osprey_list_append(incr, i * 7);
  OspreyList *via = make_range_list(0, 300);
  /* via has 0..299, multiply by 7 → use direct builder for fairness */
  OspreyListBuilder *b = osprey_list_builder_new();
  for (int64_t i = 0; i < 300; i++) osprey_list_builder_push(b, i * 7);
  OspreyList *via_b = osprey_list_builder_seal(b);
  assert(osprey_list_length(incr) == 300);
  assert(osprey_list_length(via_b) == 300);
  for (int64_t i = 0; i < 300; i++) {
    assert(osprey_list_get(incr, i) == i * 7);
    assert(osprey_list_get(via_b, i) == i * 7);
  }
  /* via without multiplication is also correct */
  (void)via;
  printf("  list_builder_matches_incremental OK\n");
}

static void test_iter_full_coverage(void) {
  OspreyList *l = make_range_list(0, 500);
  OspreyListIter *it = osprey_list_iter_new(l);
  int64_t expected = 0;
  int64_t v = 0;
  int64_t sum = 0;
  while (osprey_list_iter_next(it, &v)) {
    assert(v == expected);
    sum += v;
    expected++;
  }
  assert(expected == 500);
  assert(sum == 500 * 499 / 2); /* Gauss */
  free(it);

  /* Iter over empty produces no values. */
  OspreyListIter *it_e = osprey_list_iter_new(osprey_list_empty());
  assert(osprey_list_iter_next(it_e, &v) == 0);
  free(it_e);
  printf("  list_iter_full_coverage (sum check + empty) OK\n");
}

static void test_stress_10k(void) {
  /* 10k elements forces level-2 internal nodes (> 32*32 = 1024). */
  OspreyList *l = osprey_list_empty();
  for (int64_t i = 0; i < 10000; i++) l = osprey_list_append(l, i);
  assert(osprey_list_length(l) == 10000);
  /* Random-access spot checks */
  assert(osprey_list_get(l, 0) == 0);
  assert(osprey_list_get(l, 33) == 33);
  assert(osprey_list_get(l, 1024) == 1024);
  assert(osprey_list_get(l, 5000) == 5000);
  assert(osprey_list_get(l, 9999) == 9999);
  printf("  list_stress_10k OK\n");
}

static void test_get_after_drop_persistence(void) {
  /* drop returns a new list; verify both source and dropped survive
     subsequent operations. */
  OspreyList *l = make_range_list(0, 100);
  OspreyList *d = osprey_list_drop(l, 50);
  OspreyList *l2 = osprey_list_append(l, -1);
  OspreyList *d2 = osprey_list_append(d, -2);
  assert(osprey_list_get(l, 50) == 50);
  assert(osprey_list_get(d, 0) == 50);
  assert(osprey_list_get(l2, 100) == -1);
  assert(osprey_list_get(d2, 50) == -2);
  /* originals untouched */
  assert(osprey_list_length(l) == 100);
  assert(osprey_list_length(d) == 50);
  printf("  list_drop+persistence OK\n");
}

void run_list_tests(void) {
  printf("List:\n");
  test_empty();
  test_append_small();
  test_persistence_chain();
  test_branching_persistence();
  test_tree_growth_boundaries();
  test_set_deep();
  test_set_at_extremes();
  test_concat_variations();
  test_concat_crosses_tail_boundary();
  test_prepend();
  test_drop();
  test_reverse();
  test_builder_matches_incremental();
  test_iter_full_coverage();
  test_stress_10k();
  test_get_after_drop_persistence();
}
