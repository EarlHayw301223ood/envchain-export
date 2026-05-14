// Package export serialises and deserialises environment variable chains
// to human-readable text formats.
//
// Supported formats:
//
//   - FormatPosix  — POSIX shell export statements (export KEY='value')
//   - FormatDotenv — .env file key=value pairs (KEY='value')
//
// Example usage:
//
//	c, _ := chain.New("myapp")
//	c.Add("DATABASE_URL", "postgres://localhost/db")
//
//	// Write to stdout in POSIX format
//	export.Write(os.Stdout, c, export.FormatPosix)
//
//	// Write to a file in dotenv format
//	f, _ := os.Create(".env")
//	defer f.Close()
//	export.Write(f, c, export.FormatDotenv)
package export
