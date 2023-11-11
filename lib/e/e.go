package e

import "fmt"

func WrapIfErr(msg string, err error) error {
	if err == nil {
		return nil
	}
	return Wrap(msg, err)
}

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
