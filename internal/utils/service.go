package utils

import (
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func GetCmdArgs(fullCmd string) (string, string) {
	var cmd string
	var args string

	if idx := strings.IndexByte(fullCmd, ' '); idx >= 0 {
		cmd = fullCmd[:idx]
	} else {
		cmd = fullCmd
	}

	if idx := strings.IndexByte(fullCmd, ' '); idx >= 0 {
		args = fullCmd[idx+1:]
	} else {
		args = ""
	}

	return cmd, args
}

func CurrentDir(currentDir string, logger *logrus.Logger) bool {
	cwd, err := os.Getwd()
	if err != nil {
		logger.Fatalln(err)
	}

	if cwd == currentDir {
		return true
	}

	return false
}