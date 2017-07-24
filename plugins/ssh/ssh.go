package main

import "github.com/bolsunovskyi/scheduler/plugins"

type SSH struct {
}

func MakePlugin(params map[string]interface{}) plugins.JobStep {
	return SSH{}
}

func (SSH) GetName() string {
	return "ssh"
}

func (SSH) GetDescription() string {
	return "Send files or execute commands over SSH"
}

func (SSH) GetVersion() string {
	return "1.0"
}

func (SheSSHll) GetBuildParams() []plugins.StepParam {
	return []plugins.StepParam{
		{
			Name: "Command",
			Type: plugins.TypeString,
		},
	}
}
