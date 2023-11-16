package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/sivaprasadreddy/bookmarks/internal/api"
	"github.com/sivaprasadreddy/bookmarks/internal/config"
	"github.com/sivaprasadreddy/bookmarks/internal/domain"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type App struct {
	Router *gin.Engine
	Cfg    config.AppConfig
}

func NewApp(cfg config.AppConfig) *App {
	logger := config.NewLogger(cfg)
	gormDb := config.GetGormDb(cfg, logger)

	repo := domain.NewBookmarkRepository(gormDb, logger)
	bookmarkController := api.NewBookmarkController(repo, logger)

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))

	router.GET("/health", healthCheckHandler)

	router.GET("/api/bookmarks", bookmarkController.GetAll)
	router.GET("/api/bookmarks/:id", bookmarkController.GetById)
	router.POST("/api/bookmarks", bookmarkController.Create)
	router.PUT("/api/bookmarks/:id", bookmarkController.Update)
	router.DELETE("/api/bookmarks/:id", bookmarkController.Delete)

	return &App{
		Cfg:    cfg,
		Router: router,
	}
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func (app App) Run() {
	// Create a context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	port := fmt.Sprintf(":%d", app.Cfg.ServerPort)
	srv := &http.Server{
		Handler:        app.Router,
		Addr:           port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exiting")
}
