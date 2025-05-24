package graceful

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	// 其它需要跟主协程一起优雅退出的协程需要监听此ctx
	mainContext context.Context
	mainCancel  context.CancelFunc

	wg sync.WaitGroup

	signalChannel = make(chan os.Signal, 1)
)

// Config 配置
type Config struct {
	Logger             Logger        // 日志记录器
	WaitTimeout        time.Duration // 等待超时时间（秒）
	WaitingMessage     string        // 收到退出信号，等待所有协程关闭
	WaitDoneMessage    string        // 所有协程已退出
	WaitTimeoutMessage string        // 等待协程退出超时
	QuitMessage        string        // 进程退出
}

// WaitGroup 获取等待组
func WaitGroup() *sync.WaitGroup {
	return &wg
}

// Context 获取 main 上下文
func Context() context.Context {
	return mainContext
}

// Exit 手动发出退出信号
func Exit(sig os.Signal) {
	signalChannel <- sig
}

// Start 开始监听信号
func Start(cfg *Config) {
	// 监听终止信号
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞等待信号
	<-signalChannel

	// 向所有监听了 mainContext 的协程发送取消信号
	mainCancel()

	if cfg.Logger != nil && cfg.WaitingMessage != "" {
		cfg.Logger.Output(cfg.WaitingMessage)
	}

	// 监听wg的通道
	waitChannel := make(chan struct{}, 1)
	go func() {
		// 阻塞等待所有协程关闭
		wg.Wait()
		// 关闭done
		close(waitChannel)
	}()

	select {
	case <-waitChannel:
		if cfg.Logger != nil && cfg.WaitDoneMessage != "" {
			cfg.Logger.Output(cfg.WaitDoneMessage)
		}
	case <-time.After(cfg.WaitTimeout * time.Second):
		if cfg.Logger != nil && cfg.WaitTimeoutMessage != "" {
			// 超时退出
			cfg.Logger.Output(cfg.WaitTimeoutMessage)
		}
	}

	if cfg.Logger != nil && cfg.QuitMessage != "" {
		cfg.Logger.Output(cfg.QuitMessage)
	}
}
