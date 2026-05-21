// Package env provides the Run function for injecting decrypted scope
// variables into a subprocess.
//
// Variables stored in a [chain.Chain] are merged with the current process
// environment before the subprocess is launched. Chain variables take
// precedence over any identically named variables already present in the
// environment, allowing scoped secrets to shadow ambient configuration
// without permanently mutating the parent process environment.
package env
