.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/DEWH/go-DEWH/cmd/swarm
	go build -o bin/geth     github.com/DEWH/go-DEWH/cmd/geth
	go build -o bin/bootnode github.com/DEWH/go-DEWH/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/DEWH/go-DEWH/swarm/...
