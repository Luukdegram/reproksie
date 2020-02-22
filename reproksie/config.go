package reproksie

import (
	"encoding/json"
	"log"
	"os"

	"gopkg.in/yaml.v2"
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
	Domain   string
	Port     int
	Protocol Protocol
	Path     string
}

//TLS holds the certificate file paths
type TLS struct {
	CertFile string
	KeyFile  string
}

//Protocol is the web protocol used for the request. Either `http` or `https`.
type Protocol string

const (
	//Secure connection using TLS (https)
	Secure Protocol = "https"
	//NonSecure connection
	NonSecure Protocol = "http"
)

//ConfigParser is an interface that allows for parsing configurations based on their extension.
type ConfigParser interface {
	Parse(data []byte) (*Config, error)
}

//JSONParser is a JSON parser that unmarshals a json config file into Reproksie's config
type JSONParser struct{}

//YamlParser is a YAML parser that unmarshals a yaml config file into Reproksie's config
type YamlParser struct{}

//Parse parses the `json` data into a config file.
func (j *JSONParser) Parse(data []byte) (*Config, error) {
	var config Config

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil

}

//Parse parses the `yaml` data into a config file
func (y *YamlParser) Parse(data []byte) (*Config, error) {
	var config Config

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil

}

//ParseConfig parses the provided json data and sets all configuration that reproksie needs.
func ParseConfig(parser ConfigParser, data []byte) (*Config, error) {

	config, err := parser.Parse(data)
	if err != nil {
		return nil, err
	}

	if len(config.LogPath) != 0 {
		file, err := os.OpenFile(config.LogPath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		config.Logger = log.New(file, "", log.LstdFlags)
	}

	return config, nil
}
