package application

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/stuton/xm-golang-exercise/internal/application/configuration"
	companieshandler "github.com/stuton/xm-golang-exercise/internal/companies/handler"
	companiesrepository "github.com/stuton/xm-golang-exercise/internal/companies/repository"
	companiesservice "github.com/stuton/xm-golang-exercise/internal/companies/service"
	"github.com/stuton/xm-golang-exercise/internal/producer"
	"github.com/stuton/xm-golang-exercise/internal/server/http/middleware"
	usershandler "github.com/stuton/xm-golang-exercise/internal/users/handler"
	usersrepository "github.com/stuton/xm-golang-exercise/internal/users/repository"
	usersservice "github.com/stuton/xm-golang-exercise/internal/users/service"
	"github.com/stuton/xm-golang-exercise/utils/database"
	"github.com/stuton/xm-golang-exercise/utils/jwt"
)

type routers struct {
	companiesHandlers companieshandler.CompaniesHandlers
	usersHandlers     usershandler.UsersHandlers
	jwtMiddleware     middleware.JwtAuthMiddleware
}

func Run(logger *zap.SugaredLogger) error {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config, err := configuration.Load()
	if err != nil {
		logger.Fatalf("unable to load configuration: %v", err)
	}

	connection, err := database.NewDb(
		config.DatabaseURI,
		database.WithLogger(zapadapter.NewLogger(logger.Desugar())),
		database.WithMaxOpenConnections(config.DatabaseMaxOpenConn),
		database.WithMaxIdleConnectionsTime(config.DatabaseMaxIdleTime),
		database.WithMaxIdleConnections(config.DatabaseMaxIdleConn),
		database.WithMaxOpenConnectionsTime(config.DatabaseMaxOpenTime),
	)

	if err != nil {
		logger.Fatalf("unable to connect to database host: %v", err)
	}

	defer func() {
		if err := connection.Close(); err != nil {
			logger.Errorf("unable close database connections: %v", err)
		}
	}()

	kafkaConnection, err := kafka.DialLeader(ctx, "tcp", config.KafkaHost, config.KafkaTopic, config.KafkaPortition)

	if err != nil {
		logger.Fatalf("Unable to connect to kafka hosts: %v", err)
	}

	defer func() {
		if err := kafkaConnection.Close(); err != nil {
			logger.Errorf("Unable close kafka connection: %v", err)
		}
	}()

	j := jwt.New(config)

	h := routers{
		companiesHandlers: companieshandler.NewCompaniesHandlers(
			companiesservice.NewCompanyService(
				logger,
				companiesrepository.NewCompanyRepository(connection),
				producer.NewProducerProcessing(kafkaConnection, config.KafkaTopic),
			),
			logger,
		),
		usersHandlers: usershandler.NewUsersHandlers(
			usersservice.NewUserService(
				logger,
				usersrepository.NewUserRepository(connection),
				j,
			),
			logger,
		),
		jwtMiddleware: middleware.NewJwtAuthMiddleware(j),
	}

	gin.SetMode(config.GinMode)

	server := &http.Server{
		Addr:    config.ServerPort,
		Handler: initRouters(logger, h),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("unable to run server: %s\n", err)
		}
	}()

	<-ctx.Done()

	logger.Info("Shutting down server...")

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exiting")

	return nil
}

func initRouters(logger *zap.SugaredLogger, routers routers) *gin.Engine {
	r := gin.New()

	r.Use(ginzap.Ginzap(logger.Desugar(), time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger.Desugar(), true))

	v1 := r.Group("/api/v1")
	{
		userGroup := v1.Group("/login")
		{
			userGroup.POST("", routers.usersHandlers.UserLoginHandler.Login())
		}

		companiesGroup := v1.Group("/companies")
		{
			companiesGroup.GET("/:id", routers.companiesHandlers.GetCompanyByIDHandler.GetCompanyByID())

			authCompaniesGroup := companiesGroup.Group("")
			authCompaniesGroup.Use(routers.jwtMiddleware.Do())

			authCompaniesGroup.POST("", routers.companiesHandlers.CreateCompanyHandler.CreateCompany())
			authCompaniesGroup.PATCH("/:id", routers.companiesHandlers.UpdateCompanyHandler.UpdateCompany())
			authCompaniesGroup.DELETE("/:id", routers.companiesHandlers.DeleteCompanyByIDHandler.DeleteCompanyByID())
		}
	}

	return r
}
