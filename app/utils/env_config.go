package utils

import (
	"strings"
)

func GetMapEnvConfig(envVarValue string) (res map[string]string) {
	res = make(map[string]string)
	for _, confString := range strings.Split(envVarValue, "|") {
		confPair := strings.Split(confString, ":")
		if len(confPair) < 2 {
			continue
		}
		res[confPair[0]] = confPair[1]
	}
	return
}
