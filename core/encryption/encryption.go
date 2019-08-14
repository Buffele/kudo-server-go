package encryption

import "encoding/binary"

type Encryption struct {
	encryptionLime    uint64
	encryptionCurrent uint64
	decryptionLime    uint64
	decryptionCurrent uint64
}

func (self *Encryption) SeedEncrypt(seed uint32) {
	value := (uint64(seed) << 16) | 0x330E
	if self.encryptionLime != 0 {
		self.encryptionCurrent = value % self.encryptionLime
	} else {
		self.encryptionCurrent = self.encryptionLime
	}
}

func (self *Encryption) SeedDecrypt(seed uint32) {
	value := (uint64(seed) << 16) | 0x330E
	if self.decryptionLime != 0 {
		self.decryptionCurrent = value % self.decryptionLime
	} else {
		self.decryptionCurrent = self.decryptionLime
	}
}

func (self *Encryption) NextEncrypt() uint32 {
	self.encryptionCurrent = ((0x5DEECE66D * self.encryptionCurrent) + 11) & 0xFFFFFFFFFFFF
	return uint32(self.encryptionCurrent >> 17)
}

func (self *Encryption) NextDecrypt() uint32 {
	self.decryptionCurrent = ((0x5DEECE66D * self.decryptionCurrent) + 11) & 0xFFFFFFFFFFFF
	return uint32(self.decryptionCurrent >> 17)
}

func (self *Encryption) Reset() {
	self.encryptionCurrent = 0x19720212330E
	if self.encryptionLime != 0 {
		self.encryptionCurrent %= self.encryptionLime
	}

	self.decryptionCurrent = 0x19720212330E
	if self.decryptionLime != 0 {
		self.decryptionCurrent %= self.decryptionLime
	}
}

func TeaDecrypt(data []byte, key uint32) []byte {
	keys := [4]uint32{key >> 1, key << 1, key >> 7, key << 7}

	decryptedData := make([]byte, 0)

	size := 8

	for i := 0; i < len(data); i += size {
		decryptedChunk := teaDecryptInternal(data[i:(i+size)], keys)
		decryptedData = append(decryptedData, decryptedChunk...)
	}

	return decryptedData
}

func TeaEncrypt(data []byte, key uint32) []byte {
	keys := [4]uint32{key >> 1, key << 1, key >> 7, key << 7}

	encryptedData := make([]byte, 0)

	size := 8

	for i := 0; i < len(data); i += size {
		encryptedChunk := teaEncryptInternal(data[i:(i+size)], keys)
		encryptedData = append(encryptedData, encryptedChunk...)
	}

	return encryptedData
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

func teaEncryptInternal(src []byte, keys [4]uint32) []byte {
	e := binary.LittleEndian

	a, b := e.Uint32(src), e.Uint32(src[4:])
	var m uint32 = 0
	for i := 0; i < 32; i++ {
		a += (m + keys[m&3]) ^ (b + ((b << 4) ^ (b >> 5)))
		m -= 0x61C88647
		b += (m + keys[(m>>11)&3]) ^ (a + ((a << 4) ^ (a >> 5)))
	}

	dst := make([]byte, 8)
	e.PutUint32(dst, a)
	e.PutUint32(dst[4:], b)

	return dst
}

func New() *Encryption {
	return &Encryption{
		encryptionLime:    0x1000000000000,
		encryptionCurrent: 0x1330E,
		decryptionLime:    0x1000000000000,
		decryptionCurrent: 0x1330E,
	}
}
