package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/scgolang/sc"
)

func main() {
	const (
		synthName = "GrainFMExample"
	)
	server := &sc.Server{
		Network: "udp",
		Port:    57120,
	}

	stdout, stderr, err := server.Start(5 * time.Second)
	if err != nil {
		log.Fatalf("Could not start scsynth: %s", err)
	}
	defer func() { _ = server.Process.Kill() }()

	go func() {
		if _, err := io.Copy(os.Stderr, stderr); err != nil {
			log.Fatalf("Could not pipe scsynth stderr to terminal: %s", err)
		}
	}()

	go func() {
		if _, err := io.Copy(os.Stdout, stdout); err != nil {
			log.Fatalf("Could not pipe scsynth stdout to terminal: %s", err)
		}
	}()

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
			gate    = p.Add("gate", 1)
			amp     = p.Add("amp", 1)
			bus     = sc.C(0)
			mousey  = sc.MouseY{Min: sc.C(0), Max: sc.C(400)}.Rate(sc.KR)
			freqdev = sc.WhiteNoise{}.Rate(sc.KR).Mul(mousey)
		)
		env := sc.Env{
			Levels:      []sc.Input{sc.C(0), sc.C(1), sc.C(0)},
			Times:       []sc.Input{sc.C(1), sc.C(1)},
			Curve:       []string{"sine", "sine"},
			ReleaseNode: sc.C(1),
		}
		ampenv := sc.EnvGen{
			Env:        env,
			Gate:       gate,
			LevelScale: amp,
			Done:       sc.FreeEnclosing,
		}.Rate(sc.KR)

		var (
			trigFreq = sc.MouseY{Min: sc.C(8), Max: sc.C(64)}.Rate(sc.KR)
			trig     = sc.Impulse{Freq: trigFreq}.Rate(sc.KR)
			modIndex = sc.LFNoise{Interpolation: sc.NoiseLinear}.Rate(sc.KR).MulAdd(sc.C(5), sc.C(5))
			pan      = sc.MouseX{Min: sc.C(-1), Max: sc.C(1)}.Rate(sc.KR)
		)
		sig := sc.GrainFM{
			NumChannels: 2,
			Trigger:     trig,
			Dur:         sc.C(0.1),
			CarFreq:     sc.C(440),
			ModFreq:     sc.C(200).Add(freqdev),
			ModIndex:    modIndex,
			Pan:         pan,
		}.Rate(sc.AR)

		return sc.Out{bus, sig.Mul(ampenv)}.Rate(sc.AR)
	})
	if err := client.SendDef(def); err != nil {
		log.Fatal(err)
	}
	synthID := client.NextSynthID()
	if _, err := defaultGroup.Synth(synthName, synthID, sc.AddToTail, nil); err != nil {
		log.Fatal(err)
	}

	if err := server.Wait(); err != nil {
		log.Fatal(err)
	}
}
