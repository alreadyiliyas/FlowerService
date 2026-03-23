local E = require("errors.errors")

local P = {}

local function product_to_map(product)
    return {
        id = product.id,
        name = product.name,
        description = product.description,
        category_id = product.category_id,
        seller_id = product.seller_id,
        is_available = product.is_available,
        currency = product.currency,
        main_image_url = product.main_image_url,
        images = product.images,
        sizes = product.sizes,
        price_per_stem = product.price_per_stem,
        min_stems = product.min_stems,
        max_stems = product.max_stems,
        composition = product.composition,
        discount = product.discount,
        version = product.version,
        created_at = product.created_at,
        updated_at = product.updated_at,
        deleted_at = product.deleted_at,
    }
end

function P.create_product(name, description, category_id, seller_id, is_available, currency, main_image_url, images, sizes, price_per_stem, min_stems, max_stems, composition, discount)
    local now = os.time()

    return box.atomic(function()
        if name == nil or name == "" then
            error(E.NAME_IS_NULL)
        end

        if category_id == nil or box.space.categories:get(category_id) == nil then
            error(E.CATEGORY_NOT_FOUND)
        end

        if seller_id == nil or currency == nil or currency == "" or main_image_url == nil or main_image_url == "" then
            error(E.PRODUCT_INVALID)
        end

        if images == nil or sizes == nil or composition == nil then
            error(E.PRODUCT_INVALID)
        end

        if price_per_stem == nil or min_stems == nil or max_stems == nil or min_stems > max_stems then
            error(E.PRODUCT_INVALID)
        end

        if description == "" then
            description = box.NULL
        end
        if discount == nil then
            discount = box.NULL
        end

        local product = box.space.products:auto_increment{
            name,
            description,
            category_id,
            seller_id,
            is_available,
            currency,
            main_image_url,
            images,
            sizes,
            price_per_stem,
            min_stems,
            max_stems,
            composition,
            discount,
            1,
            now,
            now,
            box.NULL,
        }

        return product_to_map(product)
    end)
end

return P
