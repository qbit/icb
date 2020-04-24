package icb

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

// Packet represents a packet as described by https://www.icb.net/_jrudd/protocol.html
type Packet struct {
	Buffer bytes.Buffer
}

// Most messages have two fields, this covers that common case to a []string.
func (p *Packet) readMessage() ([]string, error) {
	a, err := p.Buffer.ReadString(1)
	if err != nil {
		return nil, err
	}

	b := p.Buffer.Bytes()

	s := []string{a}
	for _, c := range bytes.Split(b, []byte{1}) {
		s = append(s, string(c))
	}

	return s, nil

}

// Decode reads from the buffer and decodes the packet.
func (p *Packet) Decode() (*[]string, error) {
	var s []string
	t, err := p.Buffer.ReadByte()
	if err != nil {
		return nil, err
	}

	// Set our packet type
	s = append(s, string(t))

	switch t {
	// 'a' Login packet response
	case byte('a'):
		s = append(s, "")
	// 'b' Open Message packet
	case byte('b'):
		r, err := p.readMessage()
		if err != nil {
			return nil, err
		}
		s = append(s, r...)
	// 'c' Personal Message Packet
	case byte('c'):
		r, err := p.readMessage()
		if err != nil {
			return nil, err
		}
		s = append(s, r...)
	// 'd' Status Message Packet
	case byte('d'):
		r, err := p.readMessage()
		if err != nil {
			return nil, err
		}
		s = append(s, r...)
	// 'e' Error Message Packet
	case byte('e'):
		r := p.Buffer.String()
		s = append(s, r)
	// 'f' Important Message Packet
	case byte('f'):
		r, err := p.readMessage()
		if err != nil {
			return nil, err
		}
		s = append(s, r...)
	// 'g' Exit Packet
	case byte('g'):
		s = append(s, "")
	// 'i' Command Output Packet
	case byte('i'):
		var cmdType = make([]byte, 2)
		_, err := p.Buffer.Read(cmdType)
		if err != nil {
			return nil, err
		}

		switch string(cmdType) {
		case "ac", "ec":
			s = append(s, p.Buffer.String())
		default:
			return nil, fmt.Errorf("unknown command output: %s", string(cmdType))
		}
	case byte('j'):
		r, err := p.readMessage()
		if err != nil {
			return nil, err
		}
		s = append(s, r...)
	case byte('k'):
		log.Println("beeb")
	// 'l' Ping Packet
	case byte('l'):
		r := p.Buffer.String()
		s = append(s, r)
	// 'm' Pong Packet
	case byte('m'):
		r := p.Buffer.String()
		s = append(s, r)
	case byte('n'):
		log.Println("nop")
	default:
		//return nil, fmt.Errorf("unknown packet: %s", string(t))
		log.Printf("unknown packet: %q", string(t))
	}

	return &s, nil
}

// Encode writes an ICB formatted packet to the Buffer.
func (p *Packet) Encode(params []string) error {
	l := len(params) - 1
	// TODO check for valid types ?
	p.Buffer.Write([]byte(params[0]))
	for i, param := range params[1:] {
		p.Buffer.WriteString(param)
		if l != i+1 {
			p.Buffer.WriteByte(1)
		}
	}
	p.Buffer.WriteByte(0)
	return nil
}

// Send takes a connection and sends the ICB packet over it.
func (p *Packet) Send(c net.Conn) error {
	var err error
	_, err = c.Write([]byte{byte(p.Buffer.Len())})
	if err != nil {
		return err
	}
	_, err = c.Write(p.Buffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}
