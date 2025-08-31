package main

import (
	"bytes"
	"context"
	"time"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/log"
	"github.com/code-by-meal/go-rdp/stack/tpkt"
)

func main() {
	log.Dbg("<s>[TEST]</> TPKT")

	host := "192.168.64.3"
	port := uint16(3389)
	ctx := context.Background()
	stream, err := core.NewStream(host, port, 5*time.Second, ctx)

	if err != nil {
		log.Err("failed connect to tcp", err)

		return
	}

	data := []byte("Hello wrld!")
	buff := bytes.NewBuffer(data)

	tpkt.Write(stream, buff)
}
