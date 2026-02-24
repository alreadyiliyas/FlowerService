package utils

import (
	"errors"

	"regexp"
	"strings"

	"github.com/ilyas/flower/services/auth/internal/dto"
)

// Проверка на валидацию регистрации нового пользвателя
func ValidateAuthData(in dto.RegistrationRequest) error {
	if strings.TrimSpace(in.FirstName) == "" {
		return errors.New("имя не должно быть пустым")
	}
	if len(in.FirstName) < 3 || len(in.FirstName) > 25 {
		return errors.New("имя должно содержать от 2 до 25 символов")
	}
	if strings.TrimSpace(in.LastName) == "" {
		return errors.New("фамилия не  должно быть пустым")
	}
	if len(in.LastName) < 2 || len(in.LastName) > 25 {
		return errors.New("фамилия должна содержать от 2 до 25 символов")
	}
	if !IsValidPhoneNumber(in.PhoneNumber) {
		return errors.New("проверьте номер телефона")
	}
	if in.Role != "admin" && in.Role != "user" && in.Role != "seller" {
		return errors.New("не правильно выбрана роль")
	}
	return nil
}

// Проверка на валидацию номера телефона
func IsValidPhoneNumber(phone string) bool {
	res := regexp.MustCompile(`^\+\d{11,15}$`)
	return res.MatchString(phone)
}

// Проверка на валидацию пароля
func IsValidatePassword(password string) error {
	if password == "" {
		return errors.New("пароль не должен быть пустым")
	}
	if len(password) < 8 || len(password) > 32 {
		return errors.New("пароль должен соблюдать от 8 до 32 символов")
	}

	upperCase := regexp.MustCompile(`[A-Z]`)
	lowerCase := regexp.MustCompile(`[a-z]`)
	digit := regexp.MustCompile(`\d`)
	specialChar := regexp.MustCompile(`[!@#\$%\^&\*]`)

	if !upperCase.MatchString(password) {
		return errors.New("пароль должен содержать хотя бы одну заглавную букву")
	}
	if !lowerCase.MatchString(password) {
		return errors.New("пароль должен содержать хотя бы одну строчную букву")
	}
	if !digit.MatchString(password) {
		return errors.New("пароль должен содержать хотя бы одну цифру")
	}
	if !specialChar.MatchString(password) {
		return errors.New("пароль должен содержать хотя бы один специальный символ (!@#$%^&*)")
	}

	return nil
}
