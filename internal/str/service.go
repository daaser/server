package str

import (
	"errors"
	"strings"
)

type Service interface {
	Uppercase(string) (string, error)
	Count(string) int
}

var ErrEmpty = errors.New("Empty string")

type service struct{}

func (service) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (service) Count(s string) int {
	return len(s)
}

func NewService() Service {
	return &service{}
}
