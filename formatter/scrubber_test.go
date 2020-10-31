package formatter

import (
	"testing"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, StripLine(line), line)
	})

	t.Run("Changed", func(t *testing.T) {
		line := " server_name     abc.xyz;   "
		expected := "server_name abc.xyz;"
		assert.Equal(t, StripLine(line), expected)
	})
}
