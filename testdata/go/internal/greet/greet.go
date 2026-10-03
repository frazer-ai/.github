// Package greet is a fixture for the shared Go workflows.
package greet

import "github.com/google/uuid"

// Hello returns a greeting with a request ID.
func Hello(name string) string {
	return "hello " + name + " " + uuid.Must(uuid.NewV7()).String()
}
