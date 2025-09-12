package gcc

import (
	"bytes"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/stack/mcs"
	"github.com/code-by-meal/go-rdp/stack/mcs/per"
)

type ConferenceCreateResponse struct {
	UserData []byte
}

func NewConfernceCreateResponse() *ConferenceCreateResponse {
	return &ConferenceCreateResponse{}
}

func (c *ConferenceCreateResponse) Read(stream io.Reader) (*bytes.Buffer, error) {
	cr := mcs.NewConnectResponse()
	prefix := "gcc: ccrsp: read: %w"
	buff, err := cr.Read(stream)

	if err != nil {
		return bytes.NewBuffer([]byte{}), fmt.Errorf(prefix, err)
	}

	rdrS := bytes.NewReader(buff.Bytes())

	if _, err := per.ReadChoice(rdrS); err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	if oid, err := per.ReadOID(rdrS); err != nil || !bytes.Equal(oid, []byte{0, 0, 20, 124, 0, 1}) {
		return buff, fmt.Errorf(prefix, err)
	}

	if _, err := per.ReadLength(rdrS); err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	if _, err := per.ReadChoice(rdrS); err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	if _, err := per.ReadInteger16(rdrS, uint16(mcs.UserIDBase)); err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	if _, err := per.ReadInteger(rdrS); err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	enum, err := per.ReadEnumerated(rdrS)

	if err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	if (enum + 1) > 16 {
		return buff, fmt.Errorf(prefix, fmt.Errorf("per invalid data"))
	}

	if _, err := per.ReadNumberOfSet(rdrS); err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	if _, err := per.ReadChoice(rdrS); err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	oStr, err := per.ReadOctetString(rdrS, 4)

	if err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	if !bytes.Equal(oStr, []byte("McDn")) {
		return buff, fmt.Errorf(prefix, fmt.Errorf("invalid octet string McDn"))
	}

	userData, err := per.ReadOctetString(rdrS, 0)

	if err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	c.UserData = userData

	return bytes.NewBuffer(c.UserData), nil
}
