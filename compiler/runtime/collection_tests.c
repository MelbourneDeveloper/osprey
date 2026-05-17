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
 *   - boundary cases at 32/33/1024/1025 element counts
 *   - mixed-key types for Map (int, string, bool)
 *
 * Each test function exits non-zero on failure (via assert).
 */

extern void run_list_tests(void);
extern void run_map_tests(void);

int main(void) {
  printf("Running collection_tests…\n");
  run_list_tests();
  run_map_tests();
  printf("All collection tests passed.\n");
  return 0;
}
