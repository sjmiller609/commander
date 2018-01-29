package api

import (
	"fmt"
	"net/http"

	"github.com/astronomerio/commander/api/routes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithField("package", "api")
)

// Client for API
type Client struct {
	handlers []routes.RouteHandler
}

// NewClient creates a new API client
func NewClient() *Client {
	return &Client{
		handlers: make([]routes.RouteHandler, 0),
	}
}

// AppendRouteHandler adds new available routes
func (c *Client) AppendRouteHandler(rh routes.RouteHandler) {
	logger := log.WithField("function", "AppendRouteHandler")
	logger.Debug("Appending new routes")

	c.handlers = append(c.handlers, rh)
}

// Serve starts the HTTP server
func (c *Client) Serve(port string) error {
	logger := log.WithField("function", "Serve")
	logger.Debug("Starting HTTP server")

	router := gin.Default()

	for _, handler := range c.handlers {
		handler.Register(router)
	}

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	return router.Run(fmt.Sprintf(":%v", port))
}
