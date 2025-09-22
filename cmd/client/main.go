package main

import (
	"context"
	"fmt"

	"github.com/code-by-meal/go-rdp/client"
	"github.com/code-by-meal/go-rdp/log"
)

type OS string

const (
	Windows OS = "Windows"
	Linux   OS = "Linux"
	MacOS   OS = "MacOS"
)

type OSVersion string

const (
	Windows10 OSVersion = "10"
	Windows11 OSVersion = "11"
	Windows8  OSVersion = "8"
	Windows7  OSVersion = "7"
	WindowsXP OSVersion = "XP"
)

type Session struct {
	Username string
	Password string
	Domain   string
	Host     string
	Port     uint16
	Hostname string
	OS
	OSVersion
}

func init() {
	log.Level = log.InfoLevel
}

func main() {
	log.Info("<s>[+]</> Start <s>RDP</> client (by <s>code-by-meal</>)")

	sessions := []Session{
		//Session{Host: "10.50.53.22", Port: 3389, Username: "admin1", Domain: "LAB1", Password: "ubuntu115!@#", OS: Windows, OSVersion: Windows7},
		Session{Host: "192.168.64.3", Port: 3389, Hostname: "TEST-MACHINE", Username: "user", Domain: "", Password: "user", OS: Windows, OSVersion: Windows11},
		//Session{Host: "172.16.0.19", Hostname: "WIN-MACH", Port: 3389, Username: "user", Domain: "", Password: "user", OS: Windows, OSVersion: Windows10},
	}
	ctx := context.Background()

	for _, session := range sessions {
		client := client.NewClient(ctx, session.Host, session.Port, session.Hostname)
		clientPrint := fmt.Sprintf("OS: <d>%s %s</> Host: <d>%s:%d</> Creds: <d>%s\\%s</>:<d>%s</>\t", string(session.OS), string(session.OSVersion), session.Host, session.Port, session.Domain, session.Username, session.Password)

		if err := client.Login(session.Domain, session.Username, session.Password); err != nil {
			log.Err("<e>[ERROR-LOGIN]</> ", clientPrint, err)

			continue
		}

		log.Info("<s>[SUCCESS-LOGIN]</> ", clientPrint)

		if err := client.Close(); err != nil {
			log.Err("<e>[ERROR-CLIENT-CONNECTION]</> ", clientPrint, err)
		}
	}
}
