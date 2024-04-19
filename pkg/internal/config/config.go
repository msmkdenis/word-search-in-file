package config

import (
	"flag"
	"os"
)

type Config struct {
	URLServer string
	FSPath    string
}

func New() *Config {
	config := &Config{}

	config.ParseFlags()
	config.ParseENV()

	return config
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.URLServer, "s", "localhost:8080", "Enter URLServer as ip_address:port Or use SERVER_ADDRESS env")
	flag.StringVar(&c.FSPath, "f", "./examples", "Enter FSPath as path or use FS_PATH env")
}

func (c *Config) ParseENV() {
	if envURLServer := os.Getenv("SERVER_ADDRESS"); envURLServer != "" {
		c.URLServer = envURLServer
	}

	if envFSPath := os.Getenv("FS_PATH"); envFSPath != "" {
		c.FSPath = envFSPath
	}
}
