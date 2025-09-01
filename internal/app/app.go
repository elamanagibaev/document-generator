package app

import (
	"log"

	"github.com/gin-gonic/gin"

	"document-generator/internal/api/httpapi"
	"document-generator/internal/config"
	"document-generator/internal/infrastructure/auth"
	"document-generator/internal/service"
	"document-generator/pkg/gotenberg"
)

func New(cfg *config.Config) *App {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(gin.Recovery())

	authenticator := auth.NewStaticAuthenticator(cfg.StaticToken)
	r.Use(httpapi.RequestIDMiddleware())
	r.Use(httpapi.LoggerMiddleware(log.Printf))
	r.Use(httpapi.AuthMiddleware(authenticator))

	gotClient := gotenberg.NewClient("http://gotenberg:3000")

	// services
	docService := service.NewDocumentService("templates", gotClient)
	excelService := service.NewExcelService("templates")
	docxService := service.NewDocxService("templates")
	generatorService := service.NewGeneratorService(docService, excelService, docxService, gotClient)

	// routes
	httpapi.RegisterRoutes(r, generatorService)

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return &App{
		engine: r,
		cfg:    cfg,
	}
}

func (a *App) Run() error {
	addr := ":" + a.cfg.Port
	log.Printf("document-generator listening on %s", addr)
	return a.engine.Run(addr)
}
