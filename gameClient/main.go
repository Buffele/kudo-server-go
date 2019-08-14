package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	bitstream "../core/bitstream"
	client "../core/client"
	encoding "../core/encoding"
	hash "../core/hash"
)

var connection *client.Connection

func login(user string, password string) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, hash.Hash("RpcAccountService_login_2"))
	encoding.WriteRShortString(buf, user)
	encoding.WriteRShortString(buf, password)
	connection.Send(buf.Bytes())
}

func requestCharacterTemplates() {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, hash.Hash("RpcLoginService_requestCharacterTemplates_0"))
	connection.Send(buf.Bytes())
}

func onLoginResult(buf *bitstream.BitReader) {
	println("[CLIENT] RpcAccountService_onLoginResult")

	errorCode := encoding.ReadRUInt16(buf)
	fullUserInfo := encoding.ReadRFullUserInfo(buf)

	fmt.Printf("Error Code: %d\n", errorCode)
	fmt.Printf("Full User Info: %+v\n", fullUserInfo)
}

func onCharacterTemplates(buf *bitstream.BitReader) {
	println("[CLIENT] RpcLoginService_onLoginResult")

	templates := encoding.ReadRCharacterTemplates(buf)
	styles := encoding.ReadRCharacterStyles(buf)
	items := encoding.ReadRRawWearItems(buf)

	fmt.Printf("Templates: %+v\n", templates)
	fmt.Printf("Styles: %+v\n", styles)
	fmt.Printf("Items: %+v\n", items)
}

func main() {
	var host = "127.0.0.1:11120"

	println("[CLIENT] Connecting to " + host + "...")
	connection = client.New(host)

	connection.RegisterRPCMethod("RpcAccountService", "onLoginResult", 2, onLoginResult)
	connection.RegisterRPCMethod("RpcLoginService", "onCharacterTemplates", 3, onCharacterTemplates)

	connection.OnOpen(func() {
		println("[CLIENT] Connected.")

		//login("test", "123456")
		requestCharacterTemplates()
	})

	connection.OnError(func(err error) {
		if !connection.Connected {
			println(err.Error())
		} else {
			println(err.Error())
		}
	})

	connection.Listen()

	fmt.Println("Press any key to exit...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}
