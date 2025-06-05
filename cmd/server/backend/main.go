package main

import (
	"log"
	"net/http"

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

	auth.FirstUser(cfg.Root.Email, cfg.Root.Username, logger)

	mux := http.NewServeMux()
	auth.AuthRoutes(mux)

	logger.Info("Server is started on http://localhost:8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		logger.Error("Server start error", "err", err)
	}
}
