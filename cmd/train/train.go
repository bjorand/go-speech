package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/bjorand/go-speech/drivers/audio"
	"github.com/bjorand/go-speech/speechrec"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
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

func xys(v []float64) plotter.XYs {
	pts := make(plotter.XYs, len(v))
	for i := range pts {
		pts[i].X = float64(i)
		pts[i].Y = v[i]
	}
	return pts
}

func main() {
	var testFileScore, wrongFileScore []float64
	flag.Parse()
	if *word == "" {
		fmt.Println("-word undefined")
		return
	}
	if *inFile == "" {
		fmt.Println("-in undefined")
		return
	}

	brain := speechrec.NewBrain(*word)
	brain.Run()

	if *testFile != "" {
		go func() {
			ticker := time.NewTicker(3 * time.Second).C
			for {
				select {
				case <-ticker:
					stream, err := audio.NewWavFileReader(*testFile)
					if err != nil {
						panic(err)
					}
					for {
						wordSamples, err := stream.WordSamples(avgAmplitudeThreshold)
						if err == io.EOF {
							break
						}
						if err != nil {
							panic(err)
						}
						for _, wordSample := range wordSamples {
							p := brain.M.Predict([][]float64{wordSample})
							testFileScore = append(testFileScore, p.RawMatrix().Data[0])

						}
					}
				}
			}
		}()
	}
	if *wrongTestFile != "" {
		go func() {
			ticker := time.NewTicker(3 * time.Second).C
			for {
				select {
				case <-ticker:
					stream, err := audio.NewWavFileReader(*wrongTestFile)
					if err != nil {
						panic(err)
					}
					for {
						wordSamples, err := stream.WordSamples(avgAmplitudeThreshold)
						if err == io.EOF {
							break
						}
						if err != nil {
							panic(err)
						}
						for _, wordSample := range wordSamples {
							p := brain.M.Predict([][]float64{wordSample})
							wrongFileScore = append(wrongFileScore, p.RawMatrix().Data[0])

						}
					}
				}
			}
		}()
	}

	stream, err := audio.NewWavFileReader(*inFile)
	if err != nil {
		panic(err)
	}
	defer stream.Close()
	for {
		wordSamples, errW := stream.WordSamples(avgAmplitudeThreshold)
		if errW == io.EOF {
			break
		}
		if errW != nil {
			panic(err)
		}
		fmt.Printf("Got word with %d sample to learn\n", len(wordSamples))
		for _, wordSample := range wordSamples {
			brain.Learn <- wordSample
		}
	}
	fmt.Println("Saving training results...")
	brain.Stop()
	brain.Save()
	fmt.Println("Training results saved.")
	rand.Seed(int64(0))

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Learning"
	p.X.Label.Text = "samples"
	p.Y.Label.Text = "signoid"

	err = plotutil.AddLinePoints(p,
		"Test file score", xys(testFileScore),
		"Wrong file score", xys(wrongFileScore))
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 10*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}
