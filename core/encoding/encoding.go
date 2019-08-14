package encoding

import (
	"bytes"
	"encoding/binary"
	"fmt"

	bitstream "../bitstream"
)

// The Byte Order.
var byteOrder binary.ByteOrder = binary.LittleEndian

type RAddressPair struct {
	Ip   string
	Port uint16
}

type RWearItem struct {
	/*
		srpc::RUInt16 uint11;
		srpc::RUInt32 uint12;
		srpc::RBool bool13;
		srpc::RpcIntType<srpc::Int32, 5> sint14;
		srpc::RUInt16 uint15;
		srpc::RInt32 uint16;
	*/

	// TODO RSubInfo3
	/*
		srpc::RUInt32 itemCode;
		srpc::RpcIntType<srpc::Int32, 5> sint2;
		srpc::RBool bool3;
		srpc::RBool bool4;
		srpc::RBool bool5;
		srpc::RBool bool6;
		srpc::RBool bool7;
		srpc::RInt32 sint8;
		srpc::RInt32 sint9;
		srpc::RUInt16 uints10[6]; // stats??
	*/
}

type RUnk1Item struct {
	/*
		srpc::RUInt32 uint11;
		srpc::RInt32 uint12;
		srpc::RUInt32 uint13;
		srpc::RBool bool14;
		srpc::RpcIntType<srpc::Int32, 5> sint15;
		srpc::RUInt16 uint16;
		srpc::RInt32 sint17;
	*/

	// TODO RSubInfo3
	/*
		srpc::RUInt32 itemCode;
		srpc::RpcIntType<srpc::Int32, 5> sint2;
		srpc::RBool bool3;
		srpc::RBool bool4;
		srpc::RBool bool5;
		srpc::RBool bool6;
		srpc::RBool bool7;
		srpc::RInt32 sint8;
		srpc::RInt32 sint9;
		srpc::RUInt16 uints10[6]; // stats??
	*/
}

type RUnkInfo4 struct {
	uint1  uint16
	uint2  uint16
	uint3  uint16
	uint4  uint16
	uint5  uint16
	uint6  uint16
	uint7  uint16
	uint81 uint16
	uint9  uint16
	uint10 uint16 // 9 bits
	uint11 uint16 // 9 bits
	uint12 uint16 // 9 bits
	uint13 uint16 // 9 bits
	uint14 uint16 // 9 bits
}

type RScore struct {
	winSingle  uint16
	loseSingle uint16
	winDouble  uint16
	loseDouble uint16
}

type RItem struct {
	uint1 uint32
	uint2 uint32
}

type RFullCharacter struct {
	id       uint32
	charType uint32
	exp      uint32
	rank     uint32
	items    []RItem
	vect10   []uint32
	scores   map[int32]RScore
	uint12   uint8
	uint13   uint8
}

type RFullUserInfo struct {
	UserId      uint32
	Nickname    string
	Str3        string
	Characters  []RFullCharacter
	SmashPoints int32
	SmashCash   int32
	WearItems   []RWearItem
	Items8      []RUnk1Item
	Data9       RUnkInfo4
}

func ReadRBool(b *bitstream.BitReader) bool {
	result, _ := b.ReadBits(1)
	return bool(result != 0)
}

func ReadRUInt8(b *bitstream.BitReader) uint8 {
	result, _ := b.ReadBits(8)
	return uint8(result)
}

func ReadRUInt16(b *bitstream.BitReader) uint16 {
	result, _ := b.ReadBits(16)
	bytes := []byte{uint8(result >> 8), uint8(result)}
	return byteOrder.Uint16(bytes)
}

func ReadRUInt32(b *bitstream.BitReader) uint32 {
	result, _ := b.ReadBits(32)
	bytes := []byte{uint8(result >> 24), uint8(result >> 16), uint8(result >> 8), uint8(result)}
	return byteOrder.Uint32(bytes)
}

func ReadRInt32(b *bitstream.BitReader) int32 {
	result, _ := b.ReadBits(32)
	bytes := []byte{uint8(result >> 24), uint8(result >> 16), uint8(result >> 8), uint8(result)}
	return int32(byteOrder.Uint32(bytes))
}

func ReadRInt32Bits(b *bitstream.BitReader, bitSize int) int32 {
	result, _ := b.ReadBits(bitSize)
	bytes := []byte{uint8(result >> 24), uint8(result >> 16), uint8(result >> 8), uint8(result)}
	return int32(byteOrder.Uint32(bytes))
}

func ReadRShortString(b *bitstream.BitReader) string {
	size := ReadRUInt8(b)
	bytes := make([]byte, size)

	for i := 0; i < int(size); i++ {
		result, _ := b.ReadBits(8)
		bytes[i] = byte(result)
	}

	return string(bytes[:size])
}

func ReadWString(buf *bitstream.BitReader) string {
	return ReadRShortString(buf)
}

func ReadRAddressPair(buf *bitstream.BitReader) RAddressPair {
	addressPair := RAddressPair{
		Ip:   ReadRShortString(buf),
		Port: ReadRUInt16(buf),
	}

	return addressPair
}

func ReadRFullUserInfo(buf *bitstream.BitReader) RFullUserInfo {
	fullUserInfo := RFullUserInfo{
		UserId:      ReadRUInt32(buf),
		Nickname:    ReadWString(buf),
		Str3:        ReadWString(buf),
		Characters:  ReadRFullCharacters(buf),
		SmashPoints: ReadRInt32(buf),
		SmashCash:   ReadRInt32(buf),
		WearItems:   ReadRWearItems(buf),
		Items8:      ReadRUnk1Items(buf),
		Data9:       ReadRUnkInfo4(buf),
	}

	return fullUserInfo
}

func ReadRFullCharacters(buf *bitstream.BitReader) []RFullCharacter {
	var characters []RFullCharacter

	charactersCount := int(ReadRUInt8(buf))

	for i := 0; i < charactersCount; i++ {
		character := RFullCharacter{}

		// TODO

		characters = append(characters, character)
	}

	return characters
}

func ReadRWearItems(buf *bitstream.BitReader) []RWearItem {
	var wearItems []RWearItem

	wearItemsCount := int(ReadRUInt8(buf))

	for i := 0; i < wearItemsCount; i++ {
		wearItem := RWearItem{}

		// TODO

		wearItems = append(wearItems, wearItem)
	}

	return wearItems
}

func ReadRUnk1Items(buf *bitstream.BitReader) []RUnk1Item {
	var unk1Items []RUnk1Item

	unk1ItemsCount := int(ReadRUInt8(buf))

	for i := 0; i < unk1ItemsCount; i++ {
		unk1Item := RUnk1Item{}

		// TODO

		unk1Items = append(unk1Items, unk1Item)
	}

	return unk1Items
}

func ReadRUnkInfo4(buf *bitstream.BitReader) RUnkInfo4 {
	unkInfo4 := RUnkInfo4{
		uint1:  ReadRUInt16(buf),
		uint2:  ReadRUInt16(buf),
		uint3:  ReadRUInt16(buf),
		uint4:  ReadRUInt16(buf),
		uint5:  ReadRUInt16(buf),
		uint6:  ReadRUInt16(buf),
		uint7:  ReadRUInt16(buf),
		uint81: ReadRUInt16(buf),
		uint9:  ReadRUInt16(buf),
	}

	// TODO 9 bits each...
	/*
		srpc::RpcUIntType<srpc::UInt16, 9> uint10;
		srpc::RpcUIntType<srpc::UInt16, 9> uint11;
		srpc::RpcUIntType<srpc::UInt16, 9> uint12;
		srpc::RpcUIntType<srpc::UInt16, 9> uint13;
		srpc::RpcUIntType<srpc::UInt16, 9> uint14;
	*/
	ReadRUInt16(buf)
	ReadRUInt16(buf)
	ReadRUInt16(buf)

	return unkInfo4
}

func WriteRBool(buf *bytes.Buffer, value bool) {
	binary.Write(buf, binary.LittleEndian, value)
}

func WriteRUInt8(buf *bytes.Buffer, value uint8) {
	binary.Write(buf, binary.LittleEndian, value)
}

func WriteRUInt16(buf *bytes.Buffer, value uint16) {
	binary.Write(buf, binary.LittleEndian, value)
}

func WriteRUInt32(buf *bytes.Buffer, value uint32) {
	binary.Write(buf, binary.LittleEndian, value)
}

func WriteRInt32(buf *bytes.Buffer, value int32) {
	binary.Write(buf, binary.LittleEndian, value)
}

func WriteRShortString(buf *bytes.Buffer, value string) {
	binary.Write(buf, binary.LittleEndian, uint8(len(value)))
	binary.Write(buf, binary.LittleEndian, []byte(value))
}

func WriteWString(buf *bytes.Buffer, value string) {
	WriteRShortString(buf, value)
}

func WriteRAddressPair(buf *bytes.Buffer, addressPair RAddressPair) {
	WriteRShortString(buf, addressPair.Ip)
	WriteRUInt16(buf, addressPair.Port)
}

func WriteRFullUserInfo(buf *bytes.Buffer, fullUserInfo RFullUserInfo) {
	WriteRUInt32(buf, fullUserInfo.UserId)
	WriteWString(buf, fullUserInfo.Nickname)
	WriteWString(buf, fullUserInfo.Str3)
	WriteRFullCharacters(buf, fullUserInfo.Characters)
	WriteRInt32(buf, fullUserInfo.SmashPoints)
	WriteRInt32(buf, fullUserInfo.SmashCash)
	WriteRWearItems(buf, fullUserInfo.WearItems)
	WriteRUnk1Items(buf, fullUserInfo.Items8)
	WriteRUnkInfo4(buf, fullUserInfo.Data9)
}

func WriteRFullCharacters(buf *bytes.Buffer, characters []RFullCharacter) {
	WriteRUInt8(buf, uint8(len(characters)))

	for i := 0; i < len(characters); i++ {
		// TODO
	}
}

func WriteRWearItems(buf *bytes.Buffer, wearItems []RWearItem) {
	WriteRUInt8(buf, uint8(len(wearItems)))

	for i := 0; i < len(wearItems); i++ {
		// TODO
	}
}

func WriteRUnk1Items(buf *bytes.Buffer, unk1Item []RUnk1Item) {
	WriteRUInt8(buf, uint8(len(unk1Item)))

	for i := 0; i < len(unk1Item); i++ {
		// TODO
	}
}

func WriteRUnkInfo4(buf *bytes.Buffer, unkInfo4 RUnkInfo4) {
	WriteRUInt16(buf, unkInfo4.uint1)
	WriteRUInt16(buf, unkInfo4.uint2)
	WriteRUInt16(buf, unkInfo4.uint3)
	WriteRUInt16(buf, unkInfo4.uint4)
	WriteRUInt16(buf, unkInfo4.uint5)
	WriteRUInt16(buf, unkInfo4.uint6)
	WriteRUInt16(buf, unkInfo4.uint7)
	WriteRUInt16(buf, unkInfo4.uint81)
	WriteRUInt16(buf, unkInfo4.uint9)
	// TODO 9 bits each...
	/*
		srpc::RpcUIntType<srpc::UInt16, 9> uint10;
		srpc::RpcUIntType<srpc::UInt16, 9> uint11;
		srpc::RpcUIntType<srpc::UInt16, 9> uint12;
		srpc::RpcUIntType<srpc::UInt16, 9> uint13;
		srpc::RpcUIntType<srpc::UInt16, 9> uint14;
	*/
	WriteRUInt16(buf, 0)
	WriteRUInt16(buf, 0)
	WriteRUInt16(buf, 0)
}

type RCharacterTemplate struct {
	// RCharacter
	CharacterID uint16 // 0 = empty, 1 = juri, 2 = sini, 3 = ai, 4 = van, 5 = hani
	Stats       [6]uint16
	Bool3       bool
	Visible     bool
	// RCharacterTemplate
	sint5  int32 // 3 Bits
	sint6  int32
	sint7  int32
	uint81 uint16
	uint9  uint16
}

type RCharacterStyle struct {
}

type RRawWearItem struct {
}

func WriteRCharacterTemplates(buf *bytes.Buffer, templates []RCharacterTemplate) {
	WriteRUInt8(buf, uint8(len(templates)))

	for i := 0; i < len(templates); i++ {
		// RCharacter
		template := templates[i]
		WriteRUInt16(buf, template.CharacterID)
		for iStats := 0; iStats < 6; iStats++ {
			WriteRUInt16(buf, template.Stats[iStats])
		}
		WriteRBool(buf, template.Bool3)
		WriteRBool(buf, template.Visible)
		// RCharacterTemplate
		WriteRInt32(buf, template.sint5) // 3 Bits
		WriteRInt32(buf, template.sint6)
		WriteRInt32(buf, template.sint7)
		WriteRUInt16(buf, template.uint81)
		WriteRUInt16(buf, template.uint9)
	}
}

func WriteRCharacterStyles(buf *bytes.Buffer, styles []RCharacterStyle) {
	WriteRUInt8(buf, uint8(len(styles)))

	for i := 0; i < len(styles); i++ {
		// TODO
	}
}

func WriteRRawWearItems(buf *bytes.Buffer, items []RRawWearItem) {
	WriteRUInt8(buf, uint8(len(items)))

	for i := 0; i < len(items); i++ {
		// TODO
	}
}

func ReadRCharacterTemplates(buf *bitstream.BitReader) []RCharacterTemplate {
	var templates []RCharacterTemplate

	templatesCount := int(ReadRUInt8(buf))

	for i := 0; i < templatesCount; i++ {
		template := RCharacterTemplate{
			CharacterID: ReadRUInt16(buf),
		}

		for iStats := 0; iStats < 6; iStats++ {
			template.Stats[iStats] = ReadRUInt16(buf)
		}

		template.Bool3 = ReadRBool(buf)
		template.Visible = ReadRBool(buf)

		println(ReadRBool(buf))
		println(ReadRBool(buf))
		println(ReadRBool(buf))

		template.sint5 = ReadRInt32Bits(buf, 3)
		template.sint6 = ReadRInt32(buf)
		template.sint7 = ReadRInt32(buf)
		template.uint81 = ReadRUInt16(buf)
		template.uint9 = ReadRUInt16(buf)

		templates = append(templates, template)
	}

	return templates
}

func ReadRCharacterStyles(buf *bitstream.BitReader) []RCharacterStyle {
	var styles []RCharacterStyle

	stylesCount := int(ReadRUInt8(buf))

	for i := 0; i < stylesCount; i++ {
		style := RCharacterStyle{}

		// TODO

		styles = append(styles, style)
	}

	return styles
}

func ReadRRawWearItems(buf *bitstream.BitReader) []RRawWearItem {
	var items []RRawWearItem

	itemsCount := int(ReadRUInt8(buf))

	for i := 0; i < itemsCount; i++ {
		item := RRawWearItem{}

		// TODO

		items = append(items, item)
	}

	return items
}

func Dump(data []byte) string {
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
