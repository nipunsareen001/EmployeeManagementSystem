package main

import (
	"Techiebulter/interview/backend/server"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	Send := make(chan os.Signal, 1)
	signal.Notify(Send, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	print("start")

	//SrvInit() initiates the PG database and  DBHandler
	srv := server.SrvInit()

	go srv.Start()

	<-Send
	//Gracefully stops all the services like Db and HTTP
	logrus.Info("Graceful shutdown")
	srv.Stop()
}
