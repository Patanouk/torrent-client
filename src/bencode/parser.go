package bencode

import (
	"bufio"
	"io"
	"strconv"
)

const separator byte = ':'

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

