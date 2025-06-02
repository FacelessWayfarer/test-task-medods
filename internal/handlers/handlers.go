package handlers

import (
	"log"
	"os"

	"github.com/FacelessWayfarer/test-task-medods/internal/service"
)

type Handler struct {
	service service.IService
	logger  Logger
}

func NewHandler(service service.Service) *Handler {
	return &Handler{
		service: &service,
		logger:  log.New(os.Stdout, "test:", log.LstdFlags),
	}
}

type Logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
}
