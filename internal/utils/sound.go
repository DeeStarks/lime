package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

const (
	SoundFormatWAV = "wav"
	SoundFormatMP3 = "mp3"
)

func PlaySound(soundUrl string) {
	s := strings.Split(soundUrl, ".")
	soundFormat := s[len(s)-1] // Get file extension

	// Add absolute path to the sound file
	_, b, _, _ := runtime.Caller(0)
	dir := filepath.Dir(b)
	// Replace last two directories with "assets/sounds"
	dirs := strings.Split(dir, "/")
	dir = strings.Join(dirs[:len(dirs)-2], "/") + "/assets/sounds/"
	// Add the sound file name
	soundUrl = dir + soundUrl

	f, err := os.Open(soundUrl)
	if err != nil {
		LogMessage("Error opening sound file: %v", err)
	}

	var (
		streamer beep.StreamSeekCloser
		format   beep.Format
	)
	switch soundFormat {
	case SoundFormatMP3:
		streamer, format, err = mp3.Decode(f)
	case SoundFormatWAV:
		streamer, format, err = wav.Decode(f)
	default:
		LogMessage("Sound format not supported. Supported formats: %v, %v", SoundFormatMP3, SoundFormatWAV)
	}

	if err != nil {
		LogMessage("Error decoding sound file: %v", err)
	}
	if streamer != nil {
		defer streamer.Close()
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(streamer)
}
