package main

import (
	"errors"
	"flag"
	"io/ioutil"
)

//App parses the given arguments and starts a new reverse proxy service
type App struct {
	appConfig
}

//appConfig holds all configurable data such as usage, name, author, etc.
type appConfig struct {
	Name    string
	Author  string
	Version string
	Usage   string
}

//NewApp creates a new App with the given config data.
func NewApp(config appConfig) *App {
	a := &App{config}
	return a
}

//Run starts a new reverse proxy service, while parsing the given arguments.
func (app *App) Run() error {
	config := flag.String("c", "", "The config file to be used when running Reproksie.")
	flag.Parse()

	if len(*config) == 0 {
		return errors.New("Missing argument. A config file is required to run Reproksie")
	}

	data, err := ioutil.ReadFile(*config)
	if err != nil {
		return err
	}

	prox := newReproksie()
	err = prox.parseConfig(data)
	if err != nil {
		return err
	}
	err = prox.start()
	if err != nil {
		return err
	}
	return nil
}
