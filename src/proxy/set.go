package proxy

import (
	"gopkg.in/mgo.v2"
	"time"
)

func NewReplSet(uri string, timeout time.Duration) *ReplSet {
	set := new(ReplSet)
	set.uri = uri
	set.Timeout = timeout
	return set
}

type ReplSet struct {
	uri     string
	Session *mgo.Session
	Timeout time.Duration
}

func (p *ReplSet) Start() error {
	session, err := mgo.Dial(p.uri)
	if err != nil {
		return err
	}

	// Save the session and return nil
	p.Session = session
	return nil
}
