package config

import (
	"github.com/joho/godotenv"
	"os"
)

const portDefault = "8080"

type ConfI interface {
	GetPort() string
}

type Conf struct {
	port string
}

func New(confPath string) (*Conf, error) {
	err := godotenv.Load(confPath)
	if err != nil {
		return nil, err
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = portDefault
	}

	return &Conf{port: port}, nil
}

func (c *Conf) GetPort() string {
	return c.port
}

