package rpc

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrCookieNotFound    = errors.New("no cookie file found")
	ErrInvalidCookieFile = errors.New("cookie file not formatted correctly")
)

func FindCookie() (username, password string, err error) {
	var homeDir string
	homeDir, err = os.UserHomeDir()
	if err != nil {
		return
	}

	username, password, err = ReadCookieFile(filepath.Join(homeDir, ".bitcoin"))
	return
}

func ReadCookieFile(dataDir string) (username, password string, err error) {
	cookieContents, err := os.ReadFile(filepath.Join(dataDir, ".cookie"))
	if errors.Is(err, os.ErrNotExist) {
		newDataDir, err := readDataDirConfigParam(dataDir)
		if err != nil {
			return "", "", err
		}

		// Try data dir for cookie file
		return ReadCookieFile(newDataDir)
	} else if err != nil {
		return "", "", err
	}

	cookieParts := strings.SplitN(strings.TrimSpace(string(cookieContents)), ":", 2)
	if len(cookieParts) != 2 {
		return "", "", ErrInvalidCookieFile
	}

	username = cookieParts[0]
	password = cookieParts[1]
	return
}

func readDataDirConfigParam(dataDir string) (string, error) {
	configFileContents, err := os.ReadFile(filepath.Join(dataDir, "bitcoin.conf"))
	if err != nil {
		return "", ErrCookieNotFound
	}

	configFileLines := strings.Split(string(configFileContents), "\n")

	for _, line := range configFileLines {
		if strings.HasPrefix(line, "datadir=") {
			newDataDir := strings.TrimSpace(strings.TrimPrefix(line, "datadir="))
			return newDataDir, nil
		}
	}

	return "", ErrCookieNotFound
}
