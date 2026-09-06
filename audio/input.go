package audio

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"

	"github.com/Wal-20/tui-tuner.git/note"
	"github.com/gen2brain/malgo"
)

const (
	minFreq         = 70.0
	maxFreq         = 400.0
	maxAperiodicity = 0.5
)

func RecordAudio() {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		fmt.Println("LOG:", message)
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)

	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = 44100 // 44.1 kHz

	// minSamples := int(deviceConfig.SampleRate / uint32(note.Notes[0].Freq))
	minSamples := 4096
	audioBuffer := make([]int16, 0, minSamples)

	onSamples := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		// pInputSamples contains microphone PCM data.

		lastFreq := 0.0
		for i := range framecount {
			offset := i * 2

			sample := int16(binary.LittleEndian.Uint16(
				// takes two pcm bytes from the byte array and converts to one 16 bit int, so sample is one pcm sample
				// LittleEndian doesn't provide Int16, so we use Uint16 and then cast to int16
				pInputSamples[offset : offset+2],
			))
			audioBuffer = append(audioBuffer, sample)

			if len(audioBuffer) == minSamples {
				// without this YIN locks onto hum, rumble, or a ringing lower string
				filtered := highPassFilter(audioBuffer, 100, 44100, 3)

				freq := detectFrequency(filtered, 44100)
				if freq-lastFreq > 100 {
					fmt.Println()
				}
				lastFreq = freq

				note, octave, cents := note.DetectNote(freq)
				fmt.Printf("\r\nDetected frequency: %.2f | Note: %v | Octave: %v | Cents: %.2f                    ", freq, note, octave, cents)

				// clear buffer
				audioBuffer = audioBuffer[:0]
			}
		}
	}

	device, err := malgo.InitDevice(ctx.Context, deviceConfig, malgo.DeviceCallbacks{
		Data: onSamples,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer device.Uninit()

	err = device.Start()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening... Press Enter to stop.")
	fmt.Scanln()
}

// highPassFilter runs `poles` passes of a one-pole high-pass over the PCM buffer.
// Each pass compounds the attenuation below the cutoff.
func highPassFilter(samples []int16, cutoff, sampleRate float64, poles int) []float64 {
	filtered := make([]float64, len(samples))
	for i, s := range samples {
		filtered[i] = float64(s)
	}

	rc := 1.0 / (2 * math.Pi * cutoff)
	alpha := rc / (rc + 1.0/sampleRate)

	for range poles {
		prevIn, prevOut := 0.0, 0.0
		for i, in := range filtered {
			out := alpha * (prevOut + in - prevIn)
			filtered[i] = out
			prevIn, prevOut = in, out
		}
	}

	return filtered
}

func detectFrequency(audioBuffer []float64, sampleRate uint32) float64 {
	tauMin := int(float64(sampleRate) / maxFreq)
	tauMax := int(float64(sampleRate) / minFreq)
	if tauMax >= len(audioBuffer) {
		tauMax = len(audioBuffer) - 1
	}

	// Step 1: Difference function
	differences := make([]float64, tauMax+1)

	for tau := 1; tau <= tauMax; tau++ {
		sqDiff := 0.0

		for i := 0; i < len(audioBuffer)-tau; i++ {
			diff := audioBuffer[i] - audioBuffer[i+tau]
			sqDiff += diff * diff
		}

		differences[tau] = sqDiff
	}

	// Step 2: Cumulative mean normalized difference
	cmndf := make([]float64, tauMax+1)
	cmndf[0] = 1.0
	runningSum := 0.0
	threshold := 0.15
	bestTau := -1

	// deepest dip so far, used when nothing crosses the threshold
	deepest, deepestTau := maxAperiodicity, -1

	for tau := 1; tau <= tauMax; tau++ {
		// The running sum accumulates from tau 1 even below tauMin, or the
		// normalization is wrong for every tau above it.
		runningSum += differences[tau]
		cmndf[tau] = differences[tau] / (runningSum / float64(tau))

		if tau < tauMin {
			continue
		}

		if cmndf[tau] < deepest {
			deepest, deepestTau = cmndf[tau], tau
		}

		if cmndf[tau] < threshold {
			// crossed the threshold — walk down to the bottom of this dip
			for tau+1 <= tauMax {
				runningSum += differences[tau+1]
				cmndf[tau+1] = differences[tau+1] / (runningSum / float64(tau+1))
				if cmndf[tau+1] >= cmndf[tau] {
					break // it turned back up; tau is the minimum
				}
				tau++
			}
			bestTau = tau
			break
		}
	}

	if bestTau < 1 {
		if deepestTau < 1 {
			return 0 // not periodic enough for a note
		}
		bestTau = deepestTau
	}

	// Parabolic interpolation
	refinedTau := float64(bestTau)

	if bestTau+1 <= tauMax {
		y1, y2, y3 := cmndf[bestTau-1], cmndf[bestTau], cmndf[bestTau+1]
		if denom := y1 - 2*y2 + y3; denom > 0 {
			if shift := (y1 - y3) / (2 * denom); shift > -1 && shift < 1 {
				refinedTau += shift
			}
		}
	}

	return float64(sampleRate) / refinedTau
}
