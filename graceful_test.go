package graceful

import (
	"log"
	"syscall"
	"testing"
	"time"
)

// 工作协程
func testWorker(t *testing.T) {
	// 添加一个计数器
	WaitGroup().Add(1)

	defer func() {
		// 将协程计数减1
		WaitGroup().Done()
		t.Log("worker exited")
	}()

	// 模拟超时
	for {
		select {
		case <-Context().Done():
			t.Log("main context done")
			return
		default:
			time.Sleep(2 * time.Second)
			t.Log("working...")
		}
	}
}

func TestGracefulExit(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	defer func() {
		logger := testNewLogger(t)

		Start(&Config{
			Logger:             logger,
			WaitTimeout:        10,
			WaitingMessage:     "waiting goroutine close",
			WaitDoneMessage:    "all goroutine have closed",
			WaitTimeoutMessage: "wait timed out",
			QuitMessage:        "exit process",
		})
	}()

	go testWorker(t)

	// 5秒钟后发送退出信号
	go func() {
		time.Sleep(5 * time.Second)
		Exit(syscall.SIGTERM)
	}()
}

func TestGracefulCancel(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	logger := testNewLogger(t)

	defer func() {
		Start(&Config{
			Logger:             logger,
			WaitTimeout:        10,
			WaitingMessage:     "waiting goroutine close",
			WaitDoneMessage:    "all goroutines have closed",
			WaitTimeoutMessage: "wait timed out",
			QuitMessage:        "exit process",
		})
	}()

	go testWorker(t)

	// 5秒钟后发送取消信号
	go func() {
		time.Sleep(5 * time.Second)
		Cancel()
	}()
}
