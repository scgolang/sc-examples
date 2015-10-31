package main

import (
	"math/rand"
	"time"

	"github.com/scgolang/sc"
)

func main() {
	const synthName = "sineTone"
	var synthID int32
	var note int
	var gain, dur float32

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
	err = client.DumpOSC(int32(1))
	if err != nil {
		panic(err)
	}
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		freq := p.Add("freq", 440)
		gain := p.Add("gain", 0.5)
		dur := p.Add("dur", 1)
		bus := sc.C(0)
		env := sc.EnvGen{
			Env:        sc.EnvPerc{Release: dur},
			LevelScale: gain,
			Done:       sc.FreeEnclosing,
		}.Rate(sc.KR)
		sig := sc.SinOsc{Freq: freq}.Rate(sc.AR).Mul(env)
		return sc.Out{bus, sig}.Rate(sc.AR)
	})
	err = client.SendDef(def)
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(125 * time.Millisecond)

	for _ = range ticker.C {
		synthID = client.NextSynthID()
		note = rand.Intn(128)
		gain = rand.Float32()
		dur = rand.Float32()
		ctls := map[string]float32{
			"freq": sc.Midicps(note),
			"gain": gain,
			"dur":  dur,
		}
		_, err = defaultGroup.Synth(synthName, synthID, sc.AddToTail, ctls)
	}
}
