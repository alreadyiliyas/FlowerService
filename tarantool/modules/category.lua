local E = require("errors.errors")
local log = require("log")

local C = {}

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

        return {
            id = category.id,
            name = category.name,
            slug = category.slug,
            description = category.description,
            image_url = category.image_url,
            created_at = category.created_at,
            updated_at = category.updated_at,
        }
    end)
end

return C
