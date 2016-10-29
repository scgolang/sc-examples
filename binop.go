package main

import (
	"log"
	"time"

	"github.com/scgolang/sc"
)

func main() {
	// create a client and connect to the server
	client, err := sc.NewClient("udp", "127.0.0.1:57121", "127.0.0.1:57120", 5*time.Second)
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
	if err := client.SendDef(def); err != nil {
		log.Fatal(err)
	}

	time.Sleep(1000 * time.Millisecond)

	id := client.NextSynthID()
	if _, err = client.Synth("Envgen1", id, sc.AddToTail, sc.DefaultGroupID, nil); err != nil {
		log.Fatal(err)
	}
	time.Sleep(5000 * time.Millisecond)
}
