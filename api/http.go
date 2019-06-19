package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/astronomer/commander/kubernetes"
	"github.com/astronomer/commander/utils"
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
	kubeClient *kubernetes.Client
	handlers []RouteHandler
}


// NewClient creates a new API client
func NewHttp(kubeClient *kubernetes.Client) *HttpServer {
	return &HttpServer{
		server: gin.Default(),
		kubeClient: kubeClient,
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
		version, err := s.kubeClient.ClientSet.Discovery().ServerVersion()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"kubeVersion" : "",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"kubeVersion" : fmt.Sprintf("%s", version),
		})
	})

	go func() {
		err := s.server.Run(port)
		if err != nil {
			panic(err)
		}
	}()
}
