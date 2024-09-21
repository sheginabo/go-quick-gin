package gin

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sheginabo/go-quick-gin/internal/presentation/handlers"
	"github.com/sheginabo/go-quick-gin/internal/presentation/middlewares"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

type Module struct {
	Router *gin.Engine
	Server *http.Server
}

type AllHandlers struct {
	InternalHandler *handlers.InternalHandler
}

func NewModule() *Module {
	//r := gin.Default()
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middlewares.Logger()) // good custom logger middleware

	gin.ForceConsoleColor()

	ginModule := &Module{
		Router: r,
	}

	ginModule.SetupRoute(ginModule.NewHandlers())

	return ginModule
}

func (module *Module) NewHandlers() AllHandlers {
	return AllHandlers{
		InternalHandler: handlers.NewInternalHandler(),
	}
}

func (module *Module) SetupRoute(allHandlers AllHandlers) {
	// basic check
	module.Router.GET("/health", handlers.HealthCheck)
	module.Router.GET("/test/ip", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"X-Forwarded-For": "ip:" + c.GetHeader("X-Forwarded-For"),
			"c_ClientIP":      "ip:" + c.ClientIP(),
		})
	})
	// basic handler
	module.Router.POST("/hello", allHandlers.InternalHandler.PostHello)
}

func enableCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ! XXX: This is a temporary solution to allow all origins
		// ! We should change this to allow only specific origins
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// Run gin
func (module *Module) Run(ctx context.Context, waitGroup *errgroup.Group) {
	module.Server = &http.Server{
		Addr:    viper.GetString("SERVER_ADDRESS"),
		Handler: module.Router,
	}

	module.Server.Handler = enableCors(module.Server.Handler)

	waitGroup.Go(func() error {
		log.Info().Msgf("Starting HTTP(Gin) server on %s\n", viper.GetString("SERVER_ADDRESS"))
		err := module.Server.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msg("HTTP(Gin) server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown HTTP(Gin) server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		shutdownErr := module.Server.Shutdown(shutdownCtx)
		if shutdownErr != nil {
			log.Error().Err(shutdownErr).Msg("failed to shutdown HTTP(Gin) server")
			return shutdownErr
		}

		log.Info().Msg("HTTP(Gin) server is stopped")
		return nil
	})
}
