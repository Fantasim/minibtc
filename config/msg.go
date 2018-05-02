package config

const (
	Protocol            = "tcp"
	NodeVersion   int32 = 1
	CommandLength       = 12
)

var (
	MainNodeIP   []byte
	MainNodePort uint16 = 4000
)
