package graceful

// Logger 用于优雅退出时的信息输出
type Logger interface {
	Output(any)
}
