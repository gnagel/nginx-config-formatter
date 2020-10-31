package formatter

import "strings"

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
