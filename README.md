```
package icb // import "suah.dev/icb"


FUNCTIONS

func DefaultHandlers() map[string]interface{}
    DefaultHandlers produces a set of handlers that will print messages to
    stdout. Ping ("l") packets will also be Pong'd ("m") if received.


TYPES

type Client struct {
	Conn     net.Conn
	Handlers map[string]interface{}
}
    Client gives us a convenient way to Read and Write ICB packets to a remote
    server.

func (c *Client) Connect(s string) (err error)
    Connect connects to specified ICB server.

func (c *Client) Read() (*Packet, error)
    Read reads the next packet from the server. The first byte contains the
    length of the packet to be read.

func (c *Client) RunHandlers(s []string) error
    RunHandlers iterates over each handler and executes on the corresponding
    packet type. If the packet is unhandled, we return an error containing the
    packet type.

func (c *Client) Write(s []string) error
    Write writes a packet to an ICB server

type Packet struct {
	Buffer bytes.Buffer
}
    Packet represents a packet as described by
    https://www.icb.net/_jrudd/protocol.html

func (p *Packet) Decode() (*[]string, error)
    Decode reads from the buffer and decodes the packet.

func (p *Packet) Encode(params []string) error
    Encode writes an ICB formatted packet to the Buffer.

func (p *Packet) Send(c net.Conn) error
    Send takes a connection and sends the ICB packet over it.
```
