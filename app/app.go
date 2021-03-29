package app

import (
	"context"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"github.com/any-lyu/go.library/logs"

)

func init() {
	app.ctx, app.cancel = context.WithCancel(context.Background())
	app.wgDone = make(chan struct{})
	go closeAppWhenAnExitSignalIsEncountered()
}

var app struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	wgDone chan struct{}

	mu     sync.Mutex
	defers []func()
}

// Go 开启一个新的 goroutine 用以执行 fn 函数, 第一个参数 ctx context.Context 会传递给 fn 的第一个参数 ctx context.Context.
//
// 当 Close 函数被调用或者收到系统的退出信号, fn 函数的 closing 会收到信号, fn 根据这个信号可以做一些逻辑.
//
// Go 函数和语言级别的 go 关键字不同点在于 Wait 会等待 fn 执行结束.
func Go(ctx context.Context, fn func(ctx context.Context, closing <-chan struct{})) {
	app.wg.Add(1)

	go func(ctx context.Context, closing <-chan struct{}, wg *sync.WaitGroup, fn func(ctx context.Context, closing <-chan struct{})) {
		defer wg.Done()

		fn(ctx, closing)
	}(ctx, Closing(), &app.wg, fn)
}

// Defer 会在 Wait 函数结束之前执行, 并且符合 FILO 的规则串行执行.
//
// Defer 可以用于在程序退出之前做资源释放, 比如关闭数据库.
func Defer(fn func()) {
	if fn == nil {
		return
	}
	app.mu.Lock()
	app.defers = append(app.defers, fn)
	app.mu.Unlock()
}

// Wait 等待所有的 Go 任务退出然后才退出.
func Wait() {
	defer close(app.wgDone)

	app.wg.Wait()

	// call defers
	app.mu.Lock()
	defer app.mu.Unlock()
	for i := len(app.defers) - 1; i >= 0; i-- {
		app.defers[i]()
	}
}

// Close 触发 closing 信号, 通知相关的 goroutine 做好退出工作.
func Close() { app.cancel() }

// Closing 返回一个信号 chan, 当 Close 被调用后这个 chan 会收到通知.
func Closing() <-chan struct{} { return app.ctx.Done() }

// closeAppWhenAnExitSignalIsEncountered 捕获退出信号关闭 app
func closeAppWhenAnExitSignalIsEncountered() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	sig := <-sigCh
	logs.Info("got-exit-signal-and-closing", "signal", sig)
	Close()

	// 等待一定时间如果程序还没有关闭则强制关闭
	timer := time.NewTimer(30 * time.Second)
	select {
	case <-timer.C:
		// 打印 panic 级别的日志, 触发告警
		if _, err := os.Stderr.WriteString("panic: graceful_shutdown_timeout"); err != nil {
			logs.Fatal("graceful_shutdown_timeout")
		}
		debug.PrintStack()
		os.Exit(1)
	case <-app.wgDone:
		timer.Stop()
	}
}
