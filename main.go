package main

import (
	"fmt"
	"os"

	"github.com/luukdegram/reproksie/reproksie"
)

func main() {
	app := reproksie.NewApp(reproksie.AppConfig{
		Name:    "Reproksie",
		Author:  "Luuk de Gram",
		Version: "0.1",
		Usage:   "Test",
	})
	err := app.Run()
	if err != nil {
		fmt.Println("An error occured: ", err)
		os.Exit(0)
	}
}
