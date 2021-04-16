package sql

import (
	"time"

	"github.com/any-lyu/go.library/stat"
)

// statistics 收集监控信息
func statistics(stats stat.Stat, name string, t time.Time, err error) {
	if stats == nil {
		return
	}
	if err != nil {
		stats.Incr(name, "breaker")
		return
	}
	stats.Timing(name, int64(time.Since(t).Seconds()))
}
