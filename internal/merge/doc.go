// Package merge provides utilities for combining multiple envchain Chain
// instances into a single unified set of environment variables.
//
// # Conflict Policies
//
// When the same key exists in more than one source chain, the caller
// must choose a ConflictPolicy:
//
//   - PolicyError     – abort with an error on the first duplicate.
//   - PolicyOverwrite – later sources win; the last value is kept.
//   - PolicySkip      – first source wins; subsequent duplicates are ignored.
//
// # Example
//
//	base, _ := chain.New("base")
//	override, _ := chain.New("override")
//	dst, _ := chain.New("merged")
//
//	if err := merge.Merge(dst, merge.PolicyOverwrite, base, override); err != nil {
//		log.Fatal(err)
//	}
package merge
