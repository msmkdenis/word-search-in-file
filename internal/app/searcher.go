package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/msmkdenis/word-search-in-file/internal/cache/memory"
	"github.com/msmkdenis/word-search-in-file/internal/config"
	"github.com/msmkdenis/word-search-in-file/internal/handler"
	"github.com/msmkdenis/word-search-in-file/internal/service"
)

func Run() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.New()
	idxCache := memory.NewIndexCache()
	searcher := service.NewSearcher(idxCache, cfg.FileWorkers)
	e := echo.New()
	handler.NewSearchHandler(e, searcher)

	// Запустили сервер HTTP
	go func() {
		errStart := e.Start(cfg.URLServer)
		if errStart != nil && !errors.Is(errStart, http.ErrServerClosed) {
			slog.Error(errStart.Error())
			os.Exit(1)
		}
	}()

	httpServerCtx, httpServerStopCtx := context.WithCancel(context.Background())

	// Канал для сигналов
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	quit := make(chan struct{})
	go func() {
		// Получили сигнал
		<-quitSignal
		// Закрыли сигнальный канал
		close(quit)
	}()

	go func() {
		// Слушаем сигнальный канал, при закрытии код идет дальше
		<-quit

		// Shutdown signal with grace period of 10 seconds
		shutdownCtx, cancel := context.WithTimeout(httpServerCtx, 10*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				slog.Error("graceful shutdown timed out.. forcing exit.")
				os.Exit(1)
			}
		}()

		// Trigger graceful shutdown
		logger.Info("Shutdown signal received, gracefully stopping http server")
		if errShutdown := e.Shutdown(shutdownCtx); errShutdown != nil {
			slog.Error("failed to shutdown http server", slog.String("error", errShutdown.Error()))
		}
		httpServerStopCtx()
	}()

	<-httpServerCtx.Done()
}
