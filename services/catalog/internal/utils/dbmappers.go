package utils

import (
	"strings"

	"github.com/ilyas/flower/services/catalog/internal/apperrors"
)

type DBRule struct {
	DBerr string
	Err   error
}

var tarantoolRules = []DBRule{
	{DBerr: "NAME_IS_NULL", Err: apperrors.ErrInvalidInput},
	{DBerr: "SLUG_IS_NULL", Err: apperrors.ErrInvalidInput},
	{DBerr: "NAME_ALREADY_EXISTS", Err: apperrors.ErrDuplicateCategoryName},
	{DBerr: "SLUG_ALREADY_EXISTS", Err: apperrors.ErrDuplicateCategorySlug},
	{DBerr: "CATEGORY_NOT_FOUND", Err: apperrors.ErrNotFound},
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
