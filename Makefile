all: scheduler

run: scheduler
	./scheduler

scheduler: plugin
	$(call blue, "Build scheduler...")
	go build -a -o scheduler .

plugin: clean
	$(call blue, "Build plugins...")
	go build -buildmode=plugin -a -o plugins/ssh/ssh.so github.com/bolsunovskyi/scheduler/plugins/ssh
	go build -buildmode=plugin -a -o plugins/shell/shell.so github.com/bolsunovskyi/scheduler/plugins/shell

clean:
	$(call blue, "Clean work tree...")
	rm -f plugins/ssh/ssh.so
	rm -f plugins/shell/shell.so
	rm -f ./scheduler

define blue
	@tput setaf 6
	@echo $1
	@tput sgr0
endef