package httpapi

import (
	"net/http"

	"document-generator/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, generator *service.GeneratorService) {
	api := r.Group("/api/v1")

	api.POST("/generate", func(c *gin.Context) {
		var req service.GenerateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		filename, contentType, data, err := generator.Generate(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if req.Format != "html" {
			c.Header("Content-Disposition", "attachment; filename="+filename)
		}
		c.Data(http.StatusOK, contentType, data)
	})
}
