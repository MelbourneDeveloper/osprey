/*
 * Implements [BUILTIN-STRING-*]
 * Scalar (string -> string / bool / int) helpers exposed to Osprey IR.
 * List-returning helpers live in string_runtime_list.c.
 *
 * Conventions: NUL-terminated UTF-8 byte sequences; outputs are malloc'd
 * and owned by the caller. All functions defend against NULL.
 */

#include <ctype.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#include "string_runtime.h"

char *osp_string_dup_internal(const char *s, size_t n) {
    char *out = (char *)malloc(n + 1);
    if (!out) return NULL;
    if (n > 0) memcpy(out, s, n);
    out[n] = '\0';
    return out;
}

char *osp_string_empty_internal(void) { return osp_string_dup_internal("", 0); }

int osp_is_ws_internal(unsigned char c) {
    return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f';
}

/* ---------- inspection / search (total) ---------- */

int64_t osp_string_is_empty(const char *s) {
    if (!s) return 1;
    return s[0] == '\0' ? 1 : 0;
}

int64_t osp_string_starts_with(const char *s, const char *prefix) {
    if (!s || !prefix) return 0;
    size_t plen = strlen(prefix);
    if (plen == 0) return 1;
    if (strlen(s) < plen) return 0;
    return strncmp(s, prefix, plen) == 0 ? 1 : 0;
}

int64_t osp_string_ends_with(const char *s, const char *suffix) {
    if (!s || !suffix) return 0;
    size_t slen = strlen(s);
    size_t suflen = strlen(suffix);
    if (suflen == 0) return 1;
    if (slen < suflen) return 0;
    return memcmp(s + slen - suflen, suffix, suflen) == 0 ? 1 : 0;
}

/* Returns first byte-index of needle in s, or -1 if absent. Empty needle = 0. */
int64_t osp_string_index_of(const char *s, const char *needle) {
    if (!s || !needle) return -1;
    if (needle[0] == '\0') return 0;
    const char *hit = strstr(s, needle);
    if (!hit) return -1;
    return (int64_t)(hit - s);
}

/* ---------- substrings (total) ---------- */

char *osp_string_take(const char *s, int64_t n) {
    if (!s) return osp_string_empty_internal();
    size_t len = strlen(s);
    if (n <= 0) return osp_string_empty_internal();
    if ((size_t)n >= len) return osp_string_dup_internal(s, len);
    return osp_string_dup_internal(s, (size_t)n);
}

char *osp_string_drop(const char *s, int64_t n) {
    if (!s) return osp_string_empty_internal();
    size_t len = strlen(s);
    if (n <= 0) return osp_string_dup_internal(s, len);
    if ((size_t)n >= len) return osp_string_empty_internal();
    return osp_string_dup_internal(s + n, len - (size_t)n);
}

/* substring: returns NULL on out-of-range or inverted indices.
 * Caller emits IndexOutOfRange when NULL. */
char *osp_string_substring(const char *s, int64_t start, int64_t end) {
    if (!s) return NULL;
    size_t len = strlen(s);
    if (start < 0 || end < start || (size_t)end > len) return NULL;
    return osp_string_dup_internal(s + start, (size_t)(end - start));
}

/* ---------- transformation (total) ---------- */

char *osp_string_to_upper(const char *s) {
    if (!s) return osp_string_empty_internal();
    size_t len = strlen(s);
    char *out = osp_string_dup_internal(s, len);
    if (!out) return NULL;
    for (size_t i = 0; i < len; i++)
        out[i] = (char)toupper((unsigned char)out[i]);
    return out;
}

char *osp_string_to_lower(const char *s) {
    if (!s) return osp_string_empty_internal();
    size_t len = strlen(s);
    char *out = osp_string_dup_internal(s, len);
    if (!out) return NULL;
    for (size_t i = 0; i < len; i++)
        out[i] = (char)tolower((unsigned char)out[i]);
    return out;
}

char *osp_string_trim_start(const char *s) {
    if (!s) return osp_string_empty_internal();
    while (*s && osp_is_ws_internal((unsigned char)*s)) s++;
    return osp_string_dup_internal(s, strlen(s));
}

char *osp_string_trim_end(const char *s) {
    if (!s) return osp_string_empty_internal();
    size_t len = strlen(s);
    while (len > 0 && osp_is_ws_internal((unsigned char)s[len - 1])) len--;
    return osp_string_dup_internal(s, len);
}

char *osp_string_trim(const char *s) {
    if (!s) return osp_string_empty_internal();
    while (*s && osp_is_ws_internal((unsigned char)*s)) s++;
    size_t len = strlen(s);
    while (len > 0 && osp_is_ws_internal((unsigned char)s[len - 1])) len--;
    return osp_string_dup_internal(s, len);
}

char *osp_string_reverse(const char *s) {
    if (!s) return osp_string_empty_internal();
    size_t len = strlen(s);
    char *out = (char *)malloc(len + 1);
    if (!out) return NULL;
    for (size_t i = 0; i < len; i++) out[i] = s[len - 1 - i];
    out[len] = '\0';
    return out;
}

/* ---------- transformation (fallible) ----------
 * NULL return = caller should emit InvalidArgument. */

char *osp_string_replace(const char *s, const char *needle, const char *replacement) {
    if (!s || !needle || !replacement || needle[0] == '\0') return NULL;
    size_t slen = strlen(s);
    size_t nlen = strlen(needle);
    size_t rlen = strlen(replacement);

    size_t count = 0;
    for (const char *p = s; (p = strstr(p, needle)) != NULL; p += nlen) count++;
    if (count == 0) return osp_string_dup_internal(s, slen);

    size_t out_len = slen + count * rlen - count * nlen;
    char *out = (char *)malloc(out_len + 1);
    if (!out) return NULL;

    char *w = out;
    const char *r = s;
    while (1) {
        const char *hit = strstr(r, needle);
        if (!hit) {
            size_t tail = strlen(r);
            memcpy(w, r, tail);
            w += tail;
            break;
        }
        size_t pre = (size_t)(hit - r);
        memcpy(w, r, pre);
        w += pre;
        memcpy(w, replacement, rlen);
        w += rlen;
        r = hit + nlen;
    }
    *w = '\0';
    return out;
}

char *osp_string_repeat(const char *s, int64_t n) {
    if (!s || n < 0) return NULL;
    if (n == 0) return osp_string_empty_internal();
    size_t len = strlen(s);
    if (len == 0) return osp_string_empty_internal();
    if ((size_t)n > (SIZE_MAX - 1) / len) return NULL;
    size_t out_len = len * (size_t)n;
    char *out = (char *)malloc(out_len + 1);
    if (!out) return NULL;
    for (int64_t i = 0; i < n; i++) memcpy(out + (size_t)i * len, s, len);
    out[out_len] = '\0';
    return out;
}

char *osp_string_pad_start(const char *s, int64_t target_length, const char *fill) {
    if (!s || !fill || fill[0] == '\0') return NULL;
    size_t slen = strlen(s);
    if (target_length <= 0 || (size_t)target_length <= slen)
        return osp_string_dup_internal(s, slen);
    size_t pad_needed = (size_t)target_length - slen;
    size_t flen = strlen(fill);
    char *out = (char *)malloc((size_t)target_length + 1);
    if (!out) return NULL;
    for (size_t i = 0; i < pad_needed; i++) out[i] = fill[i % flen];
    memcpy(out + pad_needed, s, slen);
    out[(size_t)target_length] = '\0';
    return out;
}

char *osp_string_pad_end(const char *s, int64_t target_length, const char *fill) {
    if (!s || !fill || fill[0] == '\0') return NULL;
    size_t slen = strlen(s);
    if (target_length <= 0 || (size_t)target_length <= slen)
        return osp_string_dup_internal(s, slen);
    size_t pad_needed = (size_t)target_length - slen;
    size_t flen = strlen(fill);
    char *out = (char *)malloc((size_t)target_length + 1);
    if (!out) return NULL;
    memcpy(out, s, slen);
    for (size_t i = 0; i < pad_needed; i++) out[slen + i] = fill[i % flen];
    out[(size_t)target_length] = '\0';
    return out;
}

/* ---------- parsing ---------- */

/* Returns 0 on success, 1 on failure. Strict: no whitespace, optional sign. */
int64_t osp_parse_int_strict(const char *s, int64_t *out) {
    if (!s || s[0] == '\0' || !out) return 1;
    const char *p = s;
    int neg = 0;
    if (*p == '-' || *p == '+') {
        neg = (*p == '-');
        p++;
        if (*p == '\0') return 1;
    }
    int64_t acc = 0;
    while (*p) {
        if (*p < '0' || *p > '9') return 1;
        int d = *p - '0';
        if (acc > 922337203685477580LL ||
            (acc == 922337203685477580LL && d > (neg ? 8 : 7)))
            return 1;
        acc = acc * 10 + d;
        p++;
    }
    *out = neg ? -acc : acc;
    return 0;
}

int64_t osp_parse_float_strict(const char *s, double *out) {
    if (!s || s[0] == '\0' || !out) return 1;
    char *endp = NULL;
    double v = strtod(s, &endp);
    if (!endp || *endp != '\0' || endp == s) return 1;
    *out = v;
    return 0;
}
