package main

import (
	"fmt"

	"github.com/scgolang/sc"
)

func main() {
	const synthName = "DelayCExample"

	client, err := sc.NewClient("udp", "127.0.0.1:57112", "127.0.0.1:57110")
	if err != nil {
		panic(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		panic(err)
	}
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		bus, dust := sc.C(0), sc.Dust{Density: sc.C(1)}.Rate(sc.AR).Mul(sc.C(0.5))
		noise := sc.WhiteNoise{}.Rate(sc.AR)
		decay := sc.Decay{In: dust, Decay: sc.C(0.3)}.Rate(sc.AR).Mul(noise)
		sig := sc.Delay{
			Interpolation: sc.InterpolationCubic,
			In:            decay,
			MaxDelayTime:  sc.C(0.2),
			DelayTime:     sc.C(0.2),
		}.Rate(sc.AR).Add(decay)
		return sc.Out{bus, sig}.Rate(sc.AR)
	})
	err = client.SendDef(def)
	if err != nil {
		panic(err)
	}
	synthID := client.NextSynthID()
	_, err = defaultGroup.Synth(synthName, synthID, sc.AddToTail, nil)
	fmt.Printf("created synth %d\n", synthID)
}
