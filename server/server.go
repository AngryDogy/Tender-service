package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tenderservice/handler"
)

type Server struct {
	router   *gin.Engine
	handlers []handler.Handler
}

func NewServer(router *gin.Engine, handlers []handler.Handler) *Server {
	return &Server{
		router:   router,
		handlers: handlers,
	}
}

func (s *Server) Run(serverPort string) error {
	s.router.GET("/api/ping", pingHandler)
	for _, h := range s.handlers {
		h.Setup(s.router)
	}
	return s.router.Run(serverPort)
}

func pingHandler(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "ok")
}
