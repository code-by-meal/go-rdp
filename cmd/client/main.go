package main

import (
	"context"

	"github.com/code-by-meal/go-rdp/client"
	"github.com/code-by-meal/go-rdp/log"
)

func init() {
	log.Level = log.InfoLevel
}

func main() {
	log.Info("<s>[+]</> Start <s>RDP</> client (by <s>code-by-meal</>)")

	host := "192.168.64.3"
	port := uint16(3389) // 3389 - default service MS-RDP port
	ctx := context.Background()
	domain := ""
	username := "user"
	password := "user"

	client := client.NewClient(host, port, ctx)

	if err := client.Login(domain, username, password); err != nil {
		log.Err("<e>[-]</> <e>F</>a<e>i</>l<e>e</>d login to RDP.", err)
	} else {
		log.Info("<s>[+]</> <s>S</>u<s>c</>c<s>e</>s<s>s</> login to RDP.")
	}

}
