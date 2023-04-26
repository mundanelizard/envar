package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

const saltSize = 16

func generateSalt() ([]byte, error) {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func GenHash(password string) (string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(password + hex.EncodeToString(salt)))
	return hex.EncodeToString(salt) + hex.EncodeToString(hash[:]), nil
}

func VerifyHash(password string, hashedPassword string) error {
	salt, err := hex.DecodeString(hashedPassword[:saltSize*2])
	if err != nil {
		return err
	}

	hash := sha256.Sum256([]byte(password + hex.EncodeToString(salt)))
	if hex.EncodeToString(hash[:]) != hashedPassword[saltSize*2:] {
		return errors.New("invalid password")
	}
	return nil
}
