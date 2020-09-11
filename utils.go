package configer

import (
	"strings"
)

func getFileType(file string) string {
	if strings.HasSuffix(file, ".ini") {
		return "ini"
	}
	if strings.HasSuffix(file, ".toml") {
		return "toml"
	}
	if strings.HasSuffix(file, ".yaml") || strings.HasPrefix(file, ".yml") {
		return "yaml"
	}
	if strings.HasSuffix(file, ".json") {
		return "json"
	}
	return ""
}

func IsInSlice(data string, slice []string) bool {
	for _, s := range slice {
		if s == data {
			return true
		}
	}
	return false
}
