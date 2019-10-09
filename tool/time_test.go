package tool

import (
	"testing"
	"time"
)

func TestBeginningOfDay(t *testing.T) {
	const offset = 8 * 60 * 60
	loc := time.FixedZone("Asia/Shanghai", offset)

	const year = 2018
	for month := time.January; month <= time.December; month++ {
	DayLoop:
		for day := 1; day <= 31; day++ {
			switch month {
			case time.January:
			case time.February:
				if day > 28 { // 不是闰年
					break DayLoop
				}
			case time.March:
			case time.April:
				if day > 30 {
					break DayLoop
				}
			case time.May:
			case time.June:
				if day > 30 {
					break DayLoop
				}
			case time.July:
			case time.August:
			case time.September:
				if day > 30 {
					break DayLoop
				}
			case time.October:
			case time.November:
				if day > 30 {
					break DayLoop
				}
			case time.December:
			}
			t.Log(year, month, day)
			for hour := 0; hour < 24; hour++ {
				a := time.Date(year, month, day, hour, 10, 20, 30, loc)
				b := time.Date(year, month, day, 0, 0, 0, 0, loc)

				// BeginningOfDay
				x := BeginningOfDay(a, offset)
				if !x.Equal(b) {
					t.Errorf("want equal, have:%v, want:%v", x, b)
					return
				}

				// BeginningOfDayTest
				x = BeginningOfDayTest(a, loc)
				if !x.Equal(b) {
					t.Errorf("want equal, have:%v, want:%v", x, b)
					return
				}
			}
		}
	}
}

func BenchmarkBeginningOfDay(b *testing.B) {
	const offset = 8 * 60 * 60
	t := time.Now()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		BeginningOfDay(t, offset)
	}
}

func BenchmarkBeginningOfDayTest(b *testing.B) {
	const offset = 8 * 60 * 60
	loc := time.FixedZone("Asia/Shanghai", offset)
	t := time.Now()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		BeginningOfDayTest(t, loc)
	}
}

func BeginningOfDayTest(t time.Time, loc *time.Location) time.Time {
	y, m, d := t.In(loc).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, loc)
}
