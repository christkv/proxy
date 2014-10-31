GOPATH = ${PWD}
export GOPATH

setup:
	go get gopkg.in/mgo.v2
	go get github.com/spf13/cobra

build:	
	go build mongor.go

execute: build	
	go run mongor.go
