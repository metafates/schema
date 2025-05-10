// Package schema is schema declaration and validation with static types.
// No field tags or code duplication.
//
// Schema is designed to be as developer-friendly as possible.
// The goal is to eliminate duplicative type declarations.
// You declare a schema once and it will be used as both schema and type itself.
// It's easy to compose simpler types into complex data structures.
package schema

//go:generate python3 validators.py
