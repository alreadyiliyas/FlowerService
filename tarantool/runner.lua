-- /c:/Users/ILYAS/GolandProjects/Flower/tarantool/runner.lua
local M = {}

local function ensure_migrations_space()
    local s = box.schema.space.create("schema_migrations", { if_not_exists = true })
    s:format({
        {name="id", type="string"},
        {name="applied_at", type="unsigned"},
    })
    s:create_index("primary", { parts={"id"}, if_not_exists=true })
    return s
end

function M.run(migrations)
    local ms = ensure_migrations_space()
    for _, m in ipairs(migrations) do
        if not ms:get{m.id} then
            m.up()
            ms:replace{m.id, os.time()}
        end
    end
end

return M
