package dto

import "mime/multipart"

type SizePrice struct {
	Size      string `json:"size"`
	BasePrice int    `json:"base_price"`
}

type CompositionItem struct {
	FlowerType string `json:"flower_type"`
	Stems      int    `json:"stems"`
}

type Discount struct {
	Type     string `json:"type"`
	Value    int    `json:"value"`
	StartsAt string `json:"starts_at,omitempty"`
	EndsAt   string `json:"ends_at,omitempty"`
}

type Product struct {
	ID           uint64            `json:"id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	CategoryID   uint64            `json:"category_id"`
	SellerID     uint64            `json:"seller_id"`
	MainImageURL string            `json:"main_image_url"`
	Images       []string          `json:"images"`
	IsAvailable  bool              `json:"is_available"`
	Currency     string            `json:"currency"`
	Sizes        []SizePrice       `json:"sizes"`
	PricePerStem int               `json:"price_per_stem"`
	MinStems     int               `json:"min_stems"`
	MaxStems     int               `json:"max_stems"`
	Composition  []CompositionItem `json:"composition"`
	Discount     *Discount         `json:"discount,omitempty"`
	CreatedAt    string            `json:"created_at"`
	UpdatedAt    string            `json:"updated_at"`
}

type CreateProductRequest struct {
	Product         Product                  `json:"-"`
	MainImage       multipart.File           `json:"-"`
	MainImageHeader *multipart.FileHeader    `json:"-"`
	Images          []multipart.File         `json:"-"`
	ImageHeaders    []*multipart.FileHeader  `json:"-"`
}

type UpdateProductRequest struct {
	Product         Product                 `json:"-"`
	MainImage       multipart.File          `json:"-"`
	MainImageHeader *multipart.FileHeader   `json:"-"`
	Images          []multipart.File        `json:"-"`
	ImageHeaders    []*multipart.FileHeader `json:"-"`
}

type CreateCategoryRequest struct {
	Name        string                `json:"name"`
	Slug        string                `json:"slug"`
	Description string                `json:"description,omitempty"`
	Image       multipart.File        `json:"-"`
	ImageHeader *multipart.FileHeader `json:"-"`
}

type UpdateCategoryRequest struct {
	Name        string                `json:"name"`
	Slug        string                `json:"slug"`
	Description string                `json:"description,omitempty"`
	Image       multipart.File        `json:"-"`
	ImageHeader *multipart.FileHeader `json:"-"`
}

type Category struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type ProductFilter struct {
	CategoryID  *uint64
	SellerID    *uint64
	PriceMin    *int
	PriceMax    *int
	Size        *string
	IsAvailable *bool
	Page        *int
	PageSize    *int
}

type PaginatedProducts struct {
	Items    []Product `json:"items"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
}
