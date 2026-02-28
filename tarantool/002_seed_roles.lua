return {
    id = "002_seed_roles",
    up = function()
        local r = box.space.roles
        if not r.index.name:get{"user"} then
            r:insert{1, "user", "default role"}
        end
        if not r.index.name:get{"seller"} then
            r:insert{2, "seller", "shop seller"}
        end
    end
}
