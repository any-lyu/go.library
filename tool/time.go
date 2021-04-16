package tool

import "time"

// BeginningOfHour 获取 t 这个时间点所在小时的开始时间.
//
// 返回的 time.Time 和传入的参数 t 的 *time.Location 一致!
func BeginningOfHour(t time.Time, locationOffsetSeconds int) time.Time {
	const secondsPerHour = 60 * 60
	x := t.Unix()
	x += int64(locationOffsetSeconds)
	x = (x / secondsPerHour) * secondsPerHour
	x -= int64(locationOffsetSeconds)
	return time.Unix(x, 0).In(t.Location())
}

// BeginningOfDay 获取 t 这个时间点所在天的零点时间.
//
// 返回的 time.Time 和传入的参数 t 的 *time.Location 一致!
func BeginningOfDay(t time.Time, locationOffsetSeconds int) time.Time {
	const secondsPerDay = 24 * 60 * 60
	x := t.Unix()
	x += int64(locationOffsetSeconds)
	x = (x / secondsPerDay) * secondsPerDay
	x -= int64(locationOffsetSeconds)
	return time.Unix(x, 0).In(t.Location())
}

// MondayOfWeek 获取 t 这个时间点所在星期的星期一(Monday)的零点时间.
//
// 返回的 time.Time 和传入的参数 t 的 *time.Location 一致!
func MondayOfWeek(t time.Time, locationOffsetSeconds int) time.Time {
	const secondsPerWeek = 7 * 24 * 60 * 60
	const secondsOf3Days = 3 * 24 * 60 * 60 // time.Unix(0, 0) is Thursday
	x := t.Unix()
	x += int64(locationOffsetSeconds)
	x += secondsOf3Days
	x = (x / secondsPerWeek) * secondsPerWeek
	x -= secondsOf3Days
	x -= int64(locationOffsetSeconds)
	return time.Unix(x, 0).In(t.Location())
}

// BeginningOfMonth 获取 t 这个时间点所在 month 的零点时间.
//
// 返回的 time.Time 和传入的参数 t 的 *time.Location 一致!
func BeginningOfMonth(t time.Time, loc *time.Location) time.Time {
	y, m, _ := t.In(loc).Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, loc).In(t.Location())
}

// BeginningOfNextMonth 获取 t 这个时间点的下一个 month 的零点时间.
//
// 返回的 time.Time 和传入的参数 t 的 *time.Location 一致!
func BeginningOfNextMonth(t time.Time, loc *time.Location) time.Time {
	y, m, _ := t.In(loc).Date()
	return time.Date(y, m+1, 1, 0, 0, 0, 0, loc).In(t.Location())
}

// BeginningOfYear 获取 t 这个时间点所在 year 的零点时间.
//
// 返回的 time.Time 和传入的参数 t 的 *time.Location 一致!
func BeginningOfYear(t time.Time, loc *time.Location) time.Time {
	y, _, _ := t.In(loc).Date()
	return time.Date(y, time.January, 1, 0, 0, 0, 0, loc).In(t.Location())
}

// BeginningOfNextYear 获取 t 这个时间点的下一个 year 的零点时间.
//
// 返回的 time.Time 和传入的参数 t 的 *time.Location 一致!
func BeginningOfNextYear(t time.Time, loc *time.Location) time.Time {
	y, _, _ := t.In(loc).Date()
	return time.Date(y+1, time.January, 1, 0, 0, 0, 0, loc).In(t.Location())
}

var (
	// ShanghaiLocationOffset 是东八区的 offset
	ShanghaiLocationOffset = 8 * 60 * 60
	// ShanghaiLocation 表示东八区
	ShanghaiLocation = time.FixedZone("Asia/Shanghai", 8*60*60)
)
