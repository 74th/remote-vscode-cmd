# rcode

Calling client VS Code from remote host.

If you use [remote development using ssh](https://code.visualstudio.com/docs/remote/ssh), [you can use `code` command on integrated terminal](https://code.visualstudio.com/docs/remote/ssh#_opening-a-terminal-on-a-remote-host).
But other terminals connecting to the remote host (such as iTerm2) can not use code command
`rcode` enables to call VS Code from remote host.

# Restriction

This needs to be able to connect client host machine from remote host.

# how to use

## MacOS client and Linux remote server on Virtual Machine on mac

### MacOS

Install rcode.

```
go install github.com/74th/remote-vscode-cmd/cmd/rcode@latest
```

Create rcode.command file to launch rcode.
`10.37.129.2` is MacOS IP on Host Only Network.

```
/Users/nnyn/go/bin/rcode -server -addr 10.37.129.2:5450 -command code
```

```
chmod 755 rcode.command
```

Setting to launch rcode.command. "System Preferences..." -> "Users & Groups" -> select user -> "Login Items" -> Add -> select rcode.command

### Virtual Machine

Install rcode.

```
go install github.com/74th/remote-vscode-cmd/cmd/rcode@latest
```

Create `/$HOME/.rcode.yaml` .

- client: MacOS IP on Host Only Network
- remote_hostname: host name used by Remote Development SSH

```
client: 10.37.129.2:5450
remote_hostname: linux
```

### using on Virtual Machine

```
rcode /some/dir
```

It will open /some/dir on VS Code on MacOS.

### how it works

`code` has an argument `--remote`, it launches VS Code with Remote Development.

```
code --remote ssh-remote+hostname /some/dir
```

`rcode /some/dir` on the remote host connects `rcode` server on the client, and `rcode` on the client calls `code --remote ssh-remote+hostname /some/dir`.

## how to build

```bash
go get

export GOBIN=`pwd`/.bin
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

protoc --plugin=.bin/protoc-gen-go --go_out=. --go_opt=paths=source_relative \
    --plugin=.bin/protoc-gen-go-grpc --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    rcode.proto

go build cmd/rcode
```
