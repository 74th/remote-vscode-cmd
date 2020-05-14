package rcode

import (
	"context"
	fmt "fmt"
	"net"
	"os"
	"os/exec"
	"path"

	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

type service struct {
	command string
}

func (s *service) LaunchCode(ctx context.Context, req *LaunchRequest) (*LaunchResponse, error) {
	cmd := exec.Command(s.command, req.Args...)
	err := cmd.Run()
	if err != nil {
		return &LaunchResponse{
			ReturnCode: 1,
			Message:    fmt.Sprintf("cannot launch : %s", err.Error()),
		}, nil
	}

	if cmd.ProcessState.ExitCode() > 0 {
		return &LaunchResponse{
			ReturnCode: int32(cmd.ProcessState.ExitCode()),
			Message:    fmt.Sprintf("failed command : %s", err.Error()),
		}, nil
	}

	return &LaunchResponse{}, nil
}

func NewServer(host string, command string) {
	srv := &service{
		command: command,
	}
	conn, err := net.Listen("tcp", host)
	if err != nil {
		panic(fmt.Sprintf("failed to listen %s: %s", host, err.Error()))
	}
	grpcServer := grpc.NewServer()
	RegisterRcodeServer(grpcServer, srv)
	err = grpcServer.Serve(conn)
	if err != nil {
		panic(fmt.Sprintf("failed to start: %s", err.Error()))
	}
}

type RcodeConf struct {
	host       string `yaml:host`
	remoteHost string `yaml:remote_host`
}

func CallServer(args []string) {
	home := os.Getenv("HOME")
	confPath := path.Join(home, ".rcode.yaml")
	f, err := os.Open(confPath)
	if err != nil {
		panic(fmt.Sprintf("error: %s", err.Error()))
	}
	defer f.Close()

	var conf RcodeConf
	err = yaml.NewDecoder(f).Decode(&conf)
	if err != nil {
		panic(fmt.Sprintf("%s conf decode error: %s", confPath, err.Error()))
	}

	conn, err := grpc.Dial(conf.host)
	if err != nil {
		panic(fmt.Sprintf("could not connect %s : %s", conf.host, err.Error()))
	}
	client := NewRcodeClient(conn)
	req := LaunchRequest{
		Args:       args,
		RemoteHost: conf.remoteHost,
	}
	res, err := client.LaunchCode(context.Background(), &req)
	if len(res.Message) > 0 {
		fmt.Print(res.Message)
	}
	os.Exit(int(res.ReturnCode))
}
