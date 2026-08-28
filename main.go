package main

import (
	"github.com/Wal-20/tui-tuner.git/note"
)

func main() {
	// 0.02 -> 49 calls to detectNote, reasonable start
	// note.Benchmark(5, 0.02)

	note.ListenInput()
}
