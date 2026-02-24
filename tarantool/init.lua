local log = require('log')

box.cfg{
    listen = 3301,
    memtx_memory = 256 * 1024 * 1024,
    log_level = 5,
}

log.info("Starting Tarantool...")

box.schema.user.create('app', { 
    password = 'secret',
    if_not_exists = true
})

box.schema.user.grant('app', 'read,write,execute', 'universe', nil, { if_not_exists = true })

box.once("users_space_init", function()
    log.info("Creating roles space")
    local r = box.schema.space.create("roles")
    
    log.info("Creating users space")
    local u = box.schema.space.create("users")
    
    log.info("Creating auths space")
    local a = box.schema.space.create("auths")

    r:format({
        {name="id",          type="unsigned"},
        {name="name",        type="string"},
        {name="description", type="string", is_nullable=true},
    })
    box.schema.sequence.create("roles_seq", { if_not_exists = true })
    r:create_index("primary", {
        parts = {"id"},
        sequence = "roles_seq"
    })
    r:create_index("name", {
        parts = {"name"},
        unique = true
    })
    log.info("Roles space created")

    u:format({
        {name="id",         type="unsigned"},
        {name="first_name", type="string"},
        {name="last_name",  type="string"},
        {name="role_id",    type="unsigned"},
        {name="version",    type="unsigned"},
        {name="is_active",  type="boolean"},
        {name="created_at", type="unsigned"},
        {name="updated_at", type="unsigned"},
        {name="deleted_at", type="unsigned", is_nullable=true},
    })
    box.schema.sequence.create("users_seq", { if_not_exists = true })
    u:create_index("primary", {
        parts = {"id"},
        sequence = "users_seq"
    })
    a:create_index("role_id", {
        parts = {"role_id"},
        unique = false
    })
    log.info("Users space created")

    a:format({
        {name="id",             type="unsigned"},
        {name="user_id",        type="unsigned"},
        {name="phone_number",   type="string"},
        {name="email",          type="string", is_nullable=true},
        {name="provider",       type="string", is_nullable=true},
        {name="external_id",    type="string", is_nullable=true},
        {name="password_hash",  type="string", is_nullable=true},
        {name="updated_at",     type="unsigned"},
        {name="deleted_at",     type="unsigned", is_nullable=true},
    })
    box.schema.sequence.create("auths_seq", { if_not_exists = true })
    a:create_index("primary", {
        parts = {"id"},
        sequence = "auths_seq"
    })
    a:create_index("phone_number", {
        parts = {"phone_number"},
        unique = true
    })
    a:create_index("email", {
        parts = {"email"},
        unique = false
    })
    log.info("Auths space created")
    if not r.index.name:get{"user"} then
        r:insert{1, "user", "default role"}
    end
    if not r.index.name:get{"seller"} then
        r:insert{2, "seller", "shop seller"}
    end
end)

function create_user(first_name, last_name, phone_number, roles)
    role_name = role_name or "user"
    local now = os.time()
    
    return box.atomic(function()
        if box.space.auths.index.phone_number:get{phone_number} then
            error("PHONE_ALREADY_EXISTS")
        end

        local role = box.space.roles.index.name:get{role_name}
        if role == nil then
            error("ROLE_NOT_FOUND")
        end

        local user = box.space.users:auto_increment{
            first_name, last_name, role.id, 1, false, now, now, box.NULL
        }

        box.space.auths:auto_increment{
            user.id, phone_number, box.NULL, "phone", box.NULL, box.NULL, now, box.NULL
        }

        return {user.id, user.first_name, user.last_name, role.name, now}
    end)
end

log.info("Tarantool started 🚀")
