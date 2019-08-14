package common

// PacketHeader This is the Structure for the TCP Packets.
type PacketHeader struct {
	Size          uint16
	CsMessageType uint8
	Fingerprint   uint16
	Checksum      uint32
}
