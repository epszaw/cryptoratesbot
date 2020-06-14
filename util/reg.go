package util

import "regexp"

var NumberRegexp = regexp.MustCompile(`^-?[0-9]+$`)

var TokenRegexp = regexp.MustCompile(`[A-Z]{3,4}`)
