package utils

import (
	"fmt"
	"time"

	"github.com/ilyas/flower/services/catalog/internal/dto"
	"github.com/ilyas/flower/services/catalog/internal/entities"
)

// NormalizeTarantoolMap приводит ответ Tarantool к единому виду.
// На практике Tarantool может вернуть map как в формате map[interface{}]interface{},
// так и в формате map[string]interface{}. Эта функция нужна, чтобы код репозиториев
// и парсеров работал с одной и той же структурой данных и не дублировал
// одинаковую логику приведения типов в нескольких местах.
func NormalizeTarantoolMap(raw interface{}) (map[interface{}]interface{}, error) {
	switch row := raw.(type) {
	case map[interface{}]interface{}:
		return row, nil
	case map[string]interface{}:
		out := make(map[interface{}]interface{}, len(row))
		for k, v := range row {
			out[k] = v
		}
		return out, nil
	default:
		return nil, fmt.Errorf("invalid response payload type: %T", raw)
	}
}

// ParseCategoryEntity преобразует одну запись категории из ответа Tarantool
// в entities.Category. Функция извлекает скалярные поля, корректно обрабатывает
// nullable-значения и переводит Unix timestamp в time.Time, который используется
// на уровне доменной сущности.
func ParseCategoryEntity(raw interface{}) (*entities.Category, error) {
	row, err := NormalizeTarantoolMap(raw)
	if err != nil {
		return nil, err
	}

	id, err := ToUint64(row["id"], "id")
	if err != nil {
		return nil, err
	}
	name, err := ToString(row["name"], "name")
	if err != nil {
		return nil, err
	}
	slug, err := ToString(row["slug"], "slug")
	if err != nil {
		return nil, err
	}
	description, err := ToStringNullable(row["description"], "description")
	if err != nil {
		return nil, err
	}
	imageURL, err := ToStringNullable(row["image_url"], "image_url")
	if err != nil {
		return nil, err
	}
	createdUnix, err := ToUint64(row["created_at"], "created_at")
	if err != nil {
		return nil, err
	}
	updatedUnix, err := ToUint64(row["updated_at"], "updated_at")
	if err != nil {
		return nil, err
	}

	createdAt := time.Unix(int64(createdUnix), 0)
	updatedAt := time.Unix(int64(updatedUnix), 0)

	return &entities.Category{
		ID:          &id,
		Name:        &name,
		Slug:        &slug,
		Description: description,
		ImageURL:    imageURL,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
	}, nil
}

// ParseCategoryEntities преобразует массив категорий из ответа Tarantool
// в срез entities.Category. Для каждой записи используется ParseCategoryEntity,
// чтобы правила маппинга хранились в одном месте и одинаково работали
// как для списка, так и для получения одной категории.
func ParseCategoryEntities(raw interface{}) ([]entities.Category, error) {
	rows, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response payload type: %T", raw)
	}

	result := make([]entities.Category, 0, len(rows))
	for _, item := range rows {
		category, err := ParseCategoryEntity(item)
		if err != nil {
			return nil, err
		}
		result = append(result, *category)
	}

	return result, nil
}

// ParseProductDTO преобразует одну запись продукта из ответа Tarantool в dto.Product.
// Здесь разбираются как обычные скалярные поля, так и вложенные структуры:
// изображения, размеры, состав и скидка. Дополнительно функция переводит
// Unix timestamp в строки формата RFC3339, которые ожидает транспортный слой.
func ParseProductDTO(raw interface{}) (*dto.Product, error) {
	row, err := NormalizeTarantoolMap(raw)
	if err != nil {
		return nil, err
	}

	productID, err := ToUint64(row["id"], "id")
	if err != nil {
		return nil, err
	}
	name, err := ToString(row["name"], "name")
	if err != nil {
		return nil, err
	}
	description, err := ToStringNullable(row["description"], "description")
	if err != nil {
		return nil, err
	}
	categoryID, err := ToUint64(row["category_id"], "category_id")
	if err != nil {
		return nil, err
	}
	sellerID, err := ToUint64(row["seller_id"], "seller_id")
	if err != nil {
		return nil, err
	}
	isAvailable, err := ToBool(row["is_available"], "is_available")
	if err != nil {
		return nil, err
	}
	currency, err := ToString(row["currency"], "currency")
	if err != nil {
		return nil, err
	}
	mainImageURL, err := ToString(row["main_image_url"], "main_image_url")
	if err != nil {
		return nil, err
	}
	images, err := parseImagesValue(row["images"])
	if err != nil {
		return nil, err
	}
	sizes, err := parseSizesValue(row["sizes"])
	if err != nil {
		return nil, err
	}
	pricePerStem, err := ToInt(row["price_per_stem"], "price_per_stem")
	if err != nil {
		return nil, err
	}
	minStems, err := ToInt(row["min_stems"], "min_stems")
	if err != nil {
		return nil, err
	}
	maxStems, err := ToInt(row["max_stems"], "max_stems")
	if err != nil {
		return nil, err
	}
	composition, err := parseCompositionValue(row["composition"])
	if err != nil {
		return nil, err
	}
	discount, err := parseDiscountValue(row["discount"])
	if err != nil {
		return nil, err
	}
	createdUnix, err := ToUint64(row["created_at"], "created_at")
	if err != nil {
		return nil, err
	}
	updatedUnix, err := ToUint64(row["updated_at"], "updated_at")
	if err != nil {
		return nil, err
	}

	return &dto.Product{
		ID:           productID,
		Name:         name,
		Description:  stringOrEmpty(description),
		CategoryID:   categoryID,
		SellerID:     sellerID,
		MainImageURL: mainImageURL,
		Images:       images,
		IsAvailable:  isAvailable,
		Currency:     currency,
		Sizes:        sizes,
		PricePerStem: pricePerStem,
		MinStems:     minStems,
		MaxStems:     maxStems,
		Composition:  composition,
		Discount:     discount,
		CreatedAt:    time.Unix(int64(createdUnix), 0).Format(time.RFC3339),
		UpdatedAt:    time.Unix(int64(updatedUnix), 0).Format(time.RFC3339),
	}, nil
}

// ParseProductEntity преобразует одну запись продукта из ответа Tarantool
// в entities.Product. Эта функция нужна для слоя репозитория, который должен
// работать с доменными сущностями и не зависеть от transport DTO.
func ParseProductEntity(raw interface{}) (*entities.Product, error) {
	row, err := NormalizeTarantoolMap(raw)
	if err != nil {
		return nil, err
	}

	productID, err := ToUint64(row["id"], "id")
	if err != nil {
		return nil, err
	}
	name, err := ToString(row["name"], "name")
	if err != nil {
		return nil, err
	}
	description, err := ToStringNullable(row["description"], "description")
	if err != nil {
		return nil, err
	}
	categoryID, err := ToUint64(row["category_id"], "category_id")
	if err != nil {
		return nil, err
	}
	sellerID, err := ToUint64(row["seller_id"], "seller_id")
	if err != nil {
		return nil, err
	}
	isAvailable, err := ToBool(row["is_available"], "is_available")
	if err != nil {
		return nil, err
	}
	currency, err := ToString(row["currency"], "currency")
	if err != nil {
		return nil, err
	}
	mainImageURL, err := ToString(row["main_image_url"], "main_image_url")
	if err != nil {
		return nil, err
	}
	images, err := parseImagesValue(row["images"])
	if err != nil {
		return nil, err
	}
	sizesDTO, err := parseSizesValue(row["sizes"])
	if err != nil {
		return nil, err
	}
	pricePerStem, err := ToUint64(row["price_per_stem"], "price_per_stem")
	if err != nil {
		return nil, err
	}
	minStems, err := ToUint64(row["min_stems"], "min_stems")
	if err != nil {
		return nil, err
	}
	maxStems, err := ToUint64(row["max_stems"], "max_stems")
	if err != nil {
		return nil, err
	}
	compositionDTO, err := parseCompositionValue(row["composition"])
	if err != nil {
		return nil, err
	}
	discountDTO, err := parseDiscountValue(row["discount"])
	if err != nil {
		return nil, err
	}
	createdUnix, err := ToUint64(row["created_at"], "created_at")
	if err != nil {
		return nil, err
	}
	updatedUnix, err := ToUint64(row["updated_at"], "updated_at")
	if err != nil {
		return nil, err
	}
	version, err := ToUint64(row["version"], "version")
	if err != nil {
		return nil, err
	}

	sizes := make([]entities.SizePrice, 0, len(sizesDTO))
	for _, item := range sizesDTO {
		size := item.Size
		basePrice := uint64(item.BasePrice)
		sizes = append(sizes, entities.SizePrice{
			Size:      &size,
			BasePrice: &basePrice,
		})
	}

	composition := make([]entities.CompositionItem, 0, len(compositionDTO))
	for _, item := range compositionDTO {
		flowerType := item.FlowerType
		stems := uint64(item.Stems)
		composition = append(composition, entities.CompositionItem{
			FlowerType: &flowerType,
			Stems:      &stems,
		})
	}

	var discount *entities.Discount
	if discountDTO != nil {
		discountType := discountDTO.Type
		value := uint64(discountDTO.Value)
		discount = &entities.Discount{
			Type:     &discountType,
			Value:    &value,
			StartsAt: stringOrNil(discountDTO.StartsAt),
			EndsAt:   stringOrNil(discountDTO.EndsAt),
		}
	}

	createdAt := time.Unix(int64(createdUnix), 0)
	updatedAt := time.Unix(int64(updatedUnix), 0)

	return &entities.Product{
		ID:           &productID,
		Name:         &name,
		Description:  description,
		CategoryID:   &categoryID,
		SellerID:     &sellerID,
		IsAvailable:  isAvailable,
		Currency:     &currency,
		MainImageURL: &mainImageURL,
		Images:       images,
		Sizes:        sizes,
		PricePerStem: &pricePerStem,
		MinStems:     &minStems,
		MaxStems:     &maxStems,
		Composition:  composition,
		Discount:     discount,
		Version:      version,
		CreatedAt:    &createdAt,
		UpdatedAt:    &updatedAt,
	}, nil
}

// ParsePaginatedProductEntities преобразует paginated-ответ Tarantool
// со списком продуктов и метаданными пагинации в доменную структуру.
func ParsePaginatedProductEntities(raw interface{}) (entities.PaginatedProducts, error) {
	row, err := NormalizeTarantoolMap(raw)
	if err != nil {
		return entities.PaginatedProducts{}, err
	}

	itemsRaw, ok := row["items"].([]interface{})
	if !ok {
		return entities.PaginatedProducts{}, fmt.Errorf("invalid items type: %T", row["items"])
	}

	total, err := ToInt(row["total"], "total")
	if err != nil {
		return entities.PaginatedProducts{}, err
	}
	page, err := ToInt(row["page"], "page")
	if err != nil {
		return entities.PaginatedProducts{}, err
	}
	pageSize, err := ToInt(row["page_size"], "page_size")
	if err != nil {
		return entities.PaginatedProducts{}, err
	}

	items := make([]entities.Product, 0, len(itemsRaw))
	for _, item := range itemsRaw {
		product, err := ParseProductEntity(item)
		if err != nil {
			return entities.PaginatedProducts{}, err
		}
		items = append(items, *product)
	}

	return entities.PaginatedProducts{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// parseImagesValue преобразует массив изображений из ответа Tarantool
// в срез строк с URL. Каждый элемент ожидается строкой. Если структура ответа
// не соответствует ожидаемому контракту, функция сразу возвращает ошибку.
func parseImagesValue(raw interface{}) ([]string, error) {
	rows, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid images type: %T", raw)
	}

	result := make([]string, 0, len(rows))
	for _, item := range rows {
		image, err := ToString(item, "image")
		if err != nil {
			return nil, err
		}
		result = append(result, image)
	}
	return result, nil
}

// parseSizesValue преобразует массив размеров и цен из ответа Tarantool
// в DTO-структуры. Каждый вложенный элемент сначала нормализуется через
// NormalizeTarantoolMap, чтобы функция могла безопасно читать поля независимо
// от конкретного формата map, который вернул драйвер.
func parseSizesValue(raw interface{}) ([]dto.SizePrice, error) {
	rows, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid sizes type: %T", raw)
	}

	result := make([]dto.SizePrice, 0, len(rows))
	for _, item := range rows {
		row, err := NormalizeTarantoolMap(item)
		if err != nil {
			return nil, fmt.Errorf("invalid size item type: %T", item)
		}

		size, err := ToString(row["size"], "size")
		if err != nil {
			return nil, err
		}
		basePrice, err := ToInt(row["base_price"], "base_price")
		if err != nil {
			return nil, err
		}

		result = append(result, dto.SizePrice{Size: size, BasePrice: basePrice})
	}
	return result, nil
}

// parseCompositionValue преобразует массив состава букета из ответа Tarantool
// в DTO-структуры. Функция проверяет формат каждого вложенного элемента
// и извлекает тип цветка и количество стеблей, которые нужны в контракте продукта.
func parseCompositionValue(raw interface{}) ([]dto.CompositionItem, error) {
	rows, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid composition type: %T", raw)
	}

	result := make([]dto.CompositionItem, 0, len(rows))
	for _, item := range rows {
		row, err := NormalizeTarantoolMap(item)
		if err != nil {
			return nil, fmt.Errorf("invalid composition item type: %T", item)
		}

		flowerType, err := ToString(row["flower_type"], "flower_type")
		if err != nil {
			return nil, err
		}
		stems, err := ToInt(row["stems"], "stems")
		if err != nil {
			return nil, err
		}

		result = append(result, dto.CompositionItem{FlowerType: flowerType, Stems: stems})
	}
	return result, nil
}

// parseDiscountValue преобразует поле скидки из ответа Tarantool в dto.Discount.
// Поле скидки может быть nullable, поэтому в случае отсутствия значения
// функция возвращает nil. Если скидка присутствует, дополнительно маппятся
// необязательные поля starts_at и ends_at.
func parseDiscountValue(raw interface{}) (*dto.Discount, error) {
	if raw == nil {
		return nil, nil
	}

	row, err := NormalizeTarantoolMap(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid discount type: %T", raw)
	}

	discountType, err := ToString(row["type"], "discount.type")
	if err != nil {
		return nil, err
	}
	value, err := ToInt(row["value"], "discount.value")
	if err != nil {
		return nil, err
	}

	startsAt := ""
	if rawStartsAt, ok := row["starts_at"]; ok && rawStartsAt != nil {
		startsAt, err = ToString(rawStartsAt, "discount.starts_at")
		if err != nil {
			return nil, err
		}
	}
	endsAt := ""
	if rawEndsAt, ok := row["ends_at"]; ok && rawEndsAt != nil {
		endsAt, err = ToString(rawEndsAt, "discount.ends_at")
		if err != nil {
			return nil, err
		}
	}

	return &dto.Discount{
		Type:     discountType,
		Value:    value,
		StartsAt: startsAt,
		EndsAt:   endsAt,
	}, nil
}

func stringOrEmpty(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func stringOrNil(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
