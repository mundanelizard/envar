package server

import "fmt"

func CheckAuthentication() (bool, error) {
	return true, nil
}

func CreateNewRepo(name string) (string, error) {
	return fmt.Sprintf("https://localhost:8080/repos/%s", name), nil
}
