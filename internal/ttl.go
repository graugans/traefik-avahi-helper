/*
SPDX-License-Identifier: Apache-2
Copyright © 2024 Christian Ege <ch@ege.io>
*/
package internal

import (
	"fmt"
	"strconv"
)

func TTLValueFromEnvironment() uint32 {
	const defaultTTLValue uint = 600
	ttlString := GetEnv(
		"CNAME_TTL",
		fmt.Sprint(defaultTTLValue),
	)
	ttl, err := strconv.ParseUint(ttlString, 10, 32)
	if err != nil {
		return uint32(defaultTTLValue)
	}
	return uint32(ttl)
}
