package audio

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/Wal-20/tui-tuner.git/note"
	"github.com/gen2brain/malgo"
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
				delays := make([]int, 0, 1000)

				for tau := 0; tau <= 1000; tau++ {
					delays = append(delays, tau)
				}

				freq := detectFrequency(audioBuffer, delays, 44100)
				if freq-lastFreq > 100 {
					fmt.Println()
				}
				lastFreq = freq

				note, octave, cents := note.DetectNote(freq)
				fmt.Printf("\r\nDetected frequency: %.2f | Note: %v | Octave: %v | Cents: %.2f                    ", freq, note, octave, cents)

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

func detectFrequency(audioBuffer []int16, delays []int, sampleRate uint32) float64 {
	differences := make(map[int]float64)

	// Step 1: Difference function
	for _, tau := range delays {
		sqDiff := 0.0

		for i := 0; i < len(audioBuffer)-tau; i++ {
			first := float64(audioBuffer[i])
			second := float64(audioBuffer[i+tau])

			diff := first - second
			sqDiff += diff * diff
		}

		differences[tau] = sqDiff
	}

	// Step 2: Cumulative mean normalized difference
	runningSum := 0.0
	threshold := 0.1

	bestTau := 0

	for _, tau := range delays {
		runningSum += differences[tau]

		normalized := differences[tau] /
			(runningSum / float64(tau))

		// Find the first sufficiently strong dip
		if normalized < threshold {
			bestTau = tau
			break
		}
	}

	if bestTau == 0 {
		return 0
	}

	return float64(sampleRate) / float64(bestTau)
}
