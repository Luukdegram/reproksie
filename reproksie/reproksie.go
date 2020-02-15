package reproksie

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strconv"
)

//reproksie redirects a request to the correct internal application. This allows for serving applications to the internet without opening all ports.
type reproksie struct {
	Config
	servers []*http.Server
}

//newReproksie creates a new Reproksie instance using default parameters
func newReproksie() *reproksie {
	return &reproksie{}
}

//ServeHTTP contains the proxy logic. It connects the configuration's applications and points a request towards it.
func (rep *reproksie) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Forwarded-Host", r.Host)
	defer rep.logRequest(r)

	for _, app := range rep.Applications {
		if app.match(*r.URL) {
			url, err := url.Parse(string(app.Protocol) + "://127.0.0.1:" + strconv.Itoa(app.Port))
			if err != nil {
				fmt.Println(err)
			}
			proxy := httputil.NewSingleHostReverseProxy(url)
			proxy.ServeHTTP(w, r)
			return
		}
	}
}

//start starts the reproksie service
func (rep *reproksie) start(config *Config) error {
	rep.Config = *config

	if rep.Logger == nil {
		rep.Logger = log.New(os.Stdout, "", 0)
	}

	if len(rep.EntryPoints) == 0 {
		return fmt.Errorf("No entrypoints defined")
	}

	for _, entry := range rep.EntryPoints {
		if entry.Protocol == Secure && (entry.TLS.CertFile == "" || entry.TLS.KeyFile == "") {
			return fmt.Errorf("No TLS certificates were found")
		}
	}

	rep.servers = make([]*http.Server, len(rep.EntryPoints))

	for index, entryPoint := range rep.EntryPoints {
		go func(entry *EntryPoint, index int) {
			var err error
			server := &http.Server{
				Addr:    entry.Address,
				Handler: rep,
			}
			rep.servers[index] = server
			if entry.Protocol == Secure {
				err = server.ListenAndServeTLS(entry.TLS.CertFile, entry.TLS.KeyFile)
			} else {
				err = server.ListenAndServe()
			}
			if err != nil {
				rep.Logger.Printf("An error occured: %v", err)
			}
		}(entryPoint, index)

		rep.Logger.Printf("Started proxy server on address: %s\n", entryPoint.Address)
	}

	return nil
}

//shutdown makes sure all services gracefully shutdown
func (rep *reproksie) shutdown(ctx context.Context) {
	rep.Logger.Println("Shutting down servers")

	for _, server := range rep.servers {
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			rep.Logger.Printf("An error occured while shutting down server on address: %s\n error: %v", server.Addr, err)
		}
		rep.Logger.Printf("Successfully shutdown server on address: %s", server.Addr)
	}
}

//logRequest logs the request made if a path to a logfile was given
func (rep *reproksie) logRequest(r *http.Request) {
	rep.Logger.Printf("\tHOST: %s \tPORT: %s \tMETHOD: %s \tPATH: %s \tIP: %s\n",
		r.Host,
		r.URL.Port(),
		r.Method,
		r.URL.Path,
		r.RemoteAddr)
}

//match checks if the URL of the request match the host or path of an application. Path can be a regex string
func (app *Application) match(url url.URL) bool {
	if url.Host == app.Domain {
		return true
	}

	var path string
	if path = url.Path; url.Path[0] == '/' {
		path = url.Path[1:]
	}

	if path == app.Path && len(app.Path) != 0 {
		return true
	}

	if match, err := regexp.Match(app.Path, []byte(path)); match && err == nil && len(app.Path) != 0 {
		return true
	}

	return false
}
