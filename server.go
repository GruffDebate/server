package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/bigokro/gruff-server/api"
	"github.com/bigokro/gruff-server/config"
)

func main() {

	config.Init()
	api.RW_DB_POOL = config.InitDB()
	api.RW_DB_POOL.LogMode(true)

	root := api.SetUpRouter(false, api.RW_DB_POOL)

	go func() {
		if err := root.Start(":8080"); err != nil {
			root.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := root.Shutdown(ctx); err != nil {
		root.Logger.Fatal(err)
	}
}
