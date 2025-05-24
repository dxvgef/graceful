package graceful

import (
	"log"
	"os"
	"syscall"
	"testing"
	"time"
)

// 工作协程
func testWorker(t *testing.T) {
	// 添加一个计数器
	WaitGroup().Add(1)

	// 协程退出时，将计数减1
	defer func() {
		WaitGroup().Done()
		t.Log("worker exited")
	}()

	// 模拟超时
	for i := 0; i <= 5; i++ {
		time.Sleep(2 * time.Second)
		t.Log("working...")
	}

}

func TestGraceful(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	// 启动工作协程
	go testWorker(t)

	// 5秒钟后发送退出信号
	go func() {
		time.Sleep(5 * time.Second)
		Exit(syscall.SIGTERM)
	}()

	// 创建 graceful.Logger
	logger := testNewLogger(os.Stdout, "[graceful] ", log.Ldate|log.Ltime)

	// 阻塞启动退出监听，并设置等待协议超时的时间
	Start(&Config{
		Logger:             logger,
		WaitTimeout:        10,
		WaitingMessage:     "pending coroutine shutdown",
		WaitDoneMessage:    "all coroutines have exited",
		WaitTimeoutMessage: "wait timed out",
		QuitMessage:        "exit process",
	})
}
