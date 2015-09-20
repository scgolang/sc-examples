package main

import (
	"fmt"
	"github.com/scgolang/sc"
)

func main() {
	const synthName = "Balance2Example"
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
		bus, gain := C(0), C(0.1)
		l, r := LFSaw{Freq: C(44)}.Rate(AR), Pulse{Freq: C(33)}.Rate(AR)
		pos := FSinOsc{Freq: C(0.5)}.Rate(KR)
		sig := Balance2{L: l, R: r, Pos: pos, Level: gain}.Rate(AR)
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
