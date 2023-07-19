package utils

import "os"

// EnvOrDefault is a function which can be used to retrieve an environment variable named k
// If the variable is not set, it will return d
func EnvOrDefault(k, d string) string {
	v := os.Getenv(k)
	if v != "" {
		return v
	}
	return d
}
