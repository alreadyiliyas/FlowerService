package apperrors

import "errors"

// Not found
var (
	ErrNotFound        = errors.New("запись не найдена")
	ErrRoleNotFound    = errors.New("роль не найдена")
	ErrAccountNotFound = errors.New("пользователь не найден по номеру телефона")
	ErrUserNotFound    = errors.New("пользователь не найден")
)

// Duplicate
var (
	ErrDuplicate      = errors.New("запись уже существует")
	ErrDuplicatePhone = errors.New("номер телефона уже существует")
)

// Others
var (
	ErrUnauthorized  = errors.New("не авторизован")
	ErrInvalidInput  = errors.New("неверные входные данные")
	ErrAlreadyActive = errors.New("пользователь уже активирован")
	ErrDB            = errors.New("непредвиденная ошибка сервера")
)
