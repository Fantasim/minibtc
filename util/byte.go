package util

func LenDoubleSliceByte(slice [][]byte) int {
	var size = 0
	for _, b := range slice {
		size += len(b)
	}
	return size
}