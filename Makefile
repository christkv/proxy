GOPATH = ${PWD}
export GOPATH

setup:
	go get gopkg.in/mgo.v2

build:	
	go build proxy.go

execute: build	
	go run proxy.go
