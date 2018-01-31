package config

const (
	Protocol = "tcp"
	NodeVersion int32 = 1
	CommandLength = 12
)

var (
	MainNodeIP = []byte{192,168,1,9}
	MainNodePort uint16 = 4000
)