package tui

import "github.com/Wal-20/tui-tuner.git/note"

// implements bubbletea list.Item
type tuningItem struct {
	tuning note.Tuning
}

func (item tuningItem) FilterValue() string {
	return item.tuning.Name
}

func (item tuningItem) Title() string {
	return item.tuning.Name
}

func (item tuningItem) Description() string {
	return ""
}