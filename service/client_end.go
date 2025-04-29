package service

import (
	"log"
	"net"
	"net/rpc"
)

type ClientEnd struct {
	Name string
	Addr net.Addr
}

func (e ClientEnd) Call(rpcname string, args any, reply any) error {
	c, err := rpc.Dial(e.Addr.Network(), e.Addr.String())

	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	return c.Call(rpcname, args, reply)
}
