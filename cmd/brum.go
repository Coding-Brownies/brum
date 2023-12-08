package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/beep/effects"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

//go:embed fan.mp3
var audio []byte

func volumeFunc(val float64) float64 {
	y := math.Exp(val/70) - 1
	if y > 100 {
		y = 100
	}
	return y
}

func sampleFunc(val float64) float64 {
	return math.Exp(val/40) - 0.7
}

func main() {
	// load the audio
	streamer, format, err := mp3.Decode(
		io.NopCloser(bytes.NewReader(audio)),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	// create a beep.Streamer which is looped and can be resampled, and turned up and down in volume
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30))
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}
	resampler := beep.ResampleRatio(4, 1, ctrl)
	volume := &effects.Volume{Streamer: resampler, Base: 2}

	read := make(chan bool)

	// play asyncronously
	go func() {
		// wait for the first read
		<-read
		speaker.Play(volume)
	}()

	destVal := 0.0

	// goroutine which makes followVal follow destVal
	go func() {

		followVal := destVal
		const STEP = 0.1
		const TIMEOUT = 1 * time.Millisecond

		for {
			time.Sleep(TIMEOUT)

			// if it is near stop
			if math.Abs(destVal-followVal) < STEP*2 {
				continue
			}

			if destVal > followVal {
				followVal += STEP
			} else {
				followVal -= STEP
			}

			speaker.Lock()
			resampler.SetRatio(sampleFunc(followVal))
			volume.Volume = volumeFunc(followVal)
			speaker.Unlock()
		}
	}()

	isFirstRead := true
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			log.Fatal(err)
		}

		destVal, err = strconv.ParseFloat(strings.TrimSpace(text), 64)
		if err != nil {
			continue
		}
		// if this is the first read start playing the sound
		if isFirstRead {
			read <- true
			isFirstRead = false
		}
	}
}
