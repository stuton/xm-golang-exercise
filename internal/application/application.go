package application

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"

	"github.com/rs/zerolog/log"

	"github.com/stuton/xm-golang-exercise/internal/handler"
	"github.com/stuton/xm-golang-exercise/internal/middleware"
	"github.com/stuton/xm-golang-exercise/internal/producer"
	"github.com/stuton/xm-golang-exercise/internal/repository"
	"github.com/stuton/xm-golang-exercise/internal/service"
	"github.com/stuton/xm-golang-exercise/utils"
)

type application struct {
	CompanyService service.CompanyService
	UserService    service.UserService
}

var app application

func New() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	utils.InitConfig()
	utils.InitLogger()

	connection, err := utils.NewDb(
		viper.GetString("DATABASE_URI"),
		utils.WithMaxOpenConnections(viper.GetInt("DATABASE_MAX_OPEN_CONN")),
		utils.WithMaxIdleConnectionsTime(viper.GetDuration("DATABASE_MAX_IDLE_TIME")),
		utils.WithMaxIdleConnections(viper.GetInt("DATABASE_MAX_IDLE_CONN")),
		utils.WithMaxOpenConnectionsTime(viper.GetDuration("DATABASE_MAX_OPEN_TIME")),
	)

	if err != nil {
		log.Fatal().Msgf("Unable to connect to database host: %v", err)
	}

	kafkaConnection, err := initKafkaProducer()

	if err != nil {
		log.Fatal().Msgf("Unable to connect to kafka hosts: %v", err)
	}

	app = application{
		CompanyService: service.NewCompanyService(
			repository.NewCompanyRepository(connection),
			producer.NewProducerProcessing(kafkaConnection),
		),
		UserService: service.NewUserService(
			repository.NewUserRepository(connection),
		),
	}

	srv := app.initRouters()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("Unable to run server: %s\n", err)
		}
	}()

	<-ctx.Done()

	log.Info().Msg("Shutting down server...")

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Msgf("Server forced to shutdown: %v", err)
	}

	if err := connection.Close(); err != nil {
		log.Error().Msgf("Unable close database connections: %v", err)
	}

	if err := kafkaConnection.Close(); err != nil {
		log.Error().Msgf("Unable close kafka connection: %v", err)
	}

	log.Info().Msg("Server exiting")

}

func initKafkaProducer() (*kafka.Conn, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", viper.GetString("KAFKA_HOST"), viper.GetString("KAFKA_TOPIC"), 0)

	if err != nil {
		log.Fatal().Msgf("failed to dial leader: %v", err)
	}

	return conn, err
}

func (app application) initRouters() *http.Server {
	gin.SetMode(viper.GetString("GIN_MODE"))

	r := gin.New()

	r.Use(logger.SetLogger(), gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		userGroup := v1.Group("/login")
		{
			userGroup.POST("", handler.NewUserLoginHandler(app.UserService).Login())
		}

		companiesGroup := v1.Group("/companies")
		{
			companiesGroup.GET("/:id", handler.NewGetCompanyByIDHandler(app.CompanyService).GetCompanyByID())
			// private methods
			authCompaniesGroup := companiesGroup.Group("")
			authCompaniesGroup.Use(middleware.JwtAuthMiddleware())

			authCompaniesGroup.POST("", handler.NewCreateCompanyHandler(app.CompanyService).CreateCompany())
			authCompaniesGroup.PATCH("/:id", handler.NewUpdateCompanyHandler(app.CompanyService).UpdateCompany())
			authCompaniesGroup.DELETE("/:id", handler.NewDeleteCompanyByIDHandler(app.CompanyService).DeleteCompanyByID())
		}
	}

	return &http.Server{
		Addr:    viper.GetString("APP_PORT"),
		Handler: r,
	}
}
