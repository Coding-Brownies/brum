package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
	"time"
	"unicode"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"golang.org/x/term"
)

func main() {
	// set term state
	var fd int
	if term.IsTerminal(syscall.Stdin) {
		fd = syscall.Stdin
	} else {
		tty, err := os.Open("/dev/tty")
		if err != nil {
			panic(err)
		}
		defer tty.Close()
		fd = int(tty.Fd())
	}
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// music shit

	f, err := os.Open("audio/click-2.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)

	for {
		var err error
		b := make([]byte, 1)

		read := make(chan bool)
		go func() {
			_, err = os.Stdin.Read(b)
			if b[0] == byte(4) {
				err = io.EOF
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

		shot := buffer.Streamer(0, buffer.Len())
		go speaker.Play(shot)

		if !unicode.IsPrint(rune(b[0])) {
			continue
		}

		fmt.Print(string(b[0]))
	}

}
