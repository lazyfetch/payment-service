package webhook

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type App struct {
	HTTPServer *http.Server
}

func New(port int) *App {

	router := chi.NewRouter()

	// make handler for router

	serverPort := fmt.Sprintf(":%d", port)

	httpServer := &http.Server{
		Addr:    serverPort,
		Handler: router,
	}

	return &App{
		HTTPServer: httpServer,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {

	return nil // temp
}

func (a *App) Stop() {

}
