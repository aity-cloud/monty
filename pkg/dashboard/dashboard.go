package dashboard

import (
	"context"
	"crypto/tls"
	"errors"
	"io/fs"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"log/slog"

	"github.com/aity-cloud/monty/pkg/config/v1beta1"
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/util"
	"github.com/aity-cloud/monty/web"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
)

type Server struct {
	ServerOptions
	config *v1beta1.ManagementSpec
	logger *slog.Logger
}

type extraHandler struct {
	method  string
	prefix  string
	handler []gin.HandlerFunc
}

type ServerOptions struct {
	extraHandlers []extraHandler
	assetsFS      fs.FS
}

type ServerOption func(*ServerOptions)

func (o *ServerOptions) apply(opts ...ServerOption) {
	for _, op := range opts {
		op(o)
	}
}

func WithHandler(method, prefix string, handler ...gin.HandlerFunc) ServerOption {
	return func(o *ServerOptions) {
		o.extraHandlers = append(o.extraHandlers, extraHandler{
			method:  method,
			prefix:  prefix,
			handler: handler,
		})
	}
}

func WithAssetsFS(fs fs.FS) ServerOption {
	return func(o *ServerOptions) {
		o.assetsFS = fs
	}
}

func NewServer(config *v1beta1.ManagementSpec, opts ...ServerOption) (*Server, error) {
	options := ServerOptions{
		assetsFS: web.DistFS,
	}
	options.apply(opts...)

	if !web.WebAssetsAvailable(options.assetsFS) {
		return nil, errors.New("web assets not available")
	}

	if config.WebListenAddress == "" {
		return nil, errors.New("management.webListenAddress not set in config")
	}
	return &Server{
		ServerOptions: options,
		config:        config,
		logger:        logger.New().WithGroup("dashboard"),
	}, nil
}

func (ws *Server) ListenAndServe(ctx context.Context) error {
	lg := ws.logger
	var listener net.Listener
	if ws.config.WebCerts != nil {
		certs, caPool, err := util.LoadServingCertBundle(*ws.config.WebCerts)
		if err != nil {
			return err
		}
		listener, err = tls.Listen("tcp4", ws.config.WebListenAddress, &tls.Config{
			Certificates: []tls.Certificate{*certs},
			ClientCAs:    caPool,
		})
		if err != nil {
			return err
		}
	} else {
		var err error
		listener, err = net.Listen("tcp4", ws.config.WebListenAddress)
		if err != nil {
			return err
		}
	}
	lg.With(
		"address", listener.Addr(),
	).Info("ui server starting")
	webFsTracer := otel.Tracer("webfs")
	router := gin.New()
	router.Use(
		gin.Recovery(),
		logger.GinLogger(ws.logger),
		otelgin.Middleware("monty-ui"),
	)

	// Static assets
	sub, err := fs.Sub(ws.assetsFS, "dist")
	if err != nil {
		return err
	}

	webfs := http.FS(sub)

	router.NoRoute(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		_, span := webFsTracer.Start(ctx, c.Request.URL.Path)
		defer span.End()
		path := c.Request.URL.Path
		if path[0] == '/' {
			path = path[1:]
		}
		if _, err := fs.Stat(sub, path); err == nil {
			c.FileFromFS(path, webfs)
			return
		}

		c.FileFromFS("/", webfs) // serve index.html
	})

	montyApiAddr := ws.config.HTTPListenAddress
	mgmtUrl, err := url.Parse("http://" + montyApiAddr)
	if err != nil {
		lg.With(
			"url", montyApiAddr,
			"error", err,
		).Error("failed to parse management API URL")
		panic("failed to parse management API URL")
	}
	router.Any("/monty-api/*any", gin.WrapH(http.StripPrefix("/monty-api", httputil.NewSingleHostReverseProxy(mgmtUrl))))

	for _, h := range ws.extraHandlers {
		router.Handle(h.method, h.prefix, h.handler...)
	}

	return util.ServeHandler(ctx, router, listener)
}
