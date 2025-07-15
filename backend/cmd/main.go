package main

import (
	"context"
	"log"
	"time"

	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/leaderpro/infra/config"
	db "github.com/diegoclair/leaderpro/infra/data/mysql"
	"github.com/diegoclair/leaderpro/infra/shutdown"
	"github.com/diegoclair/leaderpro/internal/application/service"
	"github.com/diegoclair/leaderpro/internal/domain"
	"github.com/diegoclair/leaderpro/internal/transport/rest"
	"github.com/diegoclair/leaderpro/migrator/mysql"
)

const (
	gracefulShutdownTimeout = 10 * time.Second
	appName                 = "leaderpro"
)

func main() {
	ctx := context.Background()

	cfg, err := config.GetConfigEnvironment(ctx, appName)
	if err != nil {
		log.Fatalf("Error to load config: %v", err)
	}
	defer cfg.Close()

	log := cfg.GetLogger()

	infra := domain.NewInfrastructureServices(
		domain.WithCacheManager(cfg.GetCacheManager()),
		domain.WithDataManager(cfg.GetDataManager()),
		domain.WithLogger(log),
		domain.WithCrypto(cfg.GetCrypto()),
		domain.WithValidator(cfg.GetValidator()),
	)

	log.Info(ctx, "Running the migrations...")
	err = mysql.Migrate(cfg.GetDataManager().(*db.MysqlConn).DB())
	if err != nil {
		log.Errorw(ctx, "error to migrate mysql", logger.Err(err))
		return
	}
	log.Info(ctx, "Migrations completed successfully")
	
	// Alternative: With custom options using the options pattern
	// err = mysql.MigrateWithOptions(cfg.GetDataManager().(*db.MysqlConn).DB(),
	// 	sqlmigrator.WithDescriptionProcessor(func(filename, instruction string) string {
	// 		return "LEADER: " + sqlmigrator.ExtractActionDescription(instruction)
	// 	}),
	// )

	apps, err := service.New(infra, cfg.App.Auth.AccessTokenDuration)
	if err != nil {
		log.Errorw(ctx, "error to get domain services", logger.Err(err))
		return
	}

	server := rest.StartRestServer(ctx, cfg, infra, apps, appName, cfg.GetHttpPort())

	shutdown.GracefulShutdown(ctx, log, shutdown.WithRestServer(server.Router.Echo()))
}
