package main

import (
	"flag"
	"github.com/romanitalian/pixel.local/img-generator/config"
	"github.com/romanitalian/pixel.local/img-generator/internal/server"
	"log"
)

var confPath = flag.String("config-file", "./config/.env", "Path to config file.")

func main() {
	conf, err := config.New(*confPath)
	if err != nil {
		log.Fatalln(err)
	}

	server.Run(conf)
}
