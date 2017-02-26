package main

import (
	"log"
	"time"

	"github.com/scgolang/sc"
	"github.com/scgolang/scids/scid"
)

func main() {
	const synthName = "CombCExample"

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
		line := sc.XLine{
			Start: sc.C(0.0001),
			End:   sc.C(0.01),
			Dur:   sc.C(20),
		}.Rate(sc.KR)
		sig := sc.Comb{
			Interpolation: sc.InterpolationCubic,
			In:            sc.WhiteNoise{}.Rate(sc.AR).Mul(sc.C(0.01)),
			MaxDelayTime:  sc.C(0.01),
			DelayTime:     line,
			DecayTime:     sc.C(0.2),
		}.Rate(sc.AR)
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
}
