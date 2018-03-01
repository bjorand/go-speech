package main

import (
	"flag"
	"fmt"
	"math/cmplx"
	"os"
	"reflect"

	"github.com/bjorand/go-speech/speechrec"
	"github.com/mjibson/go-dsp/fft"
	dspwav "github.com/mjibson/go-dsp/wav"
)

var (
	inFile                = flag.String("in", "", "Input wav file")
	word                  = flag.String("word", "", "Word")
	testFile              = flag.String("test", "", "Input testwav file")
	wrongTestFile         = flag.String("wrong", "", "Input wrong testwav file")
	avgAmplitudeThreshold = 800.0
	wordMinDuration       = 1000 // ms
	window                = 20   // ms
)

func parse(t interface{}) []float64 {
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

func main() {
	flag.Parse()
	if *word == "" {
		fmt.Println("-word undefined")
		return
	}
	if *inFile == "" {
		fmt.Println("-in undefined")
		return
	}

	// lock := sync.Mutex{}
	brain := speechrec.NewBrain(*word)
	brain.Run()
	// defer brain.Stop()
	// time.Sleep(10 * time.Second)

	f, err := os.Open(*inFile)
	if err != nil {
		panic(fmt.Sprintf("couldn't open audio file - %v", err))
	}
	defer f.Close()

	w, err := dspwav.New(f)
	if err != nil {
		panic(err)
	}
	sampleRate := int(w.SampleRate)
	// go func() {
	var isWord bool
	// var train [][]float64
	// var previousSamples []wav.Sample
	// var samples []wav.Sample
	iSample := 1
	for {
		samples, err := w.ReadSamples(sampleRate / 1000 * window)
		if err != nil {
			break
		}
		// if len(previousSamples) > 0 {
		// 	copy(samples, previousSamples[800:])
		// 	samples = append(samples, samplesR...)
		// } else {
		// 	copy(samples, samplesR)
		// }

		// copy(previousSamples, samples)
		// numSamples = numSamples + len(samples)
		s := parse(samples)
		// s := make([]float64, len(samples))
		// for i, sample := range samples {
		// 	s[i] = float64(sample.Values[0])
		//
		// }

		f := fft.FFTReal(s)
		if len(f) < 1024 {
			for i := len(f); i < 1024; i++ {
				var c complex128
				f = append(f, c)
			}
		}
		// ensure we have at least 1000 items in fftr

		// lowPassFilter := 21
		// hiPassFilter := 20000
		// for i, c := range f {
		// 	if int(math.Abs(real(c))) < lowPassFilter {
		// 		f[i] = 0
		// 	}
		// 	if int(math.Abs(real(c))) > hiPassFilter {
		// 		f[i] = 0
		// 	}
		// }
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
				isWord = false
			}
		}
		if isWord {
			iSample++
			brain.Learn <- reals

		}
	}
	fmt.Println("saving")
	brain.Stop()
	brain.Save()
	fmt.Println("saved")
	// time.Sleep(2 * time.Second)
	// }()

	// fmt.Println(numSamples, len(train))
	// for i, x := range train {
	// 	fmt.Printf("train %d/%d\n", i+1, len(train))
	// 	m.Learn([][][]float64{
	// 		{x, []float64{0.1}},
	// 	})
	// }

	// ////////////////////
	// go func() {
	// 	var isWord bool
	// 	for {
	// 		f2, err := os.Open(*testFile)
	// 		if err != nil {
	// 			panic(fmt.Sprintf("couldn't open audio file - %v", err))
	// 		}
	// 		defer f2.Close()
	//
	// 		w2, err := dspwav.New(f2)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	//
	// 		sampleRate2 := 44100
	//
	// 		isWord = false
	// 		for {
	// 			samples2, err := w2.ReadSamples(sampleRate2 / 1000 * window)
	// 			if err == io.EOF {
	// 				break
	// 			}
	// 			if err != nil {
	// 				break
	// 			}
	// 			s2 := parse(samples2)
	//
	// 			f2 := fft.FFTReal(s2)
	// 			if len(f2) < 1024 {
	// 				for i := len(f2); i < 1024; i++ {
	// 					var c complex128
	// 					f2 = append(f2, c)
	// 				}
	// 			}
	//
	// 			// lowPassFilter := 21
	// 			// hiPassFilter := 20000
	// 			// for i, c := range f2 {
	// 			// 	if int(math.Abs(real(c))) < lowPassFilter {
	// 			// 		f2[i] = 0
	// 			// 	}
	// 			// 	if int(math.Abs(real(c))) > hiPassFilter {
	// 			// 		f2[i] = 0
	// 			// 	}
	// 			// }
	// 			abs2 := make([]float64, len(f2))
	// 			reals2 := make([]float64, len(f2))
	// 			for i, x := range f2 {
	// 				abs2[i] = cmplx.Abs(x)
	// 				reals2[i] = real(x)
	// 			}
	// 			var avg float64
	// 			var sum float64
	// 			for _, x := range abs2 {
	// 				sum = sum + x
	// 			}
	// 			avg = sum / float64(len(abs2))
	// 			if avg > avgAmplitudeThreshold {
	// 				if !isWord {
	// 					isWord = true
	// 					fmt.Println("word")
	// 				}
	// 			} else {
	// 				isWord = false
	// 			}
	// 			if isWord {
	// 				// we have to train network with f value
	// 				// fmt.Println("OOOOOO", realsNormalizerd(reals))
	// 				n2 := realsNormalizerd(reals2)
	// 				fmt.Println(n2)
	// 				fmt.Println("TEST", m.Predict([][]float64{n2}))
	// 			}
	//
	// 		}
	// 		time.Sleep(3 * time.Second)
	//
	// 	}
	// }()
	// // ////////////////////
	// var isWord bool
	// for {
	// 	f3, err := os.Open(*wrongTestFile)
	// 	if err != nil {
	// 		panic(fmt.Sprintf("couldn't open audio file - %v", err))
	// 	}
	// 	defer f3.Close()
	//
	// 	w3, err := dspwav.New(f3)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	sampleRate3 := 44100
	//
	// 	isWord = false
	// 	for {
	// 		samples3, err := w3.ReadSamples(sampleRate3 / 1000 * window)
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		if err != nil {
	// 			break
	// 		}
	//
	// 		s3 := parse(samples3)
	// 		f := fft.FFTReal(s3)
	// 		if len(f) < 1024 {
	// 			for i := len(f); i < 1024; i++ {
	// 				var c complex128
	// 				f = append(f, c)
	// 			}
	// 		}
	// 		// lowPassFilter := 21
	// 		// hiPassFilter := 20000
	// 		// for i, c := range f {
	// 		// 	if int(math.Abs(real(c))) < lowPassFilter {
	// 		// 		f[i] = 0
	// 		// 	}
	// 		// 	if int(math.Abs(real(c))) > hiPassFilter {
	// 		// 		f[i] = 0
	// 		// 	}
	// 		// }
	// 		abs := make([]float64, len(f))
	// 		reals := make([]float64, len(f))
	// 		for i, x := range f {
	// 			abs[i] = cmplx.Abs(x)
	// 			reals[i] = real(x)
	// 		}
	// 		var avg float64
	// 		var sum float64
	// 		for _, x := range abs {
	// 			sum = sum + x
	// 		}
	// 		avg = sum / float64(len(abs))
	// 		if avg > avgAmplitudeThreshold {
	// 			if !isWord {
	// 				isWord = true
	// 				fmt.Println("word")
	// 			}
	// 		} else {
	// 			isWord = false
	// 		}
	// 		if isWord {
	// 			// we have to train network with f value
	// 			n3 := realsNormalizerd(reals)
	// 			fmt.Println(n3)
	// 			fmt.Println("WRONG", m.Predict([][]float64{n3}))
	// 		}
	//
	// 	}
	// 	time.Sleep(3 * time.Second)
	//
	// }

}
