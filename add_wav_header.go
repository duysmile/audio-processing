package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/mattetti/filebuffer"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	storeInMemAudio()
	storeFileWav()
}

func storeInMemAudio() io.Reader {
	pcm, err := os.Open("audio/audio.pcm")
	if err != nil {
		log.Fatal(err)
	}

	out := filebuffer.New(nil)

	// audio format: 1 <-> PCM
	e := wav.NewEncoder(out, 8000, 16, 1, 1)

	audioBuf, err := newAudioIntFromBuffer(pcm)
	if err != nil {
		log.Fatal(err)
	}

	if err := e.Write(audioBuf); err != nil {
		log.Fatal(err)
	}

	out.Seek(0, 0)
	return bytes.NewReader(out.Bytes())
}

func storeFileWav() io.Reader {
	pcm, err := os.Open("audio/audio.pcm")
	if err != nil {
		log.Fatal(err)
	}

	// generated file will be 21324254.wav
	out, err := ioutil.TempFile("audio", "*.wav")
	if err != nil {
		log.Fatal(err)
	}

	// comment this block if you want to store .wav file
	defer func() {
		if err := os.RemoveAll(out.Name()); err != nil {
			log.Fatal(err)
		}
	}()

	// audio format: 1 <-> PCM
	e := wav.NewEncoder(out, 8000, 16, 1, 1)

	audioBuf, err := newAudioIntFromBuffer(pcm)
	if err != nil {
		log.Fatal(err)
	}

	if err := e.Write(audioBuf); err != nil {
		log.Fatal(err)
	}

	// flush to disk if you want to store file
	if err := e.Close(); err != nil {
		log.Fatal(err)
	}

	content, _ := ioutil.ReadFile(out.Name())
	fmt.Println("audio: ", string(content))

	// can use temp file out for any purposes
	// this file is .wav file with following format
	// sample rate: 8000
	// bit depth: 16
	// num of channels: 1

	// Use this reader to read audio
	return bufio.NewReader(out)
}

func newAudioIntFromBuffer(r io.Reader) (*audio.IntBuffer, error) {
	buf := audio.IntBuffer{
		// custom format for your case
		Format: &audio.Format{
			NumChannels: 1,
			SampleRate:  8000,
		},
	}

	for {
		var sample int16
		err := binary.Read(r, binary.LittleEndian, &sample)
		switch {
		case err == io.EOF:
			return &buf, nil
		case err != nil:
			return nil, err
		}

		buf.Data = append(buf.Data, int(sample))
	}
}
