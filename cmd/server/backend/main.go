package main

import (
	"log"

	"github.com/Solnijko/go-knowledge-base/cmd/server/backend/auth"
	"github.com/Solnijko/go-knowledge-base/cmd/server/backend/pkg"
)

func main() {

	cfg, err := pkg.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := pkg.SetupLogger(cfg)
	logger.Info("logger initialized", "level", cfg.Logging.Level, "format", cfg.Logging.Format)

	dbconf := pkg.DBConfig{
		Host:         cfg.Database.Host,
		Port:         cfg.Database.Port,
		Name:         cfg.Database.Name,
		Username:     cfg.Database.Username,
		Password:     cfg.Database.Password,
		SSL:          cfg.Database.SSL,
		PoolMaxConns: cfg.Database.PoolMaxConns,
	}
	dbpool, err := pkg.CreateDBPool(dbconf, logger)

	if err != nil {
		logger.Error("failed to create database pool", "err", err)
		log.Fatal(err)
	}

	defer dbpool.Close()

	firstUserPassword, err := auth.GeneratePassword(12)
	if err != nil {
		logger.Error("failed to generate user password", "err", err)
		log.Fatal(err)
	}
	_, err = auth.FirstUser(cfg.Root.Email, cfg.Root.Username, firstUserPassword, logger, dbpool)
	if err != nil {
		logger.Error("failed to create first user", "err", err)
		log.Fatal(err)
	}

	// mux := http.NewServeMux()
	// auth.AuthRoutes(mux)

	// srv := &http.Server{
	// 	Addr:    ":8080",
	// 	Handler: mux,
	// }

	// logger.Info("server starting", "addr", srv.Addr)
	// if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 	logger.Error("server error", "err", err)
	// 	log.Fatal(err)
	// }
}
