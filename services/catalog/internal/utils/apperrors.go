package utils

import "errors"

// Not found
var (
	ErrNotFound = errors.New("запись не найдена")
)

// Duplicate
var (
	ErrDuplicate      = errors.New("запись уже существует")
	ErrDuplicatePhone = errors.New("номер телефона уже существует")
)
