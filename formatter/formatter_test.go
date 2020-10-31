package formatter

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

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
