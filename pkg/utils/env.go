package utils

import (
	"fmt"
	"os"
	"strconv"
)

// MustGetEnv will return the env or panic if it is not present
func MustGetEnv(key string) string {
	res := os.Getenv(key)

	if res == "" {
		panic(fmt.Sprintf("missing env key: %s", key))
	}

	return res
}

// MustGetEnvBool will return the env as boolean or panic if it is not present
func MustGetEnvBool(key string) bool {
	str := MustGetEnv(key)

	res, err := strconv.ParseBool(str)
	if err != nil {
		panic(fmt.Sprintf("can't parsing env bool key: %s", key))
	}

	return res
}
