
package tui
import (
	"fmt"
	"math"
	"strings"

	"github.com/Wal-20/tui-tuner.git/audio"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	inTuneCents = 5.0 // within this, the note counts as in tune
	meterRange  = 50.0
	holdSilent  = 12 // keep the last note on screen for ~1s of silence

	minPanel, maxPanel = 32, 72
	tightHeight        = 18 // below this, drop the blank separator lines
)

var (
	greenC = lipgloss.Color("42")
	amberC = lipgloss.Color("214")
	faintC = lipgloss.Color("240")
	plainC = lipgloss.Color("252")

	faint = lipgloss.NewStyle().Foreground(faintC)
	plain = lipgloss.NewStyle().Foreground(plainC)

	panel = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(faintC).
		Padding(1, 3)
)

// tone styles the note and everything that tracks it: green in tune, amber out.
func tone(cents float64) (lipgloss.Style, lipgloss.Style) {
	c := amberC
	if math.Abs(cents) <= inTuneCents {
		c = greenC
	}

	text := lipgloss.NewStyle().Foreground(c).Bold(true)
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(c).
		Padding(0, 3)

	return text, box
}

type readingMsg audio.Reading

type mainModel struct {
	readings       <-chan audio.Reading
	current        audio.Reading
	silent         int
	width          int
	selectedTuning tuningItem
	height         int
	stopAudio      func()
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

// / waitFor blocks in its own goroutine until the detector produces a Reading.
func waitFor(readings <-chan audio.Reading) tea.Cmd {
	return func() tea.Msg {
		r, ok := <-readings
		if !ok {
			return tea.QuitMsg{}
		}
		return readingMsg(r)
	}
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.stopAudio != nil {
				m.stopAudio()
			}
			return m, tea.Quit
		case "esc":
			selectModel, err := RunSelectModel()
			if err == nil {
				return selectModel, nil
			}
		}

	case readingMsg:
		if r := audio.Reading(msg); r.Freq > 0 {
			m.current, m.silent = r, 0
		} else {
			m.silent++
		}
		return m, waitFor(m.readings)
	}

	return m, nil
}

// dims returns the terminal size, falling back to a sane default before the
// first WindowSizeMsg arrives.
func (m mainModel) dims() (int, int) {
	w, h := m.width, m.height
	if w == 0 {
		w, h = 80, 24
	}
	return w, h
}

func (m mainModel) View() string {
	w, h := m.dims()

	// A short terminal can't fit the spaced-out layout, so drop the breathing
	// room rather than overflow.
	tight := h < tightHeight

	inner := innerWidth(w)

	card := panel
	if tight {
		card = card.Padding(0, 3)
	}

	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center,
			card.Render(m.body(inner, tight)), "", faint.Render("q quit")))
}

// innerWidth sizes the panel's content. The card costs inner+8 columns once
// padding and borders are counted, so a narrow terminal overrides minPanel.
func innerWidth(w int) int {
	inner := min(max(w-10, minPanel), maxPanel)
	return max(min(inner, w-8), 20)
}

// body renders the same number of lines in every state, so the panel doesn't
// shift under vertical centring as notes come and go.
func (m mainModel) body(inner int, tight bool) string {
	center := lipgloss.NewStyle().Width(inner).Align(lipgloss.Center)
	blank := center.Render("")

	// Groups are separated by a blank line unless space is short.
	join := func(groups ...[]string) string {
		var out []string
		for i, g := range groups {
			if i > 0 && !tight {
				out = append(out, blank)
			}
			out = append(out, g...)
		}
		return strings.Join(out, "\n")
	}

	// Each line of the note's rectangle has to be centred individually, since
	// the box is narrower than the panel.
	boxed := func(s string, box lipgloss.Style) []string {
		out := strings.Split(box.Render(s), "\n")
		for i, l := range out {
			out[i] = center.Render(l)
		}
		return out
	}

	if m.current.Freq == 0 || m.silent > holdSilent {
		idle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(faintC).
			Padding(0, 3)

		return join(
			boxed(faint.Render(" — "), idle),
			[]string{center.Render(faint.Render("listening…"))},
			[]string{blank},
			[]string{blank},
			[]string{center.Render(stringRow(audio.Reading{}, faint, inner, m.selectedTuning))},
		)
	}

	r := m.current
	text, box := tone(r.Cents)

	name := strings.ReplaceAll(r.Note, "#", "♯") + fmt.Sprint(r.Octave)

	// cents = 1200*log2(freq/target), so this inverts back to the target pitch.
	target := r.Freq / math.Pow(2, r.Cents/1200)

	// The target only fits alongside the reading on a wide enough panel; letting
	// it wrap would add a line and break the fixed height.
	info := plain.Render(fmt.Sprintf("%.2f Hz", r.Freq))
	if inner >= 34 {
		info += faint.Render(fmt.Sprintf("   target %.2f Hz", target))
	}

	return join(
		boxed(text.Render(name), box),
		[]string{center.Render(text.Render(fmt.Sprintf("%+.1f cents", r.Cents)))},
		[]string{center.Render(faint.Render("♭ ") + text.Render(meter(r.Cents, inner-8)) + faint.Render(" ♯"))},
		[]string{center.Render(stringRow(r, text, inner, m.selectedTuning))},
		[]string{center.Render(info)},
	)
}

// stringRow shows standard tuning with the string being played picked out.
func stringRow(r audio.Reading, text lipgloss.Style, inner int, tuning tuningItem) string {
	gap := "   "
	if inner < 24 {
		gap = " "
	}

	var b strings.Builder

	for i, s := range tuning.tuning.Notes {
		if i > 0 {
			b.WriteString(faint.Render(gap))
		}
		if r.Freq > 0 && s.Name == r.Note && s.Octave == r.Octave {
			b.WriteString(text.Render(s.Label))
			continue
		}
		b.WriteString(faint.Render(s.Label))
	}

	return b.String()
}

// meter draws the deviation as a dot either side of a centre line.
func meter(cents float64, width int) string {
	width = max(width, 11)
	if width%2 == 0 {
		width--
	}

	mid := width / 2
	pos := max(0, min(width-1, mid+int(math.Round(cents/meterRange*float64(mid)))))

	var b strings.Builder
	for i := range width {
		switch i {
		case pos:
			b.WriteString("●")
		case mid:
			b.WriteString("│")
		default:
			b.WriteString("─")
		}
	}

	return b.String()
}
