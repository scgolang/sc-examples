package main

import (
	"log"
	"time"

	"github.com/scgolang/sc"
	"github.com/scgolang/scids/scid"
)

func main() {
	const synthName = "FreeVerbExample"

	client, err := sc.NewClient("udp", "127.0.0.1:57110", "127.0.0.1:57120", 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		log.Fatal(err)
	}
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		var (
			mix     = p.Add("mix", 0.25)
			room    = p.Add("room", 0.15)
			damp    = p.Add("damp", 0.5)
			bus     = sc.C(0)
			impulse = sc.Impulse{Freq: sc.C(1)}.Rate(sc.AR)
			lfcub   = sc.LFCub{Freq: sc.C(1200), Iphase: sc.C(0)}.Rate(sc.AR).Mul(sc.C(0.1))
			decay   = sc.Decay{In: impulse, Decay: sc.C(0.25)}.Rate(sc.AR).Mul(lfcub)
			sig     = sc.FreeVerb{In: decay, Mix: mix, Room: room, Damp: damp}.Rate(sc.AR)
		)
		return sc.Out{bus, sc.Multi(sig, sig)}.Rate(sc.AR)
	})
	if err := client.SendDef(def); err != nil {
		log.Fatal(err)
	}
	synthID, err := scid.Next()
	if err != nil {
		log.Fatal(err)
	}
	if _, err := defaultGroup.Synth(synthName, synthID, sc.AddToTail, nil); err != nil {
		log.Fatal(err)
	}
	log.Printf("created synth %d\n", synthID)
}
