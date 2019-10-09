package tool

import "time"

// BeginningOfDay 获取 t 这个时间点所在天的零点时间.
func BeginningOfDay(t time.Time, zoneOffsetSeconds int64) time.Time {
	const secondsPerDay = 24 * 60 * 60
	x := t.Unix()
	x += zoneOffsetSeconds
	x = (x / secondsPerDay) * secondsPerDay
	x -= zoneOffsetSeconds
	return time.Unix(x, 0)
}

var (
	// ShanghaiLocation 表示东八区
	ShanghaiLocation = time.FixedZone("Asia/Shanghai", 8*60*60)
)
