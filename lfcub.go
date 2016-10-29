package main

import (
	"fmt"
	"github.com/scgolang/sc"
)

func main() {
	const synthName = "sc.LFCubExample"

	client, err := sc.NewClient("udp", "127.0.0.1:57112", "127.0.0.1:57110")
	if err != nil {
		panic(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		panic(err)
	}
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		bus, gain := sc.C(0), sc.C(0.1)
		lfo1 := sc.LFCub{Freq: sc.C(0.2)}.Rate(sc.KR).MulAdd(sc.C(8), sc.C(10))
		lfo2 := sc.LFCub{Freq: lfo1}.Rate(sc.KR).MulAdd(sc.C(400), sc.C(800))
		sig := sc.LFCub{Freq: lfo2}.Rate(sc.AR).Mul(gain)
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
