package pool

import (
	"bytes"
	"sync"
)

// BytesBufferPool16KB 是 16KB 大小的 sync.Pool
var BytesBufferPool16KB = &sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 16<<10))
	},
}
