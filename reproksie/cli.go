package reproksie

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

//App parses the given arguments and starts a new reverse proxy service
type App struct {
	AppConfig
	flagSet *flag.FlagSet
	proxy   *reproksie
}

//AppConfig holds all configurable data such as usage, name, author, etc.
type AppConfig struct {
	Name    string
	Author  string
	Version string
	Usage   string
}

//NewApp creates a new App with the given config data.
func NewApp(config AppConfig) *App {
	a := &App{
		config,
		flag.NewFlagSet(
			config.Name,
			flag.ExitOnError,
		),
		newReproksie(),
	}
	return a
}

//Run starts a new reverse proxy service, while parsing the given arguments.
func (app *App) Run(args []string) error {
	configFile := app.flagSet.String("c", "", "The config file to be used when running Reproksie.")
	background := app.flagSet.Bool("b", false, "Starts the application in the background.")

	app.flagSet.Parse(args[1:])

	if len(*configFile) == 0 {
		return errors.New("Missing argument. A config file is required to run Reproksie")
	}

	ext := filepath.Ext(*configFile)
	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return err
	}

	var parser ConfigParser
	if ext == ".json" {
		parser = &JSONParser{}
	} else if ext == ".yml" {
		parser = &YamlParser{}
	} else {
		return fmt.Errorf("Unknown config file. Config file requires to be .json or .yml")
	}

	config, err := ParseConfig(parser, data)
	if err != nil {
		return err
	}
	if *background {
		go app.proxy.start(config)
	} else {
		err = app.proxy.start(config)
		if err != nil {
			return err
		}
	}
	return nil
}

//Shutdown gracefully shuts the proxy service down
func (app *App) Shutdown(ctx context.Context) {
	app.proxy.shutdown(ctx)
}
