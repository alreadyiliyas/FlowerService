package utils

func BuildConfirmKey(phone string) string {
	return "confirm:" + phone
}

func BuildPasswordUpdateKey(phone string) string {
	return "pwd_update:" + phone
}

func BuildRefreshTokenKey(refreshKey string) string {
	return "refresh:" + refreshKey
}

func BuildRefreshTokenKeyByPhone(phone string) string {
	return "refresh_phone:" + phone
}
