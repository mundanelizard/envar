package main

import (
	"crypto/rand"
	"encoding/hex"
)

func genRandomString() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	
	if err != nil {
		panic(err)
	}

	s := hex.EncodeToString(b)

	return s
}
