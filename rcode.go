package rcode

import (
	"context"
	fmt "fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"time"

	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

type service struct {
	command string
}

func (s *service) LaunchCode(ctx context.Context, req *LaunchRequest) (*LaunchResponse, error) {
	log.Printf("called host:%s args:%v", req.RemoteHostname, req.Args)
	args := append([]string{"--remote", "ssh-remote+" + req.RemoteHostname}, req.Args...)
	cmd := exec.Command(s.command, args...)
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

// NewServer starts rcode server
func NewServer(host string, command string) {

	srv := &service{
		command: command,
	}

	var listener net.Listener
	for {
		var err error
		listener, err = net.Listen("tcp", host)
		if err == nil {
			break
		}
		fmt.Sprintf("failed to listen %s: %s\n", host, err.Error())
		time.Sleep(5 * time.Second)
	}
	grpcServer := grpc.NewServer()
	RegisterRcodeServer(grpcServer, srv)
	err := grpcServer.Serve(listener)
	if err != nil {
		panic(fmt.Sprintf("failed to start: %s", err.Error()))
	}
}

type RcodeConf struct {
	Client         string `yaml:"client"`
	RemoteHostname string `yaml:"remote_hostname"`
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

	cwd, _ := os.Getwd()
	for i := range args {
		p := path.Join(cwd, args[i])
		if _, err := os.Stat(p); err == nil {
			args[i] = p
		}
	}
	fmt.Printf("calling %s as %s with args: %s\n", conf.Client, conf.RemoteHostname, args)

	conn, err := grpc.Dial(conf.Client, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("could not connect %s : %s", conf.Client, err.Error()))
	}
	client := NewRcodeClient(conn)
	req := LaunchRequest{
		Args:           args,
		RemoteHostname: conf.RemoteHostname,
	}
	res, err := client.LaunchCode(context.Background(), &req)
	if len(res.Message) > 0 {
		fmt.Print(res.Message)
	}
	os.Exit(int(res.ReturnCode))
}
