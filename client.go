package icb

import (
	"fmt"
	"log"
	"net"
	"strings"
)

// Client gives us a convenient way to Read and Write ICB packets to a
// remote server.
type Client struct {
	Conn     net.Conn
	Handlers map[string]interface{}
}

// DefaultHandlers produces a set of handlers that will print messages to
// stdout. Ping ("l") packets will also be Pong'd ("m") if received.
func DefaultHandlers() map[string]interface{} {
	return map[string]interface{}{
		"a": func(s []string, c *Client) {
		},
		"b": func(s []string, c *Client) {
			log.Printf("%s> %s", s[1], strings.Join(s[2:], " "))
		},
		"c": func(s []string, c *Client) {
			fmt.Printf("private msg from: %s> %s\n", s[1], strings.Join(s[2:], " "))
		},
		"d": func(s []string, c *Client) {
			fmt.Println("->", strings.Join(s[1:], " "))
		},
		"e": func(s []string, c *Client) {
			fmt.Println("ERROR>", strings.Join(s[1:], " "))
		},
		"j": func(s []string, c *Client) {
			fmt.Printf("-> Connected to %s (%s)\n", s[2], s[3])
		},
		"k": func(s []string, c *Client) {
			fmt.Println("-> BEEP")
		},
		"l": func(s []string, c *Client) {
			c.Write([]string{"m"})
		},
		"n": func(s []string, c *Client) {
		},
	}
}

// RunHandlers iterates over each handler and executes on the corresponding
// packet type. If the packet is unhandled, we return an error containing
// the packet type.
func (c *Client) RunHandlers(s []string) error {
	handled := false
	for h, f := range c.Handlers {
		if h == s[0] {
			f.(func([]string, *Client))(s, c)
			handled = true
		}
	}

	if !handled {
		return fmt.Errorf("unhandled packet: %q", s[0])
	}

	return nil
}

// Connect connects to specified ICB server.
func (c *Client) Connect(s string) (err error) {
	c.Conn, err = net.Dial("tcp", s)
	if err != nil {
		return err
	}

	return nil
}

// Write writes a packet to an ICB server
func (c *Client) Write(s []string) error {
	var p Packet
	err := p.Encode(s)
	if err != nil {
		return err
	}

	return p.Send(c.Conn)
}

// Read reads the next packet from the server. The first byte contains the
// length of the packet to be read.
func (c *Client) Read() (*Packet, error) {
	var p Packet
	var l = make([]byte, 1)
	var i int

	_, err := c.Conn.Read(l)
	if err != nil {
		return nil, err
	}

	i = int(l[0])
	if err != nil {
		return nil, err
	}

	var buf = make([]byte, i)
	var read = 0

	for read < i {
		r, _ := c.Conn.Read(buf[read:])
		read += r
	}
	_, err = p.Buffer.Write(buf)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
