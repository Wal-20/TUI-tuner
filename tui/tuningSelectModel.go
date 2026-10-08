package tui

import (
	"github.com/Wal-20/tui-tuner.git/audio"
	"github.com/Wal-20/tui-tuner.git/note"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type tuningSelectModel struct {
	tuningList list.Model
}

func newTuningSelectModel() tuningSelectModel {
	items := make([]list.Item, len(note.Tunings))

	for i, tuning := range note.Tunings {
		items[i] = tuningItem{
			tuning: tuning,
		}
	}

	delegate := list.NewDefaultDelegate()

	l := list.New(items, delegate, 40, 20)
	l.Title = "Select Tuning"

	return tuningSelectModel{tuningList: l}
}

func (m tuningSelectModel) Init() tea.Cmd {
	return nil
}

func (m tuningSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			item, ok := m.tuningList.SelectedItem().(tuningItem)
			if !ok {
				return m, nil
			}

			readings, stop, err := audio.Listen()
			if err != nil {
				return m, nil
			}

			return mainModel{
				selectedTuning: item,
				readings:       readings,
				stopAudio:      stop,
			}, waitFor(readings)

		}
	}

	var cmd tea.Cmd
	m.tuningList, cmd = m.tuningList.Update(msg)

	return m, cmd
}

func (m tuningSelectModel) View() string {
	return m.tuningList.View()
}

func RunSelectModel() (tea.Model, error) {

	items := make([]list.Item, len(note.Tunings))

	for i, tuning := range note.Tunings {
		items[i] = tuningItem{
			tuning: tuning,
		}
	}

	tuningList := list.New(items, list.NewDefaultDelegate(), 40, 20)

	m, err := tea.NewProgram(tuningSelectModel{tuningList: tuningList}, tea.WithAltScreen()).Run()
	return m, err
}
