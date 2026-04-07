package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"time"
)

var startTime = time.Now()

type InfoResponse struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Pod      string `json:"pod"`
	Node     string `json:"node"`
	Uptime   string `json:"uptime"`
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()

	info := InfoResponse{
		Hostname: hostname,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Uptime:   time.Since(startTime).Round(time.Second).String(),
		Pod:      os.Getenv("POD_NAME"),
		Node:     os.Getenv("NODE_NAME"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS_ORIGIN"))

	if err := json.NewEncoder(w).Encode(info); err != nil {
		slog.Error("failed to encode info", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	slog.Info("info requested", "remote", r.RemoteAddr)
}

func handleReady(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /info", handleInfo)
	mux.HandleFunc("GET /ready", handleReady)

	slog.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
