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
