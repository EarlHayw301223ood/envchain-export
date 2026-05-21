// Package lint provides validation checks for an entire chain, reporting
// suspicious or potentially problematic environment variable entries.
package lint

import (
	"fmt"
	"strings"

	"github.com/user/envchain-export/internal/chain"
)

// Severity indicates how serious a lint finding is.
type Severity string

const (
	Warn  Severity = "WARN"
	Error Severity = "ERROR"
)

// Finding represents a single lint result for a key.
type Finding struct {
	Key      string
	Severity Severity
	Message  string
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s: %s", f.Severity, f.Key, f.Message)
}

// Lint inspects every entry in c and returns a slice of findings.
// An empty slice means no issues were detected.
func Lint(c *chain.Chain) []Finding {
	var findings []Finding

	for _, k := range c.Keys() {
		v, _ := c.Get(k)

		// Warn about keys that look like they may contain secrets in the name
		// but have suspiciously short values.
		lower := strings.ToLower(k)
		if (strings.Contains(lower, "password") ||
			strings.Contains(lower, "secret") ||
			strings.Contains(lower, "token") ||
			strings.Contains(lower, "key")) && len(v) < 8 {
			findings = append(findings, Finding{
				Key:      k,
				Severity: Warn,
				Message:  "value appears too short for a secret",
			})
		}

		// Error on values that contain literal newlines (unsafe in many shells).
		if strings.ContainsRune(v, '\n') {
			findings = append(findings, Finding{
				Key:      k,
				Severity: Error,
				Message:  "value contains a newline character",
			})
		}

		// Warn when a value contains an unquoted dollar sign (possible
		// unintended variable expansion on export).
		if strings.ContainsRune(v, '$') {
			findings = append(findings, Finding{
				Key:      k,
				Severity: Warn,
				Message:  "value contains '$' which may cause unintended shell expansion",
			})
		}
	}

	return findings
}
