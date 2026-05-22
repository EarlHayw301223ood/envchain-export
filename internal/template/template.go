// Package template renders a chain's environment variables into a
// user-supplied text template, substituting each {{ .VAR_NAME }} placeholder
// with the corresponding value from the chain.
package template

import (
	"bytes"
	"fmt"
	"io"
	"text/template"

	"github.com/nicholasgasior/envchain-export/internal/chain"
)

// ErrMissingKey is returned when the template references a key that is not
// present in the chain.
type ErrMissingKey struct {
	Key string
}

func (e *ErrMissingKey) Error() string {
	return fmt.Sprintf("template: key %q not found in scope", e.Key)
}

// Render executes tmplSrc as a Go text/template, providing all key-value pairs
// from c as the data map, and writes the result to w.
//
// Keys are accessed with {{ index . "VAR_NAME" }} or, for simple identifiers,
// {{ .VAR_NAME }}. Missing keys cause an error rather than a silent empty
// substitution.
func Render(w io.Writer, c *chain.Chain, tmplSrc string) error {
	data := make(map[string]string)
	for _, k := range c.Keys() {
		v, _ := c.Get(k)
		data[k] = v
	}

	t, err := template.New("envchain").
		Option("missingkey=error").
		Parse(tmplSrc)
	if err != nil {
		return fmt.Errorf("template: parse error: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("template: render error: %w", err)
	}

	_, err = w.Write(buf.Bytes())
	return err
}
