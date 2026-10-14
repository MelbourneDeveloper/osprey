// Implements [GPU-MATMUL] and [GPU-BACKEND-PORTABILITY]: binary64 arithmetic
// on devices without native double, with round-to-nearest/even and no FMA.
#if defined(__METAL_VERSION__)
#include <metal_stdlib>
using namespace metal;
#define GPU_GLOBAL device
#define GPU_BUFFER(N) [[buffer(N)]]
#define GPU_PARAM constant ulong &
#define GPU_THREAD , uint index [[thread_position_in_grid]]
#define GPU_INDEX (ulong)(index)
#else
#define GPU_GLOBAL global
#define GPU_BUFFER(N)
#define GPU_PARAM ulong
#define GPU_THREAD
#define GPU_INDEX get_global_id(0)
#endif
constant ulong FRAC = 0x000ffffffffffffful;
constant ulong SIGN = 0x8000000000000000ul;
constant ulong INF = 0x7ff0000000000000ul;
constant ulong QUIET = 0x0008000000000000ul;
constant ulong ONE = 0x0010000000000000ul;

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

kernel void osprey_matmul64(GPU_GLOBAL const ulong *left GPU_BUFFER(0),
                           GPU_GLOBAL const ulong *right GPU_BUFFER(1),
                           GPU_GLOBAL ulong *out GPU_BUFFER(2),
                           GPU_PARAM rows GPU_BUFFER(3),
                           GPU_PARAM inner GPU_BUFFER(4),
                           GPU_PARAM cols GPU_BUFFER(5) GPU_THREAD) {
    ulong i = GPU_INDEX;
    if (i >= rows * cols) return;
    ulong row = i / cols, col = i % cols;
    ulong total = 0;
    for (ulong k = 0; k < inner; k++) {
        ulong product = multiply64(left[row * inner + k], right[k * cols + col]);
        total = add64(total, product);
    }
    out[i] = total;
}
