package reproksie

import (
	"net/url"
	"testing"
)

func TestMatch(t *testing.T) {
	url, err := url.Parse("https://example.com/cool?id=1")
	if err != nil {
		t.Error(err)
	}

	app := &Application{
		Domain: "example.com",
	}

	if got := app.match(*url); !got {
		t.Error("Expected match on domain, but got no match.")
	}

	app = &Application{
		Path: "cool",
	}

	if got := app.match(*url); !got {
		t.Error("Expected match on path, but got no match.")
	}

	url, err = url.Parse("https://example.com/api/v1/somepath")
	app = &Application{
		Path: "api/v\\d\\/",
	}

	if got := app.match(*url); !got {
		t.Error("Expected match on regex, but got no match.")
	}

	app = &Application{
		Domain: "example1.com",
	}

	if got := app.match(*url); got {
		t.Error("Expected no match on domain, but got a match.")
	}
}

func TestStart(t *testing.T) {
	rep := newReproksie()

	err := rep.start(&Config{})
	if err == nil {
		t.Error("Expected error with no entrypoints, but got nil")
	}

	config := &Config{
		EntryPoints: []*EntryPoint{
			&EntryPoint{
				Protocol: Secure,
				Name:     "Test",
				Address:  ":8080",
			},
		},
	}

	err = rep.start(config)
	if err == nil {
		t.Error("Expected error with no tls certs, but got nil.")
	}
}
