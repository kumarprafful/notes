package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {
	if strings.Contains(err, "username") {
		return errors.New("username already taken")
	}
	if strings.Contains(err, "email") {
		return errors.New("email already taken")
	}
	return errors.New("incorrect details")
}
