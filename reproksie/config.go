package reproksie

import (
	"encoding/json"
	"log"
	"os"
)

//Config holds all configuration data needed to setup Reproksie.
type Config struct {
	EntryPoints  []*EntryPoint
	Applications []*Application
	LogPath      string `json:"log_path"`
	Logger       *log.Logger
}

//EntryPoint is the EntryPoint of a request. Reproksie will then serve the request to the correct application.
type EntryPoint struct {
	Name     string
	Address  string
	Protocol Protocol
	TLS      TLS
}

//Application can be any program running on a server that is accessable through a port. Reproksie redirects requests on the Application `host` from entrypoints to Application's `port`.
type Application struct {
	Host     string
	Port     int
	Protocol Protocol
}

//TLS holds the certificate files data
type TLS struct {
	CertFile string
	KeyFile  string
}

//Protocol is the web protocol used for the request. Either `http` or `https`.
type Protocol string

const (
	secure    Protocol = "https"
	nonSecure          = "http"
)

//ParseConfig parses the provided json data and sets all configuration that reproksie needs.
func ParseConfig(data []byte) (*Config, error) {
	var config Config
	err := json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if len(config.LogPath) != 0 {
		file, err := os.OpenFile(config.LogPath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		config.Logger = log.New(file, "Reproksie", log.LstdFlags)
	}

	return &config, nil
}
