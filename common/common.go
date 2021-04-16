package common

import (
	"context"
	"sync"

	"golang.org/x/sync/singleflight"

	"github.com/any-lyu/go.library/errgroup"
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
