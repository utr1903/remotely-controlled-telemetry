package app

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/utr1903/remotely-controlled-telemetry/apps/client/logger"
)

type httpServer struct {
	logger          *logger.Logger
	durationChannel chan time.Duration
}

func newHttpServer(
	logger *logger.Logger,
	durationChannel chan time.Duration,
) *httpServer {
	return &httpServer{
		logger:          logger,
		durationChannel: durationChannel,
	}
}

func (hs *httpServer) serve() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Run the server in a goroutine
	ctx, cancel := context.WithCancel(context.Background())
	go func(
		ctx context.Context,
	) {

		mux := http.NewServeMux()
		mux.HandleFunc("/latency", http.HandlerFunc(hs.handle))

		server := &http.Server{
			Addr:    "localhost:" + HTTP_SERVER_PORT,
			Handler: mux,
		}

		hs.logger.LogWithFields(
			logrus.InfoLevel,
			"HTTP server is running on localhost:"+HTTP_SERVER_PORT,
			map[string]string{
				"component.name": "httpserver",
			})

		err := server.ListenAndServe()
		if err != nil {
			hs.logger.LogWithFields(
				logrus.ErrorLevel,
				"HTTP server failed.",
				map[string]string{
					"component.name": "httpserver",
					"error.message":  err.Error(),
				})
		}

		<-ctx.Done()

		hs.logger.LogWithFields(
			logrus.InfoLevel,
			"Shutting down HTTP server...",
			map[string]string{
				"component.name": "httpserver",
			})
		server.Shutdown(ctx)

	}(ctx)

	// Wait for the context to be canceled or for the server to be shut down
	<-interrupt
	cancel()

	hs.logger.LogWithFields(
		logrus.InfoLevel,
		"HTTP server is shut down successfully.",
		map[string]string{
			"component.name": "httpserver",
		})
}

func (hs *httpServer) handle(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		msg := "HTTP request method is not allowed."
		hs.logger.LogWithFields(
			logrus.ErrorLevel,
			msg,
			map[string]string{
				"component.name":      "application",
				"http.request.method": r.Method,
			})

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(msg))
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		msg := "HTTP request body reading failed."
		hs.logger.LogWithFields(
			logrus.ErrorLevel,
			msg,
			map[string]string{
				"component.name": "application",
				"error.message":  err.Error(),
			})
		return
	}

	// Parse the request body
	var requestBody map[string]string
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		msg := "HTTP request body parsing failed."
		hs.logger.LogWithFields(
			logrus.ErrorLevel,
			msg,
			map[string]string{
				"component.name": "application",
				"error.message":  err.Error(),
			})
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	// Get requested duration
	durationAsString, ok := requestBody["duration"]
	if !ok {
		msg := "Duration is not defined."
		hs.logger.LogWithFields(
			logrus.ErrorLevel,
			msg,
			map[string]string{
				"component.name": "application",
				"error.message":  msg,
			})

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	// Convert to integer
	duration, err := strconv.ParseInt(durationAsString, 10, 64)
	if err != nil {
		msg := "Duration parsing failed."
		hs.logger.LogWithFields(
			logrus.ErrorLevel,
			msg,
			map[string]string{
				"component.name": "application",
				"error.message":  err.Error(),
			})

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
	}

	// Set the new duration
	hs.durationChannel <- time.Duration(duration) * time.Second

	msg := "Latency duration changed."
	hs.logger.LogWithFields(
		logrus.DebugLevel,
		msg,
		map[string]string{
			"component.name": "httpserver",
			"duration":       durationAsString,
		})
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
}
