package server

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/romanitalian/img-generate/configs"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func newServer(conf configs.Configer) *server {
	rtr := fasthttprouter.New()
	return &server{
		HTTPServer: &fasthttp.Server{
			Handler:            rtr.Handler,
			ReadTimeout:        conf.ReadTimeout(),
			WriteTimeout:       conf.WriteTimeout(),
			MaxConnsPerIP:      conf.MaxConnsPerIP(),
			MaxRequestsPerConn: conf.MaxRequestsPerConn(),
		},
		router: rtr,
	}
}

type server struct {
	HTTPServer *fasthttp.Server
	router     *fasthttprouter.Router
}

func Run(conf configs.Configer) {
	srvr := newServer(conf)
	srvr.router.GET("/img/*params", imgHandler)
	srvr.router.GET("/favicon.ico", faviconHandler)
	srvr.router.GET("/ping", pingHandler)
	srvr.router.GET("/robots.txt", robotsHandler)
	srvr.router.GET("/user.json", userHandler)

	ln, err := reuseport.Listen("tcp4", "127.0.0.1:"+conf.Port())
	if err != nil {
		log.Fatal(err)
	}
	graceful := NewGracefulListener(ln, conf.GracefullListenerMaxWaitTimeOut())
	listenErr := make(chan error, 1)
	go func() {
		listenErr <- srvr.HTTPServer.Serve(graceful)
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-listenErr:
			if err != nil {
				log.Fatalf("listener error: %s", err)
			}
			log.Printf("Server stopped. Error to stop: %v\n", err)
			os.Exit(0)
		case sig := <-osSignals:
			log.Printf("Shutdown signal (app port: %s).\n", conf.Port())
			if err := graceful.Close(); err != nil {
				log.Fatal(err)
			}
			log.Printf("Server stopped. Signal to stop: %v\n", sig)
			os.Exit(0)
		}
	}
}
