package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scul0405/saga-orchestration/cmd/order/config"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	middleware2 "github.com/scul0405/saga-orchestration/services/order/interface/http/middleware"
	"net/http"
)

type Server struct {
	config     config.HTTP
	logger     logger.Logger
	Engine     *gin.Engine
	Router     *Router
	httpServer *http.Server
}

func NewHTTPServer(config config.HTTP, logger logger.Logger, engine *gin.Engine, router *Router) *Server {
	return &Server{
		config: config,
		logger: logger,
		Engine: engine,
		Router: router,
	}
}

func NewEngine(config config.HTTP) *gin.Engine {
	gin.SetMode(config.Mode)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger()) // TODO: replace with custom logger
	engine.Use(middleware2.CORSMiddleware())

	return engine
}

func (srv *Server) InitRoutes() {
	mw := middleware2.NewJWTAuthMW(srv.Router.authSvc, srv.logger)

	apiGroup := srv.Engine.Group("/api/v1/")
	{
		orderGroup := apiGroup.Group("/orders")
		orderGroup.Use(mw.AuthMiddleware())
		{
			orderGroup.GET("/:id", srv.Router.GetDetailedOrder)
		}
	}
}

func (srv *Server) Run() error {
	srv.InitRoutes()

	srv.httpServer = &http.Server{
		Addr:    ":" + srv.config.Port,
		Handler: srv.Engine,
	}

	if err := srv.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (srv *Server) GracefulStop(ctx context.Context) error {
	return srv.httpServer.Shutdown(ctx)
}
