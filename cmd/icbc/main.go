package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"suah.dev/icb"
)

func main() {
	var c icb.Client

	var server = flag.String("server", "fake", "server to connect to")
	var port = flag.Int("port", 7326, "port")
	var nick = flag.String("nick", "dummie", "nick name to use")
	var group = flag.String("group", "suah", "group to join")

	flag.Parse()

	err := c.Connect(fmt.Sprintf("%s:%d", *server, *port))
	if err != nil {
		log.Fatal("dial", err)
	}

	defer c.Conn.Close()

	c.Handlers = icb.DefaultHandlers()

	go func() {
		for {
			p, err := c.Read()
			if err != nil {
				log.Fatal("reading: ", err)
			}

			a, err := p.Decode()
			if err != nil {
				log.Printf("error decoding: %q", p.Buffer)
				continue
			}

			c.RunHandlers(*a)
		}
	}()

	c.Write([]string{"a", *nick, *nick, *group, "login"})
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")

		err = c.Write([]string{"b", text})
		if err != nil {
			log.Println(err)
		}
	}
}
