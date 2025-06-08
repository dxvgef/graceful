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

// Context 获取 mainContext
func Context() context.Context {
	return mainContext
}

// Cancel 执行 mainCancel，关闭所有监听了 mainContext 的协程，但不会退出进程
func Cancel() {
	mainCancel()
}

// Exit 手动发出退出信号
func Exit(sig os.Signal) {
	signalChannel <- sig
}

// Start 开始监听信号
func Start(cfg *Config) {
	signal.Notify(signalChannel,
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTERM, // 系统终止信号
		syscall.SIGHUP,  // 重启或重新加载配置
		syscall.SIGQUIT, // 退出并生成core dump
		syscall.SIGUSR1, // 用户自定义信号1
		syscall.SIGUSR2, // 用户自定义信号2
	)

	// 阻塞等待退出信号
	<-signalChannel
	// 向所有监听了 mainContext 的协程发出取消信号
	mainCancel()

	// select {
	// case <-signalChannel: // 退出信号
	// 	mainCancel()
	// 	break
	// case <-mainContext.Done(): // 取消信号
	// 	break
	// }

	if cfg.Logger != nil && cfg.WaitingMessage != "" {
		cfg.Logger.Output(cfg.WaitingMessage)
	}

	waitChannel := make(chan struct{}, 1)
	go func() {
		// 阻塞等待所有协程关闭
		wg.Wait()
		// 关闭done
		close(waitChannel)
	}()

	select {
	case <-waitChannel: // 所有协程已关闭
		if cfg.Logger != nil && cfg.WaitDoneMessage != "" {
			cfg.Logger.Output(cfg.WaitDoneMessage)
		}
	case <-time.After(cfg.WaitTimeout * time.Second): // 超时退出
		if cfg.Logger != nil && cfg.WaitTimeoutMessage != "" {
			cfg.Logger.Output(cfg.WaitTimeoutMessage)
		}
	}

	if cfg.Logger != nil && cfg.QuitMessage != "" {
		cfg.Logger.Output(cfg.QuitMessage)
	}
}
