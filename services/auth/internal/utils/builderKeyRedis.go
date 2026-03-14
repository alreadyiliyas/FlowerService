package utils

func BuildConfirmKey(phone string) string {
	return "confirm:" + phone
}

func BuildPasswordUpdateKey(phone string) string {
	return "pwd_update:" + phone
}

func BuildSessionKey(sessionID string) string {
	return "session:" + sessionID
}

func BuildSessionKeyByPhone(phone string) string {
	return "session_phone:" + phone
}

func BuildUserInfoKey(phone string) string {
	return "user_info:" + phone
}
