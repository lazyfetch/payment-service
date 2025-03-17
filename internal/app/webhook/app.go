package webhookapp

import (
	"fmt"
	"net/http"
	"payment/internal/webhook/handlers"

	"github.com/go-chi/chi/v5"
)

type App struct {
	HTTPServer *http.Server
}

func New(port int) *App {

	router := chi.NewRouter()
	router.Get(handlers.RobokassaHandler())

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

	return a.HTTPServer.ListenAndServe()
}

func (a *App) Stop() {

}
