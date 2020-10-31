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
