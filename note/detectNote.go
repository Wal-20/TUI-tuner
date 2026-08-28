package note

import (
	"fmt"
	"math"
	"time"
)

func detectNote(freq float64) (string, int, float64) {
	if freq <= 0 {
		return "", 0, 0
	}

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
	fmt.Printf("Note: %v | Octave: %v | Cents: %v\n", note, octave, cents)
}

func Benchmark(seconds float64, sleepDuration float64) {
	calculations := 0

	fmt.Println("Started benchmark...")

	start := time.Now()
	for time.Since(start).Seconds() < seconds {
		for _, val := range Notes {
			if time.Since(start).Seconds() >= seconds {
				break
			}
			detectNote(val.Freq)
			calculations++
			time.Sleep(time.Duration(sleepDuration * float64(time.Second)))
		}
	}
	fmt.Printf("%v Frequency to note calculations per second\n", calculations/int(seconds))
}
