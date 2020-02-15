package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/luukdegram/reproksie/reproksie"
)

func main() {
	proxy := reproksie.NewApp(reproksie.AppConfig{
		Name:    "Reproksie",
		Author:  "Luuk de Gram",
		Version: "0.1",
		Usage:   "Test",
	})

	err := proxy.Run(os.Args)

	if err != nil {
		fmt.Println("An error occured: ", err)
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	proxy.Shutdown(ctx)
}
