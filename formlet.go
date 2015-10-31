package main

import (
	"fmt"

	"github.com/scgolang/sc"
)

func main() {
	const synthName = "FormletExample"
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
		bus, sine := sc.C(0), sc.SinOsc{Freq: sc.C(5)}.Rate(sc.KR).MulAdd(sc.C(20), sc.C(300))
		blip := sc.Blip{Freq: sine, Harm: sc.C(1000)}.Rate(sc.AR).Mul(sc.C(0.1))
		line := sc.XLine{Start: sc.C(1500), End: sc.C(700), Dur: sc.C(8)}.Rate(sc.KR)
		sig := sc.Formlet{
			In:         blip,
			Freq:       line,
			AttackTime: sc.C(0.005),
			DecayTime:  sc.C(0.4),
		}.Rate(sc.AR)
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
