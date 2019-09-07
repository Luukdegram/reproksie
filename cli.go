package reproksie

import (
	"errors"
	"flag"
	"io/ioutil"
)

//App parses the given arguments and starts a new reverse proxy service
type App struct {
	AppConfig
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
	a := &App{config}
	return a
}

//Run starts a new reverse proxy service, while parsing the given arguments.
func (app *App) Run() error {
	configFile := flag.String("c", "", "The config file to be used when running Reproksie.")
	flag.Parse()

	if len(*configFile) == 0 {
		return errors.New("Missing argument. A config file is required to run Reproksie")
	}

	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return err
	}

	prox := newReproksie()
	config, err := ParseConfig(data)
	if err != nil {
		return err
	}
	err = prox.start(config)
	if err != nil {
		return err
	}
	return nil
}
