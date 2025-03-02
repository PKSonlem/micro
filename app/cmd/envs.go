// в этом файле описана работа с переменными окружения, по надобности стоит добавлять сюда новые функции
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

func getStringFromEnv(envName string) (string, error) {
	env := os.Getenv(envName)
	if len(env) == 0 {
		return "", fmt.Errorf("env not found: %s", envName)
	}

	return env, nil
}

func getStringFromEnvOrDefault(envName string, defaultValue string) string {
	env, err := getStringFromEnv(envName)
	if err != nil {
		return defaultValue
	}

	return env
}

func getIntValueFromEnv(envName string, defaultValue int64) (int64, error) {
	sValue := os.Getenv(envName)

	if sValue != "" {
		value, err := strconv.ParseInt(sValue, 0, 32)
		if err != nil {
			return defaultValue, errors.Wrapf(err, "cant parse %s", envName)
		}

		return value, nil
	}

	return defaultValue, nil
}
