package utils

import "github.com/ilyas/flower/services/catalog/internal/entities"

// BuildImagesArgs преобразует список URL изображений в формат аргументов,
// который ожидает Tarantool при создании или обновлении продукта.
func BuildImagesArgs(images []string) []interface{} {
	result := make([]interface{}, 0, len(images))
	for _, image := range images {
		result = append(result, image)
	}
	return result
}

// BuildSizesArgs преобразует размеры и цены продукта в массив map-структур
// для передачи в Tarantool как вложенный аргумент.
func BuildSizesArgs(items []entities.SizePrice) []interface{} {
	result := make([]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"size":       *item.Size,
			"base_price": *item.BasePrice,
		})
	}
	return result
}

// BuildCompositionArgs преобразует состав букета в формат, который
// используется в Tarantool для вложенного массива компонентов продукта.
func BuildCompositionArgs(items []entities.CompositionItem) []interface{} {
	result := make([]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"flower_type": *item.FlowerType,
			"stems":       *item.Stems,
		})
	}
	return result
}

// BuildDiscountArg подготавливает скидку к передаче в Tarantool.
// Если скидка не задана, функция возвращает nil.
func BuildDiscountArg(discount *entities.Discount) interface{} {
	if discount == nil {
		return nil
	}

	return map[string]interface{}{
		"type":      *discount.Type,
		"value":     *discount.Value,
		"starts_at": ValueOrNull(discount.StartsAt),
		"ends_at":   ValueOrNull(discount.EndsAt),
	}
}
