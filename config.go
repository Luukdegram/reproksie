package main

//ReproksieConfig holds all configuration data needed to setup Reproksie.
type ReproksieConfig struct {
	EntryPoints  []*EntryPoint
	Applications []*Application
	LogPath      string `json:"log_path"`
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
