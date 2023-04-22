package main

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
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


func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
