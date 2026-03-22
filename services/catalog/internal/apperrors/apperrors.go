package apperrors

import "errors"

var (
	ErrNotFound             = errors.New("запись не найдена")
	ErrNotFoundCategoryName = errors.New("категория не найдена")
	ErrNotFoundCategorySlug = errors.New("запись 'slug' не найдена")

	ErrInvalidPhoneFormat = errors.New("неверный формат номера телефона")
	ErrCodeIsExpired      = errors.New("код подтверждения истек")

	ErrDuplicate             = errors.New("запись уже существует")
	ErrDuplicateCategoryName = errors.New("такая категория уже существует")
	ErrDuplicateCategorySlug = errors.New("такой slug уже существует")

	ErrSessionRequired = errors.New("session_id обязателен")
	ErrPasswordNull    = errors.New("пароль обязателен")

	ErrUnauthorized = errors.New("не авторизован")
	ErrForbidden    = errors.New("доступ запрещен")

	ErrInvalidInput = errors.New("неверные входные данные")

	ErrDB = errors.New("непредвиденная ошибка сервера")
)
