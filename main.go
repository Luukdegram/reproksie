package main

import (
	"fmt"
	"os"

	"github.com/luukdegram/reproksie/reproksie"
)

func main() {
	err := reproksie.NewApp(reproksie.AppConfig{
		Name:    "Reproksie",
		Author:  "Luuk de Gram",
		Version: "0.1",
		Usage:   "Test",
	}).Run(os.Args)

	if err != nil {
		fmt.Println("An error occured: ", err)
		os.Exit(2)
	}
}
