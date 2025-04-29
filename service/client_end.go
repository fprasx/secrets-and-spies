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

func (e ClientEnd) Call(rpcname string, args any, reply any) bool {
	c, err := rpc.Dial(e.Addr.Network(), e.Addr.String())

	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	err = c.Call(rpcname, args, reply)

	return err == nil
}
