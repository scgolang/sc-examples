package main

import (
	"github.com/scgolang/sc"
)

func main() {
	const synthName = "CombCExample"
	client := sc.NewClient("127.0.0.1:57112")
	err := client.Connect("127.0.0.1:57110")
	if err != nil {
		panic(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		panic(err)
	}
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		bus := C(0)
		line := XLine{
			Start: C(0.0001),
			End:   C(0.01),
			Dur:   C(20),
		}.Rate(KR)
		sig := CombC{
			In:           WhiteNoise{}.Rate(AR).Mul(C(0.01)),
			MaxDelayTime: C(0.01),
			DelayTime:    line,
			DecayTime:    C(0.2),
		}.Rate(AR)
		return Out{bus, sig}.Rate(AR)
	})
	err = client.SendDef(def)
	if err != nil {
		panic(err)
	}
	synthID := client.NextSynthID()
	_, err = defaultGroup.Synth(synthName, synthID, AddToTail, nil)
}
