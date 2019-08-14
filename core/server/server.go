package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"log"
	"net"
	"strconv"

	bitstream "../bitstream"
	common "../common"
	csMessageType "../constants/csMessageType"
	encoding "../encoding"
	encryption "../encryption"
	hash "../hash"
)

// Client holds info about connection
type Client struct {
	// Properties.
	conn       net.Conn
	Server     *Server
	RemoteAddr net.Addr
	LocalAddr  net.Addr

	// NSRPC
	encryption *encryption.Encryption
}

// Server holds info about server
type Server struct {
	// Properties.
	address string

	// Callbacks.
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)

	// NSRPC
	rpcMethods map[uint32]interface{}
}

// Read client data from channel
func (client *Client) listen() {
	reader := bufio.NewReader(client.conn)
	for {
		buf1 := make([]byte, 2)
		_, err1 := reader.Read(buf1)

		if err1 != nil {
			client.conn.Close()
			client.Server.onClientConnectionClosed(client, err1)
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
			client.conn.Close()
			client.Server.onClientConnectionClosed(client, err3)
			return
		}

		messageType := uint8(buf3[0])
		//fmt.Printf("CsMessageType: %d\n", messageType)

		if messageType != csMessageType.MtGeneral {
			client.conn.Close()
			client.Server.onClientConnectionClosed(client, fmt.Errorf("Invalid CsMessageType: %d", messageType))
			return
		}

		buf4 := make([]byte, 2)
		_, err4 := reader.Read(buf4)

		if err4 != nil {
			client.conn.Close()
			client.Server.onClientConnectionClosed(client, err4)
			return
		}

		buf5 := bytes.NewReader(buf4)

		var fingerprint uint16
		binary.Read(buf5, binary.LittleEndian, &fingerprint)

		//fmt.Printf("Fingerprint: %d\n", fingerprint)

		buf6 := make([]byte, 4)
		_, err6 := reader.Read(buf6)

		if err6 != nil {
			client.conn.Close()
			client.Server.onClientConnectionClosed(client, err6)
			return
		}

		buf7 := bytes.NewReader(buf6)

		var checksum uint32
		binary.Read(buf7, binary.LittleEndian, &checksum)

		//fmt.Printf("Checksum: %d\n", checksum)

		buf := make([]byte, size)
		num, err := reader.Read(buf)

		if err != nil {
			client.conn.Close()
			client.Server.onClientConnectionClosed(client, err)
			return
		}

		message := make([]byte, num)
		copy(message, buf)

		key := client.encryption.NextDecrypt()
		//fmt.Printf("Key: %d\n", key)

		var calculatedFingerprint uint16 = uint16(uint32(key >> 5))

		if fingerprint != calculatedFingerprint {
			client.conn.Close()
			client.Server.onClientConnectionClosed(client, fmt.Errorf("Invalid Fingerprint! Expected: %d, Calculated: %d", fingerprint, calculatedFingerprint))
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
			client.conn.Close()
			client.Server.onClientConnectionClosed(client, fmt.Errorf("Invalid Checksum! Expected: %d, Calculated: %d", checksum, calculatedChecksum))
			return
		} /* else {
			println("Checksum OK!")
		}*/

		// RPC ID
		buf10 := bytes.NewReader(decryptedData)
		bufbit := bitstream.NewReader(buf10)
		rpcID := encoding.ReadRUInt32(bufbit)

		if client.Server.rpcMethods[rpcID] != nil {
			client.Server.rpcMethods[rpcID].(func(*Client, *bitstream.BitReader))(client, bufbit)
		} else {
			println("[SERVER] Unknown RPC ID received!")
			println(rpcID)
			println(encoding.Dump(decryptedData))
		}
	}
}

// Close closes the connection to the client.
func (client *Client) Close() error {
	return client.conn.Close()
}

// Send a message to the client.
func (client *Client) Send(b []byte) {
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

	key := client.encryption.NextEncrypt()
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

	client.conn.Write(buf.Bytes())
}

// OnNewClient is called when a new Client connected.
func (server *Server) OnNewClient(callback func(c *Client)) {
	server.onNewClientCallback = callback
}

// OnClientConnectionClosed is called when a client disconnected.
func (server *Server) OnClientConnectionClosed(callback func(c *Client, err error)) {
	server.onClientConnectionClosed = callback
}

// RegisterRPCMethod Call this to register RPC Methods.
func (server *Server) RegisterRPCMethod(rpcInterface string, method string, paramCount int, methodToCall interface{}) {
	rpcName := rpcInterface + "_" + method + "_" + strconv.Itoa(paramCount)
	rpcHash := hash.Hash(rpcName)
	server.rpcMethods[rpcHash] = methodToCall
	println("[SERVER] Registered RPC Method: " + rpcName)
}

// Listen listens to clients.
func (server *Server) Listen() {
	listener, err := net.Listen("tcp", server.address)
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		client := &Client{
			conn:       conn,
			Server:     server,
			RemoteAddr: conn.RemoteAddr(),
			LocalAddr:  conn.LocalAddr(),
		}
		client.encryption = encryption.New()
		client.encryption.Reset()
		go client.listen()
		server.onNewClientCallback(client)
	}
}

// New Server.
func New(address string) *Server {
	server := &Server{
		address: address,
	}

	server.OnNewClient(func(c *Client) {})
	server.OnClientConnectionClosed(func(c *Client, err error) {})

	server.rpcMethods = make(map[uint32]interface{})

	return server
}
