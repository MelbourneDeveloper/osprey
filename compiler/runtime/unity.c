#include "unity.h"

// Unity global variables
int Unity_TestsRun = 0;
int Unity_TestsFailed = 0;
const char* Unity_TestFile = NULL;
int Unity_TestLine = 0;

// Unity helper functions
void UnityPrint(const char* string) {
    printf("%s", string);
}

void UnityPrintNumber(long number) {
    printf("%ld", number);
}
