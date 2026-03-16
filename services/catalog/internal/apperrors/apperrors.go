package apperrors

import "errors"

// Not found
var (
	ErrNotFound = errors.New("Р·Р°РїРёСЃСЊ РЅРµ РЅР°Р№РґРµРЅР°")
)

// Duplicate
var (
	ErrDuplicate      = errors.New("Р·Р°РїРёСЃСЊ СѓР¶Рµ СЃСѓС‰РµСЃС‚РІСѓРµС‚")
	ErrDuplicatePhone = errors.New("РЅРѕРјРµСЂ С‚РµР»РµС„РѕРЅР° СѓР¶Рµ СЃСѓС‰РµСЃС‚РІСѓРµС‚")
)

// Others
var (
	ErrUnauthorized = errors.New("РЅРµ Р°РІС‚РѕСЂРёР·РѕРІР°РЅ")
	ErrForbidden    = errors.New("РґРѕСЃС‚СѓРї Р·Р°РїСЂРµС‰РµРЅ")
	ErrInvalidInput = errors.New("РЅРµРІРµСЂРЅС‹Рµ РІС…РѕРґРЅС‹Рµ РґР°РЅРЅС‹Рµ")
	ErrDB           = errors.New("РЅРµРїСЂРµРґРІРёРґРµРЅРЅР°СЏ РѕС€РёР±РєР° СЃРµСЂРІРµСЂР°")
)
