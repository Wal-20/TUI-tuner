package main

import (
	"fmt"
	"os"

	"github.com/Wal-20/tui-tuner.git/tui"
)

func main() {
	if _, err := tui.RunSelectModel(); err != nil {
		fmt.Fprintln(os.Stderr, "tuner:", err)
		os.Exit(1)
	}
}
