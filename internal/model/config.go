package model

type Command struct {
	Cmd string `mapstructure:"cmd"`
}

type PathConfig struct {
	Path     string    `mapstructure:"path"`
	Commands []Command `mapstructure:"commands"`
}

type Config struct {
	Dirs []PathConfig `mapstructure:"paths"`
}
