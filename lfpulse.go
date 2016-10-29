package main

import (
	"fmt"
	"github.com/scgolang/sc"
)

func main() {
	const synthName = "sc.LFPulseExample"

	client, err := sc.NewClient("udp", "127.0.0.1:57112", "127.0.0.1:57110")
	if err != nil {
		panic(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		panic(err)
	}
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		lfoFreq, lfoPhase, lfoWidth := sc.C(3), sc.C(0), sc.C(0.3)
		bus, gain := sc.C(0), sc.C(0.1)
		freq := sc.LFPulse{lfoFreq, lfoPhase, lfoWidth}.Rate(sc.KR).MulAdd(sc.C(200), sc.C(200))
		iphase, width := sc.C(0), sc.C(0.2)
		sig := sc.LFPulse{freq, iphase, width}.Rate(sc.AR).Mul(gain)
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
