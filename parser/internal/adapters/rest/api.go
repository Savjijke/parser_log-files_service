package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/savjijke/parser-log-files-service/internal/core"
)

const locationForLogger = "adapters/rest/"

func NewParseHandler(log *slog.Logger, parser core.ParserService) http.HandlerFunc {
	var req struct {
		Path string `json:"path"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		l := log.With(
			slog.String("location", locationForLogger+"NewParseHandler"),
		)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			l.Error("failed decode body", "err", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Error("failed to r.Body.Close()")
			}
		}()

		logID, err := parser.CreateID(ctx)
		if err != nil {
			l.Error("failed create id", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"log_id": logID,
		})

		go func(path string, id int) {
			ctxForParser, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()

			err := parser.Parse(ctxForParser, path, id)
			if err != nil {
				l.Error("failed parse", "err", err)
			}
		}(req.Path, logID)
	}
}

func NewNodeHandler(log *slog.Logger, parser core.ParserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := log.With(
			slog.String("location", locationForLogger+"NewNodeHandler"),
		)
		logIDString := r.URL.Query().Get("logID")
		logID, err := strconv.Atoi(logIDString)
		if err != nil {
			l.Error("failed convert param")
			http.Error(w, "failed convert param", http.StatusBadRequest)
			return
		}
		nodeID := strings.TrimPrefix(r.URL.Path, "/api/v1/node/")

		if nodeID == "" || strings.Contains(nodeID, "/") {
			l.Error("failed get node id")
			http.Error(w, "bad node id", http.StatusBadRequest)
			return
		}

		node, err := parser.GetDetailsNode(ctx, logID, nodeID)
		if err != nil {
			if errors.Is(err, core.ErrNotFound) {
				l.Warn("node not found", "err", err)
				http.Error(w, "node not found", http.StatusNotFound)
				return
			}
			l.Error("failed get details node", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		resp := toNodeDTO(logID, node)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Error("failed to encode")
		}

	}
}

func NewPortHandler(log *slog.Logger, parser core.ParserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := log.With(
			slog.String("location", locationForLogger+"NewPortHandler"),
		)
		logIDString := r.URL.Query().Get("logID")
		logID, err := strconv.Atoi(logIDString)
		if err != nil {
			l.Error("failed convert param")
			http.Error(w, "failed convert param", http.StatusBadRequest)
			return
		}
		nodeID := strings.TrimPrefix(r.URL.Path, "/api/v1/port/")

		if nodeID == "" || strings.Contains(nodeID, "/") {
			l.Error("failed get node id")
			http.Error(w, "bad node id", http.StatusBadRequest)
			return
		}

		ports, err := parser.GetPorts(ctx, logID, nodeID)
		if err != nil {
			if errors.Is(err, core.ErrNotFound) {
				l.Warn("port not found", "err", err)
				http.Error(w, "port not found", http.StatusNotFound)
				return
			}
			l.Error("failed get ports", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(ports)
		if err != nil {
			log.Error("failed to encode")
		}

	}
}

func NewLogHandler(log *slog.Logger, parser core.ParserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := log.With(
			slog.String("location", locationForLogger+"NewPortHandler"),
		)
		logIDString := strings.TrimPrefix(r.URL.Path, "/api/v1/log/")
		if logIDString == "" || strings.Contains(logIDString, "/") {
			l.Error("failed get log id")
			http.Error(w, "bad log id", http.StatusBadRequest)
			return
		}
		logID, err := strconv.Atoi(logIDString)
		if err != nil {
			l.Error("failed convert id")
			http.Error(w, "failed convert id", http.StatusBadRequest)
			return
		}
		info, err := parser.StatsFileLog(ctx, logID)
		if err != nil {
			if errors.Is(err, core.ErrNotFound) {
				l.Warn("node not found", "err", err)
				http.Error(w, "node not found", http.StatusNotFound)
				return
			}
			l.Error("failed get stats file log", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		resp := toFileLogDTO(info)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Error("failed to encode")
		}
	}
}

func NewTopologyHandler(log *slog.Logger, parser core.ParserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := log.With(
			slog.String("location", locationForLogger+"NewPortHandler"),
		)
		logIDString := strings.TrimPrefix(r.URL.Path, "/api/v1/topology/")
		if logIDString == "" || strings.Contains(logIDString, "/") {
			l.Error("failed get log id")
			http.Error(w, "bad log id", http.StatusBadRequest)
			return
		}
		logID, err := strconv.Atoi(logIDString)
		if err != nil {
			l.Error("failed convert id")
			http.Error(w, "failed convert id", http.StatusBadRequest)
			return
		}
		topology, err := parser.GetTopology(ctx, logID)
		if err != nil {
			l.Error("failed get topology", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		resp := mapTopologyToDTO(topology)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Error("failed to encode")
		}
	}
}
