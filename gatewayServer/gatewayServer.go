package gatewayServer

import (
	"bytes"
	"encoding/binary"

	encoding "../core/encoding"
	hash "../core/hash"
	tcp_server "../core/server"
)

func onConnectionInfo(client *tcp_server.Client, addressPair encoding.RAddressPair) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, hash.Hash("RpcLoadBalancingServiceEvent_onConnectionInfo_1")) // RPC ID
	encoding.WriteRAddressPair(buf, addressPair)
	client.Send(buf.Bytes())
}

func Run() {
	println("[GATEWAY] Starting GatewayServer...")

	addr := "127.0.0.1:11111"

	server := tcp_server.New(addr)

	server.OnNewClient(func(client *tcp_server.Client) {
		println("[GATEWAY] Connection opened")

		gameserverAddress := encoding.RAddressPair{
			Ip:   "127.0.0.1",
			Port: 11120,
		}

		onConnectionInfo(client, gameserverAddress)

		println("[GATEWAY] Sending onConnectionInfo()")
	})

	server.OnClientConnectionClosed(func(client *tcp_server.Client, err error) {
		println("[GATEWAY] Connection closed")
	})

	println("[GATEWAY] Server started on " + addr)

	server.Listen()
}
