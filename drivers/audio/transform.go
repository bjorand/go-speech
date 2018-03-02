package audio

import (
	"fmt"
	"io"
	"math/cmplx"
	"reflect"

	"github.com/mjibson/go-dsp/fft"
)

func parseSamples(t interface{}) []float64 {
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)
		r := make([]float64, s.Len())
		for i := 0; i < s.Len(); i++ {
			r[i] = float64(s.Index(i).Int())
		}
		return r
	default:
		panic("unknown interface type")
	}
}

// WordSamples returns samples of a word, call multiple time to get next word
func (ws *WavStream) WordSamples(avgAmplitudeThreshold float64) ([][]float64, error) {
	window := 20 // 20ms
	var wordSamples [][]float64
	var isWord bool
	for {
		samples, err := ws.Wav.ReadSamples(ws.SampleRate / 1000 * window)
		if err != nil {
			return nil, io.EOF
		}
		s := parseSamples(samples)
		f := fft.FFTReal(s)
		if len(f) < 1024 {
			for i := len(f); i < 1024; i++ {
				var c complex128
				f = append(f, c)
			}
		}
		abs := make([]float64, len(f))
		reals := make([]float64, len(f))
		for i, x := range f {
			abs[i] = cmplx.Abs(x)
			reals[i] = real(x)
		}
		var avg float64
		var sum float64
		for _, x := range abs {
			sum = sum + x
		}
		avg = sum / float64(len(abs))
		if avg > avgAmplitudeThreshold {
			if !isWord {
				isWord = true
				fmt.Println("word")
			}
		} else {
			if isWord {
				fmt.Println("end word")
				return wordSamples, nil
			}
		}
		if isWord {
			wordSamples = append(wordSamples, reals)
		}
	}
}
