package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

func main() {
	var options struct {
		Args struct {
			Address string
		}
	}

	parser := flags.NewParser(&options, flags.Default&(^flags.PrintErrors))
	_, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	config := &tls.Config{}

	conn, err := tls.Dial("tcp", options.Args.Address, config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	netWriter := bufio.NewWriter(conn)
	go func() {
		netScanner := bufio.NewScanner(conn)

		for netScanner.Scan() {
			fmt.Println(netScanner.Text())
		}
	}()

	consoleScanner := bufio.NewScanner(os.Stdin)

	for consoleScanner.Scan() {
		netWriter.WriteString(consoleScanner.Text())
	}
}
