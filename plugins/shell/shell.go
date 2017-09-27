package main

import (
	"net/rpc/jsonrpc"

	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/natefinch/pie"
	"github.com/prometheus/common/log"
)

type Shell struct {
}

func (Shell) GetName(_ string, rsp *string) error {
	*rsp = "shell"
	return nil
}

func (Shell) GetDescription(_ string, rsp *string) error {
	*rsp = "Execute shell commands"
	return nil
}

func (Shell) GetVersion(_ string, rsp *string) error {
	*rsp = "1.0"
	return nil
}

func (Shell) GetBuildParams(_ string, rsp *[]plugins.ItemParam) error {
	*rsp = []plugins.ItemParam{
		{
			Label: "Command",
			Name:  "command",
			Type:  plugins.TypeText,
		},
	}
	return nil
}

func (Shell) HasSettings(_ string, rsp *bool) error {
	*rsp = false
	return nil
}

func main() {
	p := pie.NewProvider()
	if err := p.RegisterName("shell", Shell{}); err != nil {
		log.Fatalln(err)
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
}
