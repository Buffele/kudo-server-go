package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
)

func dump(data []byte) string {
	result := fmt.Sprintf("%x", data)

	var buffer bytes.Buffer
	for i := 0; i < len(result); i += 2 {
		if i%32 == 0 {
			buffer.WriteString("\n")
		}

		buffer.WriteString(result[i:(i+2)] + " ")
	}

	return buffer.String()
}

var encryptionLime uint64 = 0x1000000000000
var encryptionCurrent uint64 = 0x1330E
var decryptionLime uint64 = 0x1000000000000
var decryptionCurrent uint64 = 0x1330E

func seedEncrypt(seed uint32) {
	value := (uint64(seed) << 16) | 0x330E
	if encryptionLime != 0 {
		encryptionCurrent = value % encryptionLime
	} else {
		encryptionCurrent = encryptionLime
	}
}

func seedDecrypt(seed uint32) {
	value := (uint64(seed) << 16) | 0x330E
	if decryptionLime != 0 {
		decryptionCurrent = value % decryptionLime
	} else {
		decryptionCurrent = decryptionLime
	}
}

func nextEncrypt() uint32 {
	encryptionCurrent = ((0x5DEECE66D * encryptionCurrent) + 11) & 0xFFFFFFFFFFFF
	return uint32(encryptionCurrent >> 17)
}

func nextDecrypt() uint32 {
	decryptionCurrent = ((0x5DEECE66D * decryptionCurrent) + 11) & 0xFFFFFFFFFFFF
	return uint32(decryptionCurrent >> 17)
}

func reset() {
	encryptionCurrent = 0x19720212330E
	if encryptionLime != 0 {
		encryptionCurrent %= encryptionLime
	}

	decryptionCurrent = 0x19720212330E
	if decryptionLime != 0 {
		decryptionCurrent %= decryptionLime
	}

	//fmt.Printf("Reset keygen %d %d\n", encryptionCurrent, decryptionCurrent)
}

func teaDecrypt(data []byte, key uint32) []byte {
	keys := [4]uint32{key >> 1, key << 1, key >> 7, key << 7}

	decryptedData := make([]byte, 0)

	size := 8

	for i := 0; i < len(data); i += size {
		decryptedChunk := teaDecryptInternal(data[i:(i+size)], keys)
		decryptedData = append(decryptedData, decryptedChunk...)
	}

	return decryptedData
}

func teaDecryptInternal(src []byte, keys [4]uint32) []byte {
	e := binary.LittleEndian

	a, b := e.Uint32(src), e.Uint32(src[4:])
	var m uint32 = 0xC6EF3720
	for i := 0; i < 32; i++ {
		b -= (m + keys[((m>>11)&3)]) ^ (a + (16*a ^ (a >> 5)))
		m += 0x61C88647
		a -= (m + keys[(m&3)]) ^ (b + (16*b ^ (b >> 5)))
	}

	dst := make([]byte, 8)
	e.PutUint32(dst, a)
	e.PutUint32(dst[4:], b)

	return dst
}

/*
size_t P2pPacketHeader::write(void* buffer) const
{
    char* to = static_cast<char*>(buffer);
    to += AbstractPacketHeader::write(to);
    to += packet::put(static_cast<srpc::UInt8>(packetType_), to);
    to += packet::put(peerId_, to);
    to += packet::put(sequenceNumber_, to);
    to += packet::put(sentTime_, to);
    return to - static_cast<char*>(buffer);
}


size_t P2pPacketHeader::read(const void* buffer)
{
    const char* from = static_cast<const char*>(buffer);
    from += AbstractPacketHeader::read(from);
    srpc::UInt8 packetType;
    from += packet::get(packetType, from);
    packetType_ = static_cast<srpc::RpcPacketType>(packetType);
    from += packet::get(peerId_, from);
    from += packet::get(sequenceNumber_, from);
    from += packet::get(sentTime_, from);
    return from - static_cast<const char*>(buffer);
}
*/

/*
enum CsMessageType
{
    mtSystem = 0,
    mtGeneral = 1
};

// HEADER
// AbstractPacketHeader
0x10, 0x00,             // Size of Body? 16
0x01,                   // CsMessageType

0xc2, 0xe9,             // Fingerprint - 59842 - 0xE9C2
0xe6, 0xfb, 0xf1, 0xe8, // Checksum - 3908172774 - 0xe8f1fbe6
// BODY (TEA)
0xea, 0x49, 0x57, 0x97, 0x9e, 0xa6, 0xf9, 0xd7,
0x22, 0x4a, 0x1d, 0x2f, 0xda, 0x2b, 0x36, 0x13




// BODY (DECRYPTED)
db d7 43 82                   // RPC ID (uint32) - 2185484251 - 0x8243d7db
09 31 32 37 2e 30 2e 30 2e 31 // IP (String)
70 2b                         // Port (uint16) - 11120 - 0x2B70
*/

/*



enum RpcPacketType
{
    ptReliable = 0, ///< ¼ø¼­ º¸Àå. Àü¼Û º¸Àå
    ptUnreliable = 1 ///< ¼ø¼­, Àü¼Û º¸ÀåÇÏÁö ¾ÊÀ½
};



???
P2P Packets

// AbstractPacketHeader
08 00 // Size of body
01    // CsMessageType

01 RpcPacketType

00 00 00
03 00 00 00 // PeerId
00 00 00 00 // sequenceNumber_?
ff ff ff ff // PeerTime ?
3b 6e


5b d9 48 32 7b 90 bf 77
*/

func main() {
	reset()
	key := nextDecrypt()
	fmt.Printf("Decrypt Key: %d\n", key)

	var header_fgprint uint16 = uint16(uint32(key >> 5))

	fmt.Printf("Fingerprint: %d\n", header_fgprint)

	data := []byte{0xea, 0x49, 0x57, 0x97, 0x9e, 0xa6, 0xf9, 0xd7, 0x22, 0x4a, 0x1d, 0x2f, 0xda, 0x2b, 0x36, 0x13}

	decryptedData := teaDecrypt(data, key)

	fmt.Println(dump(decryptedData))

	crc := crc32.Checksum(decryptedData, crc32.IEEETable)
	fmt.Printf("CRC32: %d\n", crc)

	headerChecksum := uint32(header_fgprint) | crc

	fmt.Printf("Header Checksum: %d\n", headerChecksum)

	fmt.Println("Press enter key to exit...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}
