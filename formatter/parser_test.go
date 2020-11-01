package formatter

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNumStatmentsPerLine(t *testing.T) {
	t.Run("Single statement", func(t *testing.T) {
		line := "deny all;"
		quoteCount, lineCount := NumStatmentsPerLine(line)
		assert.Equal(t, quoteCount, 0)
		assert.Equal(t, lineCount, 1)
	})

	t.Run("2x statement", func(t *testing.T) {
		line := "allow 127.0.0.1; allow 10.0.0.0/8;"
		quoteCount, lineCount := NumStatmentsPerLine(line)
		assert.Equal(t, quoteCount, 0)
		assert.Equal(t, lineCount, 2)
	})

	t.Run("3x statement", func(t *testing.T) {
		line := "allow 127.0.0.1; allow 10.0.0.0/8; deny all;"
		quoteCount, lineCount := NumStatmentsPerLine(line)
		assert.Equal(t, quoteCount, 0)
		assert.Equal(t, lineCount, 3)
	})

	t.Run("3x statement with quotes", func(t *testing.T) {
		line := "allow \"127.0.0.1\"; allow 10.0.0.0/8; deny all;"
		quoteCount, lineCount := NumStatmentsPerLine(line)
		assert.Equal(t, quoteCount, 1)
		assert.Equal(t, lineCount, 3)
	})
}

func TestSplitStatements(t *testing.T) {
	t.Run("Single statement", func(t *testing.T) {
		line := "deny all;"
		output := SplitStatements(line)
		actual := strings.Join(IndentLines(output), "\n")
		expected := strings.Join([]string{"deny all;"}, "\n")
		assert.Equal(t, actual, expected)
	})

	t.Run("2x statement", func(t *testing.T) {
		line := "allow 127.0.0.1; allow 10.0.0.0/8;"
		output := SplitStatements(line)
		actual := strings.Join(IndentLines(output), "\n")
		expected := strings.Join([]string{"allow 127.0.0.1;", "allow 10.0.0.0/8;"}, "\n")
		assert.Equal(t, actual, expected)
	})

	t.Run("3x statement", func(t *testing.T) {
		line := "allow 127.0.0.1; allow 10.0.0.0/8; deny all;"
		output := SplitStatements(line)
		actual := strings.Join(IndentLines(output), "\n")
		expected := strings.Join([]string{"allow 127.0.0.1;", "allow 10.0.0.0/8;", "deny all;"}, "\n")
		assert.Equal(t, actual, expected)
	})

	t.Run("3x statement with quotes", func(t *testing.T) {
		line := "allow \"127.0.0.1\"; allow 10.0.0.0/8; deny all;"
		output := SplitStatements(line)
		actual := strings.Join(IndentLines(output), "\n")
		expected := strings.Join([]string{"allow \"127.0.0.1\";", "allow 10.0.0.0/8;", "deny all;"}, "\n")
		assert.Equal(t, actual, expected)
	})
}

func TestEscapeVariables(t *testing.T) {
	t.Run("No variables", func(t *testing.T) {
		line := "allow 127.0.0.1"
		output := EscapeVariables(line)
		assert.Equal(t, output, line)
	})

	t.Run("1x variable", func(t *testing.T) {
		line := "allow ${first}"
		output := EscapeVariables(line)
		expected := "allow __var_start__first__var_end__"
		assert.Equal(t, expected, output)
	})

	t.Run("2x variables", func(t *testing.T) {
		line := "allow ${first} ${second}"
		output := EscapeVariables(line)
		expected := "allow __var_start__first__var_end__ __var_start__second__var_end__"
		assert.Equal(t, expected, output)
	})
}

func TestUnescapeVariables(t *testing.T) {
	t.Run("No variables", func(t *testing.T) {
		line := "allow 127.0.0.1"
		output := UnescapeVariables(line)
		assert.Equal(t, output, line)
	})

	t.Run("1x variable", func(t *testing.T) {
		line := "allow __var_start__first__var_end__"
		output := UnescapeVariables(line)
		expected := "allow ${first}"
		assert.Equal(t, expected, output)
	})

	t.Run("2x variables", func(t *testing.T) {
		line := "allow __var_start__first__var_end__ __var_start__second__var_end__"
		output := UnescapeVariables(line)
		expected := "allow ${first} ${second}"
		assert.Equal(t, expected, output)
	})
}

func TestEscapeBlocks(t *testing.T) {
	t.Run("No blocks", func(t *testing.T) {
		body := "allow 127.0.0.1;"
		output := EscapeBlocks(body)
		assert.Equal(t, output, body)
	})

	t.Run("1x block", func(t *testing.T) {
		body := `
server {
	allow 127.0.0.1;
}
`
		output := EscapeBlocks(body)
		expected := `
server __block__start__
	allow 127.0.0.1;
__block__end__
`
		assert.Equal(t, expected, output)
	})

	t.Run("2x blocks", func(t *testing.T) {
		body := `
server {
	allow 127.0.0.1;
	location / {
		proxy_pass http://10.0.0.1;
	}
}
`
		output := EscapeBlocks(body)
		expected := `
server __block__start__
	allow 127.0.0.1;
	location / __block__start__
		proxy_pass http://10.0.0.1;
	__block__end__
__block__end__
`
		assert.Equal(t, expected, output)
	})
}

func TestUnescapeBlocks(t *testing.T) {
	t.Run("No blocks", func(t *testing.T) {
		body := "allow 127.0.0.1;"
		output := UnescapeBlocks(body)
		assert.Equal(t, output, body)
	})

	t.Run("1x block", func(t *testing.T) {
		body := `
server __block__start__
	allow 127.0.0.1;
__block__end__
`
		output := UnescapeBlocks(body)
		expected := `
server {
	allow 127.0.0.1;
}
`
		assert.Equal(t, expected, output)
	})

	t.Run("2x blocks", func(t *testing.T) {
		body := `
server __block__start__
	allow 127.0.0.1;
	location / __block__start__
		proxy_pass http://10.0.0.1;
	__block__end__
__block__end__
`
		output := UnescapeBlocks(body)
		expected := `
server {
	allow 127.0.0.1;
	location / {
		proxy_pass http://10.0.0.1;
	}
}
`
		assert.Equal(t, expected, output)
	})
}

func TestMoveOpeningBracket(t *testing.T) {
	t.Run("No changes", func(t *testing.T) {
		lines := strings.Split(
			`
server {
listen ::80;
}
`,
			"\n")
		output := MoveOpeningBracket(lines)
		expected := `
server {
listen ::80;
}
`
		assert.Equal(t, expected, strings.Join(output, "\n"))
	})

	t.Run("Trim whitespace", func(t *testing.T) {
		lines := strings.Split(
			`
server {
	listen ::80;
}
`,
			"\n")
		output := MoveOpeningBracket(lines)
		expected := `
server {
listen ::80;
}
`
		assert.Equal(t, expected, strings.Join(output, "\n"))
	})


	t.Run("Move Brackets", func(t *testing.T) {
		lines := strings.Split(
			`
server 
{
	listen ::80;
	location / 
	{
		proxy_pass http://10.0.0.1;
	}
}
`,
			"\n")
		output := MoveOpeningBracket(lines)
		expected := `
server {
listen ::80;
location / {
proxy_pass http://10.0.0.1;
}
}
`
		assert.Equal(t, expected, strings.Join(output, "\n"))
	})

}
