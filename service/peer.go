package service

import (
	"net"
	"net/rpc"
)

type Peer struct {
	Name string
	Addr net.Addr
}

func (e Peer) Call(rpcname string, args any, reply any) error {
	c, err := rpc.Dial(e.Addr.Network(), e.Addr.String())

	if err != nil {
		return err
	}

	defer c.Close()

	return c.Call(rpcname, args, reply)
}
