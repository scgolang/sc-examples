package main

import (
	"github.com/scgolang/sc"
)

func main() {
	const synthName = "CombCExample"

	client, err := sc.NewClient("udp", "127.0.0.1:57112", "127.0.0.1:57110")
	if err != nil {
		panic(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		panic(err)
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
		return sc.Out{bus, sig}.Rate(sc.AR)
	})
	err = client.SendDef(def)
	if err != nil {
		panic(err)
	}
	synthID := client.NextSynthID()
	_, err = defaultGroup.Synth(synthName, synthID, sc.AddToTail, nil)
}
