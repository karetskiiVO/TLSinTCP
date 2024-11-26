package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/jessevdk/go-flags"
)

func main() {
	var options struct {
		Args struct {
			Port     string
			CertFile string
			KeyFile  string
		} `positional-args:"yes" required:"3"`
	}

	parser := flags.NewParser(&options, flags.Default&(^flags.PrintErrors))
	_, err := parser.Parse()
	if err != nil {
		log.Panic(err)
	}

	cert, err := tls.LoadX509KeyPair(options.Args.CertFile, options.Args.KeyFile)
	if err != nil {
		log.Panic(err)
	}

	godotenv.Load("localenv.env")

	var logger io.Writer
	keylogFile, exists := os.LookupEnv("SSLKEYLOGFILE")
	if exists {
		file, err := os.Create(keylogFile)
		if err != nil {
			log.Panic(err)
		}
		defer file.Close()

		logger = file
	}

	config := &tls.Config{
		KeyLogWriter: logger,
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", ":"+options.Args.Port, config)
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
				text := scanner.Text()
				log.Printf("%v: %v", conn.RemoteAddr(), text)

				fmt.Fprintf(conn, "%v\n", text)
				if err != nil {
					log.Print(err)
				}
			}
		}()
	}
}
