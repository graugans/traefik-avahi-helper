package internal

/*
 * This is based on the Stack Overflow answer: https://stackoverflow.com/a/40326580
 * By Łukasz Wojciechowski https://stackoverflow.com/users/1223977/%c5%81ukasz-wojciechowski
 */

// SPDX-License-Identifier: CC BY-SA 3.0

import "os"

// GetEnv reads the environment given with key.
// In case the variable is not defined the value
// of fallback is returned
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
