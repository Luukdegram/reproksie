package reproksie

import (
	"testing"
)

var jsonConfig = []byte(`
{
    "entrypoints": [
        {
            "name": "http",
            "address": ":8080"
        },
        {
            "name": "https",
            "address": ":4433",
            "protocol": "https",
            "tls": {
                "certFile": "example/test.crt",
                "keyFile": "example/test.key"
            }
        }
    ],
    "applications": [
        {
            "domain": "example.com",
            "port": 8080,
            "protocol": "http"
        }
	]
}
`)

var yamlConfig = []byte(`
entrypoints:
- name: http
  address: ":8080"
- name: https
  address: ":4433"
  protocol: https
  tls:
    certfile: example/test.crt
    keyfile: example/test.key
applications:
- port: 8080
  protocol: http
  domain: example.com
`)

func TestJsonParse(t *testing.T) {
	config, err := ParseConfig(&JSONParser{}, jsonConfig)
	expectedConfig := Config{
		EntryPoints: []*EntryPoint{
			&EntryPoint{Name: "http", Address: ":8080"},
			&EntryPoint{Name: "https", Address: ":4433", Protocol: "https", TLS: TLS{CertFile: "example/test.crt", KeyFile: "example/test.key"}}},
		Applications: []*Application{&Application{Domain: "example.com", Port: 8080, Protocol: "http"}},
	}

	if err != nil {
		t.Error(err)
	}

	if len(config.EntryPoints) != len(expectedConfig.EntryPoints) {
		t.Errorf("Incorrect length of entries, got: %d, want: %d.", len(config.EntryPoints), len(expectedConfig.EntryPoints))
	}

	if len(config.Applications) != len(expectedConfig.Applications) {
		t.Errorf("Incorrect length of applications, got: %d, want: %d.", len(config.Applications), len(expectedConfig.Applications))
	}

	if config.EntryPoints[1].TLS.CertFile != expectedConfig.EntryPoints[1].TLS.CertFile {
		t.Errorf("Incorrect certFile, got: %s, want: %s", config.EntryPoints[1].TLS.CertFile, expectedConfig.EntryPoints[1].TLS.CertFile)
	}

	if config.Applications[0].Port != expectedConfig.Applications[0].Port {
		t.Errorf("Incorrect port number, got: %d, want: %d", config.Applications[0].Port, expectedConfig.Applications[0].Port)
	}
}

func TestYamlParse(t *testing.T) {
	config, err := ParseConfig(&YamlParser{}, yamlConfig)
	expectedConfig := Config{
		EntryPoints: []*EntryPoint{
			&EntryPoint{Name: "http", Address: ":8080"},
			&EntryPoint{Name: "https", Address: ":4433", Protocol: "https", TLS: TLS{CertFile: "example/test.crt", KeyFile: "example/test.key"}}},
		Applications: []*Application{&Application{Domain: "example.com", Port: 8080, Protocol: "http"}},
	}

	if err != nil {
		t.Error(err)
	}

	if len(config.EntryPoints) != len(expectedConfig.EntryPoints) {
		t.Errorf("Incorrect length of entries, got: %d, want: %d.", len(config.EntryPoints), len(expectedConfig.EntryPoints))
	}

	if len(config.Applications) != len(expectedConfig.Applications) {
		t.Errorf("Incorrect length of applications, got: %d, want: %d.", len(config.Applications), len(expectedConfig.Applications))
	}

	if config.EntryPoints[1].TLS.CertFile != expectedConfig.EntryPoints[1].TLS.CertFile {
		t.Errorf("Incorrect certFile, got: %s, want: %s", config.EntryPoints[1].TLS.CertFile, expectedConfig.EntryPoints[1].TLS.CertFile)
	}

	if config.Applications[0].Port != expectedConfig.Applications[0].Port {
		t.Errorf("Incorrect port number, got: %d, want: %d", config.Applications[0].Port, expectedConfig.Applications[0].Port)
	}
}

func TestFailParseConfig(t *testing.T) {
	_, err := ParseConfig(&JSONParser{}, []byte(`
	{
		"entrypoints",
	}
	`))

	if err == nil {
		t.Errorf("Expected error in json parse but got nil")
	}

	_, err = ParseConfig(&YamlParser{}, []byte(`
	entrypoints
	`))

	if err == nil {
		t.Errorf("Expected error in yaml parse but got nil")
	}
}

func TestProtocol(t *testing.T) {
	got := Secure
	want := "https"

	if string(got) != want {
		t.Errorf("Incorrect protocol, got: %s, want: %s.", got, want)
	}

	got = NonSecure
	want = "http"

	if string(got) != want {
		t.Errorf("Incorrect protocol, got: %s, want: %s", got, want)
	}
}
