package model

type FileEvent struct {
	EventType string
	Path string
	FileName string
}

type CmdEvent struct {
	Cmd string
	Args string
}