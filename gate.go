package main

import (
	"log"
	"time"

	"github.com/scgolang/sc"
)

func main() {
	const synthName = "GateExample"

	client, err := sc.NewClient("udp", "127.0.0.1:57110", "127.0.0.1:57120", 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		log.Fatal(err)
	}
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		bus, noise := sc.C(0), sc.WhiteNoise{}.Rate(sc.KR)
		pulse := sc.LFPulse{Freq: sc.C(1.333), Iphase: sc.C(0.5)}.Rate(sc.KR)
		sig := sc.Gate{In: noise, Trig: pulse}.Rate(sc.AR)
		return sc.Out{bus, sig}.Rate(sc.AR)
	})
	err = client.SendDef(def)
	if err != nil {
		log.Fatal(err)
	}
	synthID := client.NextSynthID()
	_, err = defaultGroup.Synth(synthName, synthID, sc.AddToTail, nil)
	log.Printf("created synth %d\n", synthID)
}
