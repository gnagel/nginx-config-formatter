package formatter

import (
	"regexp"
	"strings"
)

// Strips the line and replaces neighbouring whitespaces with single space (except when within quotation marks).
func StripLine(line string) string {
	line = strings.TrimSpace(line)
	if IsCommentLine(line) {
		return line
	}

	wrappedInQuotes := false
	parts := []string{}
	splitLine := strings.Split(line, "\"")
	for _, part := range splitLine {
		if wrappedInQuotes {
			parts = append(parts, part)
		} else {
			re := regexp.MustCompile(`[\s]+`)
			part = string(re.ReplaceAll([]byte(part), []byte(" ")))
			parts = append(parts, part)
		}
		wrappedInQuotes = !wrappedInQuotes
	}

	line = strings.Join(parts, "\"")
	return line
}

func IsCommentLine(line string) bool {
	return strings.HasPrefix(line, "#")
}

func CleanLine(line string) []string {
	line = StripLine(line)
	line = EscapeVariables(line)
	if len(line) == 0 {
		return []string{""}
	} else if IsCommentLine(line) {
		line = UnescapeVariables(line)
		return []string{line}
	}

	q, c := NumStatmentsPerLine(line)
	if q == 1 && c > 1 {
		statements := SplitStatements(line)
		output := CleanLines(statements)
		return output
	} else if q != 1 && c > 1 {
		statements := strings.Split(line, ";")
		output := make([]string, 0, len(statements))
		for _, statement := range statements {
			statement = strings.TrimSpace(statement)
			if len(statement) != 0 {
				statement += ";"
				output = append(output, CleanLines([]string{statement})...)
			}
		}
		return output
	} else if strings.HasPrefix(line, "rewrite") {
		line = UnescapeVariables(line)
		return []string{line}
	}

	output := make([]string, 0, 1)
	re := regexp.MustCompile("([{}])")
	for _, subLine := range re.Split(line, -1) {
		subLine = strings.TrimSpace(subLine)
		if len(subLine) != 0 {
			subLine = UnescapeVariables(subLine)
			output = append(output, subLine)
		}
	}
	return output
}

// Strips the lines and splits them if they contain curly brackets.
func CleanLines(lines []string) []string {
	output := make([]string, 0, len(lines))
	for _, line := range lines {
		output = append(output, CleanLine(line)...)
	}
	return output
}
