package group_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/any-lyu/go.library/errors"

	"github.com/any-lyu/go.library/group"
)

func TestZ(t *testing.T) {
	var g group.Group

	g.Add(func() error {
		fmt.Println("10")
		return nil
	}, func(e error) {
		time.Sleep(time.Second * 3)
		fmt.Println("1")
	})

	g.Add(func() error {

		fmt.Println("20")
		return nil
	}, func(e error) {

		fmt.Println("2")
	})
	g.Add(func() error {
		fmt.Println("30")
		// time.Sleep(time.Second * 3)
		return errors.New("message")
	}, func(e error) {
		fmt.Println("3")
	})
	g.Run()
	// time.Sleep(time.Second * 2)
}
