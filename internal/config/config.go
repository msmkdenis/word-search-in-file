package config

import (
	"flag"
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	URLServer   string
	FileWorkers int
}

func New() *Config {
	config := &Config{}

	config.ParseFlags()
	config.ParseENV()

	return config
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.URLServer, "s", "localhost:8080", "Enter URLServer as ip_address:port Or use SERVER_ADDRESS env")
	flag.IntVar(&c.FileWorkers, "w", 5, "Enter number of workers as int Or use FILE_WORKERS env")
}

func (c *Config) ParseENV() {
	if envURLServer := os.Getenv("SERVER_ADDRESS"); envURLServer != "" {
		c.URLServer = envURLServer
	}

	if envFileWorkers := os.Getenv("FileWorkers"); envFileWorkers != "" {
		envFileWorkersInt, err := strconv.Atoi(envFileWorkers)
		if err != nil {
			slog.Info("Bad request: FileWorkers should be int, using default value")
			return
		}
		c.FileWorkers = envFileWorkersInt
	}
}
