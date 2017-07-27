all: scheduler

run: scheduler
	./scheduler

scheduler: plugin
	$(call blue, "Build scheduler...")
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o scheduler .

plugin: clean
	$(call blue, "Build plugins...")
	CGO_ENABLED=0 GOOS=linux go build -buildmode=plugin -a -installsuffix cgo -o plugins/ssh/ssh.so github.com/bolsunovskyi/scheduler/plugins/ssh
	CGO_ENABLED=0 GOOS=linux go build -buildmode=plugin -a -installsuffix cgo -o plugins/shell/shell.so github.com/bolsunovskyi/scheduler/plugins/shell

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