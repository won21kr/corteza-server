package envoy

import (
	"os"
	"strconv"
)

func getEnvInt(name string, def int) int {
	if val, has := os.LookupEnv(name); has {
		if p, err := strconv.ParseInt(val, 10, 32); err == nil {
			return int(p)
		}
	}

	return def
}
