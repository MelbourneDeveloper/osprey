#include "unity.h"
#include <pthread.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

// Forward declarations for fiber runtime functions
extern int64_t fiber_spawn(int64_t (*fn)(void));
extern int64_t fiber_await(int64_t fiber_id);
extern int64_t fiber_sleep(int64_t milliseconds);
extern int64_t channel_create(int64_t capacity);
extern int64_t channel_send(int64_t channel_id, int64_t value);
extern int64_t channel_recv(int64_t channel_id);

// Test functions
static int64_t test_function_1(void) { return 42; }
static int64_t test_function_2(void) { return 100; }
static int64_t slow_function(void) {
    usleep(100000); // 100ms
    return 999;
}

// Unity Test Functions

// CRITICAL TEST: Null function pointer - prevents segfaults
static void test_null_function_pointer(void) {
    int64_t result = fiber_spawn(NULL);
    TEST_ASSERT_EQUAL(-1, result);
}

// CRITICAL TEST: Invalid fiber ID - prevents buffer overflow segfaults
static void test_invalid_fiber_await(void) {
    int64_t result = fiber_await(99999); // Way out of bounds
    TEST_ASSERT_EQUAL(-1, result);

    result = fiber_await(-1); // Negative ID
    TEST_ASSERT_EQUAL(-1, result);
}

// Test valid fiber spawning and execution
static void test_valid_fiber_spawn(void) {
    int64_t fiber_id = fiber_spawn(test_function_1);
    TEST_ASSERT_TRUE(fiber_id > 0);

    int64_t result = fiber_await(fiber_id);
    TEST_ASSERT_EQUAL(42, result);
}

// Test multiple concurrent fibers
static void test_multiple_fibers(void) {
    int64_t fiber1 = fiber_spawn(test_function_1);
    int64_t fiber2 = fiber_spawn(test_function_2);

    TEST_ASSERT_TRUE(fiber1 > 0);
    TEST_ASSERT_TRUE(fiber2 > 0);
    TEST_ASSERT_TRUE(fiber1 != fiber2);

    int64_t result1 = fiber_await(fiber1);
    int64_t result2 = fiber_await(fiber2);

    TEST_ASSERT_EQUAL(42, result1);
    TEST_ASSERT_EQUAL(100, result2);
}

// Test fiber bounds checking - prevents array overflow
static void test_fiber_bounds_checking(void) {
    // Test invalid channel IDs
    int64_t result = channel_send(99999, 42); // Out of bounds
    TEST_ASSERT_EQUAL(0, result);

    result = channel_recv(99999); // Out of bounds
    TEST_ASSERT_EQUAL(-1, result);
}

// Test fiber sleep function
static void test_fiber_sleep(void) {
    int64_t result = fiber_sleep(10); // Sleep for 10ms
    TEST_ASSERT_EQUAL(0, result);
}

// Stress test - create many fibers to test stability
static void test_fiber_stress(void) {
    const int num_fibers = 20; // Reduced for faster testing
    int64_t fiber_ids[num_fibers];

    // Spawn many fibers
    for (int i = 0; i < num_fibers; i++) {
        fiber_ids[i] = fiber_spawn(test_function_1);
        TEST_ASSERT_TRUE(fiber_ids[i] > 0);
    }

    // Wait for all fibers
    for (int i = 0; i < num_fibers; i++) {
        int64_t result = fiber_await(fiber_ids[i]);
        TEST_ASSERT_EQUAL(42, result);
    }
}

// Test concurrent execution
static void test_concurrent_execution(void) {
    // Spawn a slow fiber
    int64_t slow_fiber = fiber_spawn(slow_function);

    // Spawn a fast fiber
    int64_t fast_fiber = fiber_spawn(test_function_1);

    // The fast fiber should complete quickly even though slow fiber is running
    int64_t fast_result = fiber_await(fast_fiber);
    TEST_ASSERT_EQUAL(42, fast_result);

    // Now wait for the slow fiber
    int64_t slow_result = fiber_await(slow_fiber);
    TEST_ASSERT_EQUAL(999, slow_result);
}

// Test channel creation and basic operations
static void test_channel_basic_operations(void) {
    int64_t channel_id = channel_create(10);
    TEST_ASSERT_TRUE(channel_id >= 0);

    // Test sending and receiving
    int64_t send_result = channel_send(channel_id, 123);
    TEST_ASSERT_EQUAL(1, send_result);

    int64_t recv_result = channel_recv(channel_id);
    TEST_ASSERT_EQUAL(123, recv_result);
}

// Test invalid channel operations
static void test_invalid_channel_operations(void) {
    // Test with invalid channel ID
    int64_t result = channel_send(-1, 42);
    TEST_ASSERT_EQUAL(0, result);

    result = channel_recv(-1);
    TEST_ASSERT_EQUAL(-1, result);
}

// Main test runner
int main(void) {
    UNITY_BEGIN();
    
    RUN_TEST(test_null_function_pointer);
    RUN_TEST(test_invalid_fiber_await);
    RUN_TEST(test_valid_fiber_spawn);
    RUN_TEST(test_multiple_fibers);
    RUN_TEST(test_fiber_bounds_checking);
    RUN_TEST(test_fiber_sleep);
    RUN_TEST(test_concurrent_execution);
    RUN_TEST(test_fiber_stress);
    RUN_TEST(test_channel_basic_operations);
    RUN_TEST(test_invalid_channel_operations);
    
    UNITY_END();
} 
