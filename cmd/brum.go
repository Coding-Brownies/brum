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

	dest := 0.0
	actual := 0.0
	go func() {
		for {
			time.Sleep(1 * time.Millisecond)

			if math.Abs(dest-actual) < 0.1 {
				continue
			}

			if dest > actual {
				actual += 0.1
			} else {
				actual -= 0.1
			}

			speaker.Lock()
			resampler.SetRatio(math.Exp(actual/40) - 0.7)
			volume.Volume = math.Exp(actual/30) - 1
			speaker.Unlock()
		}
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

		dest, err = strconv.ParseFloat(strings.TrimSpace(text), 64)
		if err != nil {
			continue
		}
	}
}
