package audio

import (
	"fmt"
	"os"

	dspwav "github.com/mjibson/go-dsp/wav"
)

// WavStream is a wav stream
type WavStream struct {
	SampleRate  int
	Wav         *dspwav.Wav
	f           *os.File
	WordSampleC chan []float64
}

// NewWavFileReader opens a wavStream from a wav file
func NewWavFileReader(filename string) (*WavStream, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("couldn't open file - %v", err)
	}
	w, err := dspwav.New(f)
	if err != nil {
		panic(err)
	}

	ws := &WavStream{
		SampleRate:  int(w.SampleRate),
		Wav:         w,
		f:           f,
		WordSampleC: make(chan []float64),
	}
	return ws, nil
}

// Close closes the stream
func (ws *WavStream) Close() {
	ws.f.Close()
}
