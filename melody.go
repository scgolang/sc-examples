package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/scgolang/sc"
	"github.com/scgolang/scids/scid"
)

func main() {
	const synthName = "sineTone"

	// setup supercollider client
	client, err := sc.NewClient("udp", "127.0.0.1:57110", "127.0.0.1:57120", 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		log.Fatal(err)
	}
	err = client.DumpOSC(int32(1))
	if err != nil {
		log.Fatal(err)
	}
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		var (
			freq = p.Add("freq", 440)
			gain = p.Add("gain", 0.5)
			dur  = p.Add("dur", 1)
			bus  = sc.C(0)
		)
		env := sc.EnvGen{
			Env:        sc.EnvPerc{Release: dur},
			LevelScale: gain,
			Done:       sc.FreeEnclosing,
		}.Rate(sc.KR)
		sig := sc.SinOsc{Freq: freq}.Rate(sc.AR).Mul(env)
		return sc.Out{bus, sc.Multi(sig, sig)}.Rate(sc.AR)
	})
	if err := client.SendDef(def); err != nil {
		log.Fatal(err)
	}
	for _ = range time.NewTicker(125 * time.Millisecond).C {
		synthID, err := scid.Next()
		if err != nil {
			log.Fatal(err)
		}
		var (
			note = rand.Intn(128)
			gain = rand.Float32()
			dur  = rand.Float32()
		)
		_, err = defaultGroup.Synth(synthName, synthID, sc.AddToTail, map[string]float32{
			"freq": sc.Midicps(float32(note)),
			"gain": gain,
			"dur":  dur,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
