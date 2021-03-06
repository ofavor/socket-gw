package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ofavor/socket-gw/transport"

	"github.com/ofavor/socket-gw/client"
	"github.com/ofavor/socket-gw/internal/log"
)

func main() {
	client := client.NewClient(
		client.LogLevel("debug"),
	)
	if err := client.Connect(); err != nil {
		log.Fatal("Connect error:", err)
	}
	if err := client.Send(transport.NewPacket(transport.PacketTypeAuth, []byte("abcd"))); err != nil {
		log.Fatal("Auth error:", err)
	}
	log.Info("Waiting for auth ACK ...")
	ack, err := client.Recv()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("ACK:", ack)

	go func() {
		for {
			p, err := client.Recv()
			if err != nil {
				log.Error("Receive error:", err)
				os.Exit(1)
			}
			fmt.Println(string(p.Body))
		}
	}()
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		client.Send(transport.NewPacket(11, []byte(text)))
	}
}
