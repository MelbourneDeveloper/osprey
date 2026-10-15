// The device math library compiled as host C and held to the host's own IEEE
// binary64 results [GPU-BACKEND-PORTABILITY] [BUILTIN-MATH]: the four
// operations, sqrt, the conversions and the comparisons must agree bit for
// bit on every word tried (NaN payloads compare as NaN); exp, ln and tanh are
// compared with the host's long double libm rounded to binary64. Where long
// double is wider (x87 on x86-64, binary128 on aarch64 Linux) that is the
// correctly rounded value bar a rare double rounding; on Darwin it is the
// binary64 libm itself. glibc's binary64 tanh is up to two ulps off (tanh of
// 0xbfd8b6b048dce522: 1.69 ulps, where the device is 0.31 against an
// 80-digit reference), so a binary64 reference would test glibc, not the
// device. The device must land within one ulp everywhere and agree on at
// least 99% of arguments.
#include <assert.h>
#include <math.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "gpu_f64.metal"
#include "gpu_math64.metal"
#include "gpu_device_lib.metal"

#define TRIALS 400000
#define SPECIALS 12

static uint64_t state = 0x9e3779b97f4a7c15;
static unsigned long checks;

static uint64_t random_word(void) {
  state ^= state << 13;
  state ^= state >> 7;
  state ^= state << 17;
  return state;
}

static double of(uint64_t word) {
  double d;
  memcpy(&d, &word, sizeof d);
  return d;
}

static uint64_t bits(double d) {
  uint64_t word;
  memcpy(&word, &d, sizeof word);
  return word;
}

static int same(uint64_t a, uint64_t b) { return a == b || (isnan(of(a)) && isnan(of(b))); }
static uint64_t word_of(long v) { return (uint64_t)v; }
static long truncated(double d) { return (long)d; }
static uint64_t flag(bool b) { return b ? 1 : 0; }

// A word with an exponent drawn from [low, low + range), any sign and mantissa.
static uint64_t word_in(int low, int range) {
  uint64_t r = random_word();
  uint64_t e = (uint64_t)low + ((r >> 52) % (uint64_t)range);
  return (r & UINT64_C(0x800fffffffffffff)) | (e << 52);
}

static const uint64_t specials[SPECIALS] = {0, UINT64_C(0x8000000000000000), 1, UINT64_C(0x8000000000000001),
    UINT64_C(0x000fffffffffffff), UINT64_C(0x0010000000000000), UINT64_C(0x3ff0000000000000),
    UINT64_C(0xbff0000000000000), UINT64_C(0x7fefffffffffffff), UINT64_C(0x7ff0000000000000),
    UINT64_C(0xfff0000000000000), UINT64_C(0x7ff8000000000000)};

static void expect(const char *what, uint64_t a, uint64_t b, uint64_t device, uint64_t host) {
  checks++;
  if (!same(device, host)) {
    fprintf(stderr, "%s(%016llx, %016llx): device %016llx host %016llx\n", what,
        (unsigned long long)a, (unsigned long long)b, (unsigned long long)device, (unsigned long long)host);
    abort();
  }
}

typedef ulong (*device_binary)(ulong, ulong);
typedef double (*host_binary)(double, double);
static double host_add(double a, double b) { return a + b; }
static double host_sub(double a, double b) { return a - b; }
static double host_mul(double a, double b) { return a * b; }
static double host_div(double a, double b) { return a / b; }

static void binary_case(const char *what, device_binary device, host_binary host, uint64_t a, uint64_t b) {
  expect(what, a, b, device(a, b), bits(host(of(a), of(b))));
}

static void check_binary(const char *what, device_binary device, host_binary host) {
  for (unsigned i = 0; i < SPECIALS; i++)
    for (unsigned j = 0; j < SPECIALS; j++) binary_case(what, device, host, specials[i], specials[j]);
  for (unsigned i = 0; i < TRIALS; i++) {
    binary_case(what, device, host, word_in(0, 2047), word_in(0, 2047));
    binary_case(what, device, host, word_in(1000, 50), word_in(1000, 50));
    binary_case(what, device, host, word_in(0, 60), word_in(0, 60));
    binary_case(what, device, host, word_in(2000, 46), word_in(2000, 46));
  }
}

static void check_unary_and_conversions(void) {
  for (unsigned i = 0; i < SPECIALS; i++) expect("sqrt", specials[i], 0, sqrt64(specials[i]), bits(sqrt(of(specials[i]))));
  for (unsigned i = 0; i < TRIALS; i++) {
    uint64_t x = word_in(0, 2047), y = word_in(0, 2047);
    expect("sqrt", x, 0, sqrt64(x), bits(sqrt(of(x))));
    long v = (long)random_word() >> (random_word() % 64);
    expect("from_long", (uint64_t)v, 0, from_long(v), bits((double)v));
    uint64_t small = word_in(900, 185);
    if (!isnan(of(small)) && fabs(of(small)) < 9.2e18) expect("to_long", small, 0, word_of(to_long(small)), word_of(truncated(of(small))));
    expect("neg", x, 0, neg64(x), bits(-of(x)));
    expect("abs", x, 0, abs64(x), bits(fabs(of(x))));
    expect("lt", x, y, flag(lt64(x, y)), flag(of(x) < of(y)));
    expect("le", x, y, flag(le64(x, y)), flag(of(x) <= of(y)));
    expect("eq", x, y, flag(eq64(x, y)), flag(of(x) == of(y)));
  }
  expect("to_long saturates", 0, 0, word_of(to_long(UINT64_C(0x43f0000000000000))), word_of(INT64_MAX));
  expect("to_long saturates", 0, 0, word_of(to_long(UINT64_C(0xc3f0000000000000))), word_of(INT64_MIN));
  expect("to_long nan", 0, 0, word_of(to_long(UINT64_C(0x7ff8000000000000))), 0);
}

typedef ulong (*device_unary)(ulong);
typedef long double (*host_unary)(long double);

static uint64_t reference(host_unary host, uint64_t x) {
  double rounded = host(of(x));
  return bits(rounded);
}

// Distance in ulps between two finite words of the same sign.
static uint64_t ulps(uint64_t a, uint64_t b) { return a > b ? a - b : b - a; }

static void check_transcendental(const char *what, device_unary device, host_unary host, int low, int range) {
  unsigned long exact = 0, trials = 0;
  uint64_t worst = 0;
  for (unsigned i = 0; i < SPECIALS; i++) expect(what, specials[i], 0, device(specials[i]), reference(host, specials[i]));
  for (unsigned i = 0; i < TRIALS * 2; i++) {
    uint64_t x = word_in(low, range);
    if (what[0] == 'l') x &= ~UINT64_C(0x8000000000000000);
    uint64_t got = device(x), want = reference(host, x);
    trials++;
    checks++;
    if (same(got, want)) { exact++; continue; }
    uint64_t distance = ulps(got, want);
    if (distance > worst) worst = distance;
    if (distance > 1) {
      fprintf(stderr, "%s(%016llx): device %016llx host %016llx\n", what, (unsigned long long)x,
          (unsigned long long)got, (unsigned long long)want);
      abort();
    }
  }
  printf("[ok] %s: %lu of %lu match the host long double libm bit for bit, worst distance %llu ulp\n", what, exact, trials,
      (unsigned long long)worst);
  assert(trials - exact <= trials / 100);
}

static void check_checked_integers(void) {
  ulong fault = 0;
  assert(osprey_add(2, 3, &fault) == 5 && osprey_sub(2, 3, &fault) == -1 && osprey_mul(-4, 5, &fault) == -20);
  assert(osprey_div(7, -2, &fault) == -3 && osprey_rem(-7, 2, &fault) == -1 && osprey_neg(5, &fault) == -5);
  assert(osprey_abs(-9, &fault) == 9 && osprey_rem(INT64_MIN, -1, &fault) == 0 && !fault);
  assert(osprey_fdiv(from_long(1), from_long(4), &fault) == bits(0.25) && !fault);
  assert(osprey_add(INT64_MAX, 1, &fault) == 0 && fault == FAULT_OVERFLOW);
  fault = 0; assert(osprey_sub(INT64_MIN, 1, &fault) == 0 && fault == FAULT_OVERFLOW);
  fault = 0; assert(osprey_mul(INT64_MIN, -1, &fault) == 0 && fault == FAULT_OVERFLOW);
  fault = 0; assert(osprey_mul(-1, INT64_MIN, &fault) == 0 && fault == FAULT_OVERFLOW);
  fault = 0; assert(osprey_mul(INT64_C(3037000500), INT64_C(3037000500), &fault) == 0 && fault == FAULT_OVERFLOW);
  fault = 0; assert(osprey_div(INT64_MIN, -1, &fault) == 0 && fault == FAULT_OVERFLOW);
  fault = 0; assert(osprey_neg(INT64_MIN, &fault) == 0 && fault == FAULT_OVERFLOW);
  fault = 0; assert(osprey_abs(INT64_MIN, &fault) == 0 && fault == FAULT_OVERFLOW);
  fault = 0; assert(osprey_div(1, 0, &fault) == 0 && fault == FAULT_REMAINDER_BY_ZERO);
  fault = 0; assert(osprey_rem(1, 0, &fault) == 0 && fault == FAULT_REMAINDER_BY_ZERO);
  fault = 0; assert(osprey_fdiv(from_long(1), 0, &fault) == 0 && fault == FAULT_DIVIDE_BY_ZERO);
  checks += 17;
}

static void check_matmul_kernel(void) {
  ulong left[4] = {bits(0.5), bits(-1.0), bits(2.0), bits(3.0)}, out[4] = {0, 0, 0, 0};
  for (ulong i = 0; i < 4; i++) osprey_matmul64(left, left, out, 2, 2, 2, i);
  osprey_matmul64(left, left, out, 2, 2, 2, 4);
  assert(out[0] == bits(-1.75) && out[1] == bits(-3.5) && out[2] == bits(7.0) && out[3] == bits(7.0));
  checks += 4;
}

int main(void) {
  check_binary("add", add64, host_add);
  check_binary("sub", sub64, host_sub);
  check_binary("mul", multiply64, host_mul);
  check_binary("div", divide64, host_div);
  check_unary_and_conversions();
  check_checked_integers();
  check_matmul_kernel();
  check_transcendental("exp", exp64, expl, 960, 73);
  check_transcendental("ln", ln64, logl, 0, 2047);
  check_transcendental("tanh", tanh64, tanhl, 990, 38);
  printf("[ok] device binary64 library: %lu host parity checks\n", checks);
  return 0;
}
