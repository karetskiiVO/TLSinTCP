package main

import (
	"bufio"
	"crypto/tls"
	"log"

	"github.com/jessevdk/go-flags"
)

func main() {
	var options struct {
		Args struct {
			Port string
		}
	}

	parser := flags.NewParser(&options, flags.Default&(^flags.PrintErrors))
	_, err := parser.Parse()
	if err != nil {
		log.Panic(err)
	}

	config := &tls.Config{}

	listener, err := tls.Listen("tcp", options.Args.Port, config)
	if err != nil {
		log.Panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Print(err)
			continue
		}

		go func() {
			defer conn.Close()

			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				_, err := conn.Write(scanner.Bytes())
				if err != nil {
					log.Print(err)
				}
			}
		}()
	}
}
