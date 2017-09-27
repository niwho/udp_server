package main

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
)

func init() {
}

// AsyncFrame 异步写入, 考虑更通用的方式
type AsyncFrame struct {
	buf        chan []byte
	stop       chan struct{}
	flush      chan struct{}
	handler    func([]byte)
	workerNum  int
	isRunning  int32
	bufferPool *sync.Pool
}

func NewAsyncFrame(workerNum int, handler func([]byte)) *AsyncFrame {
	af := &AsyncFrame{
		buf:       make(chan []byte, 1024),
		stop:      make(chan struct{}),
		flush:     make(chan struct{}),
		handler:   handler,
		workerNum: workerNum,
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
	if workerNum <= 0 {
		af.workerNum = 1
	}
	af.Run()

	return af
}

// cleanBuf
func (af *AsyncFrame) cleanBuf() {
	for {
		select {
		case dat := <-af.buf:
			af.handler(dat)
		default:
			return
		}
	}
}

func (af *AsyncFrame) Run() {
	if !atomic.CompareAndSwapInt32(&af.isRunning, 0, 1) {
		return
	}
	if af.handler == nil {
		fmt.Fprintln(os.Stderr, "logger's handler is nil.")
		// return // if return, af.Stop() will be blocked because no goroutine consumes af.stop
	}
	for i := 0; i < af.workerNum; i++ {
		go af.runOuter()
	}
}

func (af *AsyncFrame) runOuter() {
	for {
		select {
		case <-af.stop:
			af.stop <- struct{}{}
			return
		default:
			af.runrun()
		}
	}

}

func (af *AsyncFrame) runrun() {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 20
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("AsyncFrame panic=%v\n%s\n", err, buf)
		}
	}()
	for {
		select {
		case dat, ok := <-af.buf:
			if !ok {
				fmt.Fprintln(os.Stderr, "buf channel has been closed.")
				af.stop <- struct{}{}
				return
			}
			af.handler(dat)
		case <-af.flush:
			af.cleanBuf()
			af.flush <- struct{}{}
		case <-af.stop:
			af.cleanBuf()
			af.stop <- struct{}{}
			return
		}
	}

}

func (af *AsyncFrame) Write(p []byte) (int, error) {
	// 注意拷贝
	select {
	case af.buf <- p:
		return 0, nil
	default:
		// warn write loss
		return 0, fmt.Errorf("%s", "AsyncFrame buf overflow")
	}
}

// Close safe clean
func (af *AsyncFrame) Close() {
	if !atomic.CompareAndSwapInt32(&af.isRunning, 1, 0) {
		return
	}
	af.stop <- struct{}{}
	<-af.stop
}

// Flush
func (af *AsyncFrame) Flush() {
	if atomic.LoadInt32(&af.isRunning) == 0 {
		return
	}
	af.flush <- struct{}{}
	<-af.flush
}
