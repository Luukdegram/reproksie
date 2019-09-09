package reproksie

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
)

//reproksie redirects a request to the correct internal application. This allows for serving applications to the internet without opening all ports.
type reproksie struct {
	Config
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
	errors := make(chan error)

	if len(rep.EntryPoints) == 0 {
		return fmt.Errorf("No entrypoints defined")
	}

	for _, entryPoint := range rep.EntryPoints {
		go func(entry *EntryPoint, c chan error) {
			var err error
			if entry.Protocol == Secure {
				err = http.ListenAndServeTLS(entry.Address, entry.TLS.CertFile, entry.TLS.KeyFile, rep)
			} else {
				err = http.ListenAndServe(entry.Address, rep)
			}
			if err != nil {
				c <- err
			} else {
				fmt.Printf("Started proxy server on address: %s", entry.Address)
			}
		}(entryPoint, errors)
	}

	if err := <-errors; err != nil {
		return err
	}

	return nil
}

//logRequest logs the request made if a path to a logfile was given
func (rep *reproksie) logRequest(r *http.Request) {
	if rep.Logger != nil {
		rep.Logger.Printf("\tHOST: %s \tPORT: %s \tMETHOD: %s \tPATH: %s \tIP: %s\n",
			r.Host,
			r.URL.Port(),
			r.Method,
			r.URL.Path,
			r.RemoteAddr)
	}
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
