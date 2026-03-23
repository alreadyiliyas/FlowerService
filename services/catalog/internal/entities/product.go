package entities

import "time"

type SizePrice struct {
	Size      *string `json:"size"`
	BasePrice *uint64 `json:"base_price"`
}

type CompositionItem struct {
	FlowerType *string `json:"flower_type"`
	Stems      *uint64 `json:"stems"`
}

type Discount struct {
	Type     *string `json:"type"`
	Value    *uint64 `json:"value"`
	StartsAt *string `json:"starts_at,omitempty"`
	EndsAt   *string `json:"ends_at,omitempty"`
}

type Product struct {
	ID           *uint64            `json:"id"`
	Name         *string            `json:"name"`
	Description  *string            `json:"description,omitempty"`
	CategoryID   *uint64            `json:"category_id"`
	SellerID     *uint64            `json:"seller_id"`
	IsAvailable  bool               `json:"is_available"`
	Currency     *string            `json:"currency"`
	MainImageURL *string            `json:"main_image_url"`
	Images       []string           `json:"images"`
	Sizes        []SizePrice        `json:"sizes"`
	PricePerStem *uint64            `json:"price_per_stem"`
	MinStems     *uint64            `json:"min_stems"`
	MaxStems     *uint64            `json:"max_stems"`
	Composition  []CompositionItem  `json:"composition"`
	Discount     *Discount          `json:"discount,omitempty"`
	Version      uint64             `json:"version"`
	CreatedAt    *time.Time         `json:"created_at,omitempty"`
	UpdatedAt    *time.Time         `json:"updated_at,omitempty"`
	DeletedAt    *time.Time         `json:"deleted_at,omitempty"`
}
