package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

func main() {
	var options struct {
		Args struct {
			Address  string
			CertFile string
		} `positional-args:"yes" required:"1"`
	}

	parser := flags.NewParser(&options, flags.Default&(^flags.PrintErrors))
	_, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rawCert, err := os.ReadFile(options.Args.CertFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	certs := x509.NewCertPool()
	ok := certs.AppendCertsFromPEM([]byte(rawCert))
	if !ok {
		fmt.Println("can't parse cert")
		os.Exit(1)
	}

	config := &tls.Config{
		RootCAs:            certs,
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", options.Args.Address, config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	go func() {
		netScanner := bufio.NewScanner(conn)

		for netScanner.Scan() {
			fmt.Println(netScanner.Text())
		}
	}()

	consoleScanner := bufio.NewScanner(os.Stdin)

	for consoleScanner.Scan() {
		text := consoleScanner.Text()
		conn.Write([]byte(text))
	}
}
