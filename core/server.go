package core

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	Engine *gin.Engine
	config *Config
}

func NewHTTPServer(cfg *Config) *Server {
	engine := gin.Default()
	engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	// Swagger documentation endpoint
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &Server{
		config: cfg,
		Engine: engine,
	}
}

func (s *Server) Start() {

	h := &http.Server{
		Addr:    fmt.Sprintf(":%v", s.config.PORT),
		Handler: s.Engine,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		if err := h.Close(); err != nil {
			Log.Error("failed To ShutDown Server", zap.Error(err))
		}
		Log.Info("Shut Down Server")
	}()

	if err := h.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			Log.Info("Server Closed After Interruption")
		} else {
			Log.Error("Unexpected Server Shutdown", zap.Error(err))
		}
	}
}
