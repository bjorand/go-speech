# go-speech

go-speech is proof of concept I'm writing to play with neural networks and digital signal processing.
Current goal is to do simple speech recognition for words like "yes" or "no".

Working parts are:
 - Wav recorder to generate your own corpus
 - detect words in Wav file
 - split words in 20ms samples
 - train network with Fast Fourrier transform data



## Usage

`go-speech` uses `dep` for dependencies management:
```
go get github.com/golang/dep/cmd/dep
```
Install source code:
```
go get github.com/bjorand/go-speech
cd $GOPATH/src/bjorand/go-speech
```

Install dependencies:
```
dep ensure -v
```


Record a Wav file where your repeat the same word "hello" many times separated by short silences:
```
go run cmd/recorder/recorder.go -out hellos.wav
```
Press Ctrl-c to stop and save recording.

You can also record a separate and unique occurrence of "hello" to test our network while training runs.

```
go run cmd/recorder/recorder.go -out hello-test.wav
```

For testing purpose during training, record a new file containing a different word like "crazy":

```
go run cmd/recorder/recorder.go -out crazy.wav
```

Train neural network and save results:
```
go run cmd/train/train.go -word hello -in hellos.wav -test hello-test.wav -wrong crazy.wav
```
## TODO

- cmd/recorder: start recording when a word is detected.
- filter signal (21 to 9000hz cover human voice)
