return {
    id = "001_init_spaces",
    up = function()
        local r = box.schema.space.create("roles", { if_not_exists = true })
        r:format({
            {name="id", type="unsigned"},
            {name="name", type="string"},
            {name="description", type="string", is_nullable=true},
        })
        box.schema.sequence.create("roles_seq", { if_not_exists = true })
        r:create_index("primary", { parts={{field="id", type="unsigned"}}, sequence="roles_seq", if_not_exists=true })
        r:create_index("name", { parts={{field="name", type="string"}}, unique=true, if_not_exists=true })

        local u = box.schema.space.create("users", { if_not_exists = true })
        local role_idx = u.index.role_id
        if role_idx and role_idx.parts and role_idx.parts[1] and role_idx.parts[1].type ~= "unsigned" then
            role_idx:drop()
        end
        u:format({
            {name="id", type="unsigned"},
            {name="first_name", type="string"},
            {name="last_name", type="string"},
            {name="role_id", type="unsigned"},
            {name="version", type="unsigned"},
            {name="is_active", type="boolean"},
            {name="created_at", type="unsigned"},
            {name="updated_at", type="unsigned"},
            {name="deleted_at", type="unsigned", is_nullable=true},
        })
        box.schema.sequence.create("users_seq", { if_not_exists = true })
        u:create_index("primary", { parts={{field="id", type="unsigned"}}, sequence="users_seq", if_not_exists=true })
        u:create_index("role_id", { parts={{field="role_id", type="unsigned"}}, unique=false, if_not_exists=true })

        local a = box.schema.space.create("auths", { if_not_exists = true })
        a:format({
            {name="id", type="unsigned"},
            {name="user_id", type="unsigned"},
            {name="phone_number", type="string"},
            {name="email", type="string", is_nullable=true},
            {name="provider", type="string", is_nullable=true},
            {name="external_id", type="string", is_nullable=true},
            {name="password_hash", type="string", is_nullable=true},
            {name="updated_at", type="unsigned"},
            {name="deleted_at", type="unsigned", is_nullable=true},
        })
        box.schema.sequence.create("auths_seq", { if_not_exists = true })
        a:create_index("primary", { parts={{field="id", type="unsigned"}}, sequence="auths_seq", if_not_exists=true })
        a:create_index("phone_number", { parts={{field="phone_number", type="string"}}, unique=true, if_not_exists=true })
    end
}
