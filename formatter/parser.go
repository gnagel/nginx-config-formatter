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

	output := make([]string, 0)
	buffer := make([]rune, 0, len(line))
	withinQuote := false
	for _, char := range line {
		switch char {
		case '\'': fallthrough
		case '"':
			withinQuote = !withinQuote
			buffer = append(buffer, char)
		case ';':
			buffer = append(buffer, char)
			if !withinQuote {
				line := string(buffer)
				buffer = buffer[:0]
				output = append(output, line)
			}
		case '\n':
			line := string(buffer)
			buffer = buffer[:0]
			output = append(output, line)
		default:
			buffer = append(buffer, char)
		}
	}
	if len(buffer) > 0 {
		line := string(buffer)
		buffer = buffer[:0]
		output = append(output, line)
	}

	// Make sure each line has been cleaned up
	for index, line := range output {
		output[index] = strings.TrimSpace(line)
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