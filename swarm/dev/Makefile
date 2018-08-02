.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/DEC/go-DEC/cmd/swarm
	go build -o bin/geth     github.com/DEC/go-DEC/cmd/geth
	go build -o bin/bootnode github.com/DEC/go-DEC/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/DEC/go-DEC/swarm/...
