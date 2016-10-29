package main

import (
	"github.com/scgolang/osc"
	"log"
	"net"
)

const (
	listenAddr = "127.0.0.1:57110"
)

// Send a /quit message to scsynth
func main() {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:57130")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := osc.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	errChan := make(chan error)
	doneChan := make(chan *osc.Message)

	dispatcher := osc.Dispatcher{
		"/done": func(msg *osc.Message) error {
			doneChan <- msg
			return nil
		},
	}
	go func() {
		if err := conn.Serve(dispatcher); err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("sending quit request")
	quitReq, err := osc.NewMessage("/quit")
	if err != nil {
		log.Fatal(err)
	}

	if err := conn.Send(quitReq); err != nil {
		log.Fatal(err)
	}
	select {
	case quitResp := <-doneChan:
		log.Printf("%+v\n", quitResp)
	case err := <-errChan:
		log.Fatal(err)
	}
}
