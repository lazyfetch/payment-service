package webhookapp

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Validate interface {
	ValidateWebhook(ctx context.Context, rawData []byte) error
}

func PaymentHandler(validate Validate) (pattern string, handler http.HandlerFunc) {
	return "/api/internal/govnokassa", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method)
		if r.Method != http.MethodPost {
			http.Error(w, "Only post!", http.StatusMethodNotAllowed)
			return
		}
		rawData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed!", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Валидируем чеез сервисный слой
		if err = validate.ValidateWebhook(r.Context(), rawData); err != nil {
			http.Error(w, "Failed!", http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "OK!")
		// Если ошибки нету возращаем значение OK
	}
}

type App struct {
	Log        *slog.Logger // not impl in func New()
	HTTPServer *http.Server
	Validate   Validate
}

func New(validate Validate, port int) *App {

	router := chi.NewRouter()
	router.Post(PaymentHandler(validate))

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

func (a *App) Stop(ctx context.Context) error {
	return a.HTTPServer.Shutdown(ctx)
}
