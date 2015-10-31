package main

import (
	"fmt"

	"github.com/scgolang/sc"
)

func main() {
	const synthName = "FreeVerbExample"
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
		mix := p.Add("mix", 0.25)
		room := p.Add("room", 0.15)
		damp := p.Add("damp", 0.5)
		bus := sc.C(0)
		impulse := sc.Impulse{Freq: sc.C(1)}.Rate(sc.AR)
		lfcub := sc.LFCub{Freq: sc.C(1200), Iphase: sc.C(0)}.Rate(sc.AR).Mul(sc.C(0.1))
		decay := sc.Decay{In: impulse, Decay: sc.C(0.25)}.Rate(sc.AR).Mul(lfcub)
		sig := sc.FreeVerb{In: decay, Mix: mix, Room: room, Damp: damp}.Rate(sc.AR)
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
