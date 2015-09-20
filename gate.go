package main

import (
	"fmt"
	"github.com/scgolang/sc"
)

func main() {
	const synthName = "GateExample"
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
		bus, noise := C(0), WhiteNoise{}.Rate(KR)
		pulse := LFPulse{Freq: C(1.333), Iphase: C(0.5)}.Rate(KR)
		sig := Gate{In: noise, Trig: pulse}.Rate(AR)
		return Out{bus, sig}.Rate(AR)
	})
	err = client.SendDef(def)
	if err != nil {
		panic(err)
	}
	synthID := client.NextSynthID()
	_, err = defaultGroup.Synth(synthName, synthID, AddToTail, nil)
	fmt.Printf("created synth %d\n", synthID)
}
