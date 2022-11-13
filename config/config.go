package config

import (
	"runtime"

	toml "github.com/pelletier/go-toml/v2"
)

func Default() Config {
	switch runtime.GOOS {
	case "windows":
		return Config{
			OpenCmd: []string{"Invoke-Item"},
		}
	default:
		return Config{
			OpenCmd: []string{"xdg-open"},
		}
	}
}

type Config struct {
	OpenCmd []string `toml:"open_cmd"`
}

func Unmarshal(data []byte, c *Config) error {
	return toml.Unmarshal(data, c)
}

func Marshal(c *Config) ([]byte, error) {
	return toml.Marshal(c)
}
