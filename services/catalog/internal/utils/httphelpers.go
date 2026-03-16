package utils

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ilyas/flower/services/catalog/internal/dto"
)

func ParseUint64Path(r *http.Request, key string) (uint64, error) {
	val := mux.Vars(r)[key]
	return strconv.ParseUint(val, 10, 64)
}

func ParseProductFilter(r *http.Request) dto.ProductFilter {
	q := r.URL.Query()
	filter := dto.ProductFilter{}

	if v := q.Get("category_id"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			filter.CategoryID = &n
		}
	}
	if v := q.Get("seller_id"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			filter.SellerID = &n
		}
	}
	if v := q.Get("price_min"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.PriceMin = &n
		}
	}
	if v := q.Get("price_max"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.PriceMax = &n
		}
	}
	if v := q.Get("size"); v != "" {
		filter.Size = &v
	}
	if v := q.Get("is_available"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			filter.IsAvailable = &b
		}
	}
	if v := q.Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.Page = &n
		}
	}
	if v := q.Get("page_size"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.PageSize = &n
		}
	}

	return filter
}
