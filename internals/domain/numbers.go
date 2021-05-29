package domain

import (
	"errors"
)

const lineLength = 9

// ErrInvalidNumber is used when the value provided to the Number constructor is invalid
var ErrInvalidNumber = errors.New("invalid number")

type Number struct {
	Value string
}

func NewNumber(line string) (Number, error) {
	if len(line) != lineLength {
		return Number{}, ErrInvalidNumber
	}

	for _, r := range line {
		switch r {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			// Nothing to be done
		default:
			return Number{}, ErrInvalidNumber
		}
	}

	return Number{Value: line}, nil
}
