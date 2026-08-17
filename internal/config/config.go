package config

import "os"

func Getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		v = def
	}
	return v
}
