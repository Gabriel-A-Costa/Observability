package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Gabriel-A-Costa/Observability/internal/config"
	"github.com/gin-gonic/gin"
)

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	router := gin.Default()
	router.GET("/health", health)

	router.Run(fmt.Sprintf(":%s", cfg.Port))
}
