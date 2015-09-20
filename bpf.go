package main

import (
	"fmt"
	"github.com/scgolang/sc"
)

func main() {
	const synthName = "BPFExample"
	client := sc.NewClient("127.0.0.1:57112")
	err := client.Connect("127.0.0.1:57110")
	if err != nil {
		panic(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		panic(err)
	}
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		line := XLine{C(0.7), C(300), C(20), 0}.Rate(KR)
		saw := Saw{C(200)}.Rate(AR).Mul(C(0.5))
		sine := FSinOsc{line, C(0)}.Rate(KR).MulAdd(C(3600), C(4000))
		bpf := BPF{saw, sine, C(0.3)}.Rate(AR)
		return Out{C(0), bpf}.Rate(AR)
	})
	err = client.SendDef(def)
	if err != nil {
		panic(err)
	}
	synthID := client.NextSynthID()
	_, err = defaultGroup.Synth(synthName, synthID, AddToTail, nil)
	fmt.Printf("created synth %d\n", synthID)
}
