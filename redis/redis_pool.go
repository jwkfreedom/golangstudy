package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type InitFunction func() (redis.Conn, error)

type ConnectionPoolWrapper struct {
	size int
	conn chan interface{}
}

/**
Call the init function size times. If the init function fails during any call, then
the creation of the pool is considered a failure.
We call the same function size times to make sure each connection shares the same
state.
*/
func (p *ConnectionPoolWrapper) InitPool(size int, initfn InitFunction) error {
	// Create a buffered channel allowing size senders
	p.conn = make(chan interface{}, size)
	for x := 0; x < size; x++ {
		conn, err := initfn()
		if err != nil {
			return err
		}

		// If the init function succeeded, add the connection to the channel
		fmt.Println("one connection created!")
		p.conn <- conn
	}
	p.size = size
	return nil
}

func (p *ConnectionPoolWrapper) GetConnection() interface{} {
	return <-p.conn
}

func (p *ConnectionPoolWrapper) ReleaseConnection(conn interface{}) {
	p.conn <- conn
}

var dbPool = &ConnectionPoolWrapper{}
