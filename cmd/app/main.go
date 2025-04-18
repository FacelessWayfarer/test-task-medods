package main

import (
	"context"

	"github.com/FacelessWayfarer/test-task-medods/internal/app"
	"github.com/FacelessWayfarer/test-task-medods/internal/config"
	"github.com/FacelessWayfarer/test-task-medods/pkg/logging"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logging.L(ctx).Info("config initializing")
	cfg := config.SetConfig()

	ctx = logging.ContextWithLogger(ctx, logging.NewLogger())

	Server := app.NewHTTPServer(ctx, cfg)
	logging.L(ctx).Info("Running Application")
	err := Server.ListenAndServe()
	if err != nil {
		logging.WithError(ctx, err).Fatal("main.ListenAndServe")
		return
	}
}
