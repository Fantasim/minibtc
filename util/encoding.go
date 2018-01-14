package util

import (
	"strconv"
	"fmt"
)

func IntToArrayByte(n int) []byte {
	return []byte(strconv.Itoa(n))
}

func DupByteDoubleArray(a... []byte) [][]byte{
	byteSlice := make([][]byte, len(a))
	for i, b := range a {
		byteSlice[i] = b
	}
	return byteSlice
}

func ArrayIntToArrayByte(slice []int) []byte {
	if len(slice) == 0 {
		return nil
	}
	byteSlice := make([]byte, len(slice))
	for i, b := range byteSlice {
		byteSlice[i] = byte(b)
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

func IntToHex(n int) string {
	return fmt.Sprintf("%x", n)
}

func HexToInt(hex string) int {
	i, _ := strconv.ParseInt(hex, 16, 64)
	return int(i)
}

func IntArrayToByteDoubleArray(array []int) [][]byte{
	var ret [][]byte
	for _, n := range array {
		ret = append(ret, []byte(strconv.Itoa(n)))
	}
	return ret
}