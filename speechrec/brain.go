package speechrec

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"sync"
	"time"

	mind "github.com/stevenmiller888/go-mind"
)

var (
	trainRoutines = runtime.NumCPU()
	lock          = sync.Mutex{}
)

// Brain defines a brain
type Brain struct {
	M                  *mind.Mind
	Word               string
	Learn              chan []float64
	trainingInProgress bool
	trainingWg         sync.WaitGroup
}

// NewBrain returns a new brain
func NewBrain(word string) *Brain {
	b := &Brain{
		M:     mind.New(0.3, 10000, 10, "sigmoid"),
		Word:  word,
		Learn: make(chan []float64),
	}
	return b
}

func (b *Brain) trainWorker(id int, input <-chan []float64, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range input {
		fmt.Printf("Learning with routine %d\n", id)
		s := time.Now()

		n := normalizeInput(i)
		b.M.Learn([][][]float64{
			{n, []float64{1.0}},
		})
		e := time.Now().Sub(s)
		fmt.Printf("Learning done by routine %d in %.3fs\n", id, float64(e.Nanoseconds())/1000000000.0)
	}
}

// Run runs the brain
func (b *Brain) Run() {
	for i := 0; i < trainRoutines; i++ {
		b.trainingWg.Add(1)
		go b.trainWorker(i, b.Learn, &b.trainingWg)
	}
	// for i := 0; i < trainRoutines; i++ {
	// 	b.trainingWg.Add(1)
	// 	go func(i int) {
	// 		fmt.Println("Start train routine", i)
	// 		defer b.trainingWg.Done()
	// 		defer fmt.Println("Stop train routine", i)
	// 		for {
	// 			select {
	// 			case input, ok := <-b.Learn:
	// 				if !ok {
	// 					return
	// 				}
	// 				fmt.Printf("Learning with routine %d\n", i)
	// 				s := time.Now()
	//
	// 				n := normalizeInput(input)
	// 				b.M.Learn([][][]float64{
	// 					{n, []float64{1.0}},
	// 				})
	// 				e := time.Now().Sub(s)
	// 				fmt.Printf("Learning done by routine %d in %.3fs\n", i, float64(e.Nanoseconds())/1000000000.0)
	//
	// 			}
	//
	// 		}
	// }(i)
	// }

}

// Stop stops the brain
func (b *Brain) Stop() {
	close(b.Learn)
}

// Save saves the brain to binary files
func (b *Brain) Save() {
	// wait until trainning is finished
	fmt.Println("Waiting for training before saving")
	b.trainingWg.Wait()

	o, err := os.Create(fmt.Sprintf("HiddenOutput-%s.bin", b.Word))
	if err != nil {
		panic(err)
	}
	defer o.Close()
	if _, errM := b.M.Weights.HiddenOutput.MarshalBinaryTo(o); errM != nil {
		panic(err)
	}
	i, err := os.Create(fmt.Sprintf("InputHidden-%s.bin", b.Word))
	if err != nil {
		panic(err)
	}
	defer i.Close()
	if _, err := b.M.Weights.InputHidden.MarshalBinaryTo(i); err != nil {
		panic(err)
	}
}

// Load loads binary files to the brain
func (b *Brain) Load() {
	o, err := os.Open(fmt.Sprintf("HiddenOutput-%s.bin", b.Word))
	if err != nil {
		panic(err)
	}
	defer o.Close()
	b.M.Weights.HiddenOutput.UnmarshalBinaryFrom(o)
	i, err := os.Open(fmt.Sprintf("InputHidden-%s.bin", b.Word))
	if err != nil {
		panic(err)
	}
	defer i.Close()
	b.M.Weights.InputHidden.UnmarshalBinaryFrom(i)
}

// normalizeInput normalizes values to be between 0 and 1
func normalizeInput(reals []float64) []float64 {
	ret := make([]float64, len(reals))
	var maxV float64
	for _, x := range reals {
		if x > maxV {
			maxV = x
		}
	}
	for i, x := range reals {
		ret[i] = math.Abs(x / maxV)
	}
	return ret
}
