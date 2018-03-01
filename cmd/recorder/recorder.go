package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/gordonklaus/portaudio"
	youpywav "github.com/youpy/go-wav"
)

var (
	outFile = flag.String("out", "", "Output wav file")
)

func main() {
	flag.Parse()
	if *outFile == "" {
		fmt.Println("No output file defined")
		os.Exit(1)
	}
	recording := true
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	sampleRate := 512
	in := make([]int16, sampleRate)
	out, errC := os.Create(*outFile)
	if errC != nil {
		panic(errC)
	}
	defer out.Close()
	if err := portaudio.Initialize(); err != nil {
		panic(err)
	}
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	if err != nil {
		panic(err)
	}
	defer stream.Close()
	stream.Start()
	defer stream.Stop()

	go func() {
		for _ = range signalChan {

			recording = false
		}

	}()
	wr := youpywav.NewWriter(out, 0, 1, 44100, 16)
	fmt.Println("Recording...")
	for recording {
		if err := stream.Read(); err != nil {
			fmt.Println(err)
			break
		}
		var samples []youpywav.Sample
		for _, v := range in {
			var sample youpywav.Sample
			sample.Values[0] = int(v)
			samples = append(samples, sample)
		}
		wr.WriteSamples(samples)

	}
	fmt.Println("Done")

}
