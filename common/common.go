package common

import (
	"context"
	"github.com/any-lyu/go.library/errgroup"
	"golang.org/x/sync/singleflight"
	"sync"
)

var (
	groupCommon     *errgroup.Group
	onceGroupCommon sync.Once
	singleGroup     singleflight.Group
)

// GetErrGroup get common errgroup.Group
func GetErrGroup() *errgroup.Group {
	onceGroupCommon.Do(func() {
		groupCommon = errgroup.WithContext(context.Background())
		groupCommon.GOMAXPROCS(100)
	})
	return groupCommon
}

// GetSingleGroup common singleGroup
func GetSingleGroup() *singleflight.Group {
	return &singleGroup
}
