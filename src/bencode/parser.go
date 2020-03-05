package bencode

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const separator byte = ':'
const endSeparator byte = 'e'

func decodeInteger(reader *bufio.Reader) (int64, error) {
	iByte, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	if iByte != 'i' {
		return 0, fmt.Errorf("expected i as first rune when parsing integer, got %c", iByte)
	}

	slice, err := reader.ReadSlice(endSeparator)
	if err != nil {
		return 0, err
	}

	//TODO Verify slice non empty
	intResult, err := strconv.ParseInt(string(slice[:len(slice)-1]), 10, 64)
	if err != nil {
		return 0, err
	}

	return intResult, nil
}

func decodeString(reader *bufio.Reader) (string, error) {
	slice, err := reader.ReadSlice(separator)
	if err != nil {
		return "", err
	}

	//TODO Verify slice non empty
	length, err := strconv.ParseInt(string(slice[:len(slice) - 1]), 10, 64)
	if err != nil {
		return "", err
	}

	if _, err := reader.Peek(int(length)); err != nil {
		return "", err
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(reader, buf)
	return string(buf), err
}

