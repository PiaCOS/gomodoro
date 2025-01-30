package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
)

//go:embed sounds
var staticFS embed.FS

var (
	workFlag      int
	restFlag      int
	muteFlag      bool
	startWorkPath string = "sounds/gomodoro_01.mp3"
	stopWorkPath  string = "sounds/gomodoro_02.mp3"
)

func main() {
	flag.IntVar(&workFlag, "w", 25, "Time in minutes for the work session")
	flag.IntVar(&restFlag, "r", 5, "Time in minutes for the rest session")
	flag.BoolVar(&muteFlag, "m", false, "Mute mode (-m=true)")
	flag.Parse()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	workTimer := workFlag
	restTimer := restFlag
	workPlayer := NewSoundPlayer(startWorkPath)
	restPlayer := NewSoundPlayer(stopWorkPath)
	workPlayer.StartSpeaker()

	go func() {
		<-exit
		fmt.Println("\nExiting Gomodoro Timer, bye bye nya~! :3")
		workPlayer.Streamer.Close()
		restPlayer.Streamer.Close()
		os.Exit(0)
	}()

	t := time.Now()
	fmt.Printf("Let's start our session :3\n")
	fmt.Printf("It's %v\n", t.Format("15:04:05"))
	fmt.Printf("We will work for %v minutes and then rest for %v minutes !!\n", workTimer, restTimer)

	timer := 0
	isPaused := false
	go workPlayer.PlaySound()

	for {
		time.Sleep(1 * time.Minute)
		timer++

		t := time.Now()
		if !isPaused && timer == workTimer {
			fmt.Printf("%v -- TIME TO TAKE A BREAK\n", t.Format("15:04:05"))
			isPaused = true
			timer = 0
			go restPlayer.PlaySound()
		}
		if isPaused && timer == restTimer {
			fmt.Printf("%v -- TIME TO GO BACK TO WORK\n", t.Format("15:04:05"))
			isPaused = false
			timer = 0
			go workPlayer.PlaySound()
		}
	}
}

// *bytes.Reader does not implement io.ReadCloser (missing method Close)
type ReadCloser struct {
	*bytes.Reader
}

func (rc ReadCloser) Close() error {
	return nil
}

type SoundPlayer struct {
	Streamer beep.StreamSeekCloser
	Format   beep.Format
}

func (s SoundPlayer) StartSpeaker() {
	if !muteFlag {
		speaker.Init(s.Format.SampleRate, s.Format.SampleRate.N(time.Second/10))
	}
}

func (s SoundPlayer) PlaySound() {
	if !muteFlag {
		s.Streamer.Seek(0)

		done := make(chan bool)
		speaker.Play(beep.Seq(s.Streamer, beep.Callback(func() {
			done <- true
		})))
		<-done
	}
}

func NewSoundPlayer(path string) SoundPlayer {
	data, err := staticFS.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	reader := &ReadCloser{bytes.NewReader(data)}

	streamer, format, err := mp3.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	return SoundPlayer{streamer, format}
}
