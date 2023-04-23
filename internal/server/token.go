package server

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/mundanelizard/envi/internal/lockfile"
)

func getTokenDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dir = path.Join(dir, ".envi.token")


	fmt.Println("temp dir:", dir)

	return dir, err
}

func (srv *Server) saveToken(token string) error {
	dir, err := getTokenDir()
	if err != nil {
		return err
	}
	
	lock := lockfile.New(dir) 
	err = lock.Hold()
	if err != nil {
		return err
	}

	lock.Write([]byte(token))
	lock.Commit()

	return nil
}

func (srv *Server) retrieveToken() (string, error) {
	dir, err := getTokenDir()
	if err != nil {
		return "", err
	}

	token, err := os.ReadFile(dir) 
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New("invalid auth: environment not unauthenticated")
		}
		return "", err
	}

	return string(token), nil
}
