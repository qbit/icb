package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"suah.dev/icb"
)

func main() {
	var c icb.Client
	err := c.Connect("127.0.0.1:7326")
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

	c.Write([]string{"a", "q", "q", "snakes", "login"})
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
