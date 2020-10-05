package server

import (
	"bytes"
	"github.com/romanitalian/img-generator/config"
	"github.com/romanitalian/img-generator/pkg/img"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func writeImage(w http.ResponseWriter, buffer *bytes.Buffer) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	writeImage(w, img.GenerateFavicon())
}

func robotsHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("robot"))
	if err != nil {
		log.Println(err)
	}
}

func imgHandler(w http.ResponseWriter, r *http.Request) {
	writeImage(w, img.Generate(strings.Split(r.URL.Path, "/")))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("PONG"))
	if err != nil {
		log.Println("pingHandler err: ", err)
	}
}

func Run(conf config.ConfI) {
	http.HandleFunc("/", imgHandler)
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/robots.txt", robotsHandler)

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		port := conf.GetPort()
		log.Println("Listening on " + port)

		if err := http.ListenAndServe(":"+port, nil); err != http.ErrServerClosed {
			log.Println("Error on start server: ", err)
		}
	}()

	signalValue := <-sigs
	signal.Stop(sigs)

	log.Println("stop signal: ", signalValue)
}
