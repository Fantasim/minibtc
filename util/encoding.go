package util

import (
	"strconv"
	"fmt"
)

func DupByteDoubleArray(a... []byte) [][]byte{
	byteSlice := make([][]byte, len(a))
	for i, b := range a {
		byteSlice[i] = b
	}
	return byteSlice
}

func ArrayByteToInt(slice []byte) (int, error) {
	if len(slice) == 0 {
		return 0, nil
	}
	a, err := strconv.Atoi(string(slice))
	return a, err
}

func ByteToInt(b byte) int {
	return int(b)
}

func IntToByte(n int) byte{
	return byte(n)
}

func EncodeInt(n int) []byte {
	h := fmt.Sprintf("%02x", n)
	return []byte(h)
}

func DecodeInt(d []byte) int {
	i, _ := strconv.ParseInt(string(d), 16, 64)
	return int(i)
}
