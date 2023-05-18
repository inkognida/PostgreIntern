package utils

import (
	"github.com/fsnotify/fsnotify"
	"strings"
)

func GetOpString(op fsnotify.Op) string {
	switch {
	case op&fsnotify.Create == fsnotify.Create:
		return "Create"
	case op&fsnotify.Write == fsnotify.Write:
		return "Write"
	case op&fsnotify.Remove == fsnotify.Remove:
		return "Remove"
	case op&fsnotify.Rename == fsnotify.Rename:
		return "Rename"
	case op&fsnotify.Chmod == fsnotify.Chmod:
		return "Chmod"
	default:
		return "Unknown"
	}
}

func GetLastTextBefore(s string, sep string) string {
	lastIndex := strings.LastIndex(s, sep)
	if lastIndex == -1 {
		return ""
	}
	lastText := s[lastIndex+len(sep):]
	return lastText
}