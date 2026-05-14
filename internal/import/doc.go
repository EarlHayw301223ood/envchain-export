// Package imp parses environment variable definitions from dotenv and
// POSIX shell export formats into a Chain.
//
// Supported formats:
//
//	"dotenv"  — KEY=VALUE or KEY="VALUE" lines, with optional leading "export "
//	"posix"   — export KEY='VALUE' lines as produced by the export package
//
// Blank lines and lines beginning with '#' are silently skipped. Each
// parsed key/value pair is validated via the validate package before being
// added to the destination Chain.
package imp
