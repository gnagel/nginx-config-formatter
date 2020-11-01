package formatter

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestIsCommentLine(t *testing.T) {
	t.Run("Has Comment", func(t *testing.T) {
		line := "# ..."
		assert.True(t, IsCommentLine(line))
	})

	t.Run("No Comment", func(t *testing.T) {
		line := "server_name ..."
		assert.False(t, IsCommentLine(line))
	})
}

func TestStripLine(t *testing.T) {
	t.Run("No Change", func(t *testing.T) {
		line := "server_name abc.xyz;"
		assert.Equal(t, line, StripLine(line))
	})

	t.Run("Changed", func(t *testing.T) {
		line := " server_name     abc.xyz;   "
		expected := "server_name abc.xyz;"
		assert.Equal(t, expected, StripLine(line))
	})
}

func TestCleanLine(t *testing.T) {
	t.Run("Blank Line", func(t *testing.T) {
		line := ""
		output := CleanLine(line)
		assert.Equal(t, []string{""}, output)
	})

	t.Run("1x Line", func(t *testing.T) {
		line := " server_name     abc.xyz;   "
		expected := []string{"server_name abc.xyz;"}
		output := CleanLine(line)
		assert.Equal(t, expected, output)
	})

	t.Run("2x Line", func(t *testing.T) {
		line := " server_name     abc.xyz;  listen ::80; "
		expected := []string{"server_name abc.xyz;", "listen ::80;"}
		output := CleanLine(line)
		assert.Equal(t, expected, output)
	})

	t.Run("Nested Lines", func(t *testing.T) {
		line := "server { server_name     abc.xyz;  listen ::80; location / { proxy_pass http://localhost:1234; } }"
		expected := []string{"server", "server_name abc.xyz;", "listen ::80;", "location /", "proxy_pass http://localhost:1234;", ";"}
		output := CleanLine(line)
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
		actual := strings.Join(IndentLines(lines, indentation), "\n")
		assert.Equal(t, actual, expected)
	})
}

func TestScrubBlankLines(t *testing.T) {
	t.Run("No consecutive blank lines", func(t *testing.T) {
		lines := []string{
			"server {",
			"listen ::80;",
			"}",
		}
		output := ScrubBlankLines(lines)

		expected := []string{
			"server {",
			"listen ::80;",
			"}",
			"",
		}
		assert.Equal(t, expected, output)
	})

	t.Run("Strip consecutive blank lines", func(t *testing.T) {
		lines := []string{
			"",
			"",
			"",
			"",
			"server {",
			"",
			"",
			"",
			"listen ::80;",
			"",
			"",
			"",
			"}",
			"",
			"",
			"",
			"",
		}
		output := ScrubBlankLines(lines)

		expected := []string{
			"server {",
			"",
			"listen ::80;",
			"",
			"}",
			"",
		}
		assert.Equal(t, expected, output)
	})
}
