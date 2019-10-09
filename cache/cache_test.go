package cache

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	//Set("test",[]byte("haha"))

	bytes, e := Get("test")
	fmt.Printf("%s err = %v \n", bytes, e)
}
