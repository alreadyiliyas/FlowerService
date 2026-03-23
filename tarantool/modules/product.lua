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

local function validate_product(name, category_id, seller_id, currency, main_image_url, images, sizes, price_per_stem, min_stems, max_stems, composition)
    if name == nil or name == "" then
        error(E.NAME_IS_NULL)
    end

    if category_id == nil or box.space.categories:get(category_id) == nil then
        error(E.CATEGORY_NOT_FOUND)
    end

    if seller_id == nil or currency == nil or currency == "" or main_image_url == nil or main_image_url == "" then
        error(E.PRODUCT_INVALID)
    end

    if images == nil or #images == 0 or sizes == nil or #sizes == 0 or composition == nil or #composition == 0 then
        error(E.PRODUCT_INVALID)
    end

    if price_per_stem == nil or min_stems == nil or max_stems == nil or min_stems > max_stems then
        error(E.PRODUCT_INVALID)
    end
end

local function has_size(product, required_size)
    if required_size == nil or required_size == "" then
        return true
    end

    for _, item in ipairs(product.sizes) do
        if item.size == required_size then
            return true
        end
    end

    return false
end

local function resolve_price(product, required_size)
    local selected = nil

    for _, item in ipairs(product.sizes) do
        if required_size ~= nil and required_size ~= "" and item.size == required_size then
            selected = item.base_price
            break
        end

        if selected == nil or item.base_price < selected then
            selected = item.base_price
        end
    end

    if selected == nil then
        return nil
    end

    return selected + (product.price_per_stem * product.min_stems)
end

function P.create_product(name, description, category_id, seller_id, is_available, currency, main_image_url, images, sizes, price_per_stem, min_stems, max_stems, composition, discount)
    local now = os.time()

    return box.atomic(function()
        validate_product(name, category_id, seller_id, currency, main_image_url, images, sizes, price_per_stem, min_stems, max_stems, composition)

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

function P.list_products(filter)
    filter = filter or {}

    local page = tonumber(filter.page) or 1
    local page_size = tonumber(filter.page_size) or 20
    local offset = (page - 1) * page_size
    local total = 0
    local items = {}

    for _, product in box.space.products.index.primary:pairs() do
        local include = true

        if filter.category_id ~= nil and product.category_id ~= filter.category_id then
            include = false
        end

        if include and filter.seller_id ~= nil and product.seller_id ~= filter.seller_id then
            include = false
        end

        if include and filter.is_available ~= nil and product.is_available ~= filter.is_available then
            include = false
        end

        if include and not has_size(product, filter.size) then
            include = false
        end

        if include and (filter.price_min ~= nil or filter.price_max ~= nil) then
            local price = resolve_price(product, filter.size)
            if price == nil then
                include = false
            elseif filter.price_min ~= nil and price < filter.price_min then
                include = false
            elseif filter.price_max ~= nil and price > filter.price_max then
                include = false
            end
        end

        if include then
            total = total + 1
            if total > offset and #items < page_size then
                table.insert(items, product_to_map(product))
            end
        end
    end

    return {
        items = items,
        total = total,
        page = page,
        page_size = page_size,
    }
end

function P.get_product(id)
    if id == nil then
        error(E.PRODUCT_NOT_FOUND)
    end

    local product = box.space.products:get(id)
    if not product then
        error(E.PRODUCT_NOT_FOUND)
    end

    return product_to_map(product)
end

function P.update_product(id, name, description, category_id, seller_id, is_available, currency, main_image_url, images, sizes, price_per_stem, min_stems, max_stems, composition, discount, version)
    local now = os.time()

    return box.atomic(function()
        local product = box.space.products:get(id)
        if not product then
            error(E.PRODUCT_NOT_FOUND)
        end

        validate_product(name, category_id, seller_id, currency, main_image_url, images, sizes, price_per_stem, min_stems, max_stems, composition)

        if version == nil or version ~= product.version then
            error(E.VERSION_MISMATCH)
        end

        if description == "" then
            description = box.NULL
        end
        if discount == nil then
            discount = box.NULL
        end

        local updated = box.space.products:update(product.id, {
            {"=", "name", name},
            {"=", "description", description},
            {"=", "category_id", category_id},
            {"=", "seller_id", seller_id},
            {"=", "is_available", is_available},
            {"=", "currency", currency},
            {"=", "main_image_url", main_image_url},
            {"=", "images", images},
            {"=", "sizes", sizes},
            {"=", "price_per_stem", price_per_stem},
            {"=", "min_stems", min_stems},
            {"=", "max_stems", max_stems},
            {"=", "composition", composition},
            {"=", "discount", discount},
            {"=", "version", product.version + 1},
            {"=", "updated_at", now},
        })

        return product_to_map(updated)
    end)
end

function P.delete_product(id)
    return box.atomic(function()
        local product = box.space.products:get(id)
        if not product then
            error(E.PRODUCT_NOT_FOUND)
        end

        box.space.products:delete(id)
        return true
    end)
end

return P
