package main

import (
	"log"

	"github.com/scgolang/sc"
)

func main() {
	// create a client and connect to the server
	client, err := sc.NewClient("udp", "127.0.0.1:57121", "127.0.0.1:57120")
	if err != nil {
		log.Fatal(err)
	}
	def := sc.NewSynthdef("SineTone", func(p sc.Params) sc.Ugen {
		bus, chaos := sc.C(0), sc.Line{sc.C(1.0), sc.C(2.0), sc.C(10), sc.DoNothing}.Rate(sc.KR)
		sig := sc.Crackle{chaos}.Rate(sc.AR).MulAdd(sc.C(0.5), sc.C(0.5))
		return sc.Out{bus, sig}.Rate(sc.AR)
	})
	err = client.SendDef(def)
	if err != nil {
		log.Fatal(err)
	}
	id := client.NextSynthID()
	_, err = client.Synth("SineTone", id, sc.AddToTail, sc.DefaultGroupID, nil)
	if err != nil {
		log.Fatal(err)
	}
}
