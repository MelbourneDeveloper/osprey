package codegen

import (
	"sync"
	"sync/atomic"
)

// FiberRuntime manages all fiber execution with REAL concurrency.
type FiberRuntime struct {
	fibers    map[int64]*Fiber
	channels  map[int64]*Channel
	nextID    int64
	mu        sync.RWMutex
	scheduler *Scheduler
}

// Fiber represents a concurrent execution unit.
type Fiber struct {
	ID         int64
	State      FiberState
	Value      int64
	Function   func() int64
	Done       chan struct{}
	Result     int64
	ModuleData map[string]interface{} // Isolated module state
}

// FiberState represents the execution state of a fiber.
type FiberState int

const (
	// FiberReady indicates the fiber is ready to execute.
	FiberReady FiberState = iota
	// FiberRunning indicates the fiber is currently executing.
	FiberRunning
	// FiberYielded indicates the fiber has yielded control.
	FiberYielded
	// FiberCompleted indicates the fiber has completed execution.
	FiberCompleted
)

// Channel represents a communication channel between fibers.
type Channel struct {
	ID       int64
	Capacity int
	Buffer   chan int64
	Closed   bool
}

// Scheduler manages fiber execution.
type Scheduler struct {
	readyQueue   chan *Fiber
	runningCount int32
	maxWorkers   int
}

const (
	defaultQueueSize = 1000
	defaultWorkers   = 8
	yieldValue       = 42
)

// FiberRuntime holds the global fiber runtime instance.
//
//nolint:gochecknoglobals // Global runtime required for fiber coordination
var fiberRuntime = &FiberRuntime{
	fibers:   make(map[int64]*Fiber),
	channels: make(map[int64]*Channel),
	scheduler: &Scheduler{
		readyQueue: make(chan *Fiber, defaultQueueSize),
		maxWorkers: defaultWorkers,
	},
}

// GetRuntime returns the global fiber runtime instance.
func GetRuntime() *FiberRuntime {
	return fiberRuntime
}

// init starts the fiber scheduler workers.
func init() {
	// Start worker goroutines for fiber execution
	for range defaultWorkers {
		go fiberWorker()
	}
}

// fiberWorker executes fibers from the ready queue.
func fiberWorker() {
	for fiber := range fiberRuntime.scheduler.readyQueue {
		executeFiber(fiber)
	}
}

// executeFiber runs a fiber to completion.
func executeFiber(fiber *Fiber) {
	// Mark as running
	fiber.State = FiberRunning

	atomic.AddInt32(&fiberRuntime.scheduler.runningCount, 1)

	// Execute the fiber function
	fiber.Result = fiber.Function()

	// Mark as completed
	fiber.State = FiberCompleted

	atomic.AddInt32(&fiberRuntime.scheduler.runningCount, -1)

	// Signal completion
	close(fiber.Done)
}

// FiberSpawn creates a new fiber and schedules it for execution
//
//export fiber_spawn
func FiberSpawn(fn func() int64) int64 {
	fiberRuntime.mu.Lock()
	defer fiberRuntime.mu.Unlock()

	id := atomic.AddInt64(&fiberRuntime.nextID, 1)

	fiber := &Fiber{
		ID:         id,
		State:      FiberReady,
		Function:   fn,
		Done:       make(chan struct{}),
		ModuleData: make(map[string]interface{}), // Each fiber gets isolated module state
	}

	fiberRuntime.fibers[id] = fiber

	// Schedule the fiber for execution
	fiberRuntime.scheduler.readyQueue <- fiber

	return id
}

// FiberAwait waits for a fiber to complete and returns its result
//
//export fiber_await
func FiberAwait(fiberID int64) int64 {
	fiberRuntime.mu.RLock()
	fiber, exists := fiberRuntime.fibers[fiberID]
	fiberRuntime.mu.RUnlock()

	if !exists {
		return -1 // Error: fiber not found
	}

	// Wait for fiber completion
	<-fiber.Done

	return fiber.Result
}

// FiberYield yields control and returns a value
//
//export fiber_yield
func FiberYield(value int64) int64 {
	// In a real implementation, this would:
	// 1. Save current fiber state
	// 2. Yield to scheduler
	// 3. Resume later
	// For now, just return the value
	return value
}

// ChannelCreate creates a new channel with specified capacity
//
//export channel_create
func ChannelCreate(capacity int64) int64 {
	fiberRuntime.mu.Lock()
	defer fiberRuntime.mu.Unlock()

	id := atomic.AddInt64(&fiberRuntime.nextID, 1)

	channel := &Channel{
		ID:       id,
		Capacity: int(capacity),
		Buffer:   make(chan int64, capacity),
		Closed:   false,
	}

	fiberRuntime.channels[id] = channel

	return id
}

// ChannelSend sends a value through a channel
//
//export channel_send
func ChannelSend(channelID int64, value int64) int64 {
	fiberRuntime.mu.RLock()
	channel, exists := fiberRuntime.channels[channelID]
	fiberRuntime.mu.RUnlock()

	if !exists || channel.Closed {
		return 0 // Error
	}

	// Non-blocking send
	select {
	case channel.Buffer <- value:
		return 1 // Success
	default:
		return 0 // Channel full
	}
}

// ChannelRecv receives a value from a channel
//
//export channel_recv
func ChannelRecv(channelID int64) int64 {
	fiberRuntime.mu.RLock()
	channel, exists := fiberRuntime.channels[channelID]
	fiberRuntime.mu.RUnlock()

	if !exists || channel.Closed {
		return -1 // Error
	}

	// Blocking receive
	value, ok := <-channel.Buffer
	if !ok {
		return -1 // Channel closed
	}

	return value
}

// ModuleGetState gets isolated module state for current fiber
//
//export module_get_state
func ModuleGetState(fiberID int64, key string) interface{} {
	fiberRuntime.mu.RLock()
	fiber, exists := fiberRuntime.fibers[fiberID]
	fiberRuntime.mu.RUnlock()

	if !exists {
		return nil
	}

	return fiber.ModuleData[key]
}

// ModuleSetState sets isolated module state for current fiber
//
//export module_set_state
func ModuleSetState(fiberID int64, key string, value interface{}) {
	fiberRuntime.mu.Lock()
	defer fiberRuntime.mu.Unlock()

	fiber, exists := fiberRuntime.fibers[fiberID]
	if !exists {
		return
	}

	fiber.ModuleData[key] = value
}
