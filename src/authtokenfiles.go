package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func storeAuthToken(filepath, token string) error {

	log.Printf("storeAuthToken: path: %q", filepath)

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("storeAuthToken: %s", err)
	}

	_, err = file.WriteString(token)
	if err != nil {
		return fmt.Errorf("storeAuthToken: %s", err)
	}

	return nil

}

func loadAuthToken(filepath string) (bool, string, error) {

	log.Printf("loadAuthToken: path: %q", filepath)

	file, err := os.Open(filepath)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("loadAuthToken: not found")
		return false, "", nil
	}
	if err != nil {
		return false, "", err
	}

	log.Printf("loadAuthToken: found")

	token, err := io.ReadAll(file)
	if err != nil {
		return false, "", fmt.Errorf("loadAuthToken: %s", err)
	}

	return true, string(token), nil
}
