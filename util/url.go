package util

import (
	"fmt"
	"strings"
)

func AppendQueryToUrl(url string, query map[string]string) string {
	if len(query) == 0 {
		return url
	}

	params := make([]string, 0, len(query))

	for key, value := range query {
		params = append(params, fmt.Sprintf("%s=%s", key, value))
	}

	return fmt.Sprintf("%s?%s", url, strings.Join(params, "&"))
}
