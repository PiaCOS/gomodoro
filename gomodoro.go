package main

import (
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

var (
	workFlag      int
	restFlag      int
	startWorkPath string = "gomodoro_01.mp3"
	stopWorkPath  string = "gomodoro_02.mp3"
)

func main() {
	flag.IntVar(&workFlag, "w", 25, "Time in minutes for the work session")
	flag.IntVar(&restFlag, "r", 5, "Time in minutes for the rest session")
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

	fmt.Printf("Let's start our session :3\n")
	fmt.Printf("We will work for %v minutes and then rest for %v minutes !!\n", workTimer, restTimer)

	timer := 0
	isPaused := false
	go workPlayer.PlaySound()

	for {
		time.Sleep(1 * time.Second)
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

type SoundPlayer struct {
	Streamer beep.StreamSeekCloser
	Format   beep.Format
}

func (s SoundPlayer) StartSpeaker() {
	speaker.Init(s.Format.SampleRate, s.Format.SampleRate.N(time.Second/10))
}

func (s SoundPlayer) PlaySound() {
	s.Streamer.Seek(0)

	done := make(chan bool)
	speaker.Play(beep.Seq(s.Streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}

func NewSoundPlayer(path string) SoundPlayer {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	return SoundPlayer{streamer, format}
}
