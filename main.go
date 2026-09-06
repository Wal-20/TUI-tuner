package main

import (
	"fmt"
	"os"

	"github.com/Wal-20/tui-tuner.git/tui"
)

func main() {
	if err := tui.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "tuner:", err)
		os.Exit(1)
	}
}
