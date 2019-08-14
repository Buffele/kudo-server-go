// Package bitstream is a simple wrapper around a io.Reader and io.Writer to provide bit-level access to the stream.
package bitstream

import (
	"io"
)

// A Bit is a zero or a one
type Bit bool

// A BitReader reads bits from an io.Reader
type BitReader struct {
	r     io.Reader // The internal Byte Reader.
	b     [1]byte   // The current byte read from the stream.
	count uint8     // Current Bit Position.
}

// A BitWriter writes bits to an io.Writer
type BitWriter struct {
	w     io.Writer
	b     [1]byte
	count uint8
}

// NewReader returns a BitReader that returns a single bit at a time from 'r'
func NewReader(r io.Reader) *BitReader {
	b := new(BitReader)
	b.r = r
	return b
}

// ReadBit returns the next bit from the stream, reading a new byte from the underlying reader if required.
func (bitReader *BitReader) ReadBit() (Bit, error) {
	if bitReader.count == 0 {
		if n, err := bitReader.r.Read(bitReader.b[:]); n != 1 || (err != nil && err != io.EOF) {
			return false, err
		}
		bitReader.count = 8
	}
	bitReader.count--
	d := (bitReader.b[0] & 0x80)
	bitReader.b[0] <<= 1
	return d != 0, nil
}

// NewWriter returns a BitWriter that buffers bits and write the resulting bytes to 'w'
func NewWriter(w io.Writer) *BitWriter {
	b := new(BitWriter)
	b.w = w
	b.count = 8
	return b
}

func (b *BitWriter) Pending() (byt byte, vals uint8) {
	return b.b[0], b.count
}

func (b *BitWriter) Resume(data byte, count uint8) {
	b.b[0] = data
	b.count = count
}

// WriteBit writes a single bit to the stream, writing a new byte to 'w' if required.
func (b *BitWriter) WriteBit(bit Bit) error {

	if bit {
		b.b[0] |= 1 << (b.count - 1)
	}

	b.count--

	if b.count == 0 {
		if n, err := b.w.Write(b.b[:]); n != 1 || err != nil {
			return err
		}
		b.b[0] = 0
		b.count = 8
	}

	return nil
}

// WriteByte writes a single byte to the stream, regardless of alignment
func (b *BitWriter) WriteByte(byt byte) error {

	// fill up b.b with b.count bits from byt
	b.b[0] |= byt >> (8 - b.count)

	if n, err := b.w.Write(b.b[:]); n != 1 || err != nil {
		return err
	}
	b.b[0] = byt << b.count

	return nil
}

// ReadByte reads a single byte from the stream, regardless of alignment
func (b *BitReader) ReadByte() (byte, error) {

	if b.count == 0 {
		n, err := b.r.Read(b.b[:])
		if n != 1 || (err != nil && err != io.EOF) {
			b.b[0] = 0
			return b.b[0], err
		}
		// mask io.EOF for the last byte
		if err == io.EOF {
			err = nil
		}
		return b.b[0], err
	}

	byt := b.b[0]

	var n int
	var err error
	n, err = b.r.Read(b.b[:])
	if n != 1 || (err != nil && err != io.EOF) {
		return 0, err
	}

	byt |= b.b[0] >> b.count

	b.b[0] <<= (8 - b.count)

	return byt, err
}

// ReadBits reads nbits from the stream
func (bitReader *BitReader) ReadBits(nbits int) (uint64, error) {

	var u uint64

	for nbits >= 8 {
		byt, err := bitReader.ReadByte()
		if err != nil {
			return 0, err
		}

		u = (u << 8) | uint64(byt)
		nbits -= 8
	}

	var err error
	for nbits > 0 && err != io.EOF {
		byt, err := bitReader.ReadBit()
		if err != nil {
			return 0, err
		}
		u <<= 1
		if byt {
			u |= 1
		}
		nbits--
	}

	return u, nil
}

// Flush empties the currently in-process byte by filling it with 'bit'.
func (bitWriter *BitWriter) Flush(bit Bit) error {

	for bitReader.count != 8 {
		err := bitWriter.WriteBit(bit)
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteBits writes the nbits least significant bits of u, most-significant-bit first.
func (b *BitWriter) WriteBits(u uint64, nbits int) error {
	u <<= (64 - uint(nbits))
	for nbits >= 8 {
		byt := byte(u >> 56)
		err := b.WriteByte(byt)
		if err != nil {
			return err
		}
		u <<= 8
		nbits -= 8
	}

	for nbits > 0 {
		err := b.WriteBit((u >> 63) == 1)
		if err != nil {
			return err
		}
		u <<= 1
		nbits--
	}

	return nil
}

// Reset replaces the underlying io.Reader with the provided reader and resets
// all interal state to its zero value, allowing for reuse of the BitReader
// without additional allocations
func (b *BitReader) Reset(r io.Reader) {
	b.r = r
	b.b[0] = 0x00
	b.count = 0
}

// Reset replaces the underlying io.Writer with the provided writer and resets
// all internal state to its initial value allowing for reuse of the BitWriter
// without additional allocations
func (b *BitWriter) Reset(w io.Writer) {
	b.w = w
	b.b[0] = 0x00
	b.count = 8
}



func (b *BitReader) readBitsImpl(DataType& value, int bitCount, StreamBuffer& streamBuffer, int& bitReadCount)
{
    if (! isEnoughBits<DataType>(bitCount)) {
        throw StreamingException(__FILE__, __LINE__, "invalid bitCount");
    }

    int shiftCount = CHAR_BIT - bitReadCount;
    int leftBitCount = bitCount - shiftCount;
    const size_t needBytes = static_cast<size_t>(toBytes(leftBitCount) + 1);
    const size_t streamBufferSize = streamBuffer.size();
    if (streamBufferSize < needBytes) {
        throw StreamingException(__FILE__, __LINE__, "buffer is not enough");
    }

    DataType result = (static_cast<DataType>(streamBuffer.front()) >> bitReadCount);
    for (; leftBitCount > 0; leftBitCount -= CHAR_BIT) {
        streamBuffer.pop();
        result |= (static_cast<DataType>(streamBuffer.front()) << shiftCount);
        shiftCount += CHAR_BIT;
    }

    bitReadCount = (bitReadCount + bitCount) & (CHAR_BIT - 1);
    if (bitReadCount == 0) {
        streamBuffer.pop();
    }

    const DataType mask =
        (bitCount < Bits<DataType>::size) ? (1 << bitCount) - 1 : ~0;
    value = result & mask;
}