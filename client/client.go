package client

import (
	"context"
	"fmt"
	"time"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/log"
)

type Client struct {
	Host     string
	Port     uint16
	Context  context.Context
	Stream   *core.Stream
	Domain   string
	Username string
	Password string
	Timeout  time.Duration
}

func NewClient(host string, port uint16, ctx context.Context) *Client {
	log.Dbg("Init <d>rdp-client</>.")

	return &Client{
		Context: ctx,
		Host:    host,
		Port:    port,
	}
}

func (c *Client) Login(
	domain string,
	username string,
	password string,
) error {
	log.Dbg(fmt.Sprintf("Try login with to <d>%s:%d</> (domain: <d>%s</>\tusername: <d>%s</>\tpassword: <d>%s</>)", c.Host, c.Port, domain, username, password))

	c.Domain = domain
	c.Username = username
	c.Password = password
	c.Timeout = 5 * time.Second

	// Init tcp stream
	stream, err := core.NewStream(c.Host, c.Port, c.Timeout, c.Context)

	if err != nil {
		return fmt.Errorf("login: %v", err)
	}

	c.Stream = stream

	if err := c._Negotiation(); err != nil {
		return fmt.Errorf("nego: %v", err)
	}

	return nil
}
