package formatter

import (
	"strings"
)

func FormatBody(body string) string {
	body = EscapeBlocks(body)
	lines := strings.Split(body, "\n")
	lines = CleanLines(lines)
	lines = MoveOpeningBracket(lines)
	lines = IndentLines(lines)
	body = strings.Join(lines, "\n")
	body = UnescapeBlocks(body)
	return body
}
