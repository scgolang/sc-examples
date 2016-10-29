package main

import (
	"log"
	"time"

	"github.com/scgolang/sc"
)

func main() {
	const synthName = "COscExample"

	client, err := sc.NewClient("udp", "127.0.0.1:57110", "127.0.0.1:57120", 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		log.Fatal(err)
	}
	buf, err := client.AllocBuffer(512, 1)
	if err != nil {
		log.Fatal(err)
	}

	var (
		bufRoutine = sc.BufferRoutineSine1
		bufFlags   = sc.BufferFlagNormalize | sc.BufferFlagWavetable | sc.BufferFlagClear
		partials   = []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	)
	for i, p := range partials {
		partials[i] = 1 / p
	}

	if err := buf.Gen(bufRoutine, bufFlags, partials...); err != nil {
		log.Fatal(err)
	}
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		bus, gain := sc.C(0), sc.C(0.25)
		freq, beats := sc.C(200), sc.C(0.7)
		sig := sc.COsc{
			BufNum: sc.C(float32(buf.Num)),
			Freq:   freq,
			Beats:  beats,
		}.Rate(sc.AR)
		return sc.Out{bus, sig.Mul(gain)}.Rate(sc.AR)
	})
	if err := client.SendDef(def); err != nil {
		log.Fatal(err)
	}

	synthID := client.NextSynthID()
	if _, err := defaultGroup.Synth(synthName, synthID, sc.AddToTail, nil); err != nil {
		log.Fatal(err)
	}
}
