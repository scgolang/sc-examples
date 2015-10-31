package main

import (
	"math/rand"
	"time"

	"github.com/scgolang/sc"
)

func main() {
	const synthName = "playbufExample"
	const wavFile = "kalimba_mono.wav"
	var synthID int32
	var speed, gain float32

	// setup supercollider client
	client := sc.NewClient("127.0.0.1:57111")
	err := client.Connect("127.0.0.1:57110")
	if err != nil {
		panic(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		panic(err)
	}
	// read a wav file
	buf, err := client.ReadBuffer(wavFile)
	if err != nil {
		panic(err)
	}
	// send a synthdef
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		speed := p.Add("speed", 1.0)
		gain := p.Add("gain", 0.5)
		bus := sc.C(0)
		sig := sc.PlayBuf{
			NumChannels: 1,
			BufNum:      sc.C(float32(buf.Num())),
			Speed:       speed,
			Done:        sc.FreeEnclosing,
		}.Rate(sc.AR).Mul(gain)
		return sc.Out{bus, sig}.Rate(sc.AR)
	})
	err = client.SendDef(def)
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(125 * time.Millisecond)

	for _ = range ticker.C {
		synthID = client.NextSynthID()
		speed = (float32(8.0) * rand.Float32()) + float32(0.5)
		gain = rand.Float32()
		ctls := map[string]float32{"speed": speed, "gain": gain}
		_, err = defaultGroup.Synth(synthName, synthID, sc.AddToTail, ctls)
	}
}
