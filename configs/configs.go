package configs

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

type Configer interface {
	Port() string
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	MaxConnsPerIP() int
	MaxRequestsPerConn() int
	GracefullListenerMaxWaitTimeOut() time.Duration
}

type Conf struct {
	port                            string
	readTimeout                     uint64
	writeTimeout                    uint64
	maxConnsPerIP                   int
	maxRequestsPerConn              int
	gracefullListenerMaxWaitTimeOut uint64
}

func (c *Conf) Port() string {
	return c.port
}

func (c *Conf) ReadTimeout() time.Duration {
	return time.Duration(c.readTimeout) * time.Second
}

func (c *Conf) WriteTimeout() time.Duration {
	return time.Duration(c.writeTimeout) * time.Second
}

func (c *Conf) GracefullListenerMaxWaitTimeOut() time.Duration {
	return time.Duration(c.gracefullListenerMaxWaitTimeOut) * time.Second
}

func (c *Conf) MaxConnsPerIP() int {
	return c.maxConnsPerIP
}

func (c *Conf) MaxRequestsPerConn() int {
	return c.maxRequestsPerConn
}

func New(confPath string) (*Conf, error) {
	err := godotenv.Load(confPath)
	if err != nil {
		return nil, err
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		return nil, errors.New("empty port")
	}

	readTimeout, err := strconv.ParseUint(os.Getenv("SERVER_READ_TIMEOUT"), 10, 64)
	if err != nil {
		return nil, err
	}
	writeTimeout, err := strconv.ParseUint(os.Getenv("SERVER_WRITE_TIMEOUT"), 10, 64)
	if err != nil {
		return nil, err
	}
	maxConnsPerIP, err := strconv.Atoi(os.Getenv("SERVER_MAX_CONNS_PER_IP"))
	if err != nil {
		return nil, err
	}
	maxRequestsPerConn, err := strconv.Atoi(os.Getenv("SERVER_MAX_REQUESTS_PER_CONN"))
	if err != nil {
		return nil, err
	}
	gracefullListenerMaxWaitTimeOut, err := strconv.ParseUint(os.Getenv("SERVER_GRACEFUL_LISTENER_MAX_WAIT_TIMEOUT"), 10, 64)
	if err != nil {
		return nil, err
	}

	return &Conf{
		port:                            port,
		readTimeout:                     readTimeout,
		writeTimeout:                    writeTimeout,
		maxConnsPerIP:                   maxConnsPerIP,
		maxRequestsPerConn:              maxRequestsPerConn,
		gracefullListenerMaxWaitTimeOut: gracefullListenerMaxWaitTimeOut,
	}, nil
}
