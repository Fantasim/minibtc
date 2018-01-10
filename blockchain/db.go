package blockchain

import (
	"os"
)

func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}
