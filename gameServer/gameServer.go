package gameServer

import (
	"bytes"
	"encoding/binary"

	bitstream "../core/bitstream"
	encoding "../core/encoding"
	hash "../core/hash"
	tcp_server "../core/server"
)

func onLoginResult(client *tcp_server.Client, errorCode uint16, fullUserInfo encoding.RFullUserInfo) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, hash.Hash("RpcAccountService_onLoginResult_2"))
	encoding.WriteRUInt16(buf, errorCode)
	encoding.WriteRFullUserInfo(buf, fullUserInfo)
	client.Send(buf.Bytes())
}

func onCharacterTemplates(client *tcp_server.Client, templates []encoding.RCharacterTemplate, styles []encoding.RCharacterStyle, items []encoding.RRawWearItem) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, hash.Hash("RpcLoginService_onCharacterTemplates_3"))
	encoding.WriteRCharacterTemplates(buf, templates)
	encoding.WriteRCharacterStyles(buf, styles)
	encoding.WriteRRawWearItems(buf, items)
	client.Send(buf.Bytes())
}

func onConfirmSeed(client *tcp_server.Client, buf *bitstream.BitReader) {
	println("[GAME] RpcSecurityService_onConfirmSeed")
}

func exchangeSeed(client *tcp_server.Client, buf *bitstream.BitReader) {
	println("[GAME] RpcSecurityService_exchangeSeed")
}

func heartbeat(client *tcp_server.Client, buf *bitstream.BitReader) {
	//println("[GAME] RpcHeartbeatService_heartbeat")
}

func login(client *tcp_server.Client, buf *bitstream.BitReader) {
	println("[GAME] RpcAccountService_login")

	user := encoding.ReadRShortString(buf)
	password := encoding.ReadRShortString(buf)

	println("[GAME] User: " + user)
	println("[GAME] Password: " + password)

	fullUserInfo := encoding.RFullUserInfo{
		UserId:   1,
		Nickname: "AngryBird",
		Str3:     "wstr20",
	}

	onLoginResult(client, 0, fullUserInfo)

	println("[GAME] Sending onLoginResult()")
}

func requestCharacterTemplates(client *tcp_server.Client, buf *bitstream.BitReader) {
	println("[GAME] RpcLoginService_requestCharacterTemplates")

	templates := []encoding.RCharacterTemplate{}

	template1 := encoding.RCharacterTemplate{
		CharacterID: 2,
		Stats:       [6]uint16{1120, 1040, 880, 880, 1120, 500},
		Bool3:       true,
		Visible:     true,
	}

	templates = append(templates, template1)

	styles := []encoding.RCharacterStyle{}
	items := []encoding.RRawWearItem{}

	onCharacterTemplates(client, templates, styles, items)

	println("[GAME] Sending onCharacterTemplates()")
}

func Run() {
	println("[GAME] Starting GameServer...")

	addr := "127.0.0.1:11120"

	server := tcp_server.New(addr)

	server.RegisterRPCMethod("RpcSecurityService", "onConfirmSeed", 1, onConfirmSeed)
	server.RegisterRPCMethod("RpcSecurityService", "exchangeSeed", 2, exchangeSeed)
	server.RegisterRPCMethod("RpcHeartbeatService", "heartbeat", 0, heartbeat)
	server.RegisterRPCMethod("RpcAccountService", "login", 2, login)
	server.RegisterRPCMethod("RpcLoginService", "requestCharacterTemplates", 0, requestCharacterTemplates)

	server.OnNewClient(func(client *tcp_server.Client) {
		println("[GAME] Connection opened")
	})

	server.OnClientConnectionClosed(func(client *tcp_server.Client, err error) {
		println("[GAME] Connection closed")
	})

	println("[GAME] Server started on " + addr)

	server.Listen()
}
