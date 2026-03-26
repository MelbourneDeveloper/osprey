// effects_runtime.c - Runtime handler stack for algebraic effects
// Implements dynamic handler resolution for nested effect handlers

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <pthread.h>

// Maximum handler stack depth per fiber
#define MAX_HANDLER_STACK_DEPTH 1024
#define MAX_EFFECT_NAME_LENGTH 128
#define MAX_OPERATION_NAME_LENGTH 128

// HandlerEntry represents a single handler on the stack
typedef struct {
    char effect_name[MAX_EFFECT_NAME_LENGTH];
    char operation_name[MAX_OPERATION_NAME_LENGTH];
    void *handler_func_ptr;  // Function pointer to handler
} HandlerEntry;

// HandlerStack per thread/fiber
typedef struct {
    HandlerEntry stack[MAX_HANDLER_STACK_DEPTH];
    int top;  // Index of top element (-1 means empty)
    pthread_mutex_t lock;  // Thread safety
} HandlerStack;

// Global handler stack (thread-local storage would be better for production)
static __thread HandlerStack *g_handler_stack = NULL;

// Initialize handler stack for current thread
static void ensure_handler_stack_initialized(void) {
    if (g_handler_stack == NULL) {
        g_handler_stack = (HandlerStack *)malloc(sizeof(HandlerStack));
        if (g_handler_stack == NULL) {
            fprintf(stderr, "FATAL: Failed to allocate handler stack\n");
            abort();
        }
        g_handler_stack->top = -1;
        pthread_mutex_init(&g_handler_stack->lock, NULL);
    }
}

// Push a handler onto the stack
// Returns 0 on success, -1 on stack overflow
int __osprey_handler_push(const char *effect_name, const char *operation_name, void *handler_func_ptr) {
    ensure_handler_stack_initialized();

    pthread_mutex_lock(&g_handler_stack->lock);

    if (g_handler_stack->top >= MAX_HANDLER_STACK_DEPTH - 1) {
        pthread_mutex_unlock(&g_handler_stack->lock);
        fprintf(stderr, "FATAL: Handler stack overflow (depth > %d)\n", MAX_HANDLER_STACK_DEPTH);
        return -1;
    }

    g_handler_stack->top++;
    HandlerEntry *entry = &g_handler_stack->stack[g_handler_stack->top];

    strncpy(entry->effect_name, effect_name, MAX_EFFECT_NAME_LENGTH - 1);
    entry->effect_name[MAX_EFFECT_NAME_LENGTH - 1] = '\0';

    strncpy(entry->operation_name, operation_name, MAX_OPERATION_NAME_LENGTH - 1);
    entry->operation_name[MAX_OPERATION_NAME_LENGTH - 1] = '\0';

    entry->handler_func_ptr = handler_func_ptr;

    pthread_mutex_unlock(&g_handler_stack->lock);
    return 0;
}

// Pop a handler from the stack
// Returns 0 on success, -1 on stack underflow
int __osprey_handler_pop(void) {
    ensure_handler_stack_initialized();

    pthread_mutex_lock(&g_handler_stack->lock);

    if (g_handler_stack->top < 0) {
        pthread_mutex_unlock(&g_handler_stack->lock);
        fprintf(stderr, "FATAL: Handler stack underflow\n");
        return -1;
    }

    g_handler_stack->top--;

    pthread_mutex_unlock(&g_handler_stack->lock);
    return 0;
}

// Look up handler from stack (searches from top to bottom)
// Returns handler function pointer, or NULL if not found
void *__osprey_handler_lookup(const char *effect_name, const char *operation_name) {
    ensure_handler_stack_initialized();

    pthread_mutex_lock(&g_handler_stack->lock);

    // Search from top of stack (most recent handler) to bottom
    for (int i = g_handler_stack->top; i >= 0; i--) {
        HandlerEntry *entry = &g_handler_stack->stack[i];
        if (strcmp(entry->effect_name, effect_name) == 0 &&
            strcmp(entry->operation_name, operation_name) == 0) {
            void *result = entry->handler_func_ptr;
            pthread_mutex_unlock(&g_handler_stack->lock);
            return result;
        }
    }

    pthread_mutex_unlock(&g_handler_stack->lock);
    return NULL;  // Handler not found
}

// Get current stack depth (for debugging)
int __osprey_handler_stack_depth(void) {
    ensure_handler_stack_initialized();

    pthread_mutex_lock(&g_handler_stack->lock);
    int depth = g_handler_stack->top + 1;
    pthread_mutex_unlock(&g_handler_stack->lock);

    return depth;
}

// Cleanup handler stack (call at thread exit)
void __osprey_handler_stack_cleanup(void) {
    if (g_handler_stack != NULL) {
        pthread_mutex_destroy(&g_handler_stack->lock);
        free(g_handler_stack);
        g_handler_stack = NULL;
    }
}
