package relayServer

import (
	"fmt"
	"os"

	hash "../core/hash"
)

func Run() {
	println("[RELAY] Starting RelayServer...")

	rpcMethods := []string{
		// RpcRelayServiceImpl
		"RpcRelayServiceImpl_rpcRelay_6",
		"RpcRelayServiceImpl_rpcRelayed_5",
		// SimpleRelayService
		"SimpleRelayService_rpcResolve_0",
		"SimpleRelayService_rpcResolved_2",
		"SimpleRelayService_rpcRelay_6",
		"SimpleRelayService_rpcRelayed_5",
		// RpcStunServiceImpl
		"RpcStunServiceImpl_rpcResolve_0",
		"RpcStunServiceImpl_rpcResolved_2",
		// RpcSystemService
		"RpcSystemService_rpcConnect_3",
		"RpcSystemService_rpcConnected_3",
		"RpcSystemService_rpcRequestConnectReversal_1",
		"RpcSystemService_rpcConnectReversal_1",
		"RpcSystemService_rpcNewPeerConnected_1",
		"RpcSystemService_rpcDisconnect_0",
		"RpcSystemService_rpcPing_0",
		"RpcSystemService_rpcAcknowledgement_2",
		"RpcSystemService_rpcHostMigrated_0",
	}

	for _, rpcMethod := range rpcMethods {
		fmt.Printf("RPC: %s, RPC Hash: %d, Hex: %x\n", rpcMethod, hash.Hash(rpcMethod), hash.Hash(rpcMethod))
	}

	println("Press any key to exit...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}
