package server

import (
	"bytes"
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/romanitalian/img-generate/v2/configs"
	"github.com/romanitalian/img-generate/v2/pkg/img"
	"github.com/romanitalian/img-generate/v2/pkg/logger"
)

var (
	log = logger.Get()
	// activeConnections tracks the number of active connections
	activeConnections int64
	// connectionsMu protects activeConnections
	connectionsMu sync.Mutex
)

// ConnectionTracker tracks active connections for graceful shutdown
type ConnectionTracker struct {
	handler http.Handler
}

func (ct *ConnectionTracker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	connectionsMu.Lock()
	activeConnections++
	currentConnections := activeConnections
	connectionsMu.Unlock()

	log.Debug().Int64("active_connections", currentConnections).Msg("new connection")

	defer func() {
		connectionsMu.Lock()
		activeConnections--
		currentConnections := activeConnections
		connectionsMu.Unlock()

		log.Debug().Int64("active_connections", currentConnections).Msg("connection closed")
	}()

	ct.handler.ServeHTTP(w, r)
}

func rend(w http.ResponseWriter, msg string) {
	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Error().Err(err).Msg("failed to write response")
	}
}

func rendImg(w http.ResponseWriter, buffer *bytes.Buffer) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Error().Err(err).Msg("failed to write image response")
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	buffer, err := img.GenerateFavicon(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate favicon")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	rendImg(w, buffer)
}

func imgHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	path := strings.TrimPrefix(r.URL.Path, "/img")
	buffer, err := img.Generate(ctx, strings.Split(path, "/"))
	if err != nil {
		log.Error().Err(err).Msg("failed to generate image")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	rendImg(w, buffer)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	rend(w, "OK")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	rend(w, "PONG")
}

func robotsHandler(w http.ResponseWriter, r *http.Request) {
	rend(w, "robots")
}

const (
	readTimeout     = 5 * time.Second
	writeTimeout    = 10 * time.Second
	shutdownTimeout = 30 * time.Second
	maxHeaderBytes  = 1 << 20 // 1 MB
)

func Run(ctx context.Context, conf configs.ConfI) error {
	// Create server with timeouts
	srv := &http.Server{
		Addr:           ":" + conf.GetPort(),
		Handler:        &ConnectionTracker{handler: setupRoutes()},
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	// Channel to wait for server shutdown
	serverErr := make(chan error, 1)

	// Start server
	go func() {
		log.Info().Str("port", conf.GetPort()).Msg("Server is starting...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Wait for either context cancellation or server error
	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		log.Info().Msg("Server is shutting down...")

		// Create shutdown context with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		// Disable keep-alives for new connections
		srv.SetKeepAlivesEnabled(false)

		// Wait for active connections to finish or timeout
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		go func() {
			for {
				select {
				case <-shutdownCtx.Done():
					return
				case <-ticker.C:
					connectionsMu.Lock()
					count := activeConnections
					connectionsMu.Unlock()
					log.Info().Int64("active_connections", count).Msg("waiting for connections to close")
				}
			}
		}()

		// Shutdown the server
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("Server shutdown error")
			return err
		}

		log.Info().Msg("Server stopped gracefully")
		return nil
	}
}

func setupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)
	mux.HandleFunc("/img/", imgHandler)
	mux.HandleFunc("/favicon.ico", faviconHandler)
	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/robots.txt", robotsHandler)
	return mux
}
