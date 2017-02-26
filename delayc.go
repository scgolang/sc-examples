package main

import (
	"log"
	"time"

	"github.com/scgolang/sc"
	"github.com/scgolang/scids/scid"
)

func main() {
	const synthName = "DelayCExample"

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
			bus   = sc.C(0)
			dust  = sc.Dust{Density: sc.C(1)}.Rate(sc.AR).Mul(sc.C(0.5))
			noise = sc.WhiteNoise{}.Rate(sc.AR)
			decay = sc.Decay{In: dust, Decay: sc.C(0.3)}.Rate(sc.AR).Mul(noise)
		)
		sig := sc.Delay{
			Interpolation: sc.InterpolationCubic,
			In:            decay,
			MaxDelayTime:  sc.C(0.2),
			DelayTime:     sc.C(0.2),
		}.Rate(sc.AR).Add(decay)
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
	if _, err = defaultGroup.Synth(synthName, synthID, sc.AddToTail, nil); err != nil {
		log.Fatal(err)
	}
	log.Printf("created synth %d\n", synthID)
}
