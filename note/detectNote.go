package note

import (
	"fmt"
	"math"
)

func detectNote(freq float64) (string, int, float64) {

	var noteNames = []string{
		"C", "C#", "D", "D#", "E", "F",
		"F#", "G", "G#", "A", "A#", "B",
	}
	midi := int(math.Round(69 + 12*math.Log2(freq/440)))

	note := noteNames[midi%12]
	octave := midi/12 - 1

	expectedFreq := 440 * math.Pow(2, float64(midi-69)/12)

	cents := 1200 * math.Log2(freq/expectedFreq)

	return note, octave, cents
}

func ListenInput() {
	var freq float64
	fmt.Print("Enter a frequency: ")
	fmt.Scan(&freq)
	note, octave, cents := detectNote(freq)
	fmt.Printf("\nNote: %v | Octave: %v | Cents: %v\n", note, octave, cents)
}
