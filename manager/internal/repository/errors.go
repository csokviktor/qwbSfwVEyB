package repository

import "fmt"

type NotFoundError struct{}

func (nf NotFoundError) Error() string {
	return "not found"
}

func ErrNotFoundByID(id string) error {
	return fmt.Errorf("instance %w with %s is", NotFoundError{}, id)
}
