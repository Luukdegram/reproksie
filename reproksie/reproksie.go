package reproksie

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

//reproksie redirects a request to the correct internal application. This allows for serving applications to the internet without opening all ports.
type reproksie struct {
	Config
	servers chan error
}

//newReproksie creates a new Reproksie instance using default parameters
func newReproksie() *reproksie {
	return &reproksie{}
}

//ServeHTTP contains the proxy logic. It connects the configuration's applications and points a request towards it.
func (rep *reproksie) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Forwarded-Host", r.Host)
	host := r.URL.Hostname()
	defer func(repr *reproksie, r *http.Request) {
		if repr.Logger != nil {
			repr.Logger.Printf("\tHOST: %s \tPORT: %s \tMETHOD: %s \tPATH: %s \tIP: %s \n",
				r.Host,
				r.URL.Port(),
				r.Method,
				r.URL.Path,
				r.RemoteAddr)
		}
	}(rep, r)

	for _, app := range rep.Applications {
		if app.Host == host {
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

	for _, entry := range rep.EntryPoints {
		go func(entry *EntryPoint, c chan error) {
			var err error
			if entry.Protocol == secure {
				err = http.ListenAndServeTLS(entry.Address, entry.TLS.CertFile, entry.TLS.KeyFile, rep)
			} else {
				err = http.ListenAndServe(entry.Address, rep)
			}
			if err != nil {
				c <- err
			}
		}(entry, errors)
	}

	err := <-errors
	if err != nil {
		return err
	}

	return nil
}
