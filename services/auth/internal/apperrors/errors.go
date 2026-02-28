package apperrors

import "errors"

var (
	ErrNotFound      = errors.New("запись не найдена")
	ErrDuplicate     = errors.New("запись уже существует")
	ErrUnauthorized  = errors.New("не авторизован")
	ErrInvalidInput  = errors.New("неверные входные данные")
	ErrAlreadyActive = errors.New("пользователь уже активирован")
)
