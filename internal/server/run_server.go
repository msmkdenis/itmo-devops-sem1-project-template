package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

func (s *Server) Run() {
	httpServerCtx, httpServerStopCtx := context.WithCancel(context.Background())
	defer httpServerStopCtx()

	go func() {
		if err := s.app.Listen(s.serverAddr); err != nil {
			slog.Error(err.Error())
			httpServerStopCtx()
		}
	}()

	go func() {
		<-s.quit

		slog.Info(fmt.Sprintf("Shutting down %s gracefully...", s.serverAddr))

		shutdownCtx, cancel := context.WithTimeout(httpServerCtx, s.graceTimout)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()

			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				slog.Error(fmt.Sprintf("%s graceful shutdown timed out.. forcing exit.", s.serverAddr))
				httpServerStopCtx()
			}
		}()

		if errShutdown := s.app.Shutdown(); errShutdown != nil {
			slog.Error(errShutdown.Error())
			httpServerStopCtx()
		}

		slog.Info(fmt.Sprintf("%s gracefully stopped.", s.serverAddr))
		httpServerStopCtx()
	}()

	<-httpServerCtx.Done()
	s.done <- struct{}{}
}
