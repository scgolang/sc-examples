package main

import (
	"encoding/json"
	"github.com/scgolang/sc"
	"log"
	"os"
	"time"
)

// Request status from scsynth
func main() {
	client, err := sc.NewClient("udp", "127.0.0.1:57121", "127.0.0.1:57120")
	if err != nil {
		log.Fatal(err)
	}
	status, err := client.Status(time.Minute)
	if err != nil {
		log.Fatal(err)
	}
	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(status)
	if err != nil {
		log.Fatal(err)
	}
}
