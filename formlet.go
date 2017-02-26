package main

import (
	"log"
	"time"

	"github.com/scgolang/sc"
	"github.com/scgolang/scids/scid"
)

func main() {
	const synthName = "FormletExample"

	client, err := sc.NewClient("udp", "127.0.0.1:57110", "127.0.0.1:57120", 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		log.Fatal(err)
	}
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		bus, sine := sc.C(0), sc.SinOsc{Freq: sc.C(5)}.Rate(sc.KR).MulAdd(sc.C(20), sc.C(300))
		blip := sc.Blip{Freq: sine, Harm: sc.C(1000)}.Rate(sc.AR).Mul(sc.C(0.1))
		line := sc.XLine{Start: sc.C(1500), End: sc.C(700), Dur: sc.C(8)}.Rate(sc.KR)
		sig := sc.Formlet{
			In:         blip,
			Freq:       line,
			AttackTime: sc.C(0.005),
			DecayTime:  sc.C(0.4),
		}.Rate(sc.AR)
		sig = sc.Multi(sig, sig)
		return sc.Out{bus, sig}.Rate(sc.AR)
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
