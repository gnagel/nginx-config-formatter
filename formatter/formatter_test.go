package formatter

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestFormatBody(t *testing.T) {
	const utf8Sample = `
http {
    server {
        listen 80 default_server;
        server_name example.com;

        # redirect auf https://www....
        location / {
            return 301 https://www.example.com$request_uri;
        }

        # Statusseite für Monitoring freigeben
        # line above contains german umlaut causing problems
        location /nginx_status {
            stub_status on;
            access_log off;
            allow 127.0.0.1;
            deny all;
        }
    }
}`
	const latin1Sample = `
http {
    server {
        listen 80 default_server;
        server_name example.com;

        # redirect auf https://www....
        location / {
            return 301 https://www.example.com$request_uri;
        }

        # Statusseite für Monitoring freigeben
        # line above contains german umlaut causing problems
        location /nginx_status {
            stub_status on;
            access_log off;
            allow 127.0.0.1;
            deny all;
        }
    }
}`

	t.Run("Format UTF-8 body", func(t *testing.T) {
		input := utf8Sample
		output := FormatBody(input)
		expected := strings.ReplaceAll(utf8Sample, "    ", "\t")
		assert.Equal(t, expected, output)
	})

	t.Run("Format LATIN-1 body", func(t *testing.T) {
		input := latin1Sample
		output := FormatBody(input)
		expected := strings.ReplaceAll(latin1Sample, "    ", "\t")
		assert.Equal(t, expected, output)
	})
}

func TestIndentLines(t *testing.T) {
	t.Run("Indent lines", func(t *testing.T) {
		lines := strings.Split(
			`
server {
listen 80;
location / {
proxy_pass: http://localhost:123456;
}
}
`,
			"\n")
		expected := `
server {
	listen 80;
	location / {
		proxy_pass: http://localhost:123456;
	}
}
`
		actual := strings.Join(IndentLines(lines), "\n")
		assert.Equal(t, actual, expected)
	})
}
