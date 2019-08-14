package main

import (
	"fmt"

	bitstream "../core/bitstream"
	client "../core/client"
	encoding "../core/encoding"
)

var connection *client.Connection

func onConnectionInfo(buf *bitstream.BitReader) {
	println("RpcLoadBalancingServiceEvent_onConnectionInfo")

	addressPair := encoding.ReadRAddressPair(buf)

	fmt.Printf("IP: %s\n", addressPair.Ip)
	fmt.Printf("Port: %d\n", addressPair.Port)
}

func main() {
	var host = "127.0.0.1:11111"

	println("Connecting to " + host + "...")
	connection = client.New(host)

	connection.RegisterRPCMethod("RpcLoadBalancingServiceEvent", "onConnectionInfo", 1, onConnectionInfo)

	connection.OnOpen(func() {
		println("Connected.")
	})

	connection.OnError(func(err error) {
		if !connection.Connected {
			println(err.Error())
		} else {
			println(err.Error())
		}
	})

	connection.Listen()
}
