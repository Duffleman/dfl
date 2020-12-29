package auth

import (
	"strings"
)

func Can(action string, scopes string) bool {
	blanket := map[string]struct{}{}

	for _, scope := range strings.Fields(scopes) {
		if action == scope {
			return true
		}

		if strings.HasSuffix(scope, ":*") {
			blanket[strings.TrimSuffix(scope, ":*")] = struct{}{}
		}
	}

	category := strings.Split(action, ":")[0]

	if _, ok := blanket[category]; ok {
		return true
	}

	return false
}
