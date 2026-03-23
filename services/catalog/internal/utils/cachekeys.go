package utils

import "strconv"

func BuildCategoriesListKey() string {
	return "categories:list"
}

func BuildCategoryKey(id uint64) string {
	return "category:" + strconv.FormatUint(id, 10)
}
