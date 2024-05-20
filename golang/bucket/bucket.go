package adt

import (
	"context"
	"dasht/pkg/log"
	"fmt"
	"time"
)

type Bucket struct {
	size         int
	timeoutInSec int
	name         string
	logger       *log.ZerologLogger
	onLeak       func(context.Context, []interface{}) error
	c            chan interface{}
	buffer       []interface{}
}

func NewBucket(
	size int,
	timeoutInSec int,
	name string,
	logger *log.ZerologLogger,
	onLeak func(context.Context, []interface{}) error) *Bucket {
	b := Bucket{
		name:         name,
		size:         size,
		logger:       logger,
		onLeak:       onLeak,
		timeoutInSec: timeoutInSec,
		c:            make(chan interface{}, size),
	}

	go b.watcher()
	return &b
}

func (b *Bucket) Add(data interface{}) {
	select {
	case b.c <- data:
		break
	default:
		b.logger.Error(fmt.Sprintf("bucket is full (dropping the message, buffer_size: %d, bucket_name: %s)", b.size, b.name), data)
		return
	}
}

func (b *Bucket) watcher() {
	duration := time.Duration(b.timeoutInSec) * time.Second
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			b.producer()
		case data := <-b.c:
			// runs only in single thread
			b.buffer = append(b.buffer, data)
			if len(b.buffer) >= b.size {
				b.producer()
				ticker.Reset(duration)
			}
		}
	}
	ticker.Stop()
}

func (b *Bucket) producer() {
	if len(b.buffer) <= 0 {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			b.logger.Error(fmt.Sprintf("panic happend in bucket:%v", r), r)
		}
	}()
	err := b.onLeak(context.Background(), b.buffer)
	if err != nil {
		b.logger.Error(fmt.Sprintf("calling leak handler in bucket %s returned an error", b.name), err.Error())
	}
	b.buffer = nil
}
