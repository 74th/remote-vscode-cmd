package main

import (
	"os"

	rcode "github.com/74th/remote-vscode-cmd"
	"github.com/jessevdk/go-flags"
)

var Opts struct {
	IsServer bool   `long:"server" description:"launch server"`
	Addr     string `long:"addr" description:"server addr (127.0.0.1:5653)"`
	Command  string `long:"command" description:"launch server" default:"code"`
	AllowAny bool   `long:"allow-any" description:"allow any command"`
}

func main() {
	args, err := flags.Parse(&Opts)
	if err != nil {
		os.Exit(1)
	}
	if Opts.IsServer {
		rcode.NewServer(Opts.Addr, Opts.Command, Opts.AllowAny)
	} else {
		rcode.CallServer(Opts.Command, args)
	}
}
