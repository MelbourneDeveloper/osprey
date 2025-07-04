#ifndef UNITY_H
#define UNITY_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// Unity types
typedef void (*UnityTestFunction)(void);

// Unity globals
extern int Unity_TestsRun;
extern int Unity_TestsFailed;
extern const char* Unity_TestFile;
extern int Unity_TestLine;

// Core Unity macros
#define TEST_ASSERT(condition) \
    do { \
        Unity_TestsRun++; \
        if (!(condition)) { \
            printf("FAIL: %s:%d - %s\n", Unity_TestFile, Unity_TestLine, #condition); \
            Unity_TestsFailed++; \
        } else { \
            printf("PASS: %s\n", #condition); \
        } \
    } while(0)

#define TEST_ASSERT_EQUAL(expected, actual) \
    do { \
        Unity_TestsRun++; \
        if ((expected) != (actual)) { \
            printf("FAIL: %s:%d - Expected %lld but was %lld\n", Unity_TestFile, Unity_TestLine, (long long)(expected), (long long)(actual)); \
            Unity_TestsFailed++; \
        } else { \
            printf("PASS: Expected %lld equals actual %lld\n", (long long)(expected), (long long)(actual)); \
        } \
    } while(0)

#define TEST_ASSERT_TRUE(condition) TEST_ASSERT(condition)
#define TEST_ASSERT_FALSE(condition) TEST_ASSERT(!(condition))

#define TEST_ASSERT_NOT_NULL(ptr) \
    do { \
        Unity_TestsRun++; \
        if ((ptr) == NULL) { \
            printf("FAIL: %s:%d - Expected non-null pointer\n", Unity_TestFile, Unity_TestLine); \
            Unity_TestsFailed++; \
        } else { \
            printf("PASS: Pointer is not null\n"); \
        } \
    } while(0)

// Test runner macros
#define RUN_TEST(test_func) \
    do { \
        Unity_TestFile = __FILE__; \
        Unity_TestLine = __LINE__; \
        printf("\n=== Running %s ===\n", #test_func); \
        test_func(); \
    } while(0)

#define UNITY_BEGIN() \
    do { \
        Unity_TestsRun = 0; \
        Unity_TestsFailed = 0; \
        printf("Unity Test Framework Starting...\n"); \
    } while(0)

#define UNITY_END() \
    do { \
        printf("\n=== TEST RESULTS ===\n"); \
        printf("Tests Run: %d\n", Unity_TestsRun); \
        printf("Tests Failed: %d\n", Unity_TestsFailed); \
        printf("Tests Passed: %d\n", Unity_TestsRun - Unity_TestsFailed); \
        if (Unity_TestsFailed == 0) { \
            printf("🎉 ALL TESTS PASSED! 🎉\n"); \
        } else { \
            printf("💥 %d TESTS FAILED! 💥\n", Unity_TestsFailed); \
        } \
        return Unity_TestsFailed; \
    } while(0)

// Function declarations
void UnityPrint(const char* string);
void UnityPrintNumber(long number);

#endif // UNITY_H
