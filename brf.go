package main

import (
	"fmt"

	"github.com/scgolang/sc"
)

func main() {
	const synthName = "BRFExample"
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
		line := sc.XLine{sc.C(0.7), sc.C(300), sc.C(20), 0}.Rate(sc.KR)
		saw := sc.Saw{sc.C(200)}.Rate(sc.AR).Mul(sc.C(0.5))
		sine := sc.SinOsc{line, sc.C(0)}.Rate(sc.KR).MulAdd(sc.C(3800), sc.C(4000))
		bpf := sc.BRF{saw, sine, sc.C(0.3)}.Rate(sc.AR)
		return sc.Out{sc.C(0), bpf}.Rate(sc.AR)
	})
	err = client.SendDef(def)
	if err != nil {
		panic(err)
	}
	synthID := client.NextSynthID()
	_, err = defaultGroup.Synth(synthName, synthID, sc.AddToTail, nil)
	fmt.Printf("created synth %d\n", synthID)
}
