package telegram

import (
	"strconv"
	"strings"
)

func int64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func parseCallbackPlantID(data, prefix string) (int64, bool) {
	if !strings.HasPrefix(data, prefix) {
		return 0, false
	}

	id, err := strconv.ParseInt(strings.TrimPrefix(data, prefix), 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}
