package handler

import "github.com/gin-gonic/gin"

type Handler interface {
	Setup(router *gin.Engine)
}
