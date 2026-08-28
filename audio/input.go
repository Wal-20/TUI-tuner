package audio

import (
	"encoding/binary"
	"fmt"
	"log"

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
	var audioBuffer []int16

	onSamples := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		// pInputSamples contains microphone PCM data.
		fmt.Printf("Received %d frames\n", framecount)
		samples := make([]int16, framecount)

		for i := range framecount {
			offset := i * 2

			samples[i] = int16(binary.LittleEndian.Uint16(
				// takes two pcm bytes from the byte array and converts to one 16 bit int, so samples[i] is one pcm sample
				// LittleEndian doesn't provide Int16, so we use Uint16 and then cast to int16
				pInputSamples[offset : offset+2],
			))
			audioBuffer = append(audioBuffer, samples[i])

			if len(audioBuffer) >= minSamples {
				fmt.Println("Ready to detect frequency!")
				// pitch := detectFrequency(audioBuffer[:minSamples])
				audioBuffer = audioBuffer[minSamples:]
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

func detectFrequency(audioBuffer []int16, interval int, delay int) {
	// buffer contains the minimum needed pcm samples to estimate frequency
	// YIN: compare pairs of samples, squared difference between each sample of the pair, sum the differences

}
