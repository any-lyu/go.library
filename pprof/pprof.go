package pprof

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof" // justifying
	"time"

	"github.com/any-lyu/go.library/logs"
)

// Config pprof config
type Config struct {
	Host string
	Port int
}

// Start start pprof
func Start(toExit <-chan struct{}, config *Config) {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	server := &http.Server{Addr: addr, Handler: nil}
	go func() {
		logs.Info("pprof start...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.Fatal(err)
		}
	}()
	<-toExit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logs.Error("failed-to-shutdown-pprof-server", "addr", addr, "error", err.Error())
		return
	}
}
