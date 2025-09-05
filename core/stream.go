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

func NewStream(ctx context.Context, host string, port uint16, timeout time.Duration) (*Stream, error) {
	dialer := net.Dialer{}
	dialerCtx, cancel := context.WithTimeout(ctx, timeout)

	defer cancel()

	con, err := dialer.DialContext(dialerCtx, "tcp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		return nil, fmt.Errorf("dial tcp: %w", err)
	}

	go func() {
		<-ctx.Done()
		con.Close()
	}()

	return &Stream{
		Conn: con,
	}, nil
}

func (s *Stream) SwitchSSL(ctx context.Context, timeout time.Duration) error {
	cfg := &tls.Config{
		MinVersion:         tls.VersionTLS10,
		MaxVersion:         tls.VersionTLS13,
		InsecureSkipVerify: true,
	}
	tlsClient := tls.Client(s.Conn, cfg)
	tlsContext, cancel := context.WithTimeout(ctx, timeout)

	defer cancel()

	if err := tlsClient.HandshakeContext(tlsContext); err != nil {
		return fmt.Errorf("switch ssl: %w", err)
	}

	s.Conn = tlsClient

	return nil
}

func (s *Stream) Read(d []byte) (int, error) {
	n, err := s.Conn.Read(d)

	if err != nil {
		return n, fmt.Errorf("stream: read %w", err)
	}

	return n, nil
}

func (s *Stream) Write(d []byte) (int, error) {
	n, err := s.Conn.Write(d)

	if err != nil {
		return n, fmt.Errorf("stream: write %w", err)
	}

	return n, nil
}
