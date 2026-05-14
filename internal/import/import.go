// Package imp provides functionality for parsing environment variable
// definitions from various file formats into a Chain.
package imp

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/user/envchain-export/internal/chain"
	"github.com/user/envchain-export/internal/validate"
)

// Format represents the input file format to parse.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatPosix   Format = "posix"
)

// Read parses environment variable definitions from r in the given format
// and adds them to dst. Lines that are blank or start with '#' are ignored.
func Read(dst *chain.Chain, r io.Reader, format Format) error {
	scanner := bufio.NewScanner(r)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		var key, value string
		var err error
		switch format {
		case FormatDotenv:
			key, value, err = parseDotenvLine(line)
		case FormatPosix:
			key, value, err = parsePosixLine(line)
		default:
			return fmt.Errorf("unknown format: %q", format)
		}
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNum, err)
		}
		if err := validate.Pair(key, value); err != nil {
			return fmt.Errorf("line %d: %w", lineNum, err)
		}
		dst.Add(key, value)
	}
	return scanner.Err()
}

// parseDotenvLine parses a line in KEY=VALUE or KEY="VALUE" format.
func parseDotenvLine(line string) (string, string, error) {
	// Strip optional leading "export "
	line = strings.TrimPrefix(line, "export ")
	idx := strings.IndexByte(line, '=')
	if idx < 0 {
		return "", "", fmt.Errorf("missing '=' in %q", line)
	}
	key := line[:idx]
	value := unquote(line[idx+1:])
	return key, value, nil
}

// parsePosixLine parses a line in export KEY='VALUE' format.
func parsePosixLine(line string) (string, string, error) {
	line = strings.TrimPrefix(line, "export ")
	return parseDotenvLine(line)
}

// unquote strips surrounding single or double quotes and unescapes inner quotes.
func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '\'' && s[len(s)-1] == '\'') ||
			(s[0] == '"' && s[len(s)-1] == '"') {
			inner := s[1 : len(s)-1]
			if s[0] == '\'' {
				return strings.ReplaceAll(inner, "'\''", "'")
			}
			return strings.ReplaceAll(inner, `\"`, `"`)
		}
	}
	return s
}
