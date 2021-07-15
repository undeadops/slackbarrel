help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

## build: Build binary
build:
	go mod download
	mkdir -p build
	GOOS=linux go build -o build/barrel-agent.linux cmd/agent/*.go
	GOOS=linux go build -o build/barrel-server.linux cmd/server/*.go
	GOOS=darwin go build -o build/barrel-server.darwin cmd/server/*.go

## run-server: Run the Management server from the local directory with default options
run-server:
	go run cmd/server/*.go

## clean: Delete previous builds
clean:
	rm -rf build

## ngrok: Run ngrok to start port-forwarding from behind a firewall with ngrok service
ngrok:
	ngrok http 5000

