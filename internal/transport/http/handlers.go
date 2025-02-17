package http

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/invopop/ctxi18n/i18n"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"{{gitlab_url}}/internal/transport/http/middleware"
)

//go:embed static
var staticFiles embed.FS

type Exampler interface {
	Add(ctx context.Context, field string) error
}

type Server struct {
	Logger         *slog.Logger
	Config         ServerConfig
	Server         *http.Server
	ExampleService Exampler
}

type ServerConfig struct {
	Host          string
	Port          int
	Environment   string
	DefaultLocale i18n.Code
	AuthDisabled  bool
}

func NewServer(
	exampler Exampler,
	logger *slog.Logger,
	keyfunc jwt.Keyfunc,
	config ServerConfig,
) (*Server, error) {
	s := &Server{
		ExampleService: exampler,
		Logger:         logger,
		Config:         config,
	}

	fsys, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return nil, fmt.Errorf("failed to create embed file system: %w", err)
	}

	handler := s.setupHTTPRoutes(config, keyfunc, fsys)
	writeTimeout := 10
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port),
		ReadTimeout:  time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		Handler:      handler,
	}
	s.Server = httpServer

	return s, nil
}

func (s *Server) setupHTTPRoutes(config ServerConfig, keyfunc jwt.Keyfunc, fsys fs.FS) http.Handler {
	m := middleware.Middleware{
		DefaultLocale: config.DefaultLocale.String(),
		Logger:        s.Logger,
		Keyfunc:       keyfunc,
		DisableAuth:   config.AuthDisabled,
	}

	mux := http.NewServeMux()
	mux.Handle("/", m.Locale(http.HandlerFunc(s.indexHandler)))

	static := http.FileServer(http.FS(fsys))
	mux.Handle("/static/", http.StripPrefix("/static", static))

	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/readiness", s.readinessHandler)

	if config.Environment == "local" {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	httpSpanName := func(_ string, r *http.Request) string {
		return fmt.Sprintf("HTTP %s %s", r.Method, r.URL.Path)
	}

	otelFilters := func(r *http.Request) bool {
		return r.URL.Path != "/health" && r.URL.Path != "/readiness" && strings.HasPrefix(r.URL.Path, "/static")
	}

	routes := m.Logging(mux)

	handler := otelhttp.NewHandler(
		routes,
		"/",
		otelhttp.WithFilter(otelFilters),
		otelhttp.WithSpanNameFormatter(httpSpanName),
	)
	return handler
}

func (s *Server) Serve(ctx context.Context) error {
	s.Logger.InfoContext(ctx, "starting server")
	err := s.Server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.Logger.InfoContext(ctx, "shutting down server")
	err := s.Server.Shutdown(ctx)
	return err
}
