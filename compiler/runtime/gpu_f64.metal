// Implements [GPU-MATMUL] and [GPU-BACKEND-PORTABILITY]: binary64 arithmetic
// on devices without native double, with round-to-nearest/even and no FMA.
// One source compiles as Metal, as OpenCL C and as host C (gpu_shader_tests.c
// holds every function to the host's own IEEE results, bit for bit).
#if defined(__METAL_VERSION__)
#include <metal_stdlib>
using namespace metal;
#define GPU_GLOBAL device
#define GPU_BUFFER(N) [[buffer(N)]]
#define GPU_PARAM constant ulong &
#define GPU_THREAD , uint index [[thread_position_in_grid]]
#define GPU_INDEX (ulong)(index)
#elif defined(__OPENCL_VERSION__) || defined(__OPENCL_C_VERSION__)
#define GPU_GLOBAL global
#define GPU_BUFFER(N)
#define GPU_PARAM ulong
#define GPU_THREAD
#define GPU_INDEX get_global_id(0)
#else
#include <stdbool.h>
#include <stdint.h>
typedef uint64_t ulong;
typedef unsigned int uint;
#define GPU_GLOBAL
#define GPU_BUFFER(N)
#define GPU_PARAM ulong
#define GPU_THREAD , ulong thread_index
#define GPU_INDEX thread_index
#define constant static const
#define kernel
static uint clz(ulong x) { return x ? (uint)__builtin_clzll(x) : 64u; }
#endif
constant ulong FRAC = 0x000ffffffffffffful;
constant ulong SIGN = 0x8000000000000000ul;
constant ulong INF = 0x7ff0000000000000ul;
constant ulong QUIET = 0x0008000000000000ul;
constant ulong ONE = 0x0010000000000000ul;
constant ulong F64_ONE = 0x3ff0000000000000ul;
constant long LONG_MIN_WORD = (long)0x8000000000000000ul;

ulong jam(ulong x, uint shift) {
    if (!shift) return x;
    if (shift >= 64) return (ulong)(x != 0);
    return (x >> shift) | (ulong)((x << (64 - shift)) != 0);
}

bool is_nan(ulong x) { return (x & INF) == INF && (x & FRAC) != 0; }
bool is_inf(ulong x) { return (x & ~SIGN) == INF; }
bool is_zero(ulong x) { return (x & ~SIGN) == 0; }
ulong nan_result(ulong a, ulong b) {
    if (is_nan(a)) return a | QUIET;
    if (is_nan(b)) return b | QUIET;
    return INF | QUIET;
}

typedef struct { ulong mantissa; int exponent; } Parts;
Parts parts(ulong x) {
    uint e = (uint)((x >> 52) & 2047);
    ulong m = (x & FRAC) | (e ? ONE : 0ul);
    int power = e ? (int)(e) - 1023 : -1022;
    if (m && m < ONE) {
        uint shift = (uint)(clz(m)) - 11;
        m <<= shift;
        power -= (int)(shift);
    }
    return (Parts){m, power};
}

// `extended` carries 53 significant bits above three rounding bits, the top
// bit standing for 2^exponent; the result is rounded to nearest, ties to even.
ulong rounded(ulong sign, int exponent, ulong extended) {
    if (!extended) return sign;
    if (exponent < -1022) {
        extended = jam(extended, (uint)(-1022 - exponent));
        exponent = -1022;
    }
    ulong bits = extended & 7;
    ulong m = (extended >> 3) + (ulong)(bits > 4 || (bits == 4 && ((extended >> 3) & 1)));
    if (m >= (ONE << 1)) { m >>= 1; exponent++; }
    if (exponent > 1023) return sign | INF;
    ulong e = m >= ONE ? (ulong)(exponent + 1023) : 0ul;
    return sign | (e << 52) | (m & FRAC);
}

typedef struct { ulong high; ulong low; } Wide;
Wide wide_product(ulong a, ulong b) {
    ulong a0 = a & 0xfffffffful, a1 = a >> 32;
    ulong b0 = b & 0xfffffffful, b1 = b >> 32;
    ulong w0 = a0 * b0;
    ulong t = a1 * b0 + (w0 >> 32);
    ulong w1 = (t & 0xfffffffful) + a0 * b1;
    return (Wide){a1 * b1 + (t >> 32) + (w1 >> 32),
            (w1 << 32) | (w0 & 0xfffffffful)};
}

ulong multiply64(ulong a, ulong b) {
    ulong sign = (a ^ b) & SIGN;
    if (is_nan(a) || is_nan(b)) return nan_result(a, b);
    if (is_inf(a) || is_inf(b)) {
        return is_zero(a) || is_zero(b) ? nan_result(0, 0) : sign | INF;
    }
    if (is_zero(a) || is_zero(b)) return sign;
    Parts x = parts(a), y = parts(b);
    Wide p = wide_product(x.mantissa, y.mantissa);
    uint carry = (uint)((p.high >> 41) & 1);
    uint shift = 49 + carry;
    ulong extended = (p.high << (64 - shift)) | (p.low >> shift);
    extended |= (ulong)((p.low << (64 - shift)) != 0);
    return rounded(sign, x.exponent + y.exponent + (int)(carry), extended);
}

ulong add_finite(ulong a, ulong b) {
    Parts x = parts(a), y = parts(b);
    if (x.exponent < y.exponent || (x.exponent == y.exponent && x.mantissa < y.mantissa)) {
        ulong swap = a; a = b; b = swap;
        Parts temp = x; x = y; y = temp;
    }
    ulong large = x.mantissa << 3;
    ulong small = jam(y.mantissa << 3, (uint)(x.exponent - y.exponent));
    bool same = ((a ^ b) & SIGN) == 0;
    ulong value = same ? large + small : large - small;
    if (!value) return same ? (a & SIGN) : 0ul;
    if (value & (ONE << 4)) { value = jam(value, 1); x.exponent++; }
    uint shift = (uint)(clz(value)) - 8;
    return rounded(a & SIGN, x.exponent - (int)(shift), value << shift);
}

ulong add64(ulong a, ulong b) {
    if (is_nan(a) || is_nan(b)) return nan_result(a, b);
    if (is_inf(a) && is_inf(b) && ((a ^ b) & SIGN)) return nan_result(0, 0);
    if (is_inf(a)) return a;
    if (is_inf(b)) return b;
    if (is_zero(a) && is_zero(b)) return (a & b) & SIGN;
    if (is_zero(a)) return b;
    if (is_zero(b)) return a;
    return add_finite(a, b);
}

ulong neg64(ulong a) { return a ^ SIGN; }
ulong sub64(ulong a, ulong b) { return add64(a, b ^ SIGN); }
ulong abs64(ulong a) { return a & ~SIGN; }

// The quotient's 57 leading bits by restoring division plus one sticky bit,
// and the exponent adjustment (-1 when the dividend's mantissa is smaller).
Parts quotient_bits(ulong dividend, ulong divisor) {
    ulong rem = dividend, q = 0;
    for (uint i = 0; i < 57; i++) {
        q <<= 1;
        if (rem >= divisor) { rem -= divisor; q |= 1; }
        rem <<= 1;
    }
    int adjust = q >> 56 ? 0 : -1;
    if (!adjust) q = jam(q, 1);
    return (Parts){q | (ulong)(rem != 0), adjust};
}

ulong divide64(ulong a, ulong b) {
    ulong sign = (a ^ b) & SIGN;
    if (is_nan(a) || is_nan(b)) return nan_result(a, b);
    if (is_inf(a)) return is_inf(b) ? nan_result(0, 0) : sign | INF;
    if (is_inf(b)) return sign;
    if (is_zero(b)) return is_zero(a) ? nan_result(0, 0) : sign | INF;
    if (is_zero(a)) return sign;
    Parts x = parts(a), y = parts(b);
    Parts q = quotient_bits(x.mantissa, y.mantissa);
    return rounded(sign, x.exponent - y.exponent + q.exponent, q.mantissa);
}

// Signed order of a non-NaN word; both zeros map to 0 so they compare equal.
long ordered_key(ulong x) { return (x & SIGN) ? -(long)(x & ~SIGN) : (long)x; }
bool lt64(ulong a, ulong b) { return !is_nan(a) && !is_nan(b) && ordered_key(a) < ordered_key(b); }
bool le64(ulong a, ulong b) { return !is_nan(a) && !is_nan(b) && ordered_key(a) <= ordered_key(b); }
bool eq64(ulong a, ulong b) { return !is_nan(a) && !is_nan(b) && ordered_key(a) == ordered_key(b); }

// Exact int -> float widening with round-to-nearest/even [GPU-CONVERT].
ulong from_long(long v) {
    if (!v) return 0;
    ulong sign = v < 0 ? SIGN : 0ul;
    ulong magnitude = v < 0 ? 0ul - (ulong)v : (ulong)v;
    uint shift = (uint)clz(magnitude);
    return rounded(sign, 63 - (int)shift, jam(magnitude << shift, 8));
}

// Truncation toward zero; out-of-range magnitudes saturate and NaN is 0, the
// arm64 fcvtzs contract the host oracle runs under.
long to_long(ulong x) {
    if (is_nan(x) || is_zero(x)) return 0;
    Parts p = parts(x);
    if (p.exponent >= 63) return (x & SIGN) ? LONG_MIN_WORD : (long)~SIGN;
    int shift = p.exponent - 52;
    ulong magnitude = shift >= 0 ? p.mantissa << shift
                    : (shift <= -64 ? 0ul : p.mantissa >> (uint)(-shift));
    return (x & SIGN) ? -(long)magnitude : (long)magnitude;
}

// Correctly rounded square root: digit-by-digit over the mantissa's 27 bit
// pairs and 29 zero pairs, giving 56 root bits plus the remainder as sticky.
ulong sqrt64(ulong a) {
    if (is_nan(a)) return a | QUIET;
    if (is_zero(a) || a == INF) return a;
    if (a & SIGN) return nan_result(0, 0);
    Parts x = parts(a);
    ulong m = x.mantissa;
    int e = x.exponent;
    if (e & 1) { m <<= 1; e -= 1; }
    ulong root = 0, rem = 0;
    for (int i = 55; i >= 0; i--) {
        ulong pair = i >= 29 ? (m >> (uint)(2 * (i - 29))) & 3 : 0ul;
        rem = (rem << 2) | pair;
        ulong trial = (root << 2) | 1;
        if (rem >= trial) { rem -= trial; root = (root << 1) | 1; } else root <<= 1;
    }
    return rounded(0, e / 2, root | (ulong)(rem != 0));
}
