package formatter

import (
	"strings"
)

// Indents the lines according to their nesting level determined by curly brackets.
func IndentLines(lines []string) []string {
	output := make([]string, 0, len(lines))
	indentDepth := 0
	for _, line := range lines {
		if !IsCommentLine(line) && (strings.HasSuffix(line, "}") || strings.HasSuffix(line, blockEnd)) && indentDepth > 0 {
			indentDepth -= 1
		}

		if len(line) != 0 {
			output = append(output, strings.Repeat(indentation, indentDepth)+line)
		} else {
			output = append(output, "")
		}

		if !IsCommentLine(line) && (strings.HasSuffix(line, "{") || strings.HasSuffix(line, blockStart)) {
			indentDepth += 1
		}
	}

	return output
}

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
