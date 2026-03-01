package utils

import (
	"strings"

	"github.com/ilyas/flower/services/auth/internal/apperrors"
)

type DBRule struct {
	DBerr string
	Err   error
}

var tarantoolRules = []DBRule{
	{DBerr: "PHONE_ALREADY_EXISTS", Err: apperrors.ErrDuplicatePhone},
	{DBerr: "ROLE_NOT_FOUND", Err: apperrors.ErrRoleNotFound},
	{DBerr: "ACCOUNT_NOT_FOUND", Err: apperrors.ErrAccountNotFound},
	{DBerr: "USER_NOT_FOUND", Err: apperrors.ErrUserNotFound},
	{DBerr: "ALREADY_ACTIVE", Err: apperrors.ErrInvalidInput},
}

func MapDBError(err error, rules []DBRule) error {
	if err == nil {
		return nil
	}

	msg := err.Error()
	for _, r := range rules {
		if strings.Contains(msg, r.DBerr) {
			return r.Err
		}
	}
	return err
}

func MapTarantoolError(err error) error {
	return MapDBError(err, tarantoolRules)
}
