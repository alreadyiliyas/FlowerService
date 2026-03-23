local E = require("errors.errors")
local log = require("log")

local C = {}

local function category_to_map(category)
    return {
        id = category.id,
        name = category.name,
        slug = category.slug,
        description = category.description,
        image_url = category.image_url,
        created_at = category.created_at,
        updated_at = category.updated_at,
    }
end

function C.create_category(name, slug, description, image_url)
    local now = os.time()

    return box.atomic(function()
        if name == nil or name == "" then
            log.info("name is required parameter")
            error(E.NAME_IS_NULL)
        end

        if slug == nil or slug == "" then
            log.info("slug is required parameter")
            error(E.SLUG_IS_NULL)
        end

        if not string.match(slug, "^[a-z0-9-]+$") then
            log.info("slug format invalid")
            error(E.SLUG_IS_NULL)
        end

        if box.space.categories.index.name:get{name} then
            log.info("name_already_exists")
            error(E.NAME_ALREADY_EXISTS)
        end

        if box.space.categories.index.slug:get{slug} then
            log.info("slug_already_exists")
            error(E.SLUG_ALREADY_EXISTS)
        end

        if description == "" then
            description = box.NULL
        end
        if image_url == "" then
            image_url = box.NULL
        end

        local category = box.space.categories:auto_increment{
            name, slug, description, image_url, now, now
        }

        return category_to_map(category)
    end)
end

function C.list_categories()
    local result = {}

    for _, category in box.space.categories.index.primary:pairs() do
        table.insert(result, category_to_map(category))
    end

    return result
end

function C.get_category(id)
    if id == nil then
        error(E.CATEGORY_NOT_FOUND)
    end

    local category = box.space.categories:get(id)
    if not category then
        error(E.CATEGORY_NOT_FOUND)
    end

    return category_to_map(category)
end

function C.update_category(id, name, slug, description, image_url)
    local now = os.time()

    return box.atomic(function()
        local category = box.space.categories:get(id)
        if not category then
            error(E.CATEGORY_NOT_FOUND)
        end

        if name == nil or name == "" then
            error(E.NAME_IS_NULL)
        end
        if slug == nil or slug == "" then
            error(E.SLUG_IS_NULL)
        end
        if not string.match(slug, "^[a-z0-9-]+$") then
            error(E.SLUG_IS_NULL)
        end

        local found_name = box.space.categories.index.name:get{name}
        if found_name and found_name.id ~= category.id then
            error(E.NAME_ALREADY_EXISTS)
        end

        local found_slug = box.space.categories.index.slug:get{slug}
        if found_slug and found_slug.id ~= category.id then
            error(E.SLUG_ALREADY_EXISTS)
        end

        if description == "" then
            description = box.NULL
        end
        if image_url == "" then
            image_url = box.NULL
        end

        local updated = box.space.categories:update(category.id, {
            {"=", "name", name},
            {"=", "slug", slug},
            {"=", "description", description},
            {"=", "image_url", image_url},
            {"=", "updated_at", now},
        })

        return category_to_map(updated)
    end)
end

function C.delete_category(id)
    return box.atomic(function()
        local category = box.space.categories:get(id)
        if not category then
            error(E.CATEGORY_NOT_FOUND)
        end

        box.space.categories:delete(id)
        return true
    end)
end

return C
