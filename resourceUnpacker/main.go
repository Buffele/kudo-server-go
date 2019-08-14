package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/rasky/go-lzo"
)

type FileInfos struct {
	CRC                  uint32
	Flag                 int8
	UncompressedFileSize int32
	CompressedFileSize   int32
	FolderIndex          int32
	FilePosition         int32
	Name                 string
}

/*
uint32_t CRC32C(uint32_t crc, const byte_t* buf, size_t len)
{
	crc = ~crc;
	while (len--)
	{
		crc ^= *buf++;
		for (int k = 0; k < 8; k++)
			crc = crc & 1 ? (crc >> 1) ^ 0xEDB88320 : crc >> 1;
	}
	return ~crc;
}

func CRC32C(uint32_t crc, const byte_t* buf, size_t len)
{
	crc = ~crc;
	while (len--)
	{
		crc ^= *buf++;
		for (int k = 0; k < 8; k++)
			crc = crc & 1 ? (crc >> 1) ^ 0xEDB88320 : crc >> 1;
	}
	return ~crc;
}
*/

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readSmashString(stream *bytes.Reader) string {
	var size uint16
	binary.Read(stream, binary.LittleEndian, &size)

	var data = make([]byte, size)
	binary.Read(stream, binary.LittleEndian, &data)

	return string(data)[:len(data)-1]
}

func main() {
	start := time.Now()

	file, err := os.Open("./input/Resource.npk")
	check(err)
	defer file.Close()

	stat, _ := file.Stat()

	reader := bufio.NewReader(file)

	data := make([]byte, stat.Size())
	_, err2 := reader.Read(data)
	check(err2)

	stream := bytes.NewReader(data)

	var magic uint16
	binary.Read(stream, binary.LittleEndian, &magic)

	var majorVersion uint8
	binary.Read(stream, binary.LittleEndian, &majorVersion)

	var minorVersion uint8
	binary.Read(stream, binary.LittleEndian, &minorVersion)

	var unk1 uint32
	binary.Read(stream, binary.LittleEndian, &unk1)

	var unk2 uint32
	binary.Read(stream, binary.LittleEndian, &unk2)

	date := readSmashString(stream)

	fmt.Println("--- HEADER START ---")
	fmt.Printf("Magic: %d\n", magic)
	fmt.Printf("Version: %d.%d\n", majorVersion, minorVersion)
	fmt.Printf("Unk 1: %d\n", unk1)
	fmt.Printf("Unk 2: %d\n", unk2)
	fmt.Printf("Date: %s\n", date)
	fmt.Println("--- HEADER END ---")

	fmt.Println("Reading Folder List...")

	var folderCount uint16
	binary.Read(stream, binary.LittleEndian, &folderCount)

	fmt.Printf("Folder Count: %d\n", folderCount)

	folderCount2 := int(folderCount)

	fmt.Println("--- FOLDER LIST START ---")

	paths := make([]string, folderCount2)
	for i := 0; i < folderCount2; i++ {
		path := readSmashString(stream)
		fmt.Printf("Path: %s\n", path)
		paths[i] = path
		os.MkdirAll("./output/"+path, os.ModePerm)
	}

	fmt.Println("--- FOLDER LIST END ---")

	var fileCount uint32
	binary.Read(stream, binary.LittleEndian, &fileCount)

	fmt.Printf("File Count: %d\n", fileCount)

	fmt.Println("--- FILE LIST START ---")

	fileCount2 := int(fileCount)

	files := make([]FileInfos, fileCount2)
	for i := 0; i < fileCount2; i++ {
		var crc uint32
		binary.Read(stream, binary.LittleEndian, &crc)

		var flag int8
		binary.Read(stream, binary.LittleEndian, &flag)

		var uncompressedFileSize int32
		binary.Read(stream, binary.LittleEndian, &uncompressedFileSize)

		var compressedFileSize int32
		binary.Read(stream, binary.LittleEndian, &compressedFileSize)

		var folderIndex int32
		binary.Read(stream, binary.LittleEndian, &folderIndex)

		var filePosition int32
		binary.Read(stream, binary.LittleEndian, &filePosition)

		name := readSmashString(stream)

		//fmt.Printf("CRC: %d, Flag: %d, FolderIndex: %d, UC: %d, Name: %s\n", crc, flag, folderIndex, uncompressedFileSize, name)

		var fileSize int = 0
		if flag == 1 {
			fileSize = int(compressedFileSize)
		} else if flag == 9 {
			fileSize = int(compressedFileSize)
		}

		file := make([]byte, fileSize)
		stream.ReadAt(file, int64(filePosition))

		files[i] = FileInfos{
			CRC:                  crc,
			Flag:                 flag,
			UncompressedFileSize: uncompressedFileSize,
			CompressedFileSize:   compressedFileSize,
			FolderIndex:          folderIndex,
			FilePosition:         filePosition,
			Name:                 name,
		}

		decoded, err := lzo.Decompress1X(bytes.NewReader(file), 0, 0)
		check(err)

		fullPath := ""
		if folderIndex == -1 {
			fullPath = "./output/" + name
		} else {
			fullPath = "./output/" + paths[folderIndex] + name
		}
		ioutil.WriteFile(fullPath, decoded, 0777)
	}

	fmt.Println("--- FILE LIST END ---")

	elapsed := time.Since(start)
	fmt.Printf("Time elapsed %s\n", elapsed)

	fmt.Println("Press enter key to exit...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}
