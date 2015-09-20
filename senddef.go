package main

import (
	"github.com/scgolang/sc"
	"log"
)

func main() {
	// create a client and connect to the server
	client := sc.NewClient("127.0.0.1:57121")
	err := client.Connect("127.0.0.1:57120")
	if err != nil {
		log.Fatal(err)
	}
	def := sc.NewSynthdef("SineTone", func(p sc.Params) sc.Ugen {
		bus, chaos := C(0), Line{C(1.0), C(2.0), C(10), DoNothing}.Rate(KR)
		sig := Crackle{chaos}.Rate(AR).MulAdd(C(0.5), C(0.5))
		return Out{bus, sig}.Rate(AR)
	})
	err = client.SendDef(def)
	if err != nil {
		log.Fatal(err)
	}
	id := client.NextSynthID()
	_, err = client.Synth("SineTone", id, AddToTail, DefaultGroupID, nil)
	if err != nil {
		log.Fatal(err)
	}
}
