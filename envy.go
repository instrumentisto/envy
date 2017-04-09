package envy

import (
	"errors"
)

type Parser struct{}

func (p Parser) Parse(obj interface{}) error {
	return errors.New("unimplemented")
}
