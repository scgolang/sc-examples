package main

import (
	"log"
	"time"

	"github.com/scgolang/sc"
)

func main() {
	// create a client and connect to the server
	client := sc.NewClient("127.0.0.1:57121")
	err := client.Connect("127.0.0.1:57120")
	if err != nil {
		log.Fatal(err)
	}
	// create a synthdef
	def := sc.NewSynthdef("Envgen1", func(p sc.Params) sc.Ugen {
		bus := sc.C(0)
		attack, release := sc.C(0.01), sc.C(1)
		level, curveature := sc.C(1), sc.C(-4)
		perc := sc.EnvPerc{attack, release, level, curveature}
		gate, levelScale, levelBias, timeScale := sc.C(1), sc.C(1), sc.C(0), sc.C(1)
		ampEnv := sc.EnvGen{perc, gate, levelScale, levelBias, timeScale, sc.FreeEnclosing}.Rate(sc.KR)
		noise := sc.PinkNoise{}.Rate(sc.AR).Mul(ampEnv)
		return sc.Out{bus, noise}.Rate(sc.AR)
	})
	err = client.SendDef(def)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1000 * time.Millisecond)
	id := client.NextSynthID()
	_, err = client.Synth("Envgen1", id, sc.AddToTail, sc.DefaultGroupID, nil)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(5000 * time.Millisecond)
}
