package main

import (
	"context"
	"time"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/log"
	"github.com/code-by-meal/go-rdp/stack/pdu/conn"
	"github.com/code-by-meal/go-rdp/stack/x224"
)

func main() {
	log.Dbg("<s>[TEST]</> TPKT")

	//host := "192.168.64.3"
	host := "172.16.0.19"
	port := uint16(3389)
	ctx := context.Background()
	username := "user"
	//password := "user"
	stream, err := core.NewStream(host, port, 5*time.Second, ctx)

	if err != nil {
		log.Err("failed connect to tcp", err)

		return
	}

	// connection request
	cr := conn.NewConnectionRequest(username)
	if err := cr.Write(stream); err != nil {
		log.Err(err)

		return
	}

	// confirm responsi
	buff, err := x224.Read(stream)

	if err != nil {
		log.Err(err)

		return
	}

	_ = buff
}
