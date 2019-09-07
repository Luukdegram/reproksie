package reproksie

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

//reproksie redirects a request to the correct internal application. This allows for serving applications to the internet without opening all ports.
type reproksie struct {
	ReproksieConfig
	servers chan error
	logger  *log.Logger
}

//newReproksie creates a new Reproksie instance using default parameters
func newReproksie() *reproksie {
	return &reproksie{}
}

//ServeHTTP contains the proxy logic. It connects the configuration's applications and points a request towards it.
func (reproksie *reproksie) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Forwarded-Host", r.Header.Get("Host"))
	host := r.URL.Hostname()
	for _, app := range reproksie.Applications {
		if app.Host == host {
			url, err := url.Parse(string(app.Protocol) + "://127.0.0.1:" + strconv.Itoa(app.Port))
			if err != nil {
				fmt.Println(err)
			}
			proxy := httputil.NewSingleHostReverseProxy(url)
			proxy.ServeHTTP(w, r)
			break
		}
	}
}

//start starts the reproksie service
func (reproksie *reproksie) start(config *ReproksieConfig) error {
	reproksie.ReproksieConfig = *config
	errors := make(chan error)

	if len(reproksie.EntryPoints) == 0 {
		return fmt.Errorf("No entrypoints defined")
	}

	for _, entry := range reproksie.EntryPoints {
		go func(entry *EntryPoint, c chan error) {
			var err error
			if entry.Protocol == secure {
				err = http.ListenAndServeTLS(entry.Address, entry.TLS.CertFile, entry.TLS.KeyFile, reproksie)
			} else {
				err = http.ListenAndServe(entry.Address, reproksie)
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

//LogHandler implements the ServeHTTP function to log each request to the given writer. I.E. a file writer.
type LogHandler struct {
	writer io.Writer
	next   http.Handler
}

//ServeHTTP is a http.Handler interface function that serves the request. In this case, adds a log before calling the next handler.
func (lh *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(lh.writer, "%s \tHOST: %s \tPORT: %s \tMETHOD: %s \tPATH: %s \tIP: %s \n",
		time.Now().UTC().Format("2006-01-02 15:04:05"),
		r.Host,
		r.URL.Port(),
		r.Method,
		r.URL.Path,
		r.RemoteAddr)
	lh.next.ServeHTTP(w, r)
}

/*
func test() {
	file, err := os.OpenFile(config.LogPath,
        os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(fmt.Errorf("Could not open file. Error: %s", err))
	}

	defer file.Close()

	remote, err := url.Parse("http://google.com")
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	w := io.MultiWriter(os.Stdout, file)
	http.Handle("/", &LogHandler{w, &Reproksie{proxy}})
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}
}
*/
