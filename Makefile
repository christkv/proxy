GOPATH = ${PWD}
export GOPATH

setup:
	go get gopkg.in/mgo.v2

build:	
	go build mongor.go

execute: build	
	go run mongor.go
