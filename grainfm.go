package main

import (
	"fmt"

	"github.com/scgolang/sc"
)

func main() {
	const synthName = "GrainFMExample"

	client, err := sc.NewClient("udp", "127.0.0.1:57112", "127.0.0.1:57110")
	if err != nil {
		panic(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		panic(err)
	}
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		gate := p.Add("gate", 1)
		amp := p.Add("amp", 1)
		bus := sc.C(0)
		mousey := sc.MouseY{Min: sc.C(0), Max: sc.C(400)}.Rate(sc.KR)
		freqdev := sc.WhiteNoise{}.Rate(sc.KR).Mul(mousey)
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
		trig := sc.Impulse{Freq: sc.C(10)}.Rate(sc.KR)
		modIndex := sc.LFNoise{Interpolation: sc.NoiseLinear}.Rate(sc.KR).MulAdd(sc.C(5), sc.C(5))
		pan := sc.MouseX{Min: sc.C(-1), Max: sc.C(1)}.Rate(sc.KR)
		sig := sc.GrainFM{
			NumChannels: 2,
			Trigger:     trig,
			Dur:         sc.C(0.1),
			CarFreq:     sc.C(440).Add(freqdev),
			ModFreq:     sc.C(200),
			ModIndex:    modIndex,
			Pan:         pan,
		}.Rate(sc.AR)
		return sc.Out{bus, sig.Mul(ampenv)}.Rate(sc.AR)
	})
	err = client.SendDef(def)
	if err != nil {
		panic(err)
	}
	synthID := client.NextSynthID()
	_, err = defaultGroup.Synth(synthName, synthID, sc.AddToTail, nil)
	fmt.Printf("created synth %d\n", synthID)
}
