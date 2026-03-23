package utils

import (
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ilyas/flower/services/catalog/internal/dto"
)

const CategoryMaxUploadSize int64 = 5 << 20

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

func OpenMultipartFiles(form *multipart.Form, field string) ([]multipart.File, []*multipart.FileHeader, error) {
	if form == nil || form.File == nil {
		return nil, nil, nil
	}

	headers := form.File[field]
	files := make([]multipart.File, 0, len(headers))
	resultHeaders := make([]*multipart.FileHeader, 0, len(headers))

	for _, header := range headers {
		file, err := header.Open()
		if err != nil {
			CloseMultipartFiles(files)
			return nil, nil, err
		}
		files = append(files, file)
		resultHeaders = append(resultHeaders, header)
	}

	return files, resultHeaders, nil
}

func CloseMultipartFiles(files []multipart.File) {
	for _, file := range files {
		_ = file.Close()
	}
}
