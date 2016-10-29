package main

import (
	"fmt"

	"github.com/scgolang/sc"
)

func main() {
	const synthName = "GateExample"

	client, err := sc.NewClient("udp", "127.0.0.1:57112", "127.0.0.1:57110")
	if err != nil {
		panic(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		panic(err)
	}
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		bus, noise := sc.C(0), sc.WhiteNoise{}.Rate(sc.KR)
		pulse := sc.LFPulse{Freq: sc.C(1.333), Iphase: sc.C(0.5)}.Rate(sc.KR)
		sig := sc.Gate{In: noise, Trig: pulse}.Rate(sc.AR)
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
