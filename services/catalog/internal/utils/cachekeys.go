package utils

import (
	"strconv"
	"strings"

	"github.com/ilyas/flower/services/catalog/internal/dto"
)

func BuildCategoriesListKey() string {
	return "categories:list"
}

func BuildCategoryKey(id uint64) string {
	return "category:" + strconv.FormatUint(id, 10)
}

func BuildProductKey(id uint64) string {
	return "product:" + strconv.FormatUint(id, 10)
}

func BuildProductsListVersionKey() string {
	return "products:list:version"
}

func BuildProductsListKey(filter dto.ProductFilter, version string) string {
	var b strings.Builder
	b.WriteString("products:list:")
	b.WriteString(version)
	b.WriteString(":")

	if filter.CategoryID != nil {
		b.WriteString("category=")
		b.WriteString(strconv.FormatUint(*filter.CategoryID, 10))
	}
	b.WriteString("|")
	if filter.SellerID != nil {
		b.WriteString("seller=")
		b.WriteString(strconv.FormatUint(*filter.SellerID, 10))
	}
	b.WriteString("|")
	if filter.PriceMin != nil {
		b.WriteString("price_min=")
		b.WriteString(strconv.Itoa(*filter.PriceMin))
	}
	b.WriteString("|")
	if filter.PriceMax != nil {
		b.WriteString("price_max=")
		b.WriteString(strconv.Itoa(*filter.PriceMax))
	}
	b.WriteString("|")
	if filter.Size != nil {
		b.WriteString("size=")
		b.WriteString(*filter.Size)
	}
	b.WriteString("|")
	if filter.IsAvailable != nil {
		b.WriteString("available=")
		b.WriteString(strconv.FormatBool(*filter.IsAvailable))
	}
	b.WriteString("|")
	if filter.Page != nil {
		b.WriteString("page=")
		b.WriteString(strconv.Itoa(*filter.Page))
	}
	b.WriteString("|")
	if filter.PageSize != nil {
		b.WriteString("page_size=")
		b.WriteString(strconv.Itoa(*filter.PageSize))
	}

	return b.String()
}
