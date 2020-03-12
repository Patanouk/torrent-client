package bencode

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const separator byte = ':'
const endSeparator byte = 'e'

//TODO Change return type
func Parse(reader *bufio.Reader) ([]interface{}, error) {
	var result []interface{}
	var err error

	for _, err = reader.Peek(1); err == nil; _, err = reader.Peek(1) {
		if err := parse(reader, &result); err != nil {
			return nil, err
		}

		if reader.Buffered() == 0 {
			return result, err
		}
	}

	return nil, err
}

func parse(reader *bufio.Reader, result *[]interface{}) error {
	c, err := reader.ReadByte()
	if err != nil {
		return err
	}

	switch {
	case c == 'i':
		i, err := decodeInteger(reader)
		if err != nil {
			return err
		}
		
		*result = append(*result, i)
	case c >= '0' && c <= '9':
		if err := reader.UnreadByte(); err != nil {
			return err
		}

		s, err := decodeString(reader)
		if err != nil {
			return err
		}
		*result = append(*result, s)
	case c == 'l':
		var list []interface{}
		for {
			if c, err = reader.ReadByte(); err != nil {
				return err
			}

			if c == 'e' {
				break
			}

			if err := reader.UnreadByte(); err != nil {
				return err
			}

			if err := parse(reader, &list); err != nil {
				return err
			}
		}
		*result = append(*result, list)
	default:
		return fmt.Errorf("unrecognized character : %c", c)
	}
	return nil
}


func decodeInteger(reader *bufio.Reader) (int64, error) {
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

