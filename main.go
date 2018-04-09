package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"
)

var base = "https://app.kanjialive.com"

func main() {
	logger := log.New()
	logger.Level, _ = log.ParseLevel("debug")
	log.SetLevel(logger.Level)

	logger.Formatter = &log.TextFormatter{ForceColors: false, FullTimestamp: true, TimestampFormat: "Jan 2 15:04:05.000"}
	log.SetFormatter(logger.Formatter)

	r := chi.NewRouter()
	r.Use(NewStructuredLogger(logger))
	r.Use(middleware.Recoverer)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Get("/api/*", func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", base+r.RequestURI, nil)

		req.Header.Set("accept", "application/json, text/plain, */*")
		req.Header.Set("referer", "https://app.kanjialive.com/")
		req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36")

		resp, err := client.Do(req)

		if err != nil {
			logger.Error(err)
			return
		}
		defer func() {
			if resp.Body != nil {
				resp.Body.Close()
			}
		}()

		if _, err := io.Copy(w, resp.Body); err != nil {
			logger.Error(err)
			return
		}
	})

	r.Method(http.MethodGet, "/*", http.FileServer(rice.MustFindBox("dist").HTTPBox()))

	http.ListenAndServe(":7780", r)
}

type StructuredLogger struct {
	Logger *log.Logger
}

func NewStructuredLogger(logger *log.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: log.NewEntry(l.Logger)}
	logFields := log.Fields{}

	logFields["ts"] = time.Now().UTC().Format(time.RFC1123)

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logFields["http_method"] = r.Method

	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	entry.Logger = entry.Logger.WithFields(logFields)

	entry.Logger.Infoln("request started")

	return entry
}

type StructuredLoggerEntry struct {
	Logger log.FieldLogger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, elapsed time.Duration) {
	l.Logger = l.Logger.WithFields(log.Fields{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
	})

	l.Logger.Infoln("request complete")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(log.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

func GetLogEntry(r *http.Request) log.FieldLogger {
	entry := middleware.GetLogEntry(r).(*StructuredLoggerEntry)
	return entry.Logger
}

func LogEntrySetField(r *http.Request, key string, value interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithField(key, value)
	}
}

func LogEntrySetFields(r *http.Request, fields map[string]interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithFields(fields)
	}
}
