package core

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"time"
)

type Stream struct {
	Conn net.Conn
}

func NewStream(host string, port uint16, timeout time.Duration, ctx context.Context) (*Stream, error) {
	dialer := net.Dialer{}
	dialerCtx, _ := context.WithTimeout(ctx, timeout)
	con, err := dialer.DialContext(dialerCtx, "tcp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		return nil, fmt.Errorf("dial tcp: %v", err)
	}

	go func() {
		<-ctx.Done()
		con.Close()
	}()

	return &Stream{
		Conn: con,
	}, nil
}

func (s *Stream) SwitchSSL() error {
	cfg := &tls.Config{
		MinVersion:         tls.VersionTLS10,
		MaxVersion:         tls.VersionTLS13,
		InsecureSkipVerify: true,
	}
	tlsClient := tls.Client(s.Conn, cfg)

	if err := tlsClient.Handshake(); err != nil {
		return fmt.Errorf("switch ssl: %v", err)
	}

	s.Conn = tlsClient

	return nil
}

func (s *Stream) Read(d []byte) (int, error) {
	return s.Conn.Read(d)
}

func (s *Stream) Write(d []byte) (int, error) {
	return s.Conn.Write(d)
}

func ReadFull(stream io.Reader, length int) ([]byte, error) {
	buff := make([]byte, length)

	if _, err := io.ReadFull(stream, buff); err != nil {
		return buff, fmt.Errorf("stream: read full: %v", err)
	}

	return buff, nil
}
