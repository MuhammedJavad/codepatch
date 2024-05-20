package adt

import (
	"context"
	"dasht/pkg/log"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

// Tests that the bucket leaks once when called with a single size
func TestBucket_WithSingleLeak_ShouldLeakOnce(t *testing.T) {
	// Arrange
	const size int = 2
	const timeout = 10
	wg := new(sync.WaitGroup)
	wg.Add(1)
	var counter int
	var leakedLength int
	b := configureTestBucket(size, timeout, func(ctx context.Context, i []interface{}) error {
		defer wg.Done()
		leakedLength = len(i)
		counter++
		return nil
	})
	// Act
	for i := 0; i < size; i++ {
		b.Add(struct{}{})
	}
	// Assert
	wg.Wait()
	assert.True(t, counter == 1)
	assert.True(t, leakedLength == size)
}

// Tests that the bucket leaks twice when called with double the size
func TestBucket_WithDoubleLeak_ShouldLeakTwice(t *testing.T) {
	// Arrange
	const size int = 2
	const timeout = 10
	wg := new(sync.WaitGroup)
	wg.Add(2)
	var counter int
	var leakedLength int
	b := configureTestBucket(size, timeout, func(ctx context.Context, i []interface{}) error {
		defer wg.Done()
		leakedLength += len(i)
		counter++
		return nil
	})
	// Act
	for i := 0; i < size*2; i++ {
		b.Add(struct{}{})
		time.Sleep(100 * time.Millisecond)
	}
	// Assert
	wg.Wait()
	assert.True(t, counter == 2)
	assert.True(t, leakedLength == size*2)
}

// Tests that the bucket leaks all its contents when a timeout occurs.
func TestBucket_WithTimeout_ShouldLeakEverything(t *testing.T) {
	// Arrange
	const publishNumber int = 2
	const size int = 200
	const timeout = 1
	wg := new(sync.WaitGroup)
	wg.Add(1)
	var counter int
	var leakedLength int
	b := configureTestBucket(size, timeout, func(ctx context.Context, i []interface{}) error {
		defer wg.Done()
		leakedLength += len(i)
		counter++
		return nil
	})
	// Act
	for i := 0; i < publishNumber; i++ {
		b.Add(struct{}{})
		time.Sleep(100 * time.Millisecond)
	}
	// Assert
	wg.Wait()
	assert.True(t, counter == 1)
	assert.True(t, leakedLength == publishNumber)
}

// Tests that the bucket leaks only once when a timeout and a size call happen simultaneously.
func TestBucket_WithTimeoutAndSize_ShouldOnlyLeakOnce(t *testing.T) {
	// Arrange
	const size int = 2
	const timeout = 1
	wg := new(sync.WaitGroup)
	wg.Add(1)
	var counter int
	var leakedLength int
	b := configureTestBucket(size, timeout, func(ctx context.Context, i []interface{}) error {
		defer wg.Done()
		leakedLength += len(i)
		counter++
		return nil
	})
	// Act
	for i := 0; i < size; i++ {
		b.Add(struct{}{})
	}
	// Assert
	wg.Wait()
	time.Sleep(time.Second)
	assert.True(t, counter == 1)
	assert.True(t, leakedLength == size)
}

// Tests that the bucket handles gracefully (without crashing) when a panic occurs during a leak
func TestBucket_WithPanicOnLeak_ShouldHandleGracefully(t *testing.T) {
	// Arrange
	// Arrange
	const size int = 2
	const timeout = 100
	wg := new(sync.WaitGroup)
	wg.Add(1)
	b := configureTestBucket(size, timeout, func(ctx context.Context, i []interface{}) error {
		defer wg.Done()
		panic("absurd")
	})
	// Act
	for i := 0; i < size; i++ {
		b.Add(struct{}{})
	}
	// Assert
	wg.Wait()
	time.Sleep(time.Second)
	assert.True(t, true)
}

func configureTestBucket(
	size int,
	timeoutInSec int,
	onLeak func(context.Context, []interface{}) error) *Bucket {

	logger := log.InitTestLogger()
	return NewBucket(size, timeoutInSec, "name", logger, onLeak)
}
