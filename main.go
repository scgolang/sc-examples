package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/scgolang/sc"
	"github.com/scgolang/scid"
)

const usage = `
sc-examples [Options]

Options:

-l                  List the available examples.
-play EXAMPLE       Play the specified example.
`

var defs = map[string]*sc.Synthdef{}

func main() {
	var (
		list      bool
		selection string
	)
	flag.BoolVar(&list, "l", false, "list the available example sounds")
	flag.StringVar(&selection, "play", "", "example to play")
	flag.Parse()

	if list {
		doList()
		os.Exit(0)
	}
	os.Exit(doPlay(selection))
}

func doList() {
	for name := range defs {
		fmt.Println(name)
	}
}

func doPlay(selection string) int {
	def, ok := defs[selection]
	if !ok {
		fmt.Fprintf(os.Stderr, "unrecognized selection: %s\n", selection)
	}
	if err := scid.Play(def, nil); err != nil {
		fmt.Fprintln(os.Stderr, err.Error)
		return 1
	}
	return 0
}

func init() {
	flag.Usage = func() { fmt.Println(usage) }

	add("Balance2Example", func(p sc.Params) sc.Ugen {
		bus, gain := sc.C(0), sc.C(0.1)
		l, r := sc.LFSaw{Freq: sc.C(44)}.Rate(sc.AR), sc.Pulse{Freq: sc.C(33)}.Rate(sc.AR)
		pos := sc.SinOsc{Freq: sc.C(0.5)}.Rate(sc.KR)
		sig := sc.Balance2{L: l, R: r, Pos: pos, Level: gain}.Rate(sc.AR)
		return sc.Out{bus, sig}.Rate(sc.AR)
	})

	add("Envgen1", func(p sc.Params) sc.Ugen {
		ampEnv := sc.EnvGen{
			Env: sc.EnvPerc{
				Attack:  sc.C(0.01),
				Release: sc.C(1),
				Level:   sc.C(1),
				Curve:   sc.C(-4),
			},
			Gate:       sc.C(1),
			LevelScale: sc.C(1),
			LevelBias:  sc.C(0),
			TimeScale:  sc.C(1),
			Done:       sc.FreeEnclosing,
		}.Rate(sc.KR)

		return sc.Out{
			Bus:      sc.C(0),
			Channels: sc.PinkNoise{}.Rate(sc.AR).Mul(ampEnv),
		}.Rate(sc.AR)
	})

	add("BPFExample", func(p sc.Params) sc.Ugen {
		return sc.Out{
			Bus: sc.C(0),
			Channels: sc.BPF{
				In: sc.Saw{
					Freq: sc.C(200),
				}.Rate(sc.AR).Mul(sc.C(0.5)),

				Freq: sc.SinOsc{
					Freq: sc.XLine{
						Start: sc.C(0.7),
						End:   sc.C(300),
						Dur:   sc.C(20),
						Done:  0,
					}.Rate(sc.KR),
				}.Rate(sc.KR).MulAdd(sc.C(3600), sc.C(4000)),

				RQ: sc.C(0.3),
			}.Rate(sc.AR),
		}.Rate(sc.AR)
	})

	add("BRFExample", func(p sc.Params) sc.Ugen {
		return sc.Out{
			Bus: sc.C(0),
			Channels: sc.BRF{
				In: sc.Saw{
					Freq: sc.C(200),
				}.Rate(sc.AR).Mul(sc.C(0.5)),

				Freq: sc.SinOsc{
					Freq: sc.XLine{
						Start: sc.C(0.7),
						End:   sc.C(300),
						Dur:   sc.C(20),
						Done:  0,
					}.Rate(sc.KR),
				}.Rate(sc.KR).MulAdd(sc.C(3800), sc.C(4000)),

				RQ: sc.C(0.3),
			}.Rate(sc.AR),
		}.Rate(sc.AR)
	})

	add("CombCExample", func(p sc.Params) sc.Ugen {
		line := sc.XLine{
			Start: sc.C(0.0001),
			End:   sc.C(0.01),
			Dur:   sc.C(20),
		}.Rate(sc.KR)

		sig := sc.Comb{
			Interpolation: sc.InterpolationCubic,
			In:            sc.WhiteNoise{}.Rate(sc.AR).Mul(sc.C(0.01)),
			MaxDelayTime:  sc.C(0.01),
			DelayTime:     line,
			DecayTime:     sc.C(0.2),
		}.Rate(sc.AR)

		return sc.Out{sc.C(0), sc.Multi(sig, sig)}.Rate(sc.AR)
	})

	add("Decay2Example", func(p sc.Params) sc.Ugen {
		bus := sc.C(0)
		line := sc.XLine{Start: sc.C(1), End: sc.C(50), Dur: sc.C(20)}.Rate(sc.KR)
		pulse := sc.Impulse{Freq: line, Phase: sc.C(0.25)}.Rate(sc.AR)
		sig := sc.Decay2{In: pulse, Attack: sc.C(0.01), Decay: sc.C(0.2)}.Rate(sc.AR)
		gain := sc.SinOsc{Freq: sc.C(600)}.Rate(sc.AR)
		return sc.Out{bus, sig.Mul(gain)}.Rate(sc.AR)
	})

	add("DelayCExample", func(p sc.Params) sc.Ugen {
		bus, dust := sc.C(0), sc.Dust{Density: sc.C(1)}.Rate(sc.AR).Mul(sc.C(0.5))
		noise := sc.WhiteNoise{}.Rate(sc.AR)
		decay := sc.Decay{In: dust, Decay: sc.C(0.3)}.Rate(sc.AR).Mul(noise)
		sig := sc.Delay{
			Interpolation: sc.InterpolationCubic,
			In:            decay,
			MaxDelayTime:  sc.C(0.2),
			DelayTime:     sc.C(0.2),
		}.Rate(sc.AR).Add(decay)
		return sc.Out{bus, sig}.Rate(sc.AR)
	})

	add("FormletExample", func(p sc.Params) sc.Ugen {
		bus, sine := sc.C(0), sc.SinOsc{Freq: sc.C(5)}.Rate(sc.KR).MulAdd(sc.C(20), sc.C(300))
		blip := sc.Blip{Freq: sine, Harm: sc.C(1000)}.Rate(sc.AR).Mul(sc.C(0.1))
		line := sc.XLine{Start: sc.C(1500), End: sc.C(700), Dur: sc.C(8)}.Rate(sc.KR)
		sig := sc.Formlet{
			In:         blip,
			Freq:       line,
			AttackTime: sc.C(0.005),
			DecayTime:  sc.C(0.4),
		}.Rate(sc.AR)
		return sc.Out{bus, sig}.Rate(sc.AR)
	})

	add("FreeverbExample", func(p sc.Params) sc.Ugen {
		mix := p.Add("mix", 0.25)
		room := p.Add("room", 0.15)
		damp := p.Add("damp", 0.5)
		bus := sc.C(0)
		impulse := sc.Impulse{Freq: sc.C(1)}.Rate(sc.AR)
		lfcub := sc.LFCub{Freq: sc.C(1200), Iphase: sc.C(0)}.Rate(sc.AR).Mul(sc.C(0.1))
		decay := sc.Decay{In: impulse, Decay: sc.C(0.25)}.Rate(sc.AR).Mul(lfcub)
		sig := sc.FreeVerb{In: decay, Mix: mix, Room: room, Damp: damp}.Rate(sc.AR)
		return sc.Out{bus, sig}.Rate(sc.AR)
	})

	add("FsinoscExample", func(p sc.Params) sc.Ugen {
		bus := sc.C(0)
		line := sc.XLine{sc.C(4), sc.C(401), sc.C(8), 0}.Rate(sc.KR)
		sin1 := sc.SinOsc{line, sc.C(0)}.Rate(sc.AR).MulAdd(sc.C(200), sc.C(800))
		sin2 := sc.SinOsc{Freq: sin1}.Rate(sc.AR).Mul(sc.C(0.2))
		return sc.Out{bus, sin2}.Rate(sc.AR)
	})

	add("GateExample", func(p sc.Params) sc.Ugen {
		bus, noise := sc.C(0), sc.WhiteNoise{}.Rate(sc.KR)
		pulse := sc.LFPulse{Freq: sc.C(1.333), Iphase: sc.C(0.5)}.Rate(sc.KR)
		sig := sc.Gate{In: noise, Trig: pulse}.Rate(sc.AR)
		return sc.Out{bus, sig}.Rate(sc.AR)
	})

	add("GrainFMExample", func(p sc.Params) sc.Ugen {
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

	add("LFCubExample", func(p sc.Params) sc.Ugen {
		bus, gain := sc.C(0), sc.C(0.1)
		lfo1 := sc.LFCub{Freq: sc.C(0.2)}.Rate(sc.KR).MulAdd(sc.C(8), sc.C(10))
		lfo2 := sc.LFCub{Freq: lfo1}.Rate(sc.KR).MulAdd(sc.C(400), sc.C(800))
		sig := sc.LFCub{Freq: lfo2}.Rate(sc.AR).Mul(gain)
		return sc.Out{bus, sig}.Rate(sc.AR)
	})

	add("LFPulseExample", func(p sc.Params) sc.Ugen {
		lfoFreq, lfoPhase, lfoWidth := sc.C(3), sc.C(0), sc.C(0.3)
		bus, gain := sc.C(0), sc.C(0.1)
		freq := sc.LFPulse{lfoFreq, lfoPhase, lfoWidth}.Rate(sc.KR).MulAdd(sc.C(200), sc.C(200))
		iphase, width := sc.C(0), sc.C(0.2)
		sig := sc.LFPulse{freq, iphase, width}.Rate(sc.AR).Mul(gain)
		return sc.Out{bus, sig}.Rate(sc.AR)
	})
}

func add(name string, f sc.UgenFunc) {
	defs[name] = sc.NewSynthdef(name, f)
}
