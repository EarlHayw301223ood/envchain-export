// Package export provides functionality to export and import environment
// variable chains to and from shell-compatible formats.
package export

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envchain-export/internal/chain"
)

// Format represents an export output format.
type Format string

const (
	// FormatPosix outputs POSIX-compatible export statements.
	FormatPosix Format = "posix"
	// FormatDotenv outputs .env file compatible key=value pairs.
	FormatDotenv Format = "dotenv"
)

// Write serialises all variables in c to w using the given format.
func Write(w io.Writer, c *chain.Chain, format Format) error {
	vars := c.All()
	for _, kv := range vars {
		var line string
		switch format {
		case FormatDotenv:
			line = fmt.Sprintf("%s=%s\n", kv.Key, quoteValue(kv.Value))
		case FormatPosix:
			line = fmt.Sprintf("export %s=%s\n", kv.Key, quoteValue(kv.Value))
		default:
			return fmt.Errorf("unknown format: %q", format)
		}
		if _, err := io.WriteString(w, line); err != nil {
			return fmt.Errorf("write: %w", err)
		}
	}
	return nil
}

// quoteValue wraps v in single quotes, escaping any existing single quotes.
func quoteValue(v string) string {
	escaped := strings.ReplaceAll(v, "'", `'\''`)
	return "'" + escaped + "'"
}
