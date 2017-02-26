package main

import (
	"log"
	"time"

	"github.com/scgolang/sc"
	"github.com/scgolang/scids/scid"
)

func main() {
	const synthName = "BRFExample"

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
			line = sc.XLine{sc.C(0.7), sc.C(300), sc.C(20), 0}.Rate(sc.KR)
			saw  = sc.Saw{sc.C(200)}.Rate(sc.AR).Mul(sc.C(0.5))
			sine = sc.SinOsc{line, sc.C(0)}.Rate(sc.KR).MulAdd(sc.C(3800), sc.C(4000))
			bpf  = sc.BRF{saw, sine, sc.C(0.3)}.Rate(sc.AR)
			sig  = sc.Multi(bpf, bpf)
		)
		return sc.Out{sc.C(0), sig}.Rate(sc.AR)
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
