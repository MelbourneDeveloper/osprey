// Native unit tests for the mock C library [SYMEX-MOCK-STREAM]
// [SYMEX-MOCK-STATELESS] [SYMEX-MOCK-UNSUPPORTED]. The doubles are compiled
// against their own headers and linked here beside the host libc; this file
// is compiled against the host headers and reports through `dprintf`, a
// host function the mock does not define (its own `write` goes to the sink).
// A shim stands in for KLEE: symbolic
// bytes come from a queue the test fills, so each test is one concrete path.
//
// The runtime units under test (`osp_input`, `readFile`) are linked too, so
// the suite proves the runtime reaches the double, not only that the double
// works [SYMEX-MOCK-CLOSED].
#include "hunt_mock.h"

#include <setjmp.h>
#include <stdarg.h>
#include <stddef.h>
#include <stdint.h>
#include <string.h>

// ---- host-side reporting -------------------------------------------------

int dprintf(int fd, const char *format, ...);

static int failures;
static int checks;

static void say(const char *text) { (void)dprintf(2, "%s", text); }

static void check(int condition, const char *what) {
  checks++;
  if (!condition) {
    failures++;
    say("FAIL: ");
    say(what);
    say("\n");
  }
}

// ---- KLEE shim --------------------------------------------------------------

static unsigned char queue[OSP_HUNT_STREAM_MAX];
static size_t queue_head;
static size_t queue_tail;
static int range_values[16];
static size_t range_head;
static size_t range_tail;
static char last_symbolic[OSP_HUNT_NAME_MAX];
static char last_error[OSP_HUNT_NAME_MAX];
static jmp_buf on_error;

static void queue_bytes(const char *bytes, size_t count) {
  memcpy(queue + queue_tail, bytes, count);
  queue_tail += count;
}

static void queue_range(int value) { range_values[range_tail++] = value; }

void klee_make_symbolic(void *address, size_t count, const char *name) {
  strncpy(last_symbolic, name, sizeof(last_symbolic) - 1);
  for (size_t i = 0; i < count; i++) {
    ((unsigned char *)address)[i] = queue_head < queue_tail ? queue[queue_head++] : 0;
  }
}

int klee_range(int begin, int end, const char *name) {
  (void)name;
  int value = range_head < range_tail ? range_values[range_head++] : begin;
  return value < begin ? begin : value >= end ? end - 1 : value;
}

void klee_assume(uintptr_t condition) { (void)condition; }

void klee_report_error(const char *file, int line, const char *message, const char *suffix) {
  (void)file;
  (void)line;
  (void)suffix;
  strncpy(last_error, message, sizeof(last_error) - 1);
  longjmp(on_error, 1);
}

void __assert_fail(const char *expression, const char *file, unsigned line, const char *function) {
  (void)expression;
  (void)file;
  (void)line;
  (void)function;
  longjmp(on_error, 2);
}

static void reset(void) {
  queue_head = queue_tail = 0;
  range_head = range_tail = 0;
  last_symbolic[0] = '\0';
  last_error[0] = '\0';
  osp_hunt_reset();
}

// ---- the doubles under test, through their libc names ----------------------

typedef struct osp_hunt_file FILE_;
FILE_ *fopen(const char *path, const char *mode);
int fclose(FILE_ *stream);
size_t fread(void *out, size_t size, size_t count, FILE_ *stream);
char *fgets(char *out, int cap, FILE_ *stream);
int fseeko(FILE_ *stream, long offset, int whence);
long ftello(FILE_ *stream);
int fprintf(FILE_ *stream, const char *format, ...);
int snprintf(char *out, size_t cap, const char *format, ...);
int printf(const char *format, ...);
int puts(const char *text);
int getchar(void);
void qsort(void *base, size_t count, size_t size, int (*compare)(const void *, const void *));
char *getenv(const char *name);
double strtod(const char *text, char **end);
char *strstr(const char *haystack, const char *needle);
size_t strcspn(const char *text, const char *reject);
void arc4random_buf(void *out, size_t count);
int fork(void);
double sqrt(double x);
double exp(double x);
double log(double x);
double tanh(double x);
double pow(double base, double exponent);
double floor(double x);
double fmod(double x, double y);
extern FILE_ *osp_hunt_stdout;
extern FILE_ *osp_hunt_stdin;
extern int osp_hunt_errno;

// The runtime units that must reach the doubles.
char *osp_input(void);
char *read_file(char *path);

// ---- captured output -------------------------------------------------------

static char captured[4096];
static size_t captured_length;

static void capture(const char *data, size_t count) {
  if (captured_length + count < sizeof(captured)) {
    memcpy(captured + captured_length, data, count);
    captured_length += count;
    captured[captured_length] = '\0';
  }
}

static void clear_capture(void) {
  captured_length = 0;
  captured[0] = '\0';
}

// ---- tests -----------------------------------------------------------------

static void stdin_stream_feeds_the_runtime_line_reader(void) {
  reset();
  osp_hunt_input_bytes = 16;
  queue_bytes("hello\nrest", 10);
  queue_range(10);
  char *line = osp_input();
  check(line != NULL && strcmp(line, "hello") == 0, "osp_input returns the first line");
  check(osp_hunt_calls[OSP_HUNT_CALL_READ] == 6, "osp_input read one byte per call through the double");
  check(strcmp(last_symbolic, "stdin") == 0, "the stdin bytes are the symbolic object named stdin");
  char *second = osp_input();
  check(second != NULL && strcmp(second, "rest") == 0, "the cursor continues after the newline");
  check(getchar() == -1, "end of the symbolic length is EOF");
}

static void stdin_stream_honours_a_short_symbolic_length(void) {
  reset();
  osp_hunt_input_bytes = 8;
  queue_bytes("abcdefgh", 8);
  queue_range(3);
  char *line = osp_input();
  check(line != NULL && strcmp(line, "abc") == 0, "stdin.len cuts the stream short");
}

static void read_file_opens_a_symbolic_stream_named_by_path(void) {
  reset();
  osp_hunt_input_bytes = 8;
  queue_range(1);
  queue_bytes("payload!", 8);
  queue_range(7);
  char path[] = "data.txt";
  char *text = read_file(path);
  check(text != NULL && strcmp(text, "payload") == 0, "readFile returns the symbolic bytes up to the length");
  check(strcmp(last_symbolic, "data.txt") == 0, "the file bytes are named by their path");
  check(osp_hunt_calls[OSP_HUNT_CALL_FOPEN] == 1 && osp_hunt_calls[OSP_HUNT_CALL_FREAD] >= 1, "readFile went through fopen and fread");
  check(osp_hunt_calls[OSP_HUNT_CALL_FCLOSE] == 1, "readFile closed the stream");
}

static void fopen_explores_the_failure_branch(void) {
  reset();
  queue_range(0);
  FILE_ *stream = fopen("missing.txt", "r");
  check(stream == NULL, "a zero on <path>.ok is an open failure");
  check(osp_hunt_errno == 2, "the failure sets ENOENT");
}

static void fseek_and_fgets_walk_the_stream(void) {
  reset();
  osp_hunt_input_bytes = 12;
  queue_range(1);
  queue_bytes("one\ntwo\n3333", 12);
  queue_range(12);
  FILE_ *stream = fopen("lines", "r");
  char line[8];
  check(fgets(line, sizeof(line), stream) != NULL && strcmp(line, "one\n") == 0, "fgets stops at the newline");
  check(fseeko(stream, 0, 2) == 0 && ftello(stream) == 12, "SEEK_END lands on the symbolic length");
  check(fseeko(stream, 4, 0) == 0 && fgets(line, sizeof(line), stream) != NULL && strcmp(line, "two\n") == 0, "SEEK_SET then fgets reads the second line");
  check(fseeko(stream, 99, 0) == -1, "seeking past the length fails");
  fclose(stream);
}

static void written_bytes_reach_the_sink(void) {
  reset();
  clear_capture();
  osp_hunt_sink = capture;
  printf("n=%d s=%s c=%c", 42, "ok", 'x');
  puts("!");
  fprintf(osp_hunt_stdout, "%5.2f|%-4d|%03d|%x|%lld", 3.14159, 7, 5, 255, -9000000000LL);
  check(strcmp(captured, "n=42 s=ok c=x!\n 3.14|7   |005|ff|-9000000000") == 0, captured);
  check(osp_hunt_calls[OSP_HUNT_CALL_PRINTF] == 3, "each print went through the double");
}

static void snprintf_formats_like_the_host(void) {
  char out[64];
  check(snprintf(out, sizeof(out), "%.16e", 1234.5) == 22 && strcmp(out, "1.2345000000000000e+03") == 0, out);
  check(snprintf(out, sizeof(out), "%.17g", 0.1) == 19 && strcmp(out, "0.10000000000000001") == 0, out);
  check(snprintf(out, sizeof(out), "%g", 100000.0) == 6 && strcmp(out, "100000") == 0, out);
  check(snprintf(out, sizeof(out), "%g", 1000000.0) == 5 && strcmp(out, "1e+06") == 0, out);
  check(snprintf(out, sizeof(out), "%g", 0.0001) == 6 && strcmp(out, "0.0001") == 0, out);
  check(snprintf(out, sizeof(out), "%.3f", -2.0005) == 6 && strcmp(out, "-2.001") == 0, out);
  check(snprintf(out, sizeof(out), "%s", "truncated") == 9 && snprintf(out, 4, "%s", "truncated") == 9 && strcmp(out, "tru") == 0, "snprintf reports the full length and clips");
  check(snprintf(out, sizeof(out), "%zu %ld %u %%", (size_t)5, -3L, 9u) == 8 && strcmp(out, "5 -3 9 %") == 0, out);
  check(snprintf(out, sizeof(out), "%f", 1.0 / 0.0) == 3 && strcmp(out, "inf") == 0, out);
}

static void string_and_number_doubles_match_libc(void) {
  char *end;
  check(strtod("  -12.5e2xyz", &end) == -1250.0 && strcmp(end, "xyz") == 0, "strtod parses sign, fraction and exponent");
  check(strtod("nope", &end) == 0.0 && strcmp(end, "nope") == 0, "strtod leaves end on no digits");
  check(strcmp(strstr("haystack", "st"), "stack") == 0 && strstr("hay", "z") == NULL, "strstr finds and misses");
  check(strcspn("key=value", "=") == 3, "strcspn counts to the first reject");
  check(getenv("HOME") == NULL && osp_hunt_calls[OSP_HUNT_CALL_GETENV] == 1, "the environment is empty");
}

static int by_value(const void *a, const void *b) { return *(const int *)a - *(const int *)b; }

static void qsort_sorts(void) {
  int values[] = {5, -1, 3, 3, 0, 9};
  qsort(values, 6, sizeof(int), by_value);
  int sorted = values[0] == -1 && values[1] == 0 && values[2] == 3 && values[3] == 3 && values[4] == 5 && values[5] == 9;
  check(sorted && osp_hunt_calls[OSP_HUNT_CALL_QSORT] == 1, "qsort orders the elements");
}

static void entropy_is_a_named_symbolic_object(void) {
  reset();
  queue_bytes("\x01\x02\x03\x04", 4);
  unsigned char bytes[4];
  arc4random_buf(bytes, 4);
  check(bytes[0] == 1 && bytes[3] == 4, "entropy comes from the symbolic draw");
  check(strncmp(last_symbolic, "entropy.", 8) == 0, last_symbolic);
  check(osp_hunt_calls[OSP_HUNT_CALL_ENTROPY] == 1, "the entropy double counted the draw");
}

static int close_to(double a, double b) { return a - b < 1e-9 && b - a < 1e-9; }

static void math_doubles_are_accurate(void) {
  check(close_to(sqrt(2.0), 1.4142135623730951), "sqrt");
  check(close_to(exp(1.0), 2.718281828459045), "exp");
  check(close_to(log(10.0), 2.302585092994046), "log");
  check(close_to(tanh(0.5), 0.46211715726000974), "tanh");
  check(close_to(pow(2.0, 10.0), 1024.0) && close_to(pow(2.0, 0.5), 1.4142135623730951), "pow");
  check(floor(-1.5) == -2.0 && close_to(fmod(7.5, 2.0), 1.5), "floor and fmod");
  check(sqrt(-1.0) != sqrt(-1.0), "sqrt of a negative is nan");
}

static void unsupported_calls_stop_the_path(void) {
  reset();
  if (setjmp(on_error) == 0) {
    (void)fork();
    check(0, "fork must not return");
  }
  check(strcmp(last_error, "unmodelled: fork") == 0, last_error);
  check(osp_hunt_calls[OSP_HUNT_CALL_UNSUPPORTED] == 1, "the stub counted the call");
}

int main(void) {
  stdin_stream_feeds_the_runtime_line_reader();
  stdin_stream_honours_a_short_symbolic_length();
  read_file_opens_a_symbolic_stream_named_by_path();
  fopen_explores_the_failure_branch();
  fseek_and_fgets_walk_the_stream();
  written_bytes_reach_the_sink();
  snprintf_formats_like_the_host();
  string_and_number_doubles_match_libc();
  qsort_sorts();
  entropy_is_a_named_symbolic_object();
  math_doubles_are_accurate();
  unsupported_calls_stop_the_path();
  char summary[64];
  int length = snprintf(summary, sizeof(summary), "hunt_mock_tests: %d checks, %d failures\n", checks, failures);
  (void)length;
  say(summary);
  return failures == 0 ? 0 : 1;
}
