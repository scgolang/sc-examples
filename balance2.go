package main

import (
	"fmt"

	"github.com/scgolang/sc"
)

func main() {
	const synthName = "Balance2Example"

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
		l, r := sc.LFSaw{Freq: sc.C(44)}.Rate(sc.AR), sc.Pulse{Freq: sc.C(33)}.Rate(sc.AR)
		pos := sc.SinOsc{Freq: sc.C(0.5)}.Rate(sc.KR)
		sig := sc.Balance2{L: l, R: r, Pos: pos, Level: gain}.Rate(sc.AR)
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
