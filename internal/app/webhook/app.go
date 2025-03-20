package webhookapp

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Validate interface {
	ValidateWebhook() error
}

func RobokassaHandler() (pattern string, handler http.HandlerFunc) {
	return "/api/robokassa", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Only post!", http.StatusMethodNotAllowed)
		}

	}
}

type App struct {
	HTTPServer *http.Server
	Validate   Validate
}

func New(validate Validate, port int) *App {

	router := chi.NewRouter()
	router.Get(RobokassaHandler())

	serverPort := fmt.Sprintf(":%d", port)

	httpServer := &http.Server{
		Addr:    serverPort,
		Handler: router,
	}

	return &App{
		HTTPServer: httpServer,
		Validate:   validate,
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
	// temp
}
