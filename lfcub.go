package main

import (
	"log"
	"time"

	"github.com/scgolang/sc"
	"github.com/scgolang/scids/scid"
)

func main() {
	const synthName = "sc.LFCubExample"

	client, err := sc.NewClient("udp", "127.0.0.1:57110", "127.0.0.1:57120", 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		log.Fatal(err)
	}
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		var (
			bus  = sc.C(0)
			gain = sc.C(0.1)
			lfo1 = sc.LFCub{Freq: sc.C(0.2)}.Rate(sc.KR).MulAdd(sc.C(8), sc.C(10))
			lfo2 = sc.LFCub{Freq: lfo1}.Rate(sc.KR).MulAdd(sc.C(400), sc.C(800))
			sig  = sc.LFCub{Freq: lfo2}.Rate(sc.AR).Mul(gain)
		)
		return sc.Out{bus, sc.Multi(sig, sig)}.Rate(sc.AR)
	})
	if err := client.SendDef(def); err != nil {
		log.Fatal(err)
	}
	synthID, err := scid.Next()
	if err != nil {
		log.Fatal(err)
	}
	if _, err := defaultGroup.Synth(synthName, synthID, sc.AddToTail, nil); err != nil {
		log.Fatal(err)
	}
	log.Printf("created synth %d\n", synthID)
}
