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
	router.Use(middleware.Metrics())

	router.GET("/health", health)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	logger.Info("Server started", zap.String("port", cfg.Port))

	router.Run(fmt.Sprintf(":%s", cfg.Port))
}
