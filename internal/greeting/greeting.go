package greeting

import (
	"errors"
	"fmt"
	"strings"
)

func Build(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", errors.New("name cannot be empty")
	}

	return fmt.Sprintf("Hello, %s!", trimmed), nil
}
