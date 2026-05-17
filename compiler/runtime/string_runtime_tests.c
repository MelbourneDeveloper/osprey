/*
 * Implements [BUILTIN-STRING-*] verification.
 *
 * Strict assertion-driven tests for every helper in string_runtime.c and
 * string_runtime_list.c. Each test exercises both the happy path AND
 * every documented error/edge case. A failure aborts (assert) — the test
 * binary's exit status is the verdict.
 *
 * Wired into compiler/Makefile under the `c-test` target.
 */

#include <assert.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "string_runtime.h"

/* ---------- scalar predicates ---------- */

static void test_is_empty(void) {
    assert(osp_string_is_empty("") == 1);
    assert(osp_string_is_empty("a") == 0);
    assert(osp_string_is_empty("hello world") == 0);
    assert(osp_string_is_empty(NULL) == 1); /* NULL defended */
    printf("  ok  is_empty\n");
}

static void test_starts_with(void) {
    assert(osp_string_starts_with("hello world", "hello") == 1);
    assert(osp_string_starts_with("hello world", "world") == 0);
    assert(osp_string_starts_with("hello world", "") == 1);
    assert(osp_string_starts_with("", "") == 1);
    assert(osp_string_starts_with("", "x") == 0);
    assert(osp_string_starts_with("hi", "hello") == 0); /* prefix longer than s */
    assert(osp_string_starts_with("GET /api", "GET ") == 1);
    assert(osp_string_starts_with(NULL, "x") == 0);
    assert(osp_string_starts_with("x", NULL) == 0);
    printf("  ok  starts_with\n");
}

static void test_ends_with(void) {
    assert(osp_string_ends_with("hello world", "world") == 1);
    assert(osp_string_ends_with("hello world", "hello") == 0);
    assert(osp_string_ends_with("hello world", "") == 1);
    assert(osp_string_ends_with("", "") == 1);
    assert(osp_string_ends_with("", "x") == 0);
    assert(osp_string_ends_with("hi", "hello") == 0); /* suffix longer than s */
    assert(osp_string_ends_with("image.png", ".png") == 1);
    assert(osp_string_ends_with("image.PNG", ".png") == 0); /* case-sensitive */
    printf("  ok  ends_with\n");
}

static void test_index_of(void) {
    assert(osp_string_index_of("hello", "ell") == 1);
    assert(osp_string_index_of("hello", "h") == 0);
    assert(osp_string_index_of("hello", "o") == 4);
    assert(osp_string_index_of("hello", "xyz") == -1);     /* not found */
    assert(osp_string_index_of("hello", "") == 0);          /* empty needle */
    assert(osp_string_index_of("foo=bar=baz", "=") == 3);   /* first occurrence */
    assert(osp_string_index_of(NULL, "x") == -1);
    assert(osp_string_index_of("x", NULL) == -1);
    printf("  ok  index_of\n");
}

/* ---------- substring helpers ---------- */

static void test_take(void) {
    char *out;
    out = osp_string_take("hello", 3);  assert(strcmp(out, "hel") == 0);   free(out);
    out = osp_string_take("hello", 0);  assert(strcmp(out, "")    == 0);   free(out);
    out = osp_string_take("hello", -5); assert(strcmp(out, "")    == 0);   free(out);
    out = osp_string_take("hello", 5);  assert(strcmp(out, "hello") == 0); free(out);
    out = osp_string_take("hello", 99); assert(strcmp(out, "hello") == 0); free(out); /* clamp */
    out = osp_string_take("", 3);       assert(strcmp(out, "")    == 0);   free(out);
    out = osp_string_take(NULL, 3);     assert(strcmp(out, "")    == 0);   free(out);
    printf("  ok  take\n");
}

static void test_drop(void) {
    char *out;
    out = osp_string_drop("hello", 3);  assert(strcmp(out, "lo")   == 0); free(out);
    out = osp_string_drop("hello", 0);  assert(strcmp(out, "hello") == 0); free(out);
    out = osp_string_drop("hello", -5); assert(strcmp(out, "hello") == 0); free(out);
    out = osp_string_drop("hello", 5);  assert(strcmp(out, "")    == 0);  free(out); /* exact */
    out = osp_string_drop("hello", 99); assert(strcmp(out, "")    == 0);  free(out); /* clamp */
    out = osp_string_drop("", 3);       assert(strcmp(out, "")    == 0);  free(out);
    out = osp_string_drop(NULL, 3);     assert(strcmp(out, "")    == 0);  free(out);
    printf("  ok  drop\n");
}

static void test_substring(void) {
    char *out;
    out = osp_string_substring("hello", 1, 4); assert(strcmp(out, "ell") == 0); free(out);
    out = osp_string_substring("hello", 0, 5); assert(strcmp(out, "hello") == 0); free(out);
    out = osp_string_substring("hello", 0, 0); assert(strcmp(out, "") == 0); free(out);
    out = osp_string_substring("hello", 5, 5); assert(strcmp(out, "") == 0); free(out);
    /* error cases: must return NULL */
    assert(osp_string_substring("hello", -1, 3) == NULL); /* start < 0 */
    assert(osp_string_substring("hello", 0, 99) == NULL); /* end > len */
    assert(osp_string_substring("hello", 4, 2)  == NULL); /* end < start */
    assert(osp_string_substring(NULL, 0, 1)     == NULL);
    printf("  ok  substring\n");
}

/* ---------- transformation ---------- */

static void test_to_upper_lower(void) {
    char *out;
    out = osp_string_to_upper("hello"); assert(strcmp(out, "HELLO") == 0); free(out);
    out = osp_string_to_upper("");      assert(strcmp(out, "")      == 0); free(out);
    out = osp_string_to_upper("AbC");   assert(strcmp(out, "ABC")   == 0); free(out);
    out = osp_string_to_upper("1!@a");  assert(strcmp(out, "1!@A")  == 0); free(out);
    out = osp_string_to_upper(NULL);    assert(strcmp(out, "")      == 0); free(out);

    out = osp_string_to_lower("HELLO"); assert(strcmp(out, "hello") == 0); free(out);
    out = osp_string_to_lower("");      assert(strcmp(out, "")      == 0); free(out);
    out = osp_string_to_lower("AbC");   assert(strcmp(out, "abc")   == 0); free(out);
    printf("  ok  to_upper/to_lower\n");
}

static void test_trim(void) {
    char *out;
    out = osp_string_trim("  hello  ");      assert(strcmp(out, "hello") == 0); free(out);
    out = osp_string_trim("\t\n hi \r\n");   assert(strcmp(out, "hi")    == 0); free(out);
    out = osp_string_trim("hello");          assert(strcmp(out, "hello") == 0); free(out);
    out = osp_string_trim("");               assert(strcmp(out, "")      == 0); free(out);
    out = osp_string_trim("   ");            assert(strcmp(out, "")      == 0); free(out); /* all ws */
    out = osp_string_trim(NULL);             assert(strcmp(out, "")      == 0); free(out);

    out = osp_string_trim_start("  hello  "); assert(strcmp(out, "hello  ") == 0); free(out);
    out = osp_string_trim_start("hello");     assert(strcmp(out, "hello")   == 0); free(out);
    out = osp_string_trim_start("   ");       assert(strcmp(out, "")        == 0); free(out);

    out = osp_string_trim_end("  hello  ");   assert(strcmp(out, "  hello") == 0); free(out);
    out = osp_string_trim_end("hello");       assert(strcmp(out, "hello")   == 0); free(out);
    out = osp_string_trim_end("   ");         assert(strcmp(out, "")        == 0); free(out);
    printf("  ok  trim/trim_start/trim_end\n");
}

static void test_reverse(void) {
    char *out;
    out = osp_string_reverse("abc");   assert(strcmp(out, "cba")    == 0); free(out);
    out = osp_string_reverse("a");     assert(strcmp(out, "a")      == 0); free(out);
    out = osp_string_reverse("");      assert(strcmp(out, "")       == 0); free(out);
    out = osp_string_reverse("12345"); assert(strcmp(out, "54321")  == 0); free(out);
    out = osp_string_reverse(NULL);    assert(strcmp(out, "")       == 0); free(out);
    printf("  ok  reverse\n");
}

static void test_replace(void) {
    char *out;
    out = osp_string_replace("a-b-c", "-", "_");        assert(strcmp(out, "a_b_c") == 0); free(out);
    out = osp_string_replace("aaa", "a", "bb");         assert(strcmp(out, "bbbbbb") == 0); free(out);
    out = osp_string_replace("hello", "xyz", "Q");      assert(strcmp(out, "hello") == 0); free(out); /* no match */
    out = osp_string_replace("hello", "l", "");         assert(strcmp(out, "heo") == 0); free(out); /* shrink */
    out = osp_string_replace("", "x", "y");             assert(strcmp(out, "") == 0); free(out);
    /* error cases */
    assert(osp_string_replace("hello", "", "x")    == NULL); /* empty needle */
    assert(osp_string_replace(NULL,    "x", "y")   == NULL);
    assert(osp_string_replace("h",     NULL, "y")  == NULL);
    assert(osp_string_replace("h",     "x", NULL)  == NULL);
    printf("  ok  replace\n");
}

static void test_repeat(void) {
    char *out;
    out = osp_string_repeat("ab", 3);   assert(strcmp(out, "ababab") == 0); free(out);
    out = osp_string_repeat("x", 5);    assert(strcmp(out, "xxxxx") == 0); free(out);
    out = osp_string_repeat("ab", 0);   assert(strcmp(out, "") == 0); free(out); /* n=0 -> empty */
    out = osp_string_repeat("ab", 1);   assert(strcmp(out, "ab") == 0); free(out);
    out = osp_string_repeat("", 99);    assert(strcmp(out, "") == 0); free(out); /* empty source */
    /* error cases */
    assert(osp_string_repeat("ab", -1) == NULL);
    assert(osp_string_repeat(NULL,  3) == NULL);
    printf("  ok  repeat\n");
}

static void test_pad(void) {
    char *out;
    out = osp_string_pad_start("7",   3, "0");  assert(strcmp(out, "007") == 0); free(out);
    out = osp_string_pad_start("42",  5, "ab"); assert(strcmp(out, "aba42") == 0); free(out);
    out = osp_string_pad_start("hi",  2, "x");  assert(strcmp(out, "hi") == 0); free(out); /* no pad needed */
    out = osp_string_pad_start("hi",  1, "x");  assert(strcmp(out, "hi") == 0); free(out); /* target < len */
    out = osp_string_pad_start("",    3, "0");  assert(strcmp(out, "000") == 0); free(out);

    out = osp_string_pad_end  ("7",   3, ".");  assert(strcmp(out, "7..") == 0); free(out);
    out = osp_string_pad_end  ("42",  5, "ab"); assert(strcmp(out, "42aba") == 0); free(out);
    out = osp_string_pad_end  ("hi",  2, "x");  assert(strcmp(out, "hi") == 0); free(out);

    /* error: empty fill */
    assert(osp_string_pad_start("hi", 5, "") == NULL);
    assert(osp_string_pad_end  ("hi", 5, "") == NULL);
    assert(osp_string_pad_start("hi", 5, NULL) == NULL);
    printf("  ok  pad_start/pad_end\n");
}

/* ---------- parsing ---------- */

static void test_parse_int(void) {
    int64_t out;
    assert(osp_parse_int_strict("42", &out) == 0);     assert(out == 42);
    assert(osp_parse_int_strict("-42", &out) == 0);    assert(out == -42);
    assert(osp_parse_int_strict("+42", &out) == 0);    assert(out == 42);
    assert(osp_parse_int_strict("0", &out) == 0);      assert(out == 0);
    assert(osp_parse_int_strict("9223372036854775807", &out) == 0); /* INT64_MAX */
    assert(out == 9223372036854775807LL);
    assert(osp_parse_int_strict("-9223372036854775808", &out) == 0); /* INT64_MIN */
    assert(out == (-9223372036854775807LL - 1));

    /* rejections */
    assert(osp_parse_int_strict("",        &out) != 0);
    assert(osp_parse_int_strict("abc",     &out) != 0);
    assert(osp_parse_int_strict("12abc",   &out) != 0);
    assert(osp_parse_int_strict("abc12",   &out) != 0);
    assert(osp_parse_int_strict(" 42",     &out) != 0); /* leading space */
    assert(osp_parse_int_strict("42 ",     &out) != 0); /* trailing space */
    assert(osp_parse_int_strict("-",       &out) != 0); /* sign with no digits */
    assert(osp_parse_int_strict("+",       &out) != 0);
    assert(osp_parse_int_strict("9223372036854775808",  &out) != 0); /* overflow */
    assert(osp_parse_int_strict("-9223372036854775809", &out) != 0); /* underflow */
    assert(osp_parse_int_strict(NULL,      &out) != 0);
    printf("  ok  parse_int_strict\n");
}

static void test_parse_float(void) {
    double out;
    assert(osp_parse_float_strict("3.14", &out) == 0);   assert(out > 3.13 && out < 3.15);
    assert(osp_parse_float_strict("0", &out) == 0);      assert(out == 0.0);
    assert(osp_parse_float_strict("-2.5", &out) == 0);   assert(out > -2.51 && out < -2.49);
    assert(osp_parse_float_strict("1e3", &out) == 0);    assert(out > 999.9 && out < 1000.1);

    /* rejections */
    assert(osp_parse_float_strict("",      &out) != 0);
    assert(osp_parse_float_strict("abc",   &out) != 0);
    assert(osp_parse_float_strict("3.14x", &out) != 0); /* trailing junk */
    assert(osp_parse_float_strict(NULL,    &out) != 0);
    printf("  ok  parse_float_strict\n");
}

/* ---------- list-returning ---------- */

static void test_split(void) {
    osp_string_list *list;

    list = osp_string_split("a,b,c", ",");
    assert(list != NULL); assert(list->length == 3);
    assert(strcmp(list->items[0], "a") == 0);
    assert(strcmp(list->items[1], "b") == 0);
    assert(strcmp(list->items[2], "c") == 0);
    osp_string_list_free(list);

    list = osp_string_split("hello", ",");
    assert(list != NULL); assert(list->length == 1);
    assert(strcmp(list->items[0], "hello") == 0);
    osp_string_list_free(list);

    list = osp_string_split(",,", ",");
    assert(list != NULL); assert(list->length == 3); /* "", "", "" */
    assert(strcmp(list->items[0], "") == 0);
    assert(strcmp(list->items[1], "") == 0);
    assert(strcmp(list->items[2], "") == 0);
    osp_string_list_free(list);

    list = osp_string_split("foo::bar::baz", "::");
    assert(list != NULL); assert(list->length == 3);
    assert(strcmp(list->items[0], "foo") == 0);
    assert(strcmp(list->items[1], "bar") == 0);
    assert(strcmp(list->items[2], "baz") == 0);
    osp_string_list_free(list);

    /* error: empty separator */
    assert(osp_string_split("hello", "") == NULL);
    assert(osp_string_split(NULL,    ",") == NULL);
    printf("  ok  split\n");
}

static void test_lines(void) {
    osp_string_list *list;

    list = osp_string_lines("a\nb\nc");
    assert(list != NULL); assert(list->length == 3);
    assert(strcmp(list->items[0], "a") == 0);
    assert(strcmp(list->items[1], "b") == 0);
    assert(strcmp(list->items[2], "c") == 0);
    osp_string_list_free(list);

    /* trailing newline does NOT produce an empty final entry */
    list = osp_string_lines("a\nb\n");
    assert(list != NULL); assert(list->length == 2);
    assert(strcmp(list->items[0], "a") == 0);
    assert(strcmp(list->items[1], "b") == 0);
    osp_string_list_free(list);

    list = osp_string_lines("");
    assert(list != NULL); assert(list->length == 0);
    osp_string_list_free(list);

    list = osp_string_lines("single");
    assert(list != NULL); assert(list->length == 1);
    assert(strcmp(list->items[0], "single") == 0);
    osp_string_list_free(list);

    list = osp_string_lines(NULL);
    assert(list != NULL); assert(list->length == 0);
    osp_string_list_free(list);
    printf("  ok  lines\n");
}

static void test_words(void) {
    osp_string_list *list;

    list = osp_string_words("a b c");
    assert(list != NULL); assert(list->length == 3);
    assert(strcmp(list->items[0], "a") == 0);
    assert(strcmp(list->items[1], "b") == 0);
    assert(strcmp(list->items[2], "c") == 0);
    osp_string_list_free(list);

    /* runs of whitespace collapse; empties dropped */
    list = osp_string_words("  hello\t\tworld  \n  goodbye  ");
    assert(list != NULL); assert(list->length == 3);
    assert(strcmp(list->items[0], "hello") == 0);
    assert(strcmp(list->items[1], "world") == 0);
    assert(strcmp(list->items[2], "goodbye") == 0);
    osp_string_list_free(list);

    list = osp_string_words("");
    assert(list != NULL); assert(list->length == 0);
    osp_string_list_free(list);

    list = osp_string_words("   \t\n\r   ");
    assert(list != NULL); assert(list->length == 0);
    osp_string_list_free(list);

    list = osp_string_words(NULL);
    assert(list != NULL); assert(list->length == 0);
    osp_string_list_free(list);
    printf("  ok  words\n");
}

static void test_join(void) {
    /* Build a list manually and join it. */
    osp_string_list *list = osp_string_split("a,b,c", ",");
    char *out = osp_string_join(list, "-");
    assert(strcmp(out, "a-b-c") == 0);
    free(out);
    osp_string_list_free(list);

    /* join round-trips split */
    list = osp_string_split("foo::bar::baz", "::");
    out = osp_string_join(list, "::");
    assert(strcmp(out, "foo::bar::baz") == 0);
    free(out);
    osp_string_list_free(list);

    /* empty list → "" */
    list = osp_string_split("", ",");
    out = osp_string_join(list, "x");
    assert(strcmp(out, "") == 0);
    free(out);
    osp_string_list_free(list);

    printf("  ok  join\n");
}

/* ---------- entry point ---------- */

int main(void) {
    printf("🧪 string_runtime_tests\n");
    test_is_empty();
    test_starts_with();
    test_ends_with();
    test_index_of();
    test_take();
    test_drop();
    test_substring();
    test_to_upper_lower();
    test_trim();
    test_reverse();
    test_replace();
    test_repeat();
    test_pad();
    test_parse_int();
    test_parse_float();
    test_split();
    test_lines();
    test_words();
    test_join();
    printf("✅ all string_runtime tests passed\n");
    return 0;
}
