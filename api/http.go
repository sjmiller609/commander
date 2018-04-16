package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/astronomerio/commander/utils"
)

var (
	log = logrus.WithField("package", "api")
)

type RouteHandler interface {
	Register(router *gin.Engine)
}

// Server for HTTP API
type HttpServer struct {
	server *gin.Engine
	handlers []RouteHandler
}


// NewClient creates a new API client
func NewHttp() *HttpServer {
	return &HttpServer{
		server: gin.Default(),
		handlers: make([]RouteHandler, 0),
	}
}

// AppendRouteHandler adds new available routes
func (s *HttpServer) AppendRouteHandler(rh RouteHandler) {
	logger := log.WithField("function", "AppendRouteHandler")
	logger.Debug("Appending new routes")

	s.handlers = append(s.handlers, rh)
}

// Serve starts the HTTP server
func (s *HttpServer) Serve(port string) {
	logger := log.WithField("function", "Serve")
	logger.Debug("Starting HTTP server")

	port = utils.EnsurePrefix(port, ":")


	for _, handler := range s.handlers {
		handler.Register(s.server)
	}

	s.server.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	go func() {
		err := s.server.Run(port)
		if err != nil {
			panic(err)
		}
	}()
}
