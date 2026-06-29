package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Gabriel-A-Costa/Observability/internal/config"
	"github.com/Gabriel-A-Costa/Observability/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := config.NewLogger(cfg.Env)
	defer logger.Sync()

	router := gin.Default()

	// Middleware para geração de metricas no prometheus
	router.Use(middleware.Metrics())
	// Midlleware para geração de logs no loki
	router.Use(middleware.Logger(logger))

	router.GET("/health", health)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/error-500", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "simulated error 500"})
	})
	router.GET("/error-400", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "simulated error 404"})
	})

	logger.Info("Server started", zap.String("port", cfg.Port))

	router.Run(fmt.Sprintf(":%s", cfg.Port))
}
