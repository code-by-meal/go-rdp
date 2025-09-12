package client

import (
	"context"
	"fmt"
	"time"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/log"
	"github.com/code-by-meal/go-rdp/stack/rdp/nego"
)

type Client struct {
	Host             string
	Port             uint16
	Context          context.Context
	Stream           *core.Stream
	Domain           string
	Username         string
	Password         string
	Timeout          time.Duration
	SelectedProtocol nego.NegoProtocol
	Width            uint32
	Height           uint32
	Hostname         string
	UserID           uint16
	ChannelIDs       []uint16
}

func NewClient(ctx context.Context, host string, port uint16, hostname string) *Client {
	log.Dbg("Init <d>rdp-client</>.")

	return &Client{
		Context:  ctx,
		Host:     host,
		Port:     port,
		Hostname: hostname,
	}
}

func (c *Client) Login(
	domain string,
	username string,
	password string,
) error {
	log.Dbg(fmt.Sprintf("Try login with to <d>%s:%d</> (domain: <d>%s</>\tusername: <d>%s</>\tpassword: <d>%s</>)", c.Host, c.Port, domain, username, password))

	prefix := "login: %w"

	c.Domain = domain
	c.Username = username
	c.Password = password
	c.Timeout = 5 * time.Second

	// Init tcp stream
	stream, err := core.NewStream(c.Context, c.Host, c.Port, c.Timeout)

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	c.Stream = stream

	if err := c._Negotiation(); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := c._BasicSettingExchange(); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := c._ChannelConnection(); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}

func (c *Client) Close() error {
	if err := c.Stream.Conn.Close(); err != nil {
		return fmt.Errorf("close: %w", err)
	}

	return nil
}
