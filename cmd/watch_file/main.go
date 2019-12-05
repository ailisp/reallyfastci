package main

import (
	"fmt"
	"github.com/hpcloud/tail"
)

// main
func main() {
	t, _ := tail.TailFile("/tmp/b.txt", tail.Config{
		Follow: true, ReOpen: true, MustExist: false,
		Logger: tail.DiscardingLogger,
	})
	t2, _ := tail.TailFile("/tmp/exitcode", tail.Config{
		Logger:    tail.DiscardingLogger,
		MustExist: false, ReOpen: true, Follow: true,
	})
	for {
		select {
		case line := <-t.Lines:
			fmt.Println(line.Text)
		case <-t2.Lines:
			fmt.Println("bbb")
			return
		}
	}
}
