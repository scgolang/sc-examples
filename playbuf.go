package main

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/scgolang/sc"
	"github.com/scgolang/scids/scid"
)

func main() {
	const (
		synthName = "playbufExample"
		wavFile   = "flame1.wav"
	)
	client, err := sc.NewClient("udp", "127.0.0.1:57110", "127.0.0.1:57120", 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		log.Fatal(err)
	}
	// read a wav file
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	buf, err := client.ReadBuffer(filepath.Join(wd, wavFile), 1)
	if err != nil {
		log.Fatal(err)
	}
	// send a synthdef
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		var (
			speed = p.Add("speed", 1.0)
			gain  = p.Add("gain", 0.5)
			bus   = sc.C(0)
		)
		sig := sc.PlayBuf{
			NumChannels: 2,
			BufNum:      sc.C(float32(buf.Num)),
			Speed:       speed,
			Done:        sc.FreeEnclosing,
		}.Rate(sc.AR).Mul(gain)
		return sc.Out{bus, sig}.Rate(sc.AR)
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
			speed = (float32(8.0) * rand.Float32()) + float32(0.5)
			gain  = rand.Float32()
			ctls  = map[string]float32{"speed": speed, "gain": gain}
		)
		if _, err := defaultGroup.Synth(synthName, synthID, sc.AddToTail, ctls); err != nil {
			log.Fatal(err)
		}
	}
}
