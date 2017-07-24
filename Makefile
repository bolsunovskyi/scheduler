all: scheduler

run: scheduler
	./scheduler

scheduler: plugin
	$(call blue, "Build scheduler...")
	go build

plugin:
	$(call blue, "Build plugins...")
	go build -buildmode=plugin -o plugins/ssh/ssh.so github.com/bolsunovskyi/scheduler/plugins/ssh
	go build -buildmode=plugin -o plugins/shell/shell.so github.com/bolsunovskyi/scheduler/plugins/shell


define blue
	@tput setaf 6
	@echo $1
	@tput sgr0
endef