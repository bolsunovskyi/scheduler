package main

import "github.com/bolsunovskyi/scheduler/plugins"

type Shell struct {
}

func MakePlugin(params map[string]interface{}) plugins.JobStep {
	return Shell{}
}

func (Shell) GetName() string {
	return "shell"
}

func (Shell) GetDescription() string {
	return "Execute shell commands"
}

func (Shell) GetVersion() string {
	return "1.0"
}

func (Shell) GetBuildParams() []plugins.StepParam {
	return []plugins.StepParam{
		{
			Name: "Command",
			Type: plugins.TypeString,
		},
	}
}
