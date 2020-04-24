package icb

import (
	"bytes"
	"testing"
)

var (
	loginPacket = []byte{0x61, 0x74, 0x65, 0x73, 0x74, 0x1, 0x74, 0x65, 0x73, 0x74, 0x1, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x0}
	a           = &Packet{}
)

func TestEncode(t *testing.T) {
	e := a.Encode([]string{"a", "test", "test", "login"})
	if e != nil {
		t.Error(e)
	}

	result := a.Buffer.Bytes()
	success := bytes.Compare(result, loginPacket)

	if success != 0 {
		t.Errorf("Expected %q; received %q\n", loginPacket, result)
	}
}
