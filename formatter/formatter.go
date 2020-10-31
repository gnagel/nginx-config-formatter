package formatter

import "strings"


// Indents the lines according to their nesting level determined by curly brackets.
func IndentLines(lines []string) []string {
	output := make([]string, 0, len(lines))
	indentDepth := 0
	for _, line := range lines {
		if !IsCommentLine(line) && strings.HasSuffix(line, "}") && indentDepth > 0 {
			indentDepth -= 1
		}

		if len(line) != 0 {
			output = append(output, strings.Repeat(INDENTATION, indentDepth)+line)
		} else {
			output = append(output, "")
		}

		if !IsCommentLine(line) && strings.HasSuffix(line, "{") {
			indentDepth += 1
		}
	}

	return output
}
