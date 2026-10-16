// Device support for the kernels the compiler emits [GPU-BACKEND-DEVICE] and
// the builtin product [GPU-MATMUL]. A kernel cannot reach a runtime `Arith`
// policy, so a fallible integer operation or a division by zero records a
// fault word instead; the runtime aborts the program naming the kernel when a
// dispatch completes with a fault ([GPU-DEVICE-FAULT]). Follows gpu_f64.metal
// and gpu_math64.metal in the embedded shader source.

constant ulong FAULT_OVERFLOW = 1;
constant ulong FAULT_DIVIDE_BY_ZERO = 2;
constant ulong FAULT_REMAINDER_BY_ZERO = 3;

long osprey_fault(GPU_GLOBAL ulong *fault, ulong code) {
    fault[0] = code;
    return 0;
}

long osprey_add(long a, long b, GPU_GLOBAL ulong *fault) {
    ulong r = (ulong)a + (ulong)b;
    bool overflow = (((~((ulong)a ^ (ulong)b)) & ((ulong)a ^ r)) >> 63) != 0;
    return overflow ? osprey_fault(fault, FAULT_OVERFLOW) : (long)r;
}

long osprey_sub(long a, long b, GPU_GLOBAL ulong *fault) {
    ulong r = (ulong)a - (ulong)b;
    bool overflow = ((((ulong)a ^ (ulong)b) & ((ulong)a ^ r)) >> 63) != 0;
    return overflow ? osprey_fault(fault, FAULT_OVERFLOW) : (long)r;
}

// Magnitudes as unsigned words: the device backends have no signed 64-bit
// division, so every checked product, quotient and remainder works on these.
ulong magnitude(long a) { return a < 0 ? 0ul - (ulong)a : (ulong)a; }
long signed_word(ulong magnitude_word, bool negative) {
    return negative ? (long)(0ul - magnitude_word) : (long)magnitude_word;
}

long osprey_mul(long a, long b, GPU_GLOBAL ulong *fault) {
    bool negative = (a < 0) != (b < 0);
    Wide p = wide_product(magnitude(a), magnitude(b));
    ulong limit = negative ? SIGN : SIGN - 1;
    if (p.high != 0 || p.low > limit) return osprey_fault(fault, FAULT_OVERFLOW);
    return signed_word(p.low, negative);
}

// intDiv: truncating division [BUILTIN-INTDIV].
long osprey_div(long a, long b, GPU_GLOBAL ulong *fault) {
    if (b == 0) return osprey_fault(fault, FAULT_REMAINDER_BY_ZERO);
    if (a == LONG_MIN_WORD && b == -1) return osprey_fault(fault, FAULT_OVERFLOW);
    return signed_word(magnitude(a) / magnitude(b), (a < 0) != (b < 0));
}

long osprey_rem(long a, long b, GPU_GLOBAL ulong *fault) {
    if (b == 0) return osprey_fault(fault, FAULT_REMAINDER_BY_ZERO);
    return signed_word(magnitude(a) % magnitude(b), a < 0);
}

long osprey_neg(long a, GPU_GLOBAL ulong *fault) {
    return a == LONG_MIN_WORD ? osprey_fault(fault, FAULT_OVERFLOW) : -a;
}

long osprey_abs(long a, GPU_GLOBAL ulong *fault) { return a < 0 ? osprey_neg(a, fault) : a; }

ulong osprey_fdiv(ulong a, ulong b, GPU_GLOBAL ulong *fault) {
    if (!is_zero(b)) return divide64(a, b);
    long zero = osprey_fault(fault, FAULT_DIVIDE_BY_ZERO);
    return (ulong)zero;
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
