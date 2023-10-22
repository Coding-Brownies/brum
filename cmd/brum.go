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

func main() {
	embeddedReader := io.NopCloser(bytes.NewReader(audio))
	streamer, format, err := mp3.Decode(embeddedReader)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30))

	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}
	resampler := beep.ResampleRatio(4, 1, ctrl)
	volume := &effects.Volume{Streamer: resampler, Base: 2}

	read := make(chan bool)

	go func() {
		<-read
		speaker.Play(volume)
	}()

	first := true

	reader := bufio.NewReader(os.Stdin)
	for {
		var (
			text string
			err  error
		)
		go func() {

			text, err = reader.ReadString('\n')
			if first {
				first = false
				read <- true
			}
			read <- true
		}()

		<-read

		if errors.Is(err, io.EOF) {
			return
		}

		if err != nil {
			log.Fatal(err)
		}

		val, err := strconv.ParseFloat(strings.TrimSpace(text), 64)
		if err != nil {
			continue
		}

		speaker.Lock()
		resampler.SetRatio(math.Exp(val/50) - 0.7)
		volume.Volume = val / 100
		speaker.Unlock()
	}
}
