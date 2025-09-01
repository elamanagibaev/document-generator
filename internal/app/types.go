package app

import (
	"document-generator/internal/config"

	"github.com/gin-gonic/gin"
)

type App struct {
	engine *gin.Engine
	cfg    *config.Config
}
