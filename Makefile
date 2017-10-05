all: scheduler

run: scheduler
	./scheduler

scheduler: plugin
	$(call blue, "Build scheduler...")
	# go get
	go build -o scheduler .

plugin: clean
	$(call blue, "Build plugins...")
	go build -o plugins/ssh/ssh github.com/bolsunovskyi/scheduler/plugins/ssh
	go build -o plugins/shell/shell.so github.com/bolsunovskyi/scheduler/plugins/shell

clean:
	$(call blue, "Clean work tree...")
	rm -f plugins/ssh/ssh
	rm -f plugins/shell/shell
	rm -f ./scheduler

define blue
	@tput setaf 6
	@echo $1
	@tput sgr0
endef