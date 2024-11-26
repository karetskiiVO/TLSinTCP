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
		} `positional-args:"yes" required:"1"`
	}

	parser := flags.NewParser(&options, flags.Default&(^flags.PrintErrors))
	_, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	config := &tls.Config{
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
			fmt.Printf("server: %v\n", netScanner.Text())
		}
	}()

	consoleScanner := bufio.NewScanner(os.Stdin)

	for consoleScanner.Scan() {
		text := consoleScanner.Text()
		fmt.Fprintf(conn, "%v\n", text)
	}
}
