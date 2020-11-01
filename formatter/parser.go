package formatter

import (
	"regexp"
	"strings"
)

// Lines can have multiple statements seperated by a semicolon, for example: "allow 127.0.0.1; allow 10.0.0.0/8; deny all;"
// This method counts the number of statements in a line and returns the number of quotes && number of statements.
func NumStatmentsPerLine(line string) (int, int) {
	line = strings.TrimSpace(line)
	if IsCommentLine(line) {
		return 0, 0
	}

	wrappedInQuotes := false
	quoteCount := 0
	statementCount := 0
	parts := strings.Split(line, "\"")
	for _, part := range parts {
		if wrappedInQuotes {
			quoteCount = 1
		} else {
			statementCount += strings.Count(part, ";")
		}
		wrappedInQuotes = !wrappedInQuotes
	}
	return quoteCount, statementCount
}

// This is a follow upto NumStatmentsPerLine, this splits the multi-statement line into separate lines
// TODO: This is failing on quoted blocks
func SplitStatements(line string) []string {
	line = strings.TrimSpace(line)
	if IsCommentLine(line) {
		return []string{line}
	}

	wrappedInQuotes := false
	parts := strings.Split(line, "\"")
	output := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if wrappedInQuotes {
			output = append(output, part)
		} else {
			for _, innerPart := range strings.Split(part, ";") {
				innerPart = strings.TrimSpace(innerPart)
				if len(innerPart) > 0 {
					innerPart += ";"
					output = append(output, innerPart)
				}
			}
		}
		wrappedInQuotes = !wrappedInQuotes
	}
	return output
}

// Replaces variable indicators ${ and } with tags, so subsequent formatting is easier.
func EscapeVariables(line string) string {
	re := regexp.MustCompile(`\${\s*(\w+)\s*}`)
	placeholder := varStart + "${1}" + varEnd
	output := re.ReplaceAll([]byte(line), []byte(placeholder))
	return string(output)
}

// Undoes the variable replacement from EscapeVariables
func UnescapeVariables(line string) string {
	line = strings.ReplaceAll(line, varStart, "${")
	line = strings.ReplaceAll(line, varEnd, "}")
	return line
}

// Replaces bracket { and } with tags, so subsequent formatting is easier.
func EscapeBlocks(body string) string {
	output := make([]rune, 0, len(body))
	inQuotes := false
	var lastCharacter rune

	for _, char := range []rune(body) {
		if (char == '\'' || char == '"') && lastCharacter != '\\' {
			inQuotes = !inQuotes
		}
		if !inQuotes {
			switch char {
			case '{':
				output = append(output, []rune(blockStart)...)
			case '}':
				output = append(output, []rune(blockEnd)...)
			default:
				output = append(output, char)
			}
		} else {
			output = append(output, char)
		}
		lastCharacter = char
	}
	return string(output)
}

// Undoes the replacement from EscapeBlocks
func UnescapeBlocks(body string) string {
	body = strings.ReplaceAll(body, blockStart, "{")
	body = strings.ReplaceAll(body, blockEnd, "}")
	return body
}

// When opening curly bracket is in it's own line (K&R convention), it's joined with precluding line (Java).
func MoveOpeningBracket(lines []string)[]string {
	output := make([]string, 0, len(lines))
	for index, line := range lines {
		line = strings.TrimSpace(line)
		if index > 0 && line == "{" {
			output[len(output)-1] += " {"
		} else {
			output = append(output, line)
		}
	}
	return output
}