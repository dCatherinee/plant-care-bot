package domain

import "errors"

var (
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrNotFound           = errors.New("not found")
	ErrPlantNameEmpty     = errors.New("plant name is empty")
	ErrPlantAlreadyExists = errors.New("plant already exists")
	ErrInvalidPlantName   = errors.New("invalid plant name")
)
