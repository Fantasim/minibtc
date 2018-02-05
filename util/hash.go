package util

import (
	"crypto/sha256"
	"log"
	"golang.org/x/crypto/ripemd160"
)

func Sha256(data []byte) []byte {
	sha := sha256.Sum256(data)
	return sha[:]
}

func Ripemd160(data []byte) []byte {
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(data[:])
	if err != nil {
		log.Panic(err)
	}
	return RIPEMD160Hasher.Sum(nil)
}