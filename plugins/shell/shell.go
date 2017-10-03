package main

import (
	"net/rpc/jsonrpc"

	"github.com/bolsunovskyi/scheduler/plugins"
	"github.com/natefinch/pie"
	"github.com/prometheus/common/log"
)

type Shell struct {
}

func (Shell) InitDB(path string, _ *string) error {
	return nil
}

func (Shell) HandleHTTP(request plugins.HTTPRequest, rsp *plugins.HTTPResponse) error {
	rsp.Template = request.BodyStr
	rsp.Data = map[string]interface{}{
		"method": request.Method,
	}

	return nil
}

func (Shell) GetPluginParams(_ string, params *plugins.PluginParams) error {
	*params = plugins.PluginParams{
		Name:        "shell",
		Description: "Execute shell commands",
		Version:     "1.0",
		HasSettings: false,
	}
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

func main() {
	p := pie.NewProvider()
	if err := p.RegisterName("shell", Shell{}); err != nil {
		log.Fatalln(err)
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
}
