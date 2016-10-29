package main

import (
	"log"
	"time"

	"github.com/scgolang/sc"
)

func main() {
	const synthName = "Balance2Example"

	client, err := sc.NewClient("udp", "127.0.0.1:57110", "127.0.0.1:57120", 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		log.Fatal(err)
	}
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		bus, gain := sc.C(0), sc.C(0.1)
		l, r := sc.LFSaw{Freq: sc.C(44)}.Rate(sc.AR), sc.Pulse{Freq: sc.C(33)}.Rate(sc.AR)
		pos := sc.SinOsc{Freq: sc.C(0.5)}.Rate(sc.KR)
		sig := sc.Balance2{L: l, R: r, Pos: pos, Level: gain}.Rate(sc.AR)
		return sc.Out{bus, sig}.Rate(sc.AR)
	})
	if err := client.SendDef(def); err != nil {
		log.Fatal(err)
	}

	synthID := client.NextSynthID()
	if _, err := defaultGroup.Synth(synthName, synthID, sc.AddToTail, nil); err != nil {
		log.Fatal(err)
	}
	log.Printf("created synth %d\n", synthID)
}
