package main

import (
	"C"
	"github.com/bolsunovskyi/scheduler/plugins"
)

type Shell struct {
}

func MakePlugin(params map[string]interface{}) plugins.Item {
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

func (Shell) GetBuildParams() []plugins.ItemParam {
	return []plugins.ItemParam{
		{
			Label: "Command",
			Name:  "command",
			Type:  plugins.TypeText,
		},
	}
}

func (Shell) HasSettings() bool {
	return false
}
