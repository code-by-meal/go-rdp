package main

import "github.com/code-by-meal/go-rdp/log"

func init() {
	log.Level = log.InfoLevel
}

func main() {
	log.Info("<s>[+]</> Start <s>RDP</> client (by <s>code-by-meal</>)")
}
