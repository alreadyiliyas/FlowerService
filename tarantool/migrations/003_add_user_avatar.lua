return {
    id = "003_add_user_avatar",
    up = function()
        local u = box.space.users
        u:format({
            {name="id", type="unsigned"},
            {name="first_name", type="string"},
            {name="last_name", type="string"},
            {name="role_id", type="unsigned"},
            {name="version", type="unsigned"},
            {name="is_active", type="boolean"},
            {name="created_at", type="unsigned"},
            {name="updated_at", type="unsigned"},
            {name="avatar_url", type="string", is_nullable=true},
            {name="deleted_at", type="unsigned", is_nullable=true},
        })
    end
}
