package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"Alice-Seahat-Healthcare/seahat-be/config"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/database"
	"Alice-Seahat-Healthcare/seahat-be/libs/firebase"
	"Alice-Seahat-Healthcare/seahat-be/libs/mail"
	"Alice-Seahat-Healthcare/seahat-be/libs/rajaongkir"
	"Alice-Seahat-Healthcare/seahat-be/server"

	"github.com/sirupsen/logrus"
)

func main() {
	config.Load()
	logrus.SetReportCaller(true)

	db, err := database.ConnPostgres()
	if err != nil {
		logrus.Fatal(err)
	}

	firebase, err := firebase.New()
	if err != nil {
		logrus.Fatal(err)
	}

	ro, err := rajaongkir.New(config.RajaOngkir.URL, config.RajaOngkir.Key)
	if err != nil {
		logrus.Fatal(err)
	}

	dialer, err := mail.NewDialer()
	if err != nil {
		logrus.Fatal(err)
	}

	appLog, err := os.OpenFile("logs/app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatal(err)
	}

	defer appLog.Close()

	handler := server.NewServer(db, dialer, appLog, ro, firebase).SetupServer()
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", config.App.Port),
		Handler: handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	logrus.Info("Shutdown server...")
	ctx, cancel := context.WithTimeout(context.Background(), constant.TimeoutShutdown)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Server exited...")
}
