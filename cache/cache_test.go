package cache

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	//Set("server",[]byte("haha"))

	bytes, e := Get("server")
	fmt.Printf("%s err = %v \n", bytes, e)
}
