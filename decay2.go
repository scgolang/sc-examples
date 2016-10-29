package main

import (
	"log"
	"time"

	"github.com/scgolang/sc"
)

func main() {
	const synthName = "Decay2Example"

	client, err := sc.NewClient("udp", "127.0.0.1:57110", "127.0.0.1:57120", 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		log.Fatal(err)
	}
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		bus := sc.C(0)
		line := sc.XLine{Start: sc.C(1), End: sc.C(50), Dur: sc.C(20)}.Rate(sc.KR)
		pulse := sc.Impulse{Freq: line, Phase: sc.C(0.25)}.Rate(sc.AR)
		sig := sc.Decay2{In: pulse, Attack: sc.C(0.01), Decay: sc.C(0.2)}.Rate(sc.AR)
		gain := sc.SinOsc{Freq: sc.C(600)}.Rate(sc.AR)
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
