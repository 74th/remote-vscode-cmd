package main

import (
	"flag"

	rcode "github.com/74th/remote-code-cmd"
)

func main() {
	var isServer bool
	var addr string
	var command string
	flag.BoolVar(&isServer, "server", false, "launch server")
	flag.StringVar(&addr, "addr", "0.0.0.0:5450", "server addr")
	flag.StringVar(&command, "command", "code-insiders", "command")
	flag.Parse()
	if isServer {
		rcode.NewServer(addr, command)
	} else {
		rcode.CallServer(flag.Args()[1:])
	}
}
