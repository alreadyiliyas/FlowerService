package utils

func BuildConfirmKey(phone string) string {
	return "confirm:" + phone
}

func BuildPasswordUpdateKey(phone string) string {
	return "pwd_update:" + phone
}
