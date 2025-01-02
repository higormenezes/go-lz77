package lz77

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
)

const SEARCH_SIZE int = 4000

func IndexOfSubSlice(s1 []byte, s2 []byte) int {
	if len(s2) > len(s1) {
		return -1
	}

	for startIdx := range len(s1) - len(s2) + 1 {
		match := true
		for compareIdx, s1Byte := range s1[startIdx : startIdx+len(s2)] {
			if s1Byte != s2[compareIdx] {
				match = false
				break
			}
		}

		if match {
			return startIdx
		}
	}

	return -1
}

func Compress(srcPath string, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	srcFileBuf := bufio.NewReader(srcFile)
	input := []byte{}
	search := []byte{}

	for {
		// Read byte by byte from the buffer
		currentByte, err := srcFileBuf.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// while the input is present on buffer, do nothing
		input = append(input, currentByte)
		inputIndexInSearch := IndexOfSubSlice(search, input)
		if inputIndexInSearch != -1 {
			continue
		}

		// write the input block to the result file
		output := []byte{}
		length := uint16(len(input) - 1)
		distance := uint16(0)
		if length != 0 {
			prevSearchIndex := IndexOfSubSlice(search, input[0:len(input)-1])
			distance = uint16(len(search) - prevSearchIndex)
		}
		output = binary.BigEndian.AppendUint16(output, distance)
		output = binary.BigEndian.AppendUint16(output, length)
		output = append(output, input[len(input)-1])
		_, err = destFile.Write(output)
		if err != nil {
			return err
		}

		// update search
		search = append(search, input...)
		search = search[uint(max(len(search)-SEARCH_SIZE, 0)):]
		// reset input
		input = nil
	}

	if len(input) > 0 {
		output := make([]byte, 5)
		length := uint16(len(input) - 1)
		distance := uint16(0)
		if length != 0 {
			prevSearchIndex := IndexOfSubSlice(search, input[0:len(input)-1])
			distance = uint16(len(search) - prevSearchIndex)
		}
		binary.BigEndian.PutUint16(output[0:2], distance)
		binary.BigEndian.PutUint16(output[2:4], distance)
		output[4] = input[len(input)-1]
		_, err = destFile.Write(output)
		if err != nil {
			return err
		}
	}

	return nil
}

func Decompress(srcPath string, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	srcFileBuf := bufio.NewReader(srcFile)

	chunk := make([]byte, 5)

	for {
		_, err = srcFileBuf.Read(chunk)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		distance := binary.BigEndian.Uint16(chunk[0:2])
		length := binary.BigEndian.Uint16(chunk[2:4])
		b := chunk[4]

		output := []byte{}
		if length != 0 {
			temp := make([]byte, length)
			fileStat, err := destFile.Stat()
			if err != nil {
				return err
			}

			_, err = destFile.ReadAt(temp, fileStat.Size()-int64(distance))
			if err != nil {
				return err
			}
			output = append(output, temp...)
		}
		output = append(output, b)

		_, err = destFile.Write(output)
		if err != nil {
			return err
		}
	}

	return nil
}
