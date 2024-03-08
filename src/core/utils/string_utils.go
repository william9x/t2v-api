package utils

import (
	"fmt"
	"github.com/Braly-Ltd/t2v-api-core/constants"
	"github.com/google/uuid"
	"slices"
	"strings"
)

func ContainsString(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func NewUUID() (string, error) {
	taskId, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return taskId.String(), nil
}

func BuildInferenceKey(queue, taskId string) string {
	return fmt.Sprintf("%s%s%s", queue, constants.TaskIDSeparator, taskId)
}

func ExtractInferenceKey(key string) (string, string) {
	ids := strings.Split(key, constants.TaskIDSeparator)
	if len(ids) != 2 {
		return "", ""
	}
	return ids[0], ids[1]
}

func RemoveDuplicate(list []string) []string {
	slices.Sort(list)
	return slices.Compact[[]string, string](list)
}
