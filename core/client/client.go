package client

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"net"
	"strconv"

	bitstream "../bitstream"
	common "../common"
	csMessageType "../constants/csMessageType"
	encoding "../encoding"
	encryption "../encryption"
	hash "../hash"
)

// Connection The Connection Struct.
type Connection struct {
	// Properties.
	Conn      net.Conn
	Address   string
	Connected bool

	// Callbacks.
	onOpenCallback  func()
	onErrorCallback func(err error)

	// NSRPC
	encryption *encryption.Encryption
	rpcMethods map[uint32]interface{}
}

// OnOpen This is called when the TCP Connection was established successfully.
func (connection *Connection) OnOpen(f func()) {
	connection.encryption = encryption.New()
	connection.encryption.Reset()
	connection.onOpenCallback = f
}

// OnError This is called when a Connection Error occurred.
func (connection *Connection) OnError(f func(err error)) {
	connection.onErrorCallback = f
}

// Close Call this to close the connection to the server.
func (connection *Connection) Close() {
	connection.Conn.Close()
}

// Send Call this to send data to the server.
func (connection *Connection) Send(b []byte) {
	cipherBlockSize := 8
	bodySize := len(b)
	if bodySize%cipherBlockSize != 0 {
		paddingToAdd := cipherBlockSize - (bodySize % cipherBlockSize)
		//fmt.Printf("cipherBlockSize: %d, bodySize: %d, paddingToAdd: %d", cipherBlockSize, bodySize, paddingToAdd)

		writer := new(bytes.Buffer)
		binary.Write(writer, binary.LittleEndian, b)
		for i := 0; i < paddingToAdd; i++ {
			binary.Write(writer, binary.LittleEndian, uint8(0))
		}
		b = writer.Bytes()
	}

	key := connection.encryption.NextEncrypt()
	encryptedData := encryption.TeaEncrypt(b, key)

	var fingerprint uint16 = uint16(key >> 5)

	buf := new(bytes.Buffer)

	header := common.PacketHeader{
		Size:          uint16(len(encryptedData)),
		CsMessageType: csMessageType.MtGeneral,
		Fingerprint:   fingerprint,
		Checksum:      uint32(fingerprint) | crc32.Checksum(b, crc32.IEEETable),
	}

	binary.Write(buf, binary.LittleEndian, &header)
	binary.Write(buf, binary.LittleEndian, encryptedData)

	connection.Conn.Write(buf.Bytes())
}

// RegisterRPCMethod Call this to register RPC Methods.
func (connection *Connection) RegisterRPCMethod(rpcInterface string, method string, paramCount int, methodToCall interface{}) {
	rpcName := rpcInterface + "_" + method + "_" + strconv.Itoa(paramCount)
	rpcHash := hash.Hash(rpcName)
	connection.rpcMethods[rpcHash] = methodToCall
	println("[CLIENT] Registered RPC Method: " + rpcName)
}

// Listen Call this to establish a connection to the target.
func (connection *Connection) Listen() {
	conn, err := net.Dial("tcp", connection.Address)

	if err != nil {
		connection.onErrorCallback(err)
	} else {
		defer conn.Close()
		connection.Conn = conn

		connection.Connected = true
		connection.onOpenCallback()
		connection.read()
	}
}

func (connection *Connection) read() {
	reader := bufio.NewReader(connection.Conn)

	for {
		buf1 := make([]byte, 2)
		_, err1 := reader.Read(buf1)

		if err1 != nil {
			connection.Close()
			connection.onErrorCallback(err1)
			return
		}

		buf2 := bytes.NewReader(buf1)

		var size uint16
		binary.Read(buf2, binary.LittleEndian, &size)

		// TODO: Whole size.
		if uint16(reader.Buffered()) < size {
			continue
		}

		buf3 := make([]byte, 1)
		_, err3 := reader.Read(buf3)

		if err3 != nil {
			connection.Close()
			connection.onErrorCallback(err3)
			return
		}

		messageType := uint8(buf3[0])
		//fmt.Printf("Encryption Type: %d\n", messageType)

		if messageType != csMessageType.MtGeneral {
			connection.Close()
			connection.onErrorCallback(fmt.Errorf("Invalid CsMessageType: %d", messageType))
			return
		}

		buf4 := make([]byte, 2)
		_, err4 := reader.Read(buf4)

		if err4 != nil {
			connection.Close()
			connection.onErrorCallback(err4)
			return
		}

		buf5 := bytes.NewReader(buf4)

		var fingerprint uint16
		binary.Read(buf5, binary.LittleEndian, &fingerprint)

		//fmt.Printf("Fingerprint: %d\n", fingerprint)

		buf6 := make([]byte, 4)
		_, err6 := reader.Read(buf6)

		if err6 != nil {
			connection.Close()
			connection.onErrorCallback(err6)
			return
		}

		buf7 := bytes.NewReader(buf6)

		var checksum uint32
		binary.Read(buf7, binary.LittleEndian, &checksum)

		//fmt.Printf("Checksum: %d\n", checksum)

		buf := make([]byte, size)
		num, err := reader.Read(buf)

		if err != nil {
			connection.Close()
			connection.onErrorCallback(err)
			return
		}

		message := make([]byte, num)
		copy(message, buf)

		key := connection.encryption.NextDecrypt()
		//fmt.Printf("Key: %d\n", key)

		var calculatedFingerprint uint16 = uint16(uint32(key >> 5))

		if fingerprint != calculatedFingerprint {
			connection.Close()
			connection.onErrorCallback(fmt.Errorf("Invalid Fingerprint! Expected: %d, Calculated: %d", fingerprint, calculatedFingerprint))
			return
		} /* else {
			println("Fingerprint OK!")
		}*/

		decryptedData := encryption.TeaDecrypt(message, key)

		crc := crc32.Checksum(decryptedData, crc32.IEEETable)
		//fmt.Printf("CRC32: %d\n", crc)

		calculatedChecksum := uint32(fingerprint) | crc

		//fmt.Printf("Calculated Checksum: %d\n", calculatedChecksum)

		if checksum != calculatedChecksum {
			connection.Close()
			connection.onErrorCallback(fmt.Errorf("Invalid Checksum! Expected: %d, Calculated: %d", checksum, calculatedChecksum))
			return
		} /* else {
			println("Checksum OK!")
		}*/

		buf10 := bytes.NewReader(decryptedData)
		//bufbit := bitstream.NewReader(buf10)
		bufbit := &bitstream.BitReader{R: buf10}
		rpcID := encoding.ReadRUInt32(bufbit)

		if connection.rpcMethods[rpcID] != nil {
			connection.rpcMethods[rpcID].(func(*bitstream.BitReader))(bufbit)
		} else {
			println("[CLIENT] Unknown RPC ID received!")
			println(rpcID)
			println(encoding.Dump(decryptedData))
		}
		println(encoding.Dump(decryptedData))
	}
}

// New The Client Constructor.
func New(address string) *Connection {
	connection := &Connection{Address: address, Connected: false}

	connection.OnOpen(func() {})
	connection.OnError(func(err error) {})

	connection.rpcMethods = make(map[uint32]interface{})

	return connection
}
