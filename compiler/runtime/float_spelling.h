#ifndef OSPREY_FLOAT_SPELLING_H
#define OSPREY_FLOAT_SPELLING_H

#include <stddef.h>

/* The %g spelling toString(float) gives `d`: the smallest precision from 10
 * to 17 whose spelling reads back as `d` (NaN and infinities at 17). Returns
 * the length written into `buf`, or a negative value / one >= `cap` when it
 * did not fit, exactly as snprintf does. Allocates nothing. */
int osp_float_spelling(double d, char *buf, size_t cap);

#endif /* OSPREY_FLOAT_SPELLING_H */
